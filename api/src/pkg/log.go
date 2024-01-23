package pkg

import (
	"github.com/materkov/meme9/api/src/pkg/xlog"
)

func LogErr(e error) {
	if e != nil {
		xlog.Log("Logged error", xlog.Fields{
			"err": e.Error(),
		})
	}
}
