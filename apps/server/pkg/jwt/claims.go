package jwt

type Claims struct {
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
}
