
通过消息处理事件



### 字符串消息
```go
//基础
event, err := MessageEvents.NewRabbitMQ("amqp://admin:admin@127.0.0.1:5672", "basic", nil)
if(err != nil){
    panic(err)
}

ctx := context.Background()

//监听
_ = event.Subscriber(ctx, "EventName", func(body interface{}) {
    data := ""
    if err := event.Scan(body, &data); err != nil {
        println(err)
    }
    println(data)
})

//发送消息
body := "abc"
_ = event.Publish(ctx, &MessageEvents.Event{
    EventName: "EventName",
    Body: body,
})

```

### Struct消息
```go

type People struct {
    Name    string
    Age     int
    Gender   string
    Describe string
}


event, err := MessageEvents.NewRabbitMQ("amqp://admin:admin@127.0.0.1:5672", "struct", nil)
if(err != nil){
    panic(err)
}

ctx := context.Background()

//监听
_ = event.Subscriber(ctx, "EventName", func(body interface{}) {
    people := People{}
    if err := event.Scan(body, &people); err != nil {
        println(err)
    }

    fmt.Printf("%v\n", people)
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


```