package history

import (
	"context"
	"errors"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"github.com/google/go-querystring/query"
	"github.com/intchensc/qzone/api/common"
)

var (
	getQZoneHistory = "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?"
)

type HistoryAPI struct {
	*common.BaseAPI
}

// getQZoneHistoryList 获取QQ空间历史消息（限制本人），offset和count分别表示每次请求的偏移量和数目
func (h *HistoryAPI) getQZoneHistoryList(offset, count int64) ([]*QZoneHistoryItem, error) {
	if h.Err != nil {
		return nil, h.Err
	}
	// 匿名函数列表，完成子操作
	// decodeHtml 解码其中的html字符（例如\x3C）
	decodeHtml := func(dataStr string) string {
		// 1. 正则匹配 "\xHH" 的 16 进制编码部分
		re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)

		// 替换每个匹配项
		decoded := re.ReplaceAllStringFunc(dataStr, func(hex string) string {
			// 去掉 "\x" 前缀，并解析为整数
			hexValue, err := strconv.ParseInt(hex[2:], 16, 32)
			if err != nil {
				// 如果解析失败，保留原字符串
				return hex
			}
			// 转换为字符
			return string(rune(hexValue))
		})

		// 2. 去除反斜杠定义
		re2 := regexp.MustCompile(`\\+`)
		decoded = re2.ReplaceAllStringFunc(decoded, func(match string) string {
			if match == `\/` { // \/ -> /
				return `/`
			}
			return `` // 否则，去除反斜杠
		})

		return decoded
	}
	// extractHtml 提取其中的html部分
	extractHtml := func(parsed string) []string {
		//匹配 html:'(.*?)'
		//re := regexp.MustCompile(`html:'(.*?)'`)
		re := regexp.MustCompile(`html:'(.*?)',opuin`)
		matches := re.FindAllStringSubmatch(parsed, -1)
		htmls := make([]string, len(matches))
		for idx, match := range matches {
			htmls[idx] = match[1]
		}

		return htmls
	}
	// extractHistoryMsg 解析html代码，提取一条消息的数据
	extractHistoryMsg := func(html string) (*QZoneHistoryItem, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil, errors.New("parse history msg failed")
		}
		var item *QZoneHistoryItem
		doc.Find("li.f-s-s").Each(func(i int, s *goquery.Selection) {
			// sender qq
			senderQQ, _ := s.Find(".user-avatar").Attr("link")
			senderQQ = strings.TrimPrefix(senderQQ, "nameCard_")
			// 说说ID
			shuoshuoID, _ := s.Find("i[name='feed_data']").Attr("data-tid")
			// 说说消息中的图片
			var shuoshuoImgUrls []string
			s.Find(".f-ct-txtimg .img-box img").Each(func(j int, imgS *goquery.Selection) {
				attr, exists := imgS.Attr("src")
				if exists {
					shuoshuoImgUrls = append(shuoshuoImgUrls, attr)
				}
			})
			// 说说内容
			shuoshuoContent := s.Find(".f-ct-txtimg .txt-box .txt-box-title").Contents().FilterFunction(func(j int, txtS *goquery.Selection) bool {
				// 过滤掉 <a> 和 <span> 等子标签，只保留纯文本节点
				return goquery.NodeName(txtS) == "#text"
			}).Text()
			shuoshuoContent = strings.TrimSpace(shuoshuoContent)

			// createTime
			createTimeStr, _ := s.Find("i[name='feed_data']").Attr("data-abstime")
			createTime, _ := strconv.ParseInt(createTimeStr, 10, 64)
			// 互动类型
			actionType := s.Find(".f-nick .state").Text()

			// 互动内容
			comments := s.Find(".comments-content").Text()
			suffix := s.Find(".comments-content .comments-op").Text()
			if len(comments) > 0 {
				comments = strings.SplitN(comments, ": ", 2)[1]
				comments = strings.TrimSuffix(comments, suffix)
			}

			// 互动消息中的图片
			var imgUrls []string
			s.Find(".comments-content .comments-thumbnails img").Each(func(j int, imgS *goquery.Selection) {
				attrOnLoad, exists := imgS.Attr("onload")
				if exists {
					link := matchWithRegexp(attrOnLoad, `trueSrc:'(.*?)'`, true)
					if len(link) > 0 {
						imgUrls = append(imgUrls, link[0])
					}
				}
			})

			item = &QZoneHistoryItem{
				SenderQQ:        senderQQ,
				ActionType:      actionType,
				ShuoshuoID:      shuoshuoID,
				Content:         comments,
				CreateTime:      time.Unix(createTime, 0),
				ImgUrls:         imgUrls,
				ShuoshuoContent: shuoshuoContent,
				ShuoshuoImgUrls: shuoshuoImgUrls,
			}
		})
		return item, nil
	}

	// 请求历史消息数据
	data, err := h.queryQZoneHistoryList(offset, count)
	if err != nil {
		er := errors.New("QQ空间历史数据请求错误:" + err.Error())
		log.Println("QQ空间历史数据请求失败:" + er.Error())
		return nil, er
	}
	// 解码并提取HTML数据
	htmlSlice := extractHtml(decodeHtml(string(data)))
	items := make([]*QZoneHistoryItem, len(htmlSlice))
	// 分别对html切片中的每一条数据进行处理
	for idx, html := range htmlSlice {
		item, err := extractHistoryMsg(html)
		if err != nil {
			er := errors.New("QQ空间历史数据解析错误:" + err.Error())
			log.Println("QQ空间历史数据解析失败:" + er.Error())
			return nil, er
		}
		items[idx] = item
	}
	return items, nil
}

