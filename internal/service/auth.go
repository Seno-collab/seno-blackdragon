package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"seno-blackdragon/internal/keys"
	"seno-blackdragon/internal/model"
	"seno-blackdragon/internal/repository"
	"seno-blackdragon/pkg/enum"
	"seno-blackdragon/pkg/pass"

	cryptoRand "crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type JWTConfig struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration // e.g. 15 * time.Minute
	RefreshTTL    time.Duration // e.g. 30 * 24 * time.Hour
	Issuer        string        // e.g. "seno-blackdragon"
}

type AccessClaims struct {
	Email     string   `json:"email"`
	TokenType string   `json:"typ"` // "access" | "refresh"
	DeviceID  string   `json:"did"`
	SessionID string   `json:"sid"`
	Roles     []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	DeviceID  string `json:"did"`
	Uv        int    `json:"uv"`
	Fam       string `json:"fam"`
	TokenType string `json:"typ"` // "access" | "refresh"
	jwt.RegisteredClaims
}

type Session struct {
	UserID    string   `json:"user_id"`
	DeviceID  string   `json:"device_id"`
	IP        string   `json:"ip,omitempty"`
	UA        string   `json:"ua,omitempty"`
	CreatedAt string   `json:"created_at"`
	LastSeen  string   `json:"last_seen"`
	Exp       string   `json:"exp"`
	Scopes    []string `json:"scopes,omitempty"`
	MFA       bool     `json:"mfa"`
	Status    string   `json:"status"` // active | revoked
}

type Device struct {
	UserID     string `json:"user_id"`
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type,omitempty"`
	OS         string `json:"os,omitempty"`
	UA         string `json:"ua,omitempty"`
	FirstSeen  string `json:"first_seen"`
	LastSeen   string `json:"last_seen"`
	Trusted    bool   `json:"trusted"`
	Status     string `json:"status"` // active | blocked
	Name       string `json:"name,omitempty"`
}

type AuthService struct {
	userRepo *repository.UserRepo // interface
	hasher   pass.Hasher          // Argon2id/Bcrypt impl
	jwtCfg   JWTConfig
	redis    *redis.Client
	log      *zap.Logger
}

func NewAuthService(
	userRepo *repository.UserRepo,
	hasher pass.Hasher,
	redis *redis.Client,
	jwtCfg JWTConfig,
	log *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		hasher:   hasher,
		jwtCfg:   jwtCfg,
		log:      log,
		redis:    redis,
	}
}

func nowISO() string { return time.Now().UTC().Format(time.RFC3339) }
func addSecISO(sec int) string {
	return time.Now().UTC().Add(time.Duration(sec) * time.Second).Format(time.RFC3339)
}
func newID(prefix string) string {
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	return prefix + ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
func (as *AuthService) EnsureUserVersion(ctx context.Context, userID string) (int, error) {
	key := keys.UserVer(userID)
	v, err := as.redis.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if v == 0 {
		if err := as.redis.Set(ctx, key, 1, 0).Err(); err != nil {
			return 0, err
		}
		return 1, nil
	}
	return 1, nil
}

func (as *AuthService) GetUserVersion(ctx context.Context, userID string) (int, error) {
	key := keys.UserVer(userID)
	v, err := as.redis.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, err
	}
	return v, err
}

func (as *AuthService) SaveDevice(ctx context.Context, d *Device) error {
	b, _ := json.Marshal(d)
	pipe := as.redis.TxPipeline()
	pipe.Set(ctx, keys.Device(d.DeviceID), b, 0)
	pipe.SAdd(ctx, keys.UserDevice(d.UserID), d.DeviceID)
	_, err := pipe.Exec(ctx)
	return err
}

