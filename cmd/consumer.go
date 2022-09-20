package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/panjf2000/ants/v2"
	"github.com/rabbitmq/amqp091-go"
	"net/http"
	"time"
	"wechat-template-msg/consts"
	"wechat-template-msg/utility/RabbitMQ"
	"wechat-template-msg/utility/wechat"
	"wechat-template-msg/work"
)

var (
	Consumer = cConsumer{}
)

type cConsumer struct {
	g.Meta         `name:"consumer" brief:"start the consumer"`
	queueName      string
	worksNumber    int
	consumerNumber int
	msgItemChan    chan amqp091.Delivery
}

type cConsumerInput struct {
	g.Meta `name:"consumer"`
}
type cConsumerOutput struct{}

// Run 消费者执行
func (c *cConsumer) Run(ctx context.Context, in cConsumerInput) (out *cConsumerOutput, err error) {
	// 初始化
	c.init(ctx)

	gt := gtimer.New()
	// 每分钟缓存wechatConfig
	wechat.HandleWechatConfig(ctx)
	gt.Add(ctx, time.Minute, func(ctx context.Context) {
		wechat.HandleWechatConfig(ctx)
	})

	// 运行消费者监听
	_ = ants.Submit(func() {
		c.runConsumer(ctx)
	})

	// 并行执行处理消息
	_ = ants.Submit(func() {
		c.runWork()
	})
	g.Log().Info(ctx, " [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 初始化
func (c *cConsumer) init(ctx context.Context) {
	c.queueName = consts.QueueName
	worksNumber, _ := g.Cfg().Get(ctx, "server.worksNumber")
	c.worksNumber = worksNumber.Int()
	consumerNumber, _ := g.Cfg().Get(ctx, "server.consumerNumber")
	c.consumerNumber = consumerNumber.Int()
	c.msgItemChan = make(chan amqp091.Delivery, c.worksNumber*100)

	go func() {
		servePort, _ := g.Cfg().Get(ctx, "server.port")
		g.Log().Info(ctx, http.ListenAndServe(":"+servePort.String(), nil))
	}()
}

// 并行处理消息
func (c *cConsumer) runConsumer(ctx context.Context) {
	// 异步监听消息
	for i := 0; i < c.consumerNumber; i++ {
		_ = ants.Submit(func() {
			// 连接mq
			rabbitmq := RabbitMQ.NewRabbitMQSimple(ctx, c.queueName)
			rabbitmq.QueueDeclare(ctx)
			// 接收消息，写入到chan中
			rabbitmq.ConsumeSimple(ctx, c.msgItemChan)
		})
	}
}

// 并行处理消息
func (c *cConsumer) runWork() {
	for i := 0; i < c.worksNumber; i++ {
		_ = ants.Submit(func() {
			ctx := gctx.New()
			for msgItem := range c.msgItemChan {
				// 获取业务消息结果
				msgErr := work.GetMsgItemServiceResult(msgItem)
				// 处理消费者消息结果
				work.HandleConsumerMsgResult(ctx, msgItem, msgErr)
			}
		})
	}
}
