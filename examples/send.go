package main

import (
	"context"
	"github.com/xybingbing/MessageEvents"
	"github.com/xybingbing/MessageEvents/examples/lib"
)

type Obj struct {
	Name    string
	Age     int
}

func main()  {

	_, closer := lib.NewJaegerTracer("MessageEvents", "dev");
	defer closer.Close()

	const mqUrl = "amqp://admin:admin@127.0.0.1:5672"

	event, err := MessageEvents.NewRabbitMQ(mqUrl, "middleware", nil)
	if(err != nil){
		panic(err)
	}

	ctx := context.Background()

	obj := Obj{Name: "小明", Age: 18}
	event.Publish(ctx, &MessageEvents.Event{
		EventName: "add.order",
		Body: obj,
	})

	select {

	}
}
