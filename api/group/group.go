package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"regexp"

	"github.com/carlmjohnson/requests"
	"github.com/intchensc/qzone/api/common"
	"github.com/tidwall/gjson"
)

var (
	cRe          = regexp.MustCompile(`(?s)_Callback\((.*)\)`)
	userQzoneURL = "https://user.qzone.qq.com"
	ua           = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	// 获取QQ群URL
	getQQGroupURL = "https://user.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/qqgroupfriend_extend.cgi?"
	// 获取QQ群成员非好友URL
	getQQGroupMemberURL = "https://user.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/qqgroupfriend_groupinfo.cgi?"
	// 获取QQ空间历史消息
	getQZoneHistory = "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
)

type GroupAPI struct {
	*common.BaseAPI
}

// QQGroupList 群列表获取
func (g *GroupAPI) List() ([]*QQGroupResp, error) {
	if g.Err != nil {
		return nil, g.Err
	}
	// 构建请求参数结构体
	gr := &QQGroupReq{
		Uin:     g.Qq,
		Do:      "1",
		Rd:      fmt.Sprintf("%010.8f", rand.Float64()),
		Fupdate: "1",
		Clean:   "1",
		GTk:     g.Gtk2,
	}

	var data string
	url := getQQGroupURL + common.StructToStr(gr)
	err := requests.URL(url).
		UserAgent(ua).Header("cookie", g.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// url := getQQGroupURL + structToStr(gr)
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"user-agent": ua,
	// 	"cookie":     g.cookie,
	// })))
	if err != nil {
		// er := errors.New("QQ群请求错误:" + err.Error())
		log.Println("QQ群获取失败:", err)
		// return nil, er
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		er := errors.New("QQ群响应正则解析错误:" + data)
		log.Println("QQ群获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "data.group.#").Int()
	results := make([]*QQGroupResp, resLen)
	index := 0
	groups := gjson.Get(jsonStr, "data.group").Array()
	for _, group := range groups {
		gro := &QQGroupResp{
			GroupCode:   group.Get("groupcode").Int(),
			GroupName:   group.Get("groupname").String(),
			TotalMember: group.Get("total_member").Int(),
			NotFriends:  group.Get("notfriends").Int(),
		}
		results[index] = gro
		index++
	}
	return results, nil
}

// QQGroupMemberList 群友(非好友)列表获取
func (g *GroupAPI) MemberList(gid int64) ([]*QQGroupMemberResp, error) {
	if g.Err != nil {
		return nil, g.Err
	}
	gmr := &QQGroupMemberReq{
		Uin:     g.Qq,
		Gid:     gid,
		Fupdate: "1",
		Type:    "1",
		GTk:     g.Gtk2,
	}
	// url := getQQGroupMemberURL + structToStr(gmr)
	var data string
	url := getQQGroupMemberURL + common.StructToStr(gmr)
	err := requests.URL(url).
		UserAgent(ua).Header("cookie", g.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"user-agent": ua,
	// 	"cookie":     g.cookie,
	// })))
	if err != nil {
		// er := errors.New("QQ群非好友请求错误:" + err.Error())
		log.Println("QQ群非好友获取失败:", err)
		// return nil, er
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		er := errors.New("QQ群非好友正则解析错误:" + data)
		log.Println("QQ群非好友获取失败:", er.Error())
		return nil, er
	}
	jsonStr := r[1]
	resLen := gjson.Get(jsonStr, "data.notfriends").Int()
	results := make([]*QQGroupMemberResp, resLen)
	index := 0
	groupMembers := gjson.Get(jsonStr, "data.friends").Array()
	for _, groupMember := range groupMembers {
		gro := &QQGroupMemberResp{
			Uin:       groupMember.Get("fuin").Int(),
			NickName:  groupMember.Get("name").String(),
			AvatarURL: groupMember.Get("img").String(),
		}
		gro.GroupCode = gjson.Get(jsonStr, "data.groupcode").Int()
		results[index] = gro
		index++
	}
	return results, nil
}
