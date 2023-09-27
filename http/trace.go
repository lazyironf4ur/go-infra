package http

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lazyironf4ur/go-infra/conf"
	"github.com/lazyironf4ur/go-infra/util"
)

//traceid: 8位机器ip + 模块名 + timestamp + goroutineid

func InjectTrace() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		traceId := generateTraceId()
		i, _ := util.GetGoroutineID()
		util.TraceMap[i] = traceId
		ctx.Set("trace_id", traceId)
	}
}

func generateTraceId() string {
	traceId := formatIP() + "_"
	traceId += strconv.FormatInt(time.Now().UnixNano(), 16) + "_"
	moduleName := conf.GlobalConfig["module"]
	if moduleName != nil {
		for _, v := range conf.GlobalConfig["module"].(string) {
			traceId += strconv.FormatInt(int64(v), 16)
		}
	}

	i, err := util.GetGoroutineID()
	if err != nil {
		panic(err)
	}
	traceId += "_"
	traceId += strconv.FormatInt(int64(i), 16)
	return traceId
}

func formatIP() string {
	ip := getLocalIP()
	strs := strings.Split(ip, ".")
	ffip := ""
	for _, v := range strs {
		_ip, _ := strconv.Atoi(v)
		ffip += strconv.FormatInt(int64(_ip), 16)
	}

	return ffip
}

func getLocalIP() string {
	a, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, ipaddr := range a {
		if ip, ok := ipaddr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			s := ip.IP.To4().String()
			if s != "" {
				return s
			}
		}
	}
	return ""
}
