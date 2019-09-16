//go build -ldflags -H=windowsgui ido.go
package main

import (
	"ido/client"
	"ido/logger"
	"ido/model"
	"ido/push"
	"ido/session"
	"ido/socket"
	"net/http"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"strconv"
)

type Config struct {
	CertFile        string `json:"cert_file"`        // ssl cert file, /etc/letsencrypt/live/muutr.com/fullchain.pem
	KeyFile         string `json:"key_file"`         // ssl key file, /etc/letsencrypt/live/muutr.com/privkey.pem
	Debugging       bool   `json:"debugging"`        // 是否为调试状态，调试状态不作权限判断
	SessionDuration int    `json:"session_duration"` // Session的有效时间
	LogDir          string `json:"log_dir"`          // 日志的保存目录
	LogFilePrefix   string `json:"log_file_prefix"`  // 日志的保存文件名前缀
	LogMethod       int    `json:"log_method"`       // 日志输出模式，0不输出，1输出到文件，2输出到窗口
}

var (
	hub    *socket.Hub    // WebSocket
	config *Config        // IdoConfig
	root   string         // IdoRootPath
	Log    *logger.Logger // LogOut
)

func main() {

	root = getCurrPath()

	// 读取配置文件
	if err := initConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
		return
	}

	// 启动日志输出
	Log = new(logger.Logger)
	Log.Method = logger.LogMethod(config.LogMethod)
	Log.Dir = config.LogDir
	Log.FilePrefix = config.LogFilePrefix
	Log.Out(logger.RUN_LOG, "Ido Server is running!")
	model.Log = Log
	session.Log = Log
	push.Log = Log
	socket.Log = Log
	client.Log = Log

	model.Debugging = config.Debugging

	// web socket
	hub = socket.NewHub()

	// 启动其它线程, 同时监听多端口
	otherThread()

	// http
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/jpart/", partHtml)
	mux.HandleFunc("/html/admin/", adminHtml)
	mux.HandleFunc("/back.php", backHandle)
	mux.HandleFunc("/zip.php", zipHandle)
	mux.HandleFunc("/upload.php", uploadHandle)
	mux.HandleFunc("/chat.php", chatHandle)
	mux.HandleFunc("/admin.php", adminHandle)

	if len(config.CertFile) == 0 {
		config.CertFile = "/etc/letsencrypt/live/muutr.com/fullchain.pem"
	}else{
		config.CertFile = root + config.CertFile
	}
	if len(config.KeyFile) == 0 {
		config.KeyFile = "/etc/letsencrypt/live/muutr.com/privkey.pem"
	}else{
		config.KeyFile =  root + config.KeyFile
	}
	fmt.Println("ListenAndServe port is 443", config.CertFile, config.KeyFile)
	if  err := http.ListenAndServeTLS(":443", config.CertFile, config.KeyFile, mux); err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
	}
}

func otherThread() {
	// 临听80；
	go func() {
		fmt.Println("ListenAndServe port is 80")
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			url := "https://" + request.Host + request.RequestURI
			Log.Out(logger.RUN_LOG, url)
			http.Redirect(writer, request, url, http.StatusMovedPermanently)
		})
		if err := http.ListenAndServe(":80", mux); err != nil{
			Log.Out(logger.ERR_LOG, err.Error())
		}
	}()
	// 聊天室
	go hub.Run()
	// Session更新
	go session.Update(time.Duration(config.SessionDuration * 1000 * 1000 * 1000))
}

// 获取当前目录
func getCurrPath() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

// 初始配制文件
func initConfig() error {
	// 读当前目录下的ido.config
	b, err := ioutil.ReadFile(root + "/ido.config")
	if err != nil {
		return err
	}

	if len(b) == 0 {
		return errors.New("initConfig ido.config is nil")
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return err
	}
	return nil
}

