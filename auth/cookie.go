package auth

import "fmt"

type CookieAuth struct {
	Cookie string
}

func (c *CookieAuth) Login() error {
	if c.Cookie == "" {
		return fmt.Errorf("cookie 为空")
	}
	return nil
}

func (c *CookieAuth) Logout() error { return nil }

func (c *CookieAuth) IsLogin() bool { return c.Cookie != "" }

func (c *CookieAuth) GetCookie() string { return c.Cookie }
