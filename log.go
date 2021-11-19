package log


import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type mylog struct {
	name  string
	debug bool
	*zap.Logger
}

func New(name string, debug bool) *mylog {
	l:=  &mylog{
		name:  name,
		debug: debug,
	}
	l.init()
	return l
}

func (log *mylog) init()  {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	cores := []zapcore.Core{}
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	waringLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	if log.debug {
		//debug 直接输出到终端中
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), infoLevel))
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), waringLevel))
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), errorLevel))
	} else {
		// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
		infoWriter := getWriter(log.name + "_info.log")
		waringWriter := getWriter(log.name + "_waring.log")
		errorWriter := getWriter(log.name + "_error.log")
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(&infoWriter)), infoLevel))
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(&waringWriter)), waringLevel))
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(&errorWriter)), errorLevel))
	}

	// 最后创建具体的Logger
	core := zapcore.NewTee(cores..., )
	caller := zap.AddCaller()
	development := zap.Development()
	logger := zap.New(core, caller,development)
	log.Logger = logger
}


func (log *mylog) Close() error {
	return nil
}


func (log *mylog) MError(err error,opt ...option){
	if err == nil {return }
	res,causeLine := log.split(err)
	register := &Option{}
	for _,o := range opt{
		o(register)
	}
	var fileds [] zap.Field
	fileds = append(fileds,register.registerOption(log)...)
	fileds = append(fileds,zap.String("app_name",log.name), zap.String("err_line",causeLine), zap.Any("err_stack",res))
	log.Error(err.Error(), fileds...)
}



func (log *mylog) MInfo(info string,opt ...option){
	register := &Option{}
	for _,o := range opt{
		o(register)
	}
	var fileds [] zap.Field
	fileds = append(fileds,register.registerOption(log)...)
	fileds = append(fileds,zap.String("app_name",log.name))
	log.Info(info, fileds...)
}



func (log *mylog)split(err error) ([]string ,string/*cause line*/){
	str := fmt.Sprintf("%+v",err)
	tem:= strings.Split(str,"\n")
	temCause :="can't get error stack "
	for i:=0;i<len(tem);i++ {
		if i ==0 {
			tem[i] = "ERR_REASON ："+tem[i]
		}
		if i== 2 {
			temCause = tem[i]
		}
	}
	return tem,temCause
}




func getWriter(filename string) lumberjack.Logger {
	today := time.Now().Format("20060102")
	filename = fmt.Sprintf("./logs/%s/%s", today, filename)
	return lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M  128
		MaxBackups: 30,       // 日志文件最多保存多少个备份 30
		MaxAge:     7,        // 文件最多保存多少天 7
		Compress:   true,     // 是否压缩
	}
}

