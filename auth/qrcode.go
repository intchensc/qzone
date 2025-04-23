package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
)

var (
	ptqrshowURL  = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=549000912&e=2&l=M&s=3&d=72&v=4&t=0.31232733520361844&daid=5&pt_3rd_aid=0"
	ptqrloginURL = "https://xui.ptlogin2.qq.com/ssl/ptqrlogin?u1=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&ptqrtoken=%v&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-1656992258324&js_ver=22070111&js_type=1&login_sig=&pt_uistyle=40&aid=549000912&daid=5&has_onekey=1&&o1vId=1e61428d61cb5015701ad73d5fb59f73"
	checkSigURL  = "https://ptlogin2.qzone.qq.com/check_sig?pttype=1&uin=%v&service=ptqrlogin&nodirect=1&ptsigx=%v&s_url=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&f_url=&ptlang=2052&ptredirect=100&aid=549000912&daid=5&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0"
)

type QrAuth struct {
	qrsig   string // 二维码接口获取到的参数
	qrtoken string // 由qrsig计算而成
	cookie  string // 登录成功后携带的cookie
}

func NewQrAuth() *QrAuth {
	return &QrAuth{}
}

// GenerateQRCode 生成二维码，返回base64 二维码ID 用于查询扫码情况

func (q *QrAuth) Login() error {
	// 生成二维码
	b64s, err := q.GenerateQRCode()
	if err != nil {
		return errors.New("生成二维码失败: " + err.Error())
	}

	// 解码base64数据
	ddd, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		return errors.New("base64解码失败: " + err.Error())
	}
	// 保存到本地文件
	err = os.WriteFile("qrcode.png", ddd, 0666)
	if err != nil {
		return errors.New("写入二维码到文件失败: " + err.Error())
	}

	log.Println("二维码已保存到qrcode.png，请扫描登录")

	// 循环检查登录状态
	for {
		//0成功 1未扫描 2未确认 3已过期 -1系统错误
		status, err := q.CheckQRCodeStatus()
		if err != nil {
			return errors.New("检测二维码状态失败: " + err.Error())
		}

		switch status {
		case 0:
			log.Println("登录成功")
			return nil
		case 1:
			log.Println("等待扫描二维码...")
		case 2:
			log.Println("请在手机上确认登录...")
		case 3:
			return errors.New("二维码已失效或登录被拒绝")
		default:
			return errors.New("未知的登录状态: " + fmt.Sprintf("%d", status))
		}

		time.Sleep(2 * time.Second)
	}
}

func (q *QrAuth) Logout() error {
	return nil
}

func (q *QrAuth) IsLogin() bool {
	return false
}
func (q *QrAuth) GetCookie() string {
	return q.cookie
}

func (q *QrAuth) GenerateQRCode() (string, error) {
	cookiesString := ""
	q.qrsig = ""
	var data []byte
	var rc []*http.Cookie
	err := requests.URL(ptqrshowURL).Handle(func(r *http.Response) error {
		rc = r.Cookies()
		data, _ = io.ReadAll(r.Body)
		defer r.Body.Close()
		return nil
	}).Fetch(context.Background())

	// data, err := DialRequest(NewRequest(
	// 	WithUrl(ptqrshowURL),
	// 	WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	}}),
	// 	WithRespFunc(func(response *http.Response) {
	// 		for _, v := range response.Cookies() {
	// 			cookiesString = cookiesString + v.String()
	// 			if v.Name == "qrsig" {
	// 				q.qrsig = v.Value
	// 				break
	// 			}
	// 		}
	// 	})))
	if err != nil {
		// er := errors.New("空间登录二维码显示错误:" + data)
		fmt.Println(err)
		// return "", er
		// return "", er
	}
	for _, c := range rc {
		if c.Name == "qrsig" {
			cookiesString = cookiesString + c.String()
			q.qrsig = c.Value
			break
		}
	}
	if q.qrsig == "" {
		er := errors.New("空间登录二维码cookie获取错误:" + cookiesString)
		return "", er
	}
	if err != nil {
		// er := errors.New("空间登录二维码显示错误:" + err.Error())
		// return "", er
		log.Println(err)
	}

	base64 := base64.StdEncoding.EncodeToString(data)
	q.qrtoken = genderGTK(q.qrsig, 0)
	return base64, nil
}

