package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"seno-blackdragon/internal/repository"
	"seno-blackdragon/pkg/enum"
	"seno-blackdragon/pkg/pass"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JWTConfig struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration // e.g. 15 * time.Minute
	RefreshTTL    time.Duration // e.g. 30 * 24 * time.Hour
	Issuer        string        // e.g. "seno-blackdragon"
}

type Claims struct {
	UID       uuid.UUID `json:"uid"`
	Email     string    `json:"email"`
	TokenType string    `json:"typ"` // "access" | "refresh"
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo *repository.UserRepo // interface
	hasher   pass.Hasher          // Argon2id/Bcrypt impl
	jwtCfg   JWTConfig
	log      *zap.Logger
}

func NewAuthService(
	userRepo *repository.UserRepo,
	hasher pass.Hasher,
	jwtCfg JWTConfig,
	log *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		hasher:   hasher,
		jwtCfg:   jwtCfg,
		log:      log,
	}
}

// ===== token helpers =====

func (as *AuthService) makeAccessToken(ctx context.Context, u *repository.UserModel, jti string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(as.jwtCfg.AccessTTL)
	claims := &Claims{
		UID:       u.ID,
		Email:     u.Email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    as.jwtCfg.Issuer,
			Subject:   u.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString(as.jwtCfg.AccessSecret)
	return ss, exp, err
}

func (as *AuthService) makeRefreshToken(ctx context.Context, u *repository.UserModel, jti string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(as.jwtCfg.RefreshTTL)
	claims := &Claims{
		UID:       u.ID,
		Email:     u.Email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    as.jwtCfg.Issuer,
			Subject:   u.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString(as.jwtCfg.RefreshSecret)
	return ss, exp, err
}

// ===== auth flows =====

func (as *AuthService) Register(ctx context.Context, fullName string, bio string, email string, password string) (uuid.UUID, error) {
	u, err := as.userRepo.GetUserByEmail(ctx, email)
	if err == nil && !errors.Is(err, enum.ErrUserNotFound) {
		return uuid.Nil, fmt.Errorf("failed to check existing email: %w", err)
	}

	if u != nil {
		return uuid.Nil, enum.ErrEmailAlready
	}
	hashed, err := as.hasher.Hash(password)
	if err != nil {
		return uuid.Nil, err
	}
	param := &repository.UserModel{
		FullName:     fullName,
		Bio:          bio,
		Email:        email,
		PasswordHash: hashed,
	}
	id, err := as.userRepo.CreateUser(ctx, param)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (as *AuthService) Login(ctx context.Context, email string, password string) (access, refresh string, expired int64, err error) {
	u, err := as.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return "", "", 0, enum.ErrInvalidCredentials
	}
	ok, _ := as.hasher.Verify(password, u.PasswordHash)
	if !ok {
		return "", "", 0, enum.ErrInvalidCredentials
	}

	rtJTI := fmt.Sprintf("rt-%s", u.ID)
	atJTI := fmt.Sprintf("at-%s", u.ID)
	at, atExp, err := as.makeAccessToken(ctx, u, atJTI)
	rt, _, err := as.makeRefreshToken(ctx, u, rtJTI)
	if err != nil {
		return "", "", 0, err
	}

	return at, rt, int64(atExp.Unix()), nil
}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, int64, error) {
	// Parse & validate refresh token
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(as.jwtCfg.Issuer),
	)
	claims := &Claims{}
	tok, err := parser.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, enum.ErrWrongAlgorithm
		}
		return as.jwtCfg.RefreshSecret, nil
	})
	if err != nil || !tok.Valid {
		return "", "", 0, enum.ErrInvalidToken
	}
	if claims.TokenType != "refresh" {
		return "", "", 0, enum.ErrWrongType
	}

	// Check session on Redis
	// rtKey := fmt.Sprintf("RT:%s:%s", claims.UID, claims.ID)

	u, err := as.userRepo.GetUserByID(ctx, claims.UID)
	if err != nil || u == nil {
		return "", "", 0, enum.ErrUserNotFound
	}


	newRTJTI := fmt.Sprintf("rt-%d-%d", u.ID, time.Now().UnixNano())
	at, atExp, err := as.makeAccessToken(ctx, u, fmt.Sprintf("at-%d-%d", u.ID, time.Now().UnixNano()))
	if err != nil {
		return "", "", 0, err
	}
	rt, _, err := as.makeRefreshToken(ctx, u, newRTJTI)
	if err != nil {
		return "", "", 0, err
	}
	
	return at, rt, int64(atExp.Unix()), nil
}

func (as *AuthService) Logout(ctx context.Context, userID int64, refreshJTI string, accessJTI string, accessExp time.Time) error {
	// Revoke refresh session
	// key := fmt.Sprintf("RT:%d:%s", userID, refreshJTI)
	// // _ = as.tokenRdb.Del(ctx, key).Err()

	// // Optional: blacklist access token by JTI until it expires
	// if accessJTI != "" && !accessExp.IsZero() {
	// 	blKey := fmt.Sprintf("BL_AT:%s", accessJTI)
	// 	ttl := time.Until(accessExp)
	// 	if ttl > 0 {
	// 		_ = as.tokenRdb.Set(ctx, blKey, "block token", ttl).Err()
	// 	}
	// }
	return nil
}
