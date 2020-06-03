package utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func GetCallInfo(skip int) (funcName string, fileName string, lineCount int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}


	fileName = path.Base(file)               // 从file(文件全路径)中剥离出文件名
	funcName = runtime.FuncForPC(pc).Name()  // 根据pc获取函数名
	lineCount = line
	return
}

func CheckFileSize(file *os.File, maxSize int64)bool{
	fileInfo, err := file.Stat()
	if err != nil{
		panic("获取文件信息失败")
	}
	return fileInfo.Size() > maxSize
}


func CheckFileTime(file *os.File, maxSize int64)bool{
	// todo
	return false
}


func SplitLogFile(file *os.File) (newFileObj *os.File){
	fileName := file.Name()                 // 完整文件路径
	err := file.Close()
	if err != nil{
		panic("关闭文件失败")
	}
	newFileName := fmt.Sprintf("%s_%v",fileName,time.Now().Unix())
	err = os.Rename(fileName, newFileName)  // 备份原来的文件
	if err != nil{
		panic("重命名文件失败")
	}
	newFileObj, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		panic("打开新的文件失败")
	}
	return newFileObj
}