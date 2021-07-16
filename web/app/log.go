package app

import (
	"context"
	"fmt"

	"github.com/materkov/meme9/web/utils"
)

func Logf(ctx context.Context, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)

	reqID, ok := ctx.Value(utils.RequestIdKey{}).(int)
	if ok {
		msg = fmt.Sprintf("[ReqID %x] %s", reqID, msg)
	}

	go func() {
		_ = ObjectStore.WriteLog(reqID, msg)
	}()
}
