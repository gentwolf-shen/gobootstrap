package extra

import (
	"gobootstrap/logger"
	"runtime"
)

var (
	Recover = &RecoverHelper{}
)

type RecoverHelper struct{}

func (r *RecoverHelper) Process(err interface{}) {
	if err == nil {
		return
	}

	_, file, line, _ := runtime.Caller(1)
	logger.Sugar.Infof("%s : %d", file, line)
	logger.Sugar.Errorf("Recovered in %v", err)
}
