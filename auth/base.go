package auth

type BaseAuth interface {
	Login() error
	Logout() error
	IsLogin() bool
	GetCookie() string
}
