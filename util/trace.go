package util

import "context"

type traceKey string

const(
	Trace_key traceKey = "trace_id"
)

var(
	TraceMap = make(map[int]string)
)

func GetTraceCtx() context.Context{
	i, err := GetGoroutineID()
	if err != nil {
		panic(err)
	}
	traceId, ok := TraceMap[i]
	if !ok {
		return context.Background()
	}

	return context.WithValue(context.Background(), Trace_key, traceId)
}