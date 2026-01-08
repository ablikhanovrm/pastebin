package jwt

import "time"

type Manager struct {
	secret string
}

func New(secret string) *Manager {
	return &Manager{secret}
}

func (m *Manager) Generate(userID int64, ttl time.Duration) (string, error) {
	return "", nil
}
func (m *Manager) Parse(token string) (*Claims, error) { return nil, nil }
