package pass

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2idOptions struct {
	MemoryKiB uint32 // e.g. 64*1024
	TimeCost  uint32 // e.g. 3
	Threads   uint8  // e.g. 1
	KeyLen    uint32 // e.g. 32
	SaltLen   int    // e.g. 16
	Pepper    []byte // optional app-wide secret
}

type Argon2idHasher struct {
	o Argon2idOptions
}

func NewArgon2idHasher(o Argon2idOptions) *Argon2idHasher {
	// sensible defaults
	if o.MemoryKiB == 0 {
		o.MemoryKiB = 64 * 1024
	}
	if o.TimeCost == 0 {
		o.TimeCost = 3
	}
	if o.Threads == 0 {
		o.Threads = 1
	}
	if o.KeyLen == 0 {
		o.KeyLen = 32
	}
	if o.SaltLen == 0 {
		o.SaltLen = 16
	}
	return &Argon2idHasher{o: o}
}

func (h *Argon2idHasher) preimage(pw string) []byte {
	if len(h.o.Pepper) == 0 {
		return []byte(pw)
	}
	mac := hmac.New(sha256.New, h.o.Pepper)
	mac.Write([]byte(pw))
	return mac.Sum(nil)
}

func (h *Argon2idHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.o.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := argon2.IDKey(h.preimage(password), salt, h.o.TimeCost, h.o.MemoryKiB, h.o.Threads, h.o.KeyLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.o.MemoryKiB, h.o.TimeCost, h.o.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(sum),
	), nil
}

func (h *Argon2idHasher) Verify(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid argon2id hash")
	}
	var m, t uint32
	var p uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	got := argon2.IDKey(h.preimage(password), salt, t, m, p, uint32(len(want)))
	// constant-time compare
	if len(got) != len(want) {
		return false, nil
	}
	var v byte
	for i := range got {
		v |= got[i] ^ want[i]
	}
	return v == 0, nil
}

func (h *Argon2idHasher) RehashNeeded(encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return true
	}
	var m, t uint32
	var p uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return true
	}
	return m < h.o.MemoryKiB || t < h.o.TimeCost || p != h.o.Threads
}
