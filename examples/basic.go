package main

import (
	"context"
	"github.com/xybingbing/MessageEvents"
)

func main() {

	const mqUrl = "amqp://admin:admin@127.0.0.1:5672"

	event, err := MessageEvents.NewRabbitMQ(mqUrl, "basic", nil)
	if(err != nil){
		panic(err)
	}

	ctx := context.Background()

	//监听
	_ = event.Subscriber(ctx, "EventName", func(body interface{}) error {
		data := ""
		if err := event.Scan(body, &data); err != nil {
			println(err)
		}
		println(data)

		return nil
	})

	//发送消息
	body := "abc"
	_ = event.Publish(ctx, &MessageEvents.Event{
		EventName: "EventName",
		Body: body,
	})

}
