package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lazyironf4ur/go-infra/util"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"

	trace_format = " [%s] "
)

type gLogger struct {
	logger.Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

var (
	glog = new(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
)

func new(writer logger.Writer, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &gLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (g *gLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *g
	newlogger.LogLevel = level
	return &newlogger
}

func (g gLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Info {
		g.Printf(g.infoStr+formatTraceId(ctx)+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g gLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		g.Printf(g.warnStr+formatTraceId(ctx)+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g gLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.Printf(g.errStr+formatTraceId(ctx)+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g gLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			g.Printf(formatTraceId(ctx)+g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(formatTraceId(ctx)+g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			g.Printf(formatTraceId(ctx)+g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(formatTraceId(ctx)+g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.Printf(formatTraceId(ctx)+g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(formatTraceId(ctx)+g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func formatTraceId(ctx context.Context) string {
	a := ctx.Value(util.Trace_key)
	return fmt.Sprintf(trace_format, a.(string))
}
