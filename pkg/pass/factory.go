package pass

import "strings"

type Algo string

const (
	AlgoArgon2id Algo = "argon2id"
	AlgoBcrypt   Algo = "bcrypt"
)

type Config struct {
	Algo Algo

	// Common
	Pepper []byte

	// Argon2id
	MemoryKiB uint32
	TimeCost  uint32
	Threads   uint8
	KeyLen    uint32
	SaltLen   int

	// Bcrypt
	BcryptCost int
}

func New(cfg Config) Hasher {
	switch cfg.Algo {
	case AlgoArgon2id:
		return NewArgon2idHasher(Argon2idOptions{
			MemoryKiB: cfg.MemoryKiB,
			TimeCost:  cfg.TimeCost,
			Threads:   cfg.Threads,
			KeyLen:    cfg.KeyLen,
			SaltLen:   cfg.SaltLen,
			Pepper:    cfg.Pepper,
		})
	case AlgoBcrypt:
		return NewBcryptHasher(BcryptOptions{
			Cost:   cfg.BcryptCost,
			Pepper: cfg.Pepper,
		})
	default:
		return NewArgon2idHasher(Argon2idOptions{Pepper: cfg.Pepper})
	}
}

func Detect(encoded string, argon2id, bcrypt Hasher) Hasher {
	if strings.HasPrefix(encoded, "$argon2id$") {
		return argon2id
	}
	if strings.HasPrefix(encoded, "$2a$") || strings.HasPrefix(encoded, "$2b$") || strings.HasPrefix(encoded, "$2y$") {
		return bcrypt
	}
	// fallback
	return bcrypt
}
