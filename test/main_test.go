package test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/intchensc/qzone"
	"github.com/intchensc/qzone/auth"
)

func TestFriendAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试好友列表
	friends, err := qm.API.Friend().List()
	if err != nil {
		t.Fatal("获取好友列表失败:", err)
	}
	t.Logf("获取到%d个好友", len(friends))

	// 测试好友详情
	if len(friends) > 0 {
		_, err := qm.API.Friend().Detail(friends[0].Uin)
		if err != nil {
			t.Fatal("获取好友详情失败:", err)
		}
	}
}

func TestGroupAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试群列表
	groups, err := qm.API.Group().List()
	if err != nil {
		t.Fatal("获取群列表失败:", err)
	}
	t.Logf("获取到%d个群", len(groups))

	// 测试群成员
	if len(groups) > 0 {
		_, err := qm.API.Group().MemberList(groups[0].GroupCode)
		if err != nil {
			t.Fatal("获取群成员失败:", err)
		}
	}
}

// TODO: 重写历史记录
func TestHistoryAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试历史消息
	history, err := qm.API.History().List()
	if err != nil {
		t.Fatal("获取历史消息失败:", err)
	}
	t.Logf("获取到%d条历史消息", len(history))
}

func TestShuoShuoAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试说说数量
	count, err := qm.API.ShuoShuo().Count(qm.API.Qq)
	if err != nil {
		t.Fatal("获取说说数量失败:", err)
	}
	t.Logf("说说总数:%d", count)

	// 测试说说列表
	if count > 0 {
		ss, err := qm.API.ShuoShuo().List(qm.API.Qq, 3, 1000)
		if err != nil {
			t.Fatal("获取说说列表失败:", err)
		}
		t.Logf("获取到说说:%v", ss[0].Content)
	}

	// 测试最新说说
	if count > 0 {
		_, err := qm.API.ShuoShuo().Latest(qm.API.Qq)
		if err != nil {
			t.Fatal("获取最新说说失败:", err)
		}
	}

	// 测试说说评论
	if count > 0 {
		list, err := qm.API.ShuoShuo().List(qm.API.Qq, 1, 0)
		if err != nil {
			t.Fatal("获取说说列表失败:", err)
		}
		if len(list) > 0 {
			_, err := qm.API.ShuoShuo().CommentList(list[0].Tid, 10, 0)
			if err != nil {
				t.Fatal("获取说说评论失败:", err)
			}
		}
	}
}

// 测试发布纯文本说说
func TestPublishTextShuoShuo(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	resp, err := qm.API.ShuoShuo().Publish("测试纯文本说说", nil)
	if err != nil {
		t.Fatal("发布纯文本说说失败:", err)
	}
	t.Logf("发布成功, tid:%v", resp)
}

// 测试发布带图片说说
func TestPublishImageShuoShuo(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 从test目录读取图片数据
	imagePath := "./1.png"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatal("读取图片文件失败:", err)
	}
	images := []string{base64.StdEncoding.EncodeToString(imageData)}
	resp, err := qm.API.ShuoShuo().Publish("测试带图片说说", images)
	if err != nil {
		t.Fatal("发布带图片说说失败:", err)
	}
	t.Logf("发布成功, tid:%v", resp)
}
