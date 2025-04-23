package qzone

import (
	"log"

	"github.com/intchensc/qzone/api"
	"github.com/intchensc/qzone/auth"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetPrefix("[qzone]")
}

type Qzone struct {
	Auth auth.BaseAuth
	API  *api.API
}

// TODO:完善异步登录
func (q *Qzone) Login() error {
	// 使用channel等待登录完成
	done := make(chan error)
	go func() {
		err := q.Auth.Login()
		done <- err
	}()

	err := <-done
	if err != nil {
		return err
	}
	q.API.SetLogin(q.Auth)
	return nil
}

func New(auth auth.BaseAuth) *Qzone {
	return &Qzone{
		Auth: auth,
		API:  api.New(),
	}
}
