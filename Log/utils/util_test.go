package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetCallInfo(t *testing.T) {
	funcName, fileName, lineCount := GetCallInfo(1)
	fmt.Printf("funcName = %v, fileName = %v, lineCount = %v",funcName, fileName, lineCount)
}

func TestGetCurrentPath(t *testing.T) {
	path := GetCurrentPath()
	fmt.Printf("path = %v",path)
}

func TestCheckFileSize(t *testing.T) {
	path := GetCurrentPath()
	fileName := path + "/testFile.ini"
	fileObj, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	b1 := CheckFileSize(fileObj, 1)
	b2 := CheckFileSize(fileObj, 10000000000000)
	fmt.Printf("b1 = %v,b2 = %v",b1,b2)
}

func TestCheckFileTime(t *testing.T) {
	// todo
}

func TestSplitLogFile(t *testing.T) {
	path := GetCurrentPath()
	fileName := path + "/testFile.ini"
	fileObj, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	SplitLogFile(fileObj)
}

// 读取配置文件的结构体,字段必须大写
type config struct {
	Level 		 string `ini:"file_level"`
	LogPath 	 string `ini:"file_path"`
	LogNormaName string `ini:"file_name"`
	Size         int64 `ini:"file_max_size"`
	TestMap      map[int]int `ini:"file_max_map"`
	TestBool     bool `ini:"file_max_bool"`
	TestFloat    float32 `ini:"file_max_float"`
}
func TestParseConfig(t *testing.T) {
	path := GetCurrentPath()
	fileName := path + "/testFile.ini"
	c := new(config)
	err := ParseConfig(fileName, c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("config = %#v",c)
}
