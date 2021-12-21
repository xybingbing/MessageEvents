package main

import (
	"context"
	"fmt"
	"github.com/xybingbing/MessageEvents"
)

type People struct {
	Name    string
	Age     int
	Gender   string
	Describe string
}

func main() {

	const mqUrl = "amqp://admin:admin@127.0.0.1:5672"

	event, err := MessageEvents.NewRabbitMQ(mqUrl, "struct", nil)
	if(err != nil){
		panic(err)
	}

	ctx := context.Background()

	//监听
	_ = event.Subscriber(ctx, "EventName", func(body interface{}) error {
		people := People{}
		if err := event.Scan(body, &people); err != nil {
			println(err)
		}

		fmt.Printf("%v\n", people)

		return nil
	})

	//发送消息
	body := People{
		Name: "小明",
		Age: 18,
		Gender: "男",
		Describe: "是一个好人",
	}
	_ = event.Publish(ctx, &MessageEvents.Event{
		EventName: "EventName",
		Body: body,
	})

}
