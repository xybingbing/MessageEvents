package pub

import "github.com/xybingbing/MessageEvents/router"

func Fail(ctx *router.Context) error {
	msg := ctx.Message
	err := ctx.Next()
	if err == router.ErrRetryFail {
		//TODO 3次失败了消息存储起来
		println(err.Error())
		println(msg.UUID)
	}

	return err
}