func (as *AuthService) GetDevice(ctx context.Context, deviceID string) (*Device, error) {
	raw, err := as.redis.Get(ctx, keys.Device(deviceID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var d Device
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// ===== token helpers =====

func (as *AuthService) makeAccessToken(u *repository.UserModel, jti, sessionID, deviceID string, roles []string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(as.jwtCfg.AccessTTL)
	claims := &AccessClaims{
		Email:     u.Email,
		TokenType: "access",
		SessionID: sessionID,
		DeviceID:  deviceID,
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

func (as *AuthService) makeRefreshToken(u *repository.UserModel, jti, deviceID, fam string, uv int) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Duration(as.jwtCfg.RefreshTTL) * time.Second)
	claims := &RefreshClaims{
		DeviceID: deviceID,
		Fam:      fam,
		Uv:       uv,
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

func (as *AuthService) SaveSession(ctx context.Context, sid string, s *Session, ttlSec int) error {
	b, _ := json.Marshal(s)
	pipe := as.redis.TxPipeline()
	pipe.Set(ctx, keys.Session(sid), b, time.Duration(ttlSec)*time.Second)
	pipe.SAdd(ctx, keys.UserSession(s.UserID), sid)
	_, err := pipe.Exec(ctx)
	return err
}

func (as *AuthService) GetSession(ctx context.Context, sid string) (*Session, error) {
	raw, err := as.redis.Get(ctx, keys.Session(sid)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (as *AuthService) DelSession(ctx context.Context, userID, sid string) error {
	pipe := as.redis.TxPipeline()
	pipe.Del(ctx, keys.Session(sid))
	pipe.SRem(ctx, keys.UserSession(userID))
	_, err := pipe.Exec(ctx)
	return err
}

func (as *AuthService) LogoutAll(ctx context.Context, userID string) error {
	if err := as.redis.Incr(ctx, keys.UserSession(userID)).Err(); err != nil {
		return err
	}
	sids, _ := as.redis.SMembers(ctx, keys.UserSession(userID)).Result()
	if len(sids) > 0 {
		pipe := as.redis.TxPipeline()
		for _, sid := range sids {
			pipe.Del(ctx, keys.Session(sid))
			pipe.Del(ctx, keys.UserSession(userID))
		}
		_, _ = pipe.Exec(ctx)
	}
	return nil
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

func (as *AuthService) Login(ctx context.Context, cmd model.LoginCmd) (*model.TokenPair, error) {
	u, err := as.userRepo.GetUserByEmail(ctx, cmd.Email)
	if err != nil || u == nil {
		return nil, enum.ErrInvalidCredentials
	}
	ok, _ := as.hasher.Verify(cmd.Password, u.PasswordHash)
	if !ok {
		return nil, enum.ErrInvalidCredentials
	}
	did := cmd.DeviceID
	if did == "" {
		did = newID("DEV_")
	}
	dev, _ := as.GetDevice(ctx, did)
	if dev == nil {
		dev = &Device{
			UserID:    u.ID.String(),
			DeviceID:  did,
			OS:        cmd.DeviceMeta["os"],
			UA:        cmd.UA,
			FirstSeen: nowISO(),
			LastSeen:  nowISO(),
			Trusted:   false,
			Status:    model.Active,
			Name:      cmd.DeviceMeta["name"],
		}
		if err := as.SaveDevice(ctx, dev); err != nil {
			return nil, err
		}
	} else {
		dev.LastSeen = nowISO()
		_ = as.SaveDevice(ctx, dev)
	}
	sid := newID("SID_")
	session := &Session{
		UserID:    u.ID.String(),
		DeviceID:  did,
		IP:        cmd.IP,
		UA:        cmd.UA,
		CreatedAt: nowISO(),
		LastSeen:  nowISO(),
		Exp:       addSecISO(int(as.jwtCfg.AccessTTL)),
		MFA:       true,
		Status:    model.Active,
	}
	if err := as.SaveSession(ctx, sid, session, int(as.jwtCfg.AccessTTL)); err != nil {
		return nil, err
	}

	atJTI := fmt.Sprintf("at-%s-%d", u.ID, time.Now().UnixNano())
	at, atExp, err := as.makeAccessToken(u, atJTI, sid, did, []string{})
	if err != nil {
		return nil, enum.ErrInvalidToken
	}
	fam := newID("FAM_")
	rtJTI := fmt.Sprintf("rt-%s-%d", u.ID, time.Now().UnixNano())
	uv, _ := as.EnsureUserVersion(ctx, u.ID.String())

	rt, _, err := as.makeRefreshToken(u, rtJTI, did, fam, uv)
	if err != nil {
		return nil, enum.ErrInvalidToken
	}

	pipe := as.redis.TxPipeline()
	pipe.Set(ctx, keys.RTActive(rtJTI), fmt.Sprintf(`{"user_id:"%s","device_id": "%s", "fam":"%s"}`, u.ID.String(), did, fam), as.jwtCfg.RefreshTTL)
	pipe.SAdd(ctx, keys.FamActive(fam), rtJTI)
	pipe.SAdd(ctx, keys.UserDeviceFams(u.ID.String(), did), fam)
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}
	return &model.TokenPair{
		AccessToken:  at,
		RefreshToken: rt,
		Expired:      int64(atExp.Unix()),
	}, nil
}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	// Parse & validate refresh token
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(as.jwtCfg.Issuer),
	)
	claims := &RefreshClaims{}
	tok, err := parser.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, enum.ErrWrongAlgorithm
		}
		return as.jwtCfg.RefreshSecret, nil
	})
	if err != nil || !tok.Valid {
		return nil, enum.ErrInvalidToken
	}
	if claims.TokenType != "refresh" {
		return nil, enum.ErrWrongType
	}
	jti := claims.Subject
	if jti == "" {
		jti = claims.Subject
	}
	fam := claims.Fam
	if fam == "" {
		return nil, enum.ErrInvalidToken
	}
	did := claims.DeviceID
	if n, _ := as.redis.Exists(ctx, keys.FamBlack(fam)).Result(); n == 1 {
		return nil, enum.ErrFamilyBlocked
	}
	if n, _ := as.redis.Exists(ctx, keys.RTRevoked(jti)).Result(); n == 1 {
		_ = as.redis.Set(ctx, keys.FamBlack(fam), "1", as.jwtCfg.RefreshTTL).Err()
		return nil, enum.ErrRefreshRevoked
	}
	if n, _ := as.redis.Exists(ctx, keys.RTActive(jti)).Result(); n != 1 {
		return nil, enum.ErrRefreshNotActive
	}
	if claims.Uv > 0 {
		curUv, err := as.GetUserVersion(ctx, jti)
		if err != nil && curUv != claims.Uv {
			return nil, enum.ErrInvalidToken
		}
	}
	u, err := as.userRepo.GetUserByID(ctx, uuid.MustParse(claims.Subject))
	if err != nil || u == nil {
		return nil, enum.ErrUserNotFound
	}
	gotLock, _ := as.redis.SetNX(ctx, keys.RotateLock(fam), "1", 10*time.Second).Result()
	if !gotLock {
		return nil, enum.ErrRotationRace
	}
	defer as.redis.Del(ctx, keys.RotateLock(fam))
	ttlLeft := time.Until(claims.ExpiresAt.Time)
	if ttlLeft < time.Second {
		ttlLeft = time.Second
	}
	pipe := as.redis.TxPipeline()
	pipe.Del(ctx, keys.RTActive(jti))
	pipe.SRem(ctx, keys.FamActive(fam), jti)
	pipe.Set(ctx, keys.RTRevoked(jti), "1", ttlLeft)

	newRTJTI := fmt.Sprintf("rt-%d-%d", u.ID, time.Now().UnixNano())
	newRT, _, err := as.makeRefreshToken(u, newRTJTI, claims.DeviceID, fam, claims.Uv)
	if err != nil {
		return nil, err
	}
	pipe.SAdd(ctx, keys.FamActive(fam), newRTJTI)
	newSID := newID("SID_")
	newSess := &Session{
		UserID:    u.ID.String(),
		DeviceID:  did,
		IP:        "",
		UA:        "",
		CreatedAt: nowISO(),
		LastSeen:  nowISO(),
		Exp:       addSecISO(int(as.jwtCfg.AccessTTL)),
		MFA:       true,
		Status:    model.Active,
	}
	if err := as.SaveSession(ctx, newSID, newSess, int(as.jwtCfg.AccessTTL)); err != nil {
		return nil, err
	}
	newATJTI := fmt.Sprintf("at-%s-%d", u.ID, time.Now().Unix())
	newAT, newAtExp, err := as.makeAccessToken(u, newATJTI, newID("SID_"), claims.DeviceID, []string{})
	if err != nil {
		return nil, err
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}
	token := &model.TokenPair{
		AccessToken:  newAT,
		RefreshToken: newRT,
		Expired:      int64(newAtExp.Unix()),
	}
	return token, nil
}
func (as *AuthService) LogoutDevice(ctx context.Context, userID, deviceID string) error {
	fams, _ := as.redis.SMembers(ctx, keys.UserDeviceFams(userID, deviceID)).Result()
	pipe := as.redis.TxPipeline()
	for _, fam := range fams {
		pipe.Set(ctx, keys.FamBlack(fam), "1", as.jwtCfg.RefreshTTL)
		jtis, _ := as.redis.SMembers(ctx, keys.FamActive(fam)).Result()
		for _, j := range jtis {
			pipe.Del(ctx, keys.RTActive(j))
			pipe.Set(ctx, keys.RTRevoked(j), "1", time.Hour)
		}
		pipe.Del(ctx, keys.FamActive(fam))
	}
	_, _ = pipe.Exec(ctx)

	sids, _ := as.redis.SMembers(ctx, keys.UserSession(userID)).Result()
	for _, sid := range sids {
		if sess, _ := as.GetSession(ctx, sid); sess != nil && sess.DeviceID == deviceID {
			_ = as.DelSession(ctx, userID, sid)
		}
	}
	return nil
}

// func (as *AuthService) LogoutAll(ctx context.Context, userID string) error {
// 	fams, _ := as.redis.SMembers(ctx, keys.UserFam)
// }
