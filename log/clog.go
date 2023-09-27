package log

import (
	"context"
	"log"
	"os"

	"github.com/lazyironf4ur/go-infra/util"
)

type cLog struct {
	ctx context.Context
	log log.Logger
}

const trace_id = "[%s] "

var Log *cLog = new(cLog)

func init() {
	Log.log = *log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
}

func (clog *cLog) WithGinCtx() *cLog {
	clog.ctx = util.GetTraceCtx()
	return clog
}

func (clog *cLog) WithCtx(ctx context.Context) *cLog {
	clog.ctx = ctx
	return clog
}

func (clog *cLog) Info(msg string, args ...interface{}) {
	clog.log.Printf(trace_id+msg+"\n", clog.ctx.Value("trace_id").(string), args)
}
