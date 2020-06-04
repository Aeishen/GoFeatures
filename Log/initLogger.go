package Log

import (
	"GoFeatures/Log/logger"
	"GoFeatures/Log/utils"
	"fmt"
	"runtime"
)

var MyLogger logger.MyLogger

// 读取配置文件的结构体,字段必须大写
type config struct {
	Level 		 string `ini:"file_level"`
	LogPath 	 string `ini:"file_path"`
	LogNormaName string `ini:"file_name"`
	Size         int64 `ini:"file_max_size"`
	Map1        map[int]int `ini:"file_max_map"`
	Bool1        map[int]int `ini:"file_max_bool"`
	Float1        map[int]int `ini:"file_max_float"`
}

func InitLogger(){
	fileName := "GoFeatures/Log/conf.ini"
	logConf := new(config)
	err := utils.ParseConfig(fileName, logConf)
	if err != nil {
		// todo what you want
		fmt.Println(err)
		return
	}

	// TODO
	if runtime.GOOS == "windows" {
		MyLogger = logger.NewConsoleLogger(logConf.Level)
	}
	MyLogger = logger.NewFileLogger(logConf.Level, logConf.LogPath, logConf.LogNormaName, logConf.Size)
}
