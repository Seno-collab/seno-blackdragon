package pass

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, encode string) (bool, error)
	RehashNeeded(encoded string) bool
}
