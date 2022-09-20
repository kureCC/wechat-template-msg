package wechat

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var WechatConfigList gmap.IntStrMap

// GetAccessToken 获取accessToken
func GetAccessToken(wechatId int) (assessToken string, err error) {
	assessToken = WechatConfigList.Get(wechatId)
	if len(assessToken) == 0 {
		err = gerror.New("assessToken获取失败")
	}
	return
}

// HandleWechatConfig 处理微信配置
func HandleWechatConfig(ctx context.Context) {
	all, err := g.
		DB("db_wechat").
		Model("t_wechat_conf").
		Fields("wechat_id,access_token").All()
	if err != nil {
		g.Log().Error(ctx, "微信配置查询失败")
	}
	wechatConfigList := all.RecordKeyInt("wechat_id")
	for wechatId, item := range wechatConfigList {
		WechatConfigList.Set(wechatId, item["access_token"].String())
	}
}
