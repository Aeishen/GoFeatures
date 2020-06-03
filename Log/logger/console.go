// 往终端打印日志
package logger

import (
	"GoFeatures/Log/utils"
	"fmt"
	"os"
	"time"
)

type myConsoleLogger struct {
	logLevel Level
}

func NewConsoleLogger(level string) *myConsoleLogger {
	logger := new(myConsoleLogger)
	logger.logLevel = getLevelByString(level)
	return logger
}


func (l *myConsoleLogger)Debug(format string, a ...interface{}){
	l.log(DebugLevel,format, a...)
}

func (l *myConsoleLogger)Info(format string, a ...interface{}){
	l.log(InfoLevel,format, a...)
}

func (l *myConsoleLogger)Warn(format string, a ...interface{}){
	l.log(WarnLevel,format, a...)
}

func (l *myConsoleLogger)Error(format string, a ...interface{}){
	l.log(ErrorLevel,format, a...)
}

func (l *myConsoleLogger)Fatal(format string, a ...interface{}){
	l.log(FatalLevel,format, a...)
}

func (l *myConsoleLogger)log(curLevel Level,format string, a ...interface{}){
	if l.logLevel < curLevel{
		return
	}
	msg := fmt.Sprintf(format, a...)
	funcName, fileName, line := utils.GetCallInfo(3)
	nowTime := time.Now().Format("2006-01-02 15:04:05.000")
	logType := getLevelString(curLevel)
	msg = fmt.Sprintf("[%s][%s][%s:%d][%s]%s",nowTime,fileName,funcName,line,logType,msg)
	_, _ = fmt.Fprintln(os.Stdout, msg)
}