package sub

import (
	"github.com/xybingbing/MessageEvents/router"
	"strconv"
	"time"
)

func Retry(ctx *router.Context) error {
	msg := ctx.Message
	err := ctx.Next()
	if err != nil {
		runFailNum := 0
		if msg.Metadata.Get("RunFailNum") != "" {
			runFailNum, err = strconv.Atoi(msg.Metadata.Get("RunFailNum"))
			if err != nil {
				return err
			}
		}
		if runFailNum < 3 {
			runFailNum++
			msg.Metadata.Set("RunFailNum", strconv.Itoa(runFailNum))
			time.AfterFunc(time.Second*time.Duration(runFailNum*runFailNum), func() {
				ctx.MessageEngine.Run(ctx.Context, msg, ctx.Topic)
			})
			return err
		}
		//失败3次
		return router.ErrRetryFail
	}
	return err
}
