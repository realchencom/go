package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type RSLog struct {
	v *zap.SugaredLogger
}

var (
	Logger RSLog
)

func init() {
	logger := rsLog("", "info", "console")

	output := Con.GetString("log.output")
	if output == "" {
		logger.Warnf("未找到log.output配置，默认%s ", output)
		output = "console"
		Logger.v = logger
		return
	}
	//配置log.output后将重新生成日志实例，就的logger defer掉
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)
	var level, path string
	level = Con.GetString("log.level")
	if level == "" {
		level = "info"
		logger.Warnf("未找到log.level配置，默认%s ", level)
	}
	if output == "console" {
		Logger.v = rsLog("", level, output)
		goto finish
	}
	path = Con.GetString("log.path")
	if path == "" {
		path = "./main.log"
		logger.Warnf("未找到log.path配置，默认%s ", path)
	}
	Logger.v = rsLog(path, level, output)
finish:
	Logger.v.Info("Zap log init success")
}

func (rs *RSLog) Debug(template string, args ...interface{}) {
	rs.v.Debugf(template, args...)
}
func (rs *RSLog) Info(template string, args ...interface{}) {
	rs.v.Infof(template, args...)
}
func (rs *RSLog) Warn(template string, args ...interface{}) {
	rs.v.Warnf(template, args...)
}
func (rs *RSLog) Error(template string, args ...interface{}) {
	rs.v.Errorf(template, args...)
}
func (rs *RSLog) DPanic(template string, args ...interface{}) {
	rs.v.DPanicf(template, args...)
}
func (rs *RSLog) Panic(template string, args ...interface{}) {
	rs.v.Panicf(template, args...)
}
func (rs *RSLog) Fatal(template string, args ...interface{}) {
	rs.v.Fatalf(template, args...)
}

func (rs *RSLog) Sync() {
	err := rs.v.Sync()
	if err != nil {
		return
	}
}

// logpath 日志文件路径
// loglevel 日志级别
func rsLog(path string, level string, output string) *zap.SugaredLogger {
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",   //结构化（json）输出：msg的key
		LevelKey:    "level", //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:     "ts",    //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:   "file",  //结构化（json）输出：打印日志的文件对应的Key
		NameKey:     "name",
		FunctionKey: "func",
		//StacktraceKey: "trace",
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var logLevel zapcore.Level
	switch level {
	case "debug":
		//自定义日志级别：自定义Warn级别
		logLevel = zapcore.DebugLevel

	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "dpanic":
		logLevel = zapcore.DPanicLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	}
	var out zapcore.WriteSyncer
	if output == "console" {
		// 实现多个输出
		//core := zapcore.NewTee(
		//	zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel),                         //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		//	zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(warnWriter), warnLevel),                         //warn及以上写入errPath
		//	zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		//)
		out = zapcore.AddSync(os.Stdout)
	} else if output == "file" {
		hook := lumberjack.Logger{
			Filename:   path, // 日志文件路径
			MaxSize:    128,  // megabytes
			MaxBackups: 30,   // 最多保留30个备份
			MaxAge:     7,    // days
			Compress:   true, // 是否压缩 disabled by default
		}
		out = zapcore.AddSync(&hook)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		out,
		logLevel,
	)
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	return log.Sugar()
}
