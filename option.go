package log


import "go.uber.org/zap"

type Option struct {
	traceID string
	errCode int
}


func (o *Option) registerOption (log  *mylog)[]zap.Field{
	var res []zap.Field
	if o.traceID != "" {
		res = append(res, zap.String("traceID",o.traceID))
	}
	if o.errCode != 0 {
		res= append(res,zap.Int("errcode",o.errCode))
	}
	return res
}

type option func(opt *Option)



func WithTraceID(traceID string) option{
	return func(opt *Option) {
		opt.traceID =traceID
	}
}


func WithErrCode(errCode int) option {
	return func(opt *Option) {
		opt.errCode =errCode
	}
}
