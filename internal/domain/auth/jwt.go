package auth

type AuthorizationJWTPayload struct {
	Sub      string   `json:"sub"`
	DeviceID string   `json:"device_id"`
	Scopes   []string `json:"scopes"`
	Type     string   `json:"type"`
	IssuedAt int64    `json:"iat"`
	Expires  int64    `json:"exp"`
	JTI      string   `json:"jti"`
}