// 启动，解析表单及多部表单，建立Session
func start(w http.ResponseWriter, r *http.Request) int {
	// 解析传入的数据
	r.ParseForm()
	r.ParseMultipartForm(32 << 20)
	// 启用Session
	model.Session = session.Start(w, r)

	// 测试直接返回权限
	if config.Debugging  {
		return 3
	}

	// 保存用户请求记录
	role := model.Session.Get("role")
	m := map[string]interface{}{
		"ip":     r.RemoteAddr,
		"cookie": model.Session.ID,
		"time":   time.Now(),
		"user":   model.Session.Get("user"),
		"role":   role,
	}
	var data = make(map[string]string)
	for k, v := range r.Form {
		if k == "handle" {
			m[k] = v[0]
			continue
		}
		if k == "func" {
			m[k] = v[0]
			continue
		}
		data[k] = strings.Join(v, ",")
	}
	m["data"] = data

	// 返回权限
	Log.Out(logger.USER_LOG, m)
	r.Form.Set("ip", r.RemoteAddr)
	// 如果是用户登录，则不需要权限
	handle := r.Form.Get("handle")
	if handle == "User" || handle == "Sync" || handle == "Token"{
		return 3
	}
	if role == nil {
		return 0
	}
	i, _ := strconv.Atoi(role.(string))
	return i
}

// 聊天处理
func chatHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Log.Out(logger.RUN_LOG, "Into chatHandle")
	socket.ServeWs(hub, w, r)
}

// http
// 处理数据, 返回的数据不压缩
func backHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Log.Out(logger.RUN_LOG, "Into backHandle")
	if i := start(w, r); i < 1 {
		w.Write(model.Result(-1, fmt.Sprintf("backHandle.PermissionErr: %s, %s, %v", model.Session.ID, model.Session.Get("user"), i), nil))
		return
	}
	// 解析数据
	b := model.Router(r.Form)
	w.Write(b)
}

// 数据同步，返回的数据压缩
func zipHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Log.Out(logger.RUN_LOG, "Into zipHandle")
	// 解析数据
	if i := start(w, r); i < 1 {
		w.Write(model.Result(-1, fmt.Sprintf("zipHandle.PermissionErr: %s, %s, %v", model.Session.ID, model.Session.Get("user"), i), nil))
		return
	}
	b := model.Router(r.Form)

	// GZIP
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	len, err := zw.Write(b)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		w.Write(b)
		return
	}
	if len == 0 {
		Log.Out(logger.ERR_LOG, "ido.syncHandle zw.Write len = 0")
		w.Write(b)
		return
	}

	err = zw.Flush()
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		w.Write(b)
		return
	}

	err = zw.Close()
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		w.Write(b)
		return
	}
	//w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(buf.Bytes())
}

// 管理后台的Handle处理
func adminHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Log.Out(logger.RUN_LOG, "Into adminHandle")
	// 解析数据
	if i := start(w, r); i < 3 {
		w.Write(model.Result(-1, fmt.Sprintf("adminHandle.PermissionErr: %s, %s, %v", model.Session.ID, model.Session.Get("user"),i), nil))
		return
	}
	b := model.Router(r.Form)
	w.Write(b)
}

// 管理后台的静态页面
func partHtml(w http.ResponseWriter, r *http.Request) {
	Log.Out(logger.RUN_LOG, "Into part")
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	http.ServeFile(w, r, p)
}

// 管理后台的静态页面
func adminHtml(w http.ResponseWriter, r *http.Request) {
	Log.Out(logger.RUN_LOG, "Into adminHtml")
	if i := start(w, r); i < 3 {
		http.Redirect(w, r, "/alert.html", http.StatusFound)
		return
	}
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	http.ServeFile(w, r, p)
}

// 上传文件
func uploadHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Log.Out(logger.RUN_LOG, "Into uploadHandle")
	if i := start(w, r); i < 1 {
		w.Write(model.Result(-1, fmt.Sprintf("uploadHandle.PermissionErr: %s, %s, %v", model.Session.ID, model.Session.Get("user"), i), nil))
		return
	}
	b := model.Upload(r)
	r.Form.Del("path")
	w.Write(b)
}
