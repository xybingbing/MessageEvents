package main

import (
	"context"
	"fmt"
	"github.com/xybingbing/MessageEvents"
	"github.com/xybingbing/MessageEvents/examples/lib"
	"log"
	"time"
)


type Obj2 struct {
	Name    string
	Age     int
}

func main()  {

	_, closer := lib.NewJaegerTracer("MessageEvents2", "dev");
	defer closer.Close()


	const mqUrl = "amqp://admin:admin@127.0.0.1:5672"

	event, err := MessageEvents.NewRabbitMQ(mqUrl, "middleware", nil)
	if(err != nil){
		panic(err)
	}



	ctx := context.Background()

	data := Obj2{}
	event.Subscriber(ctx, "add.order", func(body interface{}) error {




		err := event.Scan(body, &data)
		if(err != nil){
			return err
		}

		log.Default().Println(data.Name)

		return fmt.Errorf("666")

		return nil

	})

	//
	//event.Subscriber(ctx, "add.order", func(body interface{}) error {
	//	err := event.Scan(body, &data)
	//	if(err != nil){
	//		return err
	//	}
	//
	//	println(data.Name)
	//
	//	return errors.New("出错了")
	//
	//})

	//var i = 0;
	//for {
	//
	//	obj := Obj{Name: "小明", Age: 18}
	//	event.Publish(ctx, &MessageEvents.Event{
	//		EventName: "add.order",
	//		Body: obj,
	//	})
	//
	//	//println(i)
	//
	//	break;
	//
	//	time.Sleep(time.Second * 2)
	//	i++
	//}

	time.Sleep(time.Second * 2)
	select {

	}
}