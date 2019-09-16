package logger

import (
	"runtime"
	"fmt"
	"log"
	"os"
	"time"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"strconv"
)

const  (
	USER_LOG = iota
	RUN_LOG
	ERR_LOG
	PRINT_LOG = 1
	FILE_LOG = 2
	MONGO_LOG = 4
)


type LogState int
type LogMethod int

type Logger struct {
	Method     LogMethod   // 无效则不输出
	Dir        string // 目录
	FilePrefix string // 日志的前缀
	current    string // 当前日期
	logger     *log.Logger
	mongo      *mgo.Session
}

type MongeInfo struct {
	//bongo.DocumentBase `bson:",inline"`
	Info interface{}
	Data interface{}
}

// 创建目录
func (e *Logger) CreateDir() bool{
	//创建目录
	err := os.MkdirAll(e.Dir, os.ModePerm)
	if err != nil {
		fmt.Println(0, "os.MkdirAll is err", err, e.Dir)
		if (e.Method | FILE_LOG ) == e.Method {
			e.Method -= FILE_LOG
		}
		return false
	}
	fmt.Println("os.MkdirAll is ok", e.Dir)
	return true
}

// 创建输出文件
func (e *Logger) CreateFile() bool {
	//创建文件
	f := fmt.Sprintf("%s//%s_%s.log", e.Dir, e.FilePrefix, e.current)
	logfile, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(0, "initLogger os.OpenFile is err", err.Error())
		if (e.Method | FILE_LOG ) == e.Method {
			e.Method -= FILE_LOG
		}
		return false
	}

	//日志绑定文件
	e.logger = log.New(logfile, "\r\n", log.Ldate|log.Ltime)
	return true
}

func (e *Logger)CreateConnection() bool{
	session, err :=  mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Println("logger.CreateConnection err", err.Error())
		if ( e.Method | MONGO_LOG ) == e.Method {
			e.Method -= MONGO_LOG
		}
		return false
	}
	if session == nil {
		log.Println("logger.CreateConnection session is nil")
		if ( e.Method | MONGO_LOG ) == e.Method {
			e.Method -= MONGO_LOG
		}
		return false
	}

	session.SetMode(mgo.Monotonic, true)
	e.mongo = session
	return true
}

//
func (e *Logger)RunInfo(skip int)(file string, line int, name string){
	//获取调用者、文件名、行信息
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}
	f := runtime.FuncForPC(pc)
	return file, line, f.Name()
}

// 输出, state
func (e *Logger) Out(state LogState, a interface{}) {
	// 无效则输出到屏
	if e.Method <= 0 {
		return
	}

	switch e.Method {
	case PRINT_LOG:
		e.Print(state, a)
	case FILE_LOG:
		e.File(state, a)
	case MONGO_LOG:
		e.Mongo(state,a)
	case PRINT_LOG + FILE_LOG:
		e.Print(state, a)
		e.File(state, a)
	case PRINT_LOG + MONGO_LOG:
		e.Print(state, a)
		e.Mongo(state,a)
	case FILE_LOG + MONGO_LOG:
		e.File(state, a)
		e.Mongo(state,a)
	default:
		e.Print(state, a)
		e.File(state, a)
		e.Mongo(state,a)
	}
}

// 输出到屏幕
func (e *Logger) Print(state LogState, a interface{}){

	tag := ""
	switch state {
	case USER_LOG:
		tag = "USER"
	case RUN_LOG:
		tag = "RUN"
	case ERR_LOG:
		tag = "ERROR"
	}

	file, line, name := e.RunInfo(3)
	log.Println(tag, file, line, name, a)
}

// 输出到文件
func (e *Logger) File(state LogState, a interface{}){

	tag := ""
	switch state {
	case USER_LOG:
		tag = "USER"
	case RUN_LOG:
		tag = "RUN"
	case ERR_LOG:
		tag = "ERROR"
	}

	// 判断创建文件夹
	if e.logger == nil {
		if b :=e.CreateDir(); !b {
			return
		}

	}

	// 按日生成文件
	s := time.Now().Format("2006_01_02")
	if s != e.current || e.logger == nil{
		e.current = s
		if b := e.CreateFile(); !b {
			return
		}
	}

	file, line, name := e.RunInfo(3)
	e.logger.Println(tag, file, line, name, a)
}

// 输出到Mongo
func (e *Logger) Mongo(state LogState, a interface{}){

	tag := ""
	file, line, name := e.RunInfo(3)

	var i interface{}
	m := map[string]interface{}{
		"file":file,
		"line":strconv.Itoa(line),
		"name":name,
		"time":time.Now().Format("2006-01-02 15:04:05.0700"),
		"info":a,
	}

	switch state {
	case USER_LOG:
		i = a
		tag = "user_log"
	case RUN_LOG:
		i = m
		tag = "run_log"
	case ERR_LOG:
		i = m
		tag = "err_log"
	}

	if b := e.CreateConnection(); !b{
		return
	}
	collection := e.mongo.DB("jiang_ido").C(tag)
	err := collection.Insert(i)
	if err != nil {
		log.Println(err.Error())
	}

}
