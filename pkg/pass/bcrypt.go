package pass

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptOptions struct {
	Cost   int    // e.g. 12
	Pepper []byte // optional app-wide secret
}

type BcryptHasher struct {
	o BcryptOptions
}

func NewBcryptHasher(o BcryptOptions) *BcryptHasher {
	if o.Cost == 0 {
		o.Cost = bcrypt.DefaultCost
	} // 10
	return &BcryptHasher{o: o}
}

func (h *BcryptHasher) preimage(pw string) []byte {
	if len(h.o.Pepper) == 0 {
		return []byte(pw)
	}
	// đơn giản: pw || pepper (hoặc HMAC như argon2id tuỳ bạn)
	var buf bytes.Buffer
	buf.WriteString(pw)
	buf.Write(h.o.Pepper)
	return buf.Bytes()
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(h.preimage(password), h.o.Cost)
	return string(b), err
}

func (h *BcryptHasher) Verify(password, encoded string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encoded), h.preimage(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return err == nil, err
}

func (h *BcryptHasher) RehashNeeded(encoded string) bool {
	cost, err := bcrypt.Cost([]byte(encoded))
	if err != nil {
		return true
	}
	return cost < h.o.Cost
}
