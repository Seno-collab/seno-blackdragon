package model

var (
	Active  = "active"
	Block   = "block"
	Revoked = "revoked"
)

type LoginCmd struct {
	Email      string
	Password   string
	DeviceID   string
	DeviceMeta map[string]string
	IP         string
	UA         string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	Expired      int64
}
