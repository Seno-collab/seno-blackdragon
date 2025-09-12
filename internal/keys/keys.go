package keys

func Device(deviceID string) string   { return "device:" + deviceID }
func UserDevice(userID string) string { return "user_devices:" + userID }

func Session(sid string) string     { return "session:" + sid }
func UserSession(uid string) string { return "user_sessions" + uid }
func RevokedSess(sid string) string { return "revoked:sess:" + sid }

func UserVer(uid string) string { return "user_ver" + uid }

func UserDeviceFams(uid, did string) string { return "user_device_fams:" + uid + ":" + did }
func UserFams(uid string) string            { return "user_fams:" + uid }
func RTActive(jti string) string            { return "rt:active" + jti }
func RTRevoked(jti string) string           { return "rt:revoked:" + jti }
func RTFamilyBlock(fam string) string       { return "rt:family:black" + fam }
func RotateLock(fam string) string          { return "rotate:lock" + fam }
func FamActive(fam string) string           { return "rt:fam:active" + fam }
func FamBlack(fam string) string            { return "rt:family:back" + fam }
