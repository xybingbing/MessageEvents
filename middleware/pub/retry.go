package pub

import (
	"github.com/xybingbing/MessageEvents/router"
	"strconv"
	"time"
)

func Retry(ctx *router.Context) error {
	msg := ctx.Message
	err := ctx.Next()
	if err != nil {
		SendFailNum := 0
		if msg.Metadata.Get("SendFailNum") != "" {
			SendFailNum, err = strconv.Atoi(msg.Metadata.Get("SendFailNum"))
			if err != nil {
				return err
			}
		}
		if SendFailNum < 3 {
			SendFailNum++
			msg.Metadata.Set("SendFailNum", strconv.Itoa(SendFailNum))
			time.AfterFunc(time.Second * time.Duration(SendFailNum), func() {
				ctx.MessageEngine.Run(ctx.Context, msg, ctx.Topic)
			})
			return err
		}
		//重试发送3次失败了
		return router.ErrRetryFail
	}
	return err
}
