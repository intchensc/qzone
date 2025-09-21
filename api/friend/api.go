package friend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/carlmjohnson/requests"
	"github.com/intchensc/qzone/api/common"
	"github.com/tidwall/gjson"
)

var (
	cRe             = regexp.MustCompile(`(?s)_Callback\((.*)\)`)
	userQzoneURL    = "https://user.qzone.qq.com"
	friendURL       = "https://h5.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/friend_show_qqfriends.cgi?g_tk=%v"
	detailFriendURL = "https://h5.qzone.qq.com/proxy/domain/base.qzone.qq.com/cgi-bin/user/cgi_userinfo_get_all?g_tk=%v"
)

type API struct {
	*common.BaseAPI
}

// List 好友列表获取 TODO:有时候显示亲密度前200好友
func (f *API) List() ([]*FriendInfoEasyResp, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	url := fmt.Sprintf(friendURL, f.Gtk) + "&uin=" + strconv.FormatInt(f.Qq, 10)
	log.Printf("请求好友列表URL: %s", url)
	var data string
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		UserAgent(common.UA).
		Header("origin", userQzoneURL).
		Header("cookie", f.Cookie).
		ToString(&data).Fetch(context.Background())
	log.Printf("原始响应数据: %s", data)

	if err != nil {

		log.Println("好友列表获取失败:", err)
		// return nil, err
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		er := errors.New("好友列表正则解析错误:" + string(data))
		log.Println("好友列表获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "items.#").Int()
	results := make([]*FriendInfoEasyResp, resLen)
	index := 0

	friends := gjson.Get(jsonStr, "items").Array()
	for _, friend := range friends {
		fie := &FriendInfoEasyResp{
			Uin:     friend.Get("uin").Int(),
			Groupid: friend.Get("groupid").Int(),
			Name:    friend.Get("name").String(),
			Remark:  friend.Get("remark").String(),
			Image:   friend.Get("image").String(),
			Online:  friend.Get("online").Int(),
		}
		results[index] = fie
		index++
	}
	return results, nil
}

// Detail 好友详细信息获取
func (f *API) Detail(uin int64) (*FriendInfoDetailResp, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	url := fmt.Sprintf(detailFriendURL, f.Gtk) + "&uin=" + strconv.FormatUint(uint64(uin), 10)
	var data string
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		Header("origin", userQzoneURL).
		Header("cookie", f.Cookie).
		ToString(&data).Fetch(context.Background())
	if err != nil {
		// er := errors.New("好友详细信息请求错误:" + err.Error())
		log.Println("好友详细信息获取失败:", err)
		// return nil, er
	}
	r := cRe.FindStringSubmatch(string(data))
	if len(r) < 2 {
		er := errors.New("好友详细信息正则解析错误:" + string(data))
		log.Println("好友详细信息获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]

	fid := &FriendInfoDetailResp{}
	if err := json.Unmarshal([]byte(jsonStr), fid); err != nil {
		er := errors.New("好友详细信息JSON绑定错误:" + err.Error())
		log.Println("好友详细信息获取失败:", er.Error())
		return nil, er
	}
	return fid, nil
}
