package RabbitMQ

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机名称
	Key       string // bind Key 名称
	Mqurl     string // 连接信息
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string, mqurl string) *RabbitMq {
	return &RabbitMq{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     mqurl,
	}
}

// Destroy 断开channel 和 connection
func (r *RabbitMq) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMq) failOnErr(ctx context.Context, err error, message string) {
	if err != nil {
		g.Log().Fatalf(ctx, "%s:%s", message, err)
	}
}

// NewRabbitMQSimple 创建简单模式下Rabbitmq实例
func NewRabbitMQSimple(ctx context.Context, queueName string) *RabbitMq {
	mqurl, _ := g.Config().Get(ctx, "rabbitmq.dns")
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "", mqurl.String())
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp091.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(ctx, err, "failed to connect rabbitmq")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(ctx, err, "failed to open a channel")

	return rabbitmq
}

// PublishSimple simple模式队列生产
func (r *RabbitMq) PublishSimple(ctx context.Context, message string) {
	// 调用channel发送消息到队列中
	err := r.channel.PublishWithContext(
		ctx,
		r.Exchange,
		r.QueueName,
		// 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		// 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp091.Persistent, // 单条消息持久化
		},
	)
	if err != nil {
		g.Log().Fatalf(ctx, "推送数据失败 ", err)
	}
}

// QueueDeclare 申请队列，如果队列不存在会自动创建，存在则跳过创建
func (r *RabbitMq) QueueDeclare(ctx context.Context) {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		true,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		g.Log().Error(ctx, "创建队列失败 ", err)
	}
}

// ConsumeSimple simple模式消费
func (r *RabbitMq) ConsumeSimple(ctx context.Context, msgItemChan chan amqp091.Delivery) {
	defer r.Destroy()
	err := r.channel.Qos(200, 0, false)
	if err != nil {
		g.Log().Error(ctx, "限流失败 ", err)
	}
	// 接收消息
	msgList, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		guid.S(),
		// 是否自动应答
		false,
		// 是否独有
		false,
		// 设置为true，表示 不能将同一个connection中生产者的消息传递给这个connection中的消费者
		false,
		// 是否阻塞
		false,
		// 额外参数
		nil,
	)
	if err != nil {
		g.Log().Error(ctx, "监听消息失败 ", err)
	}
	for msgItem := range msgList {
		msgItemChan <- msgItem
	}
}
