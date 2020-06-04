package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 获取调用信息
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


// 获取这个文件所在的路径(测试用)
func GetCurrentPath()string{
	_,filename,_,_:=runtime.Caller(1)
	return path.Dir(filename)
}

// 按大小检查文件
func CheckFileSize(file *os.File, maxSize int64)bool{
	fileInfo, err := file.Stat()
	if err != nil{
		panic("获取文件信息失败")
	}
	return fileInfo.Size() > maxSize
}

// 按时间检查文件
func CheckFileTime(file *os.File, maxSize int64)bool{
	// todo
	return false
}

// 分割文件
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


// 解析配置文件
func ParseConfig(fileName string, result interface{})error{

	t := reflect.TypeOf(result)    // 获取接口值动态类型
	tElem := t.Elem()              // 获取指针指向的元素/变量(等效于对指针类型变量做了一个*操作)
	v := reflect.ValueOf(result)   // 获取接口值动态值

	// 前提条件, result底层类型必须是个指针, 元素的底层类型必须是结构体
	if t.Kind() != reflect.Ptr || tElem.Kind() != reflect.Struct{
		return errors.New("result必须是个结构体指针")
	}

	// 打开文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil{
		err = fmt.Errorf("读取配置文件%s失败:%s", fileName,err)
		return err
	}
	// 将读取到的文件数据按行分割
	lineSlice := strings.Split(string(data),"\r\n")

	// 一行一行的解析
	for index, line := range lineSlice{
		line = strings.TrimSpace(line) // 去除字符串首尾的空白
		if len(line) == 0 || strings.HasPrefix(line,"#"){  // 空行与注释忽略
			continue
		}

		// 判断是否存在等号
		equalIndex := strings.Index(line,"=")
		if equalIndex == -1{
			err = fmt.Errorf("配置文件%s第%d行存在语法错误:%s", fileName,index+1,err)
			return err
		}
		// 安装=号分割每一行的key与value
		key := strings.TrimSpace(line[:equalIndex])
		val := strings.TrimSpace(line[equalIndex + 1:])
		if len(key) == 0 || len(val) == 0{
			err = fmt.Errorf("配置文件%s第%d行存在语法错误:%s", fileName,index+1,err)
			return err
		}
		// 利用反射给result赋值,遍历结构体每个字段与key比较,匹配上就给val赋值
		for i := 0; i < tElem.NumField(); i++{    // tElem.NumField():结构体内的字段数量
			field := tElem.Field(i)               // tElem.Field(i):根据索引获取结构体每个字段
			tag := field.Tag.Get("ini")      // 获取结构体每个字段的tag
			if key == tag{
				k := field.Type.Kind()            // 获取结构体每个字段的底层类型,判断对应类型进行赋值处理
				switch k{
				case reflect.String:
					v.Elem().Field(i).SetString(val)
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					val_, _ := strconv.ParseInt(val, 10, 64)
					v.Elem().Field(i).SetInt(val_)
				case reflect.Bool:
					val_, _ := strconv.ParseBool(val)
					v.Elem().Field(i).SetBool(val_)
				case reflect.Float32, reflect.Float64:
					val_, _ := strconv.ParseFloat(val, 64)
					v.Elem().Field(i).SetFloat(val_)
				default:
					err = fmt.Errorf("配置文件%s解析失败:结构体不可存在字段类型%s", fileName, k)
					return err
				}
			}
		}
	}
	return nil
}

