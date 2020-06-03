package main

import (
	"GoFeatures/Log/logger"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

var MyLogger logger.MyLogger
type config struct {
	level 		 string `ini:"file_level"`
	logPath 	 string `ini:"file_path"`
	logNormaName string `ini:"file_name"`
	size         int64 `ini:"file_max_size"`
}

func InitLogger(){
	logConf := new(config)

	// TODO
	if runtime.GOOS == "windows" {
		MyLogger = logger.NewConsoleLogger(logConf.level)
	}
	MyLogger = logger.NewFileLogger(logConf.level, logConf.logPath, logConf.logNormaName, logConf.size)
}

func parseConfig(fileName string, c interface{})error{
	// 打开文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil{
		err = fmt.Errorf("打开日志文件%s失败:%s", fileName,err)
		return err
	}
	// 将读取到的文件数据按行分割
	lineSlice := strings.Split(string(data),"\r\n")

	// 一行一行的解析
	for index, line := range lineSlice{
		line = strings.TrimSpace(line) // 去除字符串首尾的空白
		if len(line) == 0 && strings.HasPrefix(line,"#"){  // 空行忽略注释
			continue
		}
	}
	return nil
}

func main() {
	fileName := "GoFeatures/Log/conf.ini"
	err := parseConfig(fileName, false)
	if err != nil {
		fmt.Println(err)
	}
}