func (h *HistoryAPI) queryQZoneHistoryList(offset, count int64) (string, error) {
	if h.Err != nil {
		return "", h.Err
	}
	qzhr := QZoneHistoryReq{
		Uin:                h.Qq,
		Offset:             offset,
		Count:              count,
		BeginTime:          "",
		EndTime:            "",
		Getappnotification: "1",
		Getnotifi:          "1",
		HasGetKey:          "0",
		Useutf8:            "1",
		Outputhtmlfeed:     "1",
		Scope:              "1",
		Set:                "0",
		Format:             "json",
		Gtk:                h.Gtk,
	}
	headers := map[string]string{
		"cookie":                    h.Cookie,
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"authority":                 "user.qzone.qq.com",
		"pragma":                    "no-cache",
		"cache-control":             "no-cache",
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"sec-ch-ua":                 "\"Not A(Brand\";v=\"99\", \"Microsoft Edge\";v=\"121\", \"Chromium\";v=\"121\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"Content-Type":              "application/json; charset=utf-8",
	}
	// url_ := getQZoneHistory + structToStr(qzhr)
	var data string
	url := getQZoneHistory + common.StructToStr(qzhr)
	Vh, _ := query.Values(headers)
	err := requests.URL(url).
		Headers(Vh).UserAgent(common.UA).ContentType(common.ContentType).
		ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url_), WithHeader(headers)))
	if err != nil {
		// er := errors.New("请求历史消息数据错误:" + err.Error())
		log.Println("请求历史消息数据失败:", err)
		// return nil, err
	}
	// cookie过期或者发生了其他错误
	ans := matchWithRegexp(string(data), `"code":(.*?),`, true)
	if ans != nil {
		code, _ := strconv.Atoi(ans[0])
		if code == -3000 {
			er := errors.New("请求历史消息数据错误: cookie失效或其他错误")
			log.Println("请求历史消息数据失败:", er.Error())
			return "", er
		}
	}

	return data, nil
}

// GetQZoneHistoryList 获取本人QQ空间的所有历史消息
func (h *HistoryAPI) List() ([]*QZoneHistoryItem, error) {
	if h.Err != nil {
		return nil, h.Err
	}
	// 0. 函数中使用到的匿名函数
	// getTotal 获取历史消息总数
	getTotal := func() (int64, error) {
		var (
			low, high int64 = 0, math.MaxInt / 2
			total     int64 = 0
			count     int64 = 100
		)

		for low <= high {
			mid := (low + high) >> 1
			// 1. 请求数据
			data, err := h.queryQZoneHistoryList(mid*count, count)
			if err != nil {
				er := errors.New("QQ空间历史消息获取错误:" + err.Error())
				log.Println("QQ空间历史消息获取失败:", er.Error())
				return total, er
			}
			// 2. 解析数据
			ans := matchWithRegexp(string(data), `total_number:(.*?),`, true)
			if ans == nil {
				er := errors.New("QQ空间历史消息解析错误")
				log.Println("QQ空间历史消息解析失败:", er.Error())
				return total, er
			}
			num, _ := strconv.ParseInt(ans[0], 10, 64)
			if num <= 0 {
				high = mid - 1
			} else { // num > 0
				low = mid + 1
				total = mid*count + num
			}

			time.Sleep(2 * time.Second)
		}

		return total, nil
	}

	// 1. getTotal
	total, err := getTotal()
	if err != nil {
		return nil, err
	}
	// 2. 每次请求10条，并拼接结果
	totalItems := make([]*QZoneHistoryItem, total)

	idx := 0
	for i := 0; i <= int(total)/10; i++ {
		offset := i * 10
		items, err := h.getQZoneHistoryList(int64(offset), 10)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			if item != nil {
				totalItems[idx] = item
				idx++
			}
		}
	}

	return totalItems[:idx], nil
}

func matchWithRegexp(data, pattern string, extract bool) []string {
	re := regexp.MustCompile(pattern)
	matched := re.FindAllStringSubmatch(data, -1)
	if matched == nil {
		return nil
	}

	res := make([]string, len(matched))
	for i, match := range matched {
		if extract {
			res[i] = match[1]
		} else {
			res[i] = match[0]
		}
	}

	return res
}
