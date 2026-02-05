package auth

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type RegisterInput struct {
	Email     string
	Password  string
	Name      string
	IP        string
	UserAgent string
}