// CheckQRCodeStatus 检查二维码状态 //0成功 1未扫描 2未确认 3已过期  -1系统错误
func (q *QrAuth) CheckQRCodeStatus() (int8, error) {
	// if q.status == 0 {
	// 	return 0, nil
	// }
	qrtoken := q.qrtoken
	qrsign := q.qrsig
	qcookie := q.cookie
	urls := fmt.Sprintf(ptqrloginURL, qrtoken)
	var data []byte
	err := requests.URL(urls).Header(
		"cookie", "qrsig="+qrsign).Handle(func(r *http.Response) error {
		data, _ = io.ReadAll(r.Body)
		for _, v := range r.Cookies() {
			if v.Value != "" {
				qcookie += v.Name + "=" + v.Value + ";"
			}
		}
		defer r.Body.Close()
		return nil
	}).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(
	// 	WithUrl(urls),
	// 	WithHeader(map[string]string{
	// 		"cookie": "qrsig=" + qrsign,
	// 	}),
	// 	WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	}}),
	// 	WithRespFunc(func(response *http.Response) {
	// 		for _, v := range response.Cookies() {
	// 			if v.Value != "" {
	// 				qcookie += v.Name + "=" + v.Value + ";"
	// 			}
	// 		}
	// 	})))
	if err != nil {
		log.Println(err)
		// return -1, er
	}
	text := string(data)
	switch {
	case strings.Contains(text, "二维码未失效"):
		return 1, nil
	case strings.Contains(text, "二维码认证中"):
		return 2, nil
	case strings.Contains(text, "二维码已失效") || strings.Contains(text, "本次登录已被拒绝"):
		return 3, nil
	case strings.Contains(text, "登录成功"):
		dealedCheckText := strings.ReplaceAll(text, "'", "")
		redirectURL := strings.Split(dealedCheckText, ",")[2]
		redirectCookie, err := loginRedirect(redirectURL)
		if err != nil {
			er := errors.New("空间登录重定向失败:" + err.Error())
			return -1, er
		}
		qcookie += redirectCookie
		q.cookie = strings.ReplaceAll(qcookie, " ", "")
		return 0, nil
	}
	return 0, nil
}

// loginRedirect 登录成功回调
func loginRedirect(redirectURL string) (string, error) {
	var cookie string
	// redirectURL := strings.ReplaceAll(redirectURL, "'", "")
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", errors.New("空间登录重定向链接解析错误:" + err.Error())
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.New("空间登录重定向链接查询参数解析错误:" + err.Error())
	}
	urls := fmt.Sprintf(checkSigURL, values["uin"][0], values["ptsigx"][0])
	errr := requests.URL(urls).Handle(func(r *http.Response) error {
		for _, v := range r.Cookies() {
			if v.Value != "" {
				cookie += v.Name + "=" + v.Value + ";"
			}
		}
		defer r.Body.Close()
		return nil
	}).Fetch(context.Background())
	// _, err = DialRequest(NewRequest(
	// 	WithUrl(urls),
	// 	WithClient(&http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	}}),
	// 	WithRespFunc(func(response *http.Response) {
	// 		for _, v := range response.Cookies() {
	// 			if v.Value != "" {
	// 				cookie += v.Name + "=" + v.Value + ";"
	// 			}
	// 		}
	// 	})))
	if errr != nil {
		// return "", errors.New("空间登录重定向链接请求错误:" + errr)
	}
	return cookie, nil
}

// genderGTK 生成GTK
func genderGTK(sKey string, hash int) string {
	for _, s := range sKey {
		us, _ := strconv.Atoi(fmt.Sprintf("%d", s))
		hash += (hash << 5) + us
	}
	return fmt.Sprintf("%d", hash&0x7fffffff)
}
