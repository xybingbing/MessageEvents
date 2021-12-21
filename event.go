package MessageEvents

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/xybingbing/MessageEvents/middleware/pub"
	"github.com/xybingbing/MessageEvents/middleware/sub"
	"github.com/xybingbing/MessageEvents/router"
)


type MessageEvent struct {
	mqAddr string
	mqConn *amqp.ConnectionWrapper

	topic string
	logger watermill.LoggerAdapter

	metricsConfig *MetricsConfig
	metricsBuilder metrics.PrometheusMetricsBuilder

	pub message.Publisher
	sub message.Subscriber
	PubRouter *router.MessageEngine
	SubRouter *router.MessageEngine
}

type MetricsConfig struct {
	IsMetrics 	bool
	Addr		string
	Namespace	string
	Subsystem	string
}

func NewRabbitMQ(mqAddr string, topic string, metricsConfig *MetricsConfig) (*MessageEvent, error) {

	var logger = watermill.NewStdLogger(false, false);

	//创建mq连接
	amqpConfig := amqp.ConnectionConfig{AmqpURI: mqAddr}
	mqConn, err := amqp.NewConnection(amqpConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("create new connection: %w", err)
	}

	//监控指标
	var metrics_builder metrics.PrometheusMetricsBuilder
	if metricsConfig != nil && metricsConfig.IsMetrics == true {
		if metricsConfig.Addr == "" {
			metricsConfig.Addr = ":18081"
		}
		prometheusRegistry, closeMetricsServer := metrics.CreateRegistryAndServeHTTP(metricsConfig.Addr)
		defer closeMetricsServer()
		metrics_builder = metrics.NewPrometheusMetricsBuilder(prometheusRegistry, metricsConfig.Namespace, metricsConfig.Subsystem)
	}

	return &MessageEvent{
		mqAddr: mqAddr,
		mqConn: mqConn,
		topic: topic,
		logger: logger,
		metricsConfig: metricsConfig,
		metricsBuilder: metrics_builder,
	}, err
}

//事件
type Event struct {
	EventName string
	Body interface{}
}

func (event *MessageEvent) Publish(ctx context.Context, Event *Event) error {
	if event.pub == nil {
		//创建配置
		amqpConfig := amqp.NewDurablePubSubConfig(event.mqAddr, nil)
		amqpConfig.Publish.Transactional = true
		//创建消息发送者
		publisher, err := amqp.NewPublisherWithConnection(amqpConfig, event.logger, event.mqConn)
		if err != nil {
			return err
		}
		event.pub = publisher
		//开启metrics
		if event.metricsConfig != nil && event.metricsConfig.IsMetrics == true {
			m_publisher, err := event.metricsBuilder.DecoratePublisher(publisher)
			if err != nil {
				return err
			}
			event.pub = m_publisher
		}

		//创建路由
		event.PubRouter = router.NewMessageEngine()
		//增加生产消息的链路追踪中间价
		event.PubRouter.Use(pub.Tracer, pub.Fail, pub.Retry)
	}

	if Event != nil{

		//组装数据
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(Event.Body)
		if err != nil {
			return err
		}

		msg := message.NewMessage(watermill.NewUUID(), buf.Bytes())
		msg.Metadata.Set("EventName", Event.EventName)	//事件名称
		msg.Metadata.Set("FailNum", "0")			//发送失败次数

		event.PubRouter.AddRouteOne(Event.EventName, func(ctx *router.Context) error {
			return event.pub.Publish(event.topic, msg)
		})

		//执行
		event.PubRouter.Run(ctx, msg, event.topic)

	}

	return nil
}



//订阅事件
func (event *MessageEvent) Subscriber(ctx context.Context, EventName string, handle func(data interface{}) error) error {
	if event.sub == nil {
		//创建配置
		amqpConfig := amqp.NewDurablePubSubConfig(event.mqAddr, func(topic string) string {
			return topic
		})
		//创建消息接收者
		subscriber, err := amqp.NewSubscriberWithConnection(amqpConfig, event.logger, event.mqConn)
		if err != nil {
			return err
		}
		event.sub = subscriber

		//开启metrics
		if event.metricsConfig != nil && event.metricsConfig.IsMetrics == true {
			m_subscriber, err := event.metricsBuilder.DecorateSubscriber(subscriber)
			if err != nil {
				return err
			}
			event.sub = m_subscriber
		}
		//创建路由
		event.SubRouter = router.NewMessageEngine()
		//增加消费消息的链路追踪中间价, 错误处理
		event.SubRouter.Use(sub.Tracer, sub.Fail, sub.Retry)
		//消费消息
		messages, err := event.sub.Subscribe(ctx, event.topic)
		if err != nil {
			return err
		}
		go func() {
			for msg := range messages {
				event.SubRouter.Run(ctx, msg, event.topic)
				msg.Ack()
			}
		}()
	}

	event.SubRouter.AddRoutes(EventName, func(ctx *router.Context) error {
		return handle(ctx.Message.Payload)
	})

	return nil
}

//反序列化
func (event *MessageEvent) Scan (params interface{}, pointer interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(params.(message.Payload)))
	return dec.Decode(pointer)
}

//删除queues
func (event *MessageEvent) RemoveSubscriber(queues string) (int, error) {
	channel, err := event.mqConn.Connection().Channel()
	if err != nil {
		return 0, err
	}
	//强制删除消息队列的queues
	return channel.QueueDelete(queues, false, false, true)
}
