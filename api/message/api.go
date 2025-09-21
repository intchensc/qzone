package message

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"github.com/dop251/goja"
	"github.com/intchensc/qzone/api/common"
	"github.com/tidwall/gjson"
)

var (
	getQZoneHistory = "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
	cRe             = regexp.MustCompile(`(?s)_Callback\((.*)\)`)
)

type API struct {
	*common.BaseAPI
}

// ListN 获取最近N条
func (h *API) ListN(num int, ms int) ([]*Message, error) {
	message := make([]*Message, 0, num)
	index := 0
	mes, err := h.page(100, index)
	if err != nil {
		return nil, err
	}
	tlen := len(mes)
	for tlen < num-index {
		message = append(message, mes...)
		index = index + tlen
		mes, err = h.page(100, index)
		time.Sleep(time.Duration(ms) * time.Millisecond)
		if err != nil {
			return nil, err
		}
		tlen = len(mes)
	}
	message = append(message, mes[:num-index]...)
	return message, nil
}

// ListT 获取指定时间之后
func (h *API) ListT(t time.Time, ms int) ([]*Message, error) {
	var message []*Message
	index := 0
	mes, err := h.page(100, index)
	if err != nil {
		return nil, err
	}
	tlen := len(mes)
	for t.Before(mes[tlen-1].CreateTime) {
		message = append(message, mes...)
		index = index + tlen
		mes, err = h.page(100, index)
		time.Sleep(time.Duration(ms) * time.Millisecond)
		if err != nil {
			return nil, err
		}
		tlen = len(mes)
	}
	for _, k := range mes {
		if t.After(k.CreateTime) {
			break
		}
		message = append(message, k)
	}
	return message, err
}

// handleHTML 处理Html
func (h *API) handleHTML(message string) string {
	re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)
	ans := re.ReplaceAllStringFunc(message, func(hex string) string {
		byteValue, err := strconv.ParseUint(hex[2:], 16, 8)
		if err != nil {
			return hex
		}
		return string(rune(byteValue))
	})
	re = regexp.MustCompile(`\s+`)
	ans = re.ReplaceAllString(ans, " ")
	ans = strings.ReplaceAll(ans, "\\", "")
	return ans
}

// page 分类获取
func (h *API) page(num int, pos int) ([]*Message, error) {
	url := fmt.Sprintf(
		"https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?uin=%s&begin_time=0&end_time=0&getappnotification=1&getnotifi=1&has_get_key=0&offset=%d&set=0&count=%d&useutf8=1&outputhtmlfeed=1&scope=1&format=jsonp&g_tk=%s",
		strconv.FormatInt(h.Qq, 10), pos, num, h.Gtk,
	)
	//fmt.Println(url)
	var data string
	err := requests.
		URL(url).
		UserAgent(common.UA).
		Header("cookie", h.Cookie).
		ToString(&data).Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		log.Println("历史消息获取失败:", err)
	}
	jsonRaw := r[1]
	vm := goja.New()
	result, err := vm.RunString(fmt.Sprintf("JSON.stringify(%s)", gjson.Get(jsonRaw, "data")))
	if err != nil {
		return nil, err
	}
	totnum := gjson.Get(result.String(), "main.total_number").Int()
	if totnum == 0 {
		return nil, nil
	}
	var message []*Message
	dataS := gjson.Get(result.String(), "data").Array()
	for _, k := range dataS[:len(dataS)-1] {
		htmlRaw := k.Get("html").String()
		html := h.handleHTML(htmlRaw)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil, err
		}
		dataPub := doc.Find("i[name=feed_data]")
		tid, _ := dataPub.Attr("data-tid")
		otid, _ := dataPub.Attr("data-origtid")
		flag := doc.Find("div.f-nick>span.state").Text()
		ttime := time.Unix(k.Get("abstime").Int(), 0).Local()
		mes := &Message{
			TID:        tid,
			OriginTID:  otid,
			QQ:         k.Get("uin").String(),
			Raw:        html,
			CreateTime: ttime,
		}
		switch flag {
		case "赞了我的说说":
			mes.Type = TLike
		case "提到我":
			mes.Type = TAt
		case "评论提到我":
			mes.Type = TCommentAt
		case "评论":
			mes.Type = TComment
		default:
			mes.Type = TOther
		}
		message = append(message, mes)
		//fmt.Println(mes.TID, mes.QQ, mes.Type, mes.CreateTime)
	}
	return message, nil
}
