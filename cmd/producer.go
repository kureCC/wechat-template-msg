package cmd

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"wechat-template-msg/consts"
	"wechat-template-msg/utility/RabbitMQ"
)

var (
	Producer = cProducer{}
)

type cProducer struct {
	g.Meta    `name:"producer" brief:"start the producer"`
	queueName string
}

type cProducerInput struct {
	g.Meta `name:"producer"`
}
type cProducerOutput struct{}

func (c *cProducer) Run(ctx context.Context, in cProducerInput) (out *cProducerOutput, err error) {
	c.queueName = consts.QueueName
	// 连接mq
	rabbitmq := RabbitMQ.NewRabbitMQSimple(ctx, c.queueName)
	defer rabbitmq.Destroy()
	rabbitmq.QueueDeclare(ctx)
	for i := 0; i < 2; i++ {
		msg := getOfficialaccountTemplateUrlMessage()
		rabbitmq.PublishSimple(ctx, msg)
	}
	return
}

// 获取公众号模板url消息
func getOfficialaccountTemplateUrlMessage() (msg string) {
	msg = `
{
    "wechat_id": 1,
    "wechat_type": 1,
    "table_name": "t_wechat_template_msg",
    "officialaccount_template_message": {
        "touser": "openid",
        "template_id": "1tYCGsxxxxxxWpTBJ08Bkn17Sb6SA",
        "url": "https://www.baidu.com/",
        "data": {
            "first": {
                "value": "first！",
                "color": "#173177"
            },
            "keyword1": {
                "value": "keyword1",
                "color": "#173177"
            },
            "keyword2": {
                "value": "keyword2",
                "color": "#173177"
            },
            "remark": {
                "value": "remark！",
                "color": "#173177"
            }
        }
    }
}
`
	msg = gjson.New(msg).MustToJsonString()
	return
}

// 获取公众号模板小程序消息
func getOfficialaccountTemplateMiniprogramMessage() (msg string) {
	msg = `
{
    "wechat_id": 1,
    "wechat_type": 1,
    "table_name": "t_wechat_template_msg",
    "officialaccount_template_message": {
        "touser": "openid",
        "template_id": "1tYCGxxxxxBkn17Sb6SA",
        "miniprogram": {
            "appid": "appid",
            "pagepath": "pages/xxx/xxx/xxx"
        },
        "data": {
            "first": {
                "value": "first！",
                "color": "#173177"
            },
            "keyword1": {
                "value": "keyword1",
                "color": "#173177"
            },
            "keyword2": {
                "value": "keyword2",
                "color": "#173177"
            },
            "remark": {
                "value": "remark！",
                "color": "#173177"
            }
        }
    }
}
`
	msg = gjson.New(msg).MustToJsonString()
	return
}

// 获取小程序模板消息
func getMiniprogram15SubscribeMessage() (msg string) {
	msg = `
{
    "wechat_id": 15,
    "wechat_type": 2,
    "table_name": "t_wechat_template_msg",
    "miniprogram_subscribe_message": {
        "data": {
            "time3": {
                "color": "",
                "value": "2022-08-04 15:39:00"
            },
            "thing1": {
                "color": "",
                "value": "第二个直播间"
            },
            "thing8": {
                "color": "",
                "value": "近期有新直播啦，快去平台预约吧"
            }
        },
        "lang": "zh_CN",
        "page": "/pages/xxxx/xxx",
        "touser": "openid",
        "template_id": "yMXLZxxxxxxpu5WRsiao",
        "miniprogram_state": "trial"
    }
}
`
	msg = gjson.New(msg).MustToJsonString()
	return
}
