package work

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/rabbitmq/amqp091-go"
	"wechat-template-msg/utility/wechat/miniprogram/subscribe"
	"wechat-template-msg/utility/wechat/officialaccount/message"
)

type WechatTemplateMsg struct {
	WechatId                       int `json:"wechat_id"`
	WechatType                     int `json:"wechat_type"`
	Status                         int
	TableName                      string `json:"table_name"`
	Errmsg                         string
	AddTime                        string
	PushData                       string
	OfficialaccountTemplateMessage message.TemplateMessage `json:"officialaccount_template_message"`
	MiniprogramSubscribeMessage    subscribe.Message       `json:"miniprogram_subscribe_message"`
}

// GetMsgItemServiceResult 获取业务消息结果
func GetMsgItemServiceResult(msgItem amqp091.Delivery) (err error) {
	// 解析数据
	var wechatTemplateMsg WechatTemplateMsg
	err = gjson.Unmarshal(msgItem.Body, &wechatTemplateMsg)
	if err != nil {
		return gerror.New(" 解析消息内容错误，请重试 " + string(msgItem.Body) + err.Error())
	}
	wechatTemplateMsg.Status = 2
	wechatTemplateMsg.AddTime = gtime.Now().Format("Y-m-d H:i:s")
	// 发送请求
	switch wechatTemplateMsg.WechatType {
	case 1:
		// 公众号模板
		err = wechatTemplateMsg.officialaccountTemplateSend()
		break
	case 2:
		// 小程序模板
		err = wechatTemplateMsg.miniprogramSubscribeSend()
		break
	default:
		err = gerror.New("错误的公众号类型，请检查 " + string(msgItem.Body))
	}

	if err != nil {
		wechatTemplateMsg.Status = 3
		wechatTemplateMsg.Errmsg = err.Error()
	}

	// 记录结果到mysql
	err = wechatTemplateMsg.AddToMysql()
	return
}

// HandleConsumerMsgResult 处理消费者消息结果
func HandleConsumerMsgResult(ctx context.Context, msgItem amqp091.Delivery, msgErr error) {
	if msgErr != nil {
		g.Log().Error(ctx, msgErr)
	}

	// 确认消息
	err := msgItem.Ack(false)
	if err != nil {
		g.Log().Error(ctx, "确认消息失败：", err)
	}
}

// 公众号模板发送
func (wtm *WechatTemplateMsg) officialaccountTemplateSend() (err error) {
	wtm.PushData = gjson.New(wtm.OfficialaccountTemplateMessage).MustToJsonString()
	// 公众号模板
	_, err = message.Send(wtm.WechatId, &wtm.OfficialaccountTemplateMessage)
	return
}

// 小程序订阅发送
func (wtm *WechatTemplateMsg) miniprogramSubscribeSend() (err error) {
	wtm.PushData = gjson.New(wtm.MiniprogramSubscribeMessage).MustToJsonString()
	// 公众号模板
	err = subscribe.Send(wtm.WechatId, &wtm.MiniprogramSubscribeMessage)
	return
}

// AddToMysql 添加入库
func (wtm WechatTemplateMsg) AddToMysql() (err error) {
	if wtm.TableName == "" {
		return gerror.New("表命名不存在")
	}
	_, err = g.DB("db_wechat").Model(wtm.TableName).Insert(g.Map{
		"wechat_id":   wtm.WechatId,
		"wechat_type": wtm.WechatType,
		"status":      wtm.Status,
		"push_data":   wtm.PushData,
		"errmsg":      wtm.Errmsg,
		"add_time":    wtm.AddTime,
	})
	return
}
