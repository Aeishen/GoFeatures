// 往文件写日志
package logger

import (
	"GoFeatures/Log/utils"
	"fmt"
	"os"
	"path"
	"time"
)

type myFileLogger struct {
	logLevel Level
	logFilePath string
	logNormalFileName string
	logNormalFile *os.File
	logErrorFile *os.File
	logMaxSize int64
}

func NewFileLogger(level string, logPath string, logNormaName string, size int64) *myFileLogger {
	logger := new(myFileLogger)
	logger.logLevel = getLevelByString(level)
	logger.logNormalFileName = logNormaName
	logger.logFilePath = logPath
	logger.logMaxSize = size
	logger.initLogFile()
	return logger
}

func (l *myFileLogger)initLogFile(){
	normalFileName := path.Join(l.logFilePath,l.logNormalFileName)
	normalFile, err := os.OpenFile(normalFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("打开日志文件%s失败%v",normalFileName,err))
	}
	l.logNormalFile = normalFile

	errorFileName := normalFileName + "Err"
	errorFile, err := os.OpenFile(errorFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("打开日志文件%s失败:%v",errorFileName,err))
	}
	l.logErrorFile = errorFile
}

func (l *myFileLogger)Debug(format string, a ...interface{}){
	l.log(DebugLevel,format, a...)
}

func (l *myFileLogger)Info(format string, a ...interface{}){
	l.log(InfoLevel,format, a...)
}

func (l *myFileLogger)Warn(format string, a ...interface{}){
	l.log(WarnLevel,format, a...)
}

func (l *myFileLogger)Error(format string, a ...interface{}){
	l.log(ErrorLevel,format, a...)
}

func (l *myFileLogger)Fatal(format string, a ...interface{}){
	l.log(FatalLevel,format, a...)
}

func (l *myFileLogger)log(curLevel Level,format string, a ...interface{}){
	if l.logLevel < curLevel{
		return
	}
	msg := fmt.Sprintf(format, a...)
	funcName, fileName, line := utils.GetCallInfo(3)
	nowTime := time.Now().Format("2006-01-02 15:04:05.000")
	logType := getLevelString(curLevel)
	msg = fmt.Sprintf("[%s][%s][%s:%d][%s]%s",nowTime,fileName,funcName,line,logType,msg)

	if utils.CheckFileSize(l.logNormalFile, l.logMaxSize){
		l.logNormalFile = utils.SplitLogFile(l.logNormalFile)
	}
	_, _ = fmt.Fprintln(l.logNormalFile, msg)

	if curLevel >= ErrorLevel{
		if utils.CheckFileSize(l.logErrorFile, l.logMaxSize){
			l.logErrorFile = utils.SplitLogFile(l.logErrorFile)
		}
		_, _ = fmt.Fprintln(l.logErrorFile, msg)
	}
}