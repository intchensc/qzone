package common

import (
    "fmt"
    "net/url"
    "reflect"
    "strconv"
    "strings"
)

var (
    UA          = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
    ContentType = "application/x-www-form-urlencoded"
)

type BaseAPI struct {
    Qq     int64 // QQ号
    Gtk    string
    Gtk2   string
    Pskey  string
    Skey   string
    Uin    string
    Cookie string
    Err    error
}

// Unpack 初始化信息, 将成功扫码登录获取到的 cookie 解析
func (b *BaseAPI) Unpack(cookie string) {
    b.Err = nil
    if cookie == "" {
        b.Err = fmt.Errorf("cookie 为空")
        return
    }
    for _, v := range strings.Split(cookie, ";") {
        name, val, f := strings.Cut(v, "=")
        if f {
            name = strings.TrimSpace(name)
            val = strings.TrimSpace(val)
            switch name {
            case "uin":
                b.Uin = val
            case "skey":
                b.Skey = val
            case "p_skey":
                b.Pskey = val
            }
        }
    }
    b.Gtk = genderGTK(b.Skey, 5381)
    b.Gtk2 = genderGTK(b.Pskey, 5381)
    t, err := strconv.ParseInt(strings.TrimPrefix(b.Uin, "o"), 10, 64)
    if err != nil {
        b.Err = fmt.Errorf("uin 解析失败: %w", err)
        b.Cookie = cookie
        return
    }
    b.Qq = t
    b.Cookie = strings.TrimSpace(cookie)
    b.Err = nil
    return
}

func GetShuoShuoUnikey(uin string, tid string) (unikey string) {
    return fmt.Sprintf("http://user.qzone.qq.com/%s/mood/%s", uin, tid)
}

// genderGTK 生成GTK
func genderGTK(sKey string, hash int) string {
    for _, s := range sKey {
        us, _ := strconv.Atoi(fmt.Sprintf("%d", s))
        hash += (hash << 5) + us
    }
    return fmt.Sprintf("%d", hash&0x7fffffff)
}

func StructToStr(in interface{}) (payload string) {
    keys := make([]string, 0, 16)
    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    t := v.Type()
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        get := field.Tag.Get("json")
        if get != "" {
            var t string
            if v.Field(i).Kind() == reflect.Int64 {
                t = strconv.FormatInt(v.Field(i).Int(), 10)
            } else {
                t = v.Field(i).Interface().(string)
            }

            keys = append(keys, get+"="+url.QueryEscape(t))
        }
    }
    payload = strings.Join(keys, "&")
    return
}
