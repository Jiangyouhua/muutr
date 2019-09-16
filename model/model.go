package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"encoding/json"
	"ido/session"
	"net/http"
	"os"
	"path"
	"mime/multipart"
	"io"
	"errors"
	"ido/logger"
	"html"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"ido/push"
	"time"
)

type Conn struct {
	Name string  // 数据库名
	Info string  // 数据库连接信息
}

// Model接口方法
type Mod interface {
	Router(k string, v url.Values) []byte
}

// Model类父类
type Model struct {
	table   string     // 表名
	columns []string   // 列名
	primary string     // 主键
	values  url.Values // 传入的参数
}

var (
	db        *sql.DB
	conn          = Conn{"mysql", "root:eastism@gmail.com@tcp(localhost:3306)/jiang_ido?charset=utf8mb4"}
	models    map[string]interface{} // mod的映射
	Log       *logger.Logger
	Session   *session.Session
	paging      = 25
	datetime   = "1970-01-01 00:00:00"
	Debugging bool
)

// 初始化
func init() {
	initModels()
	initDb()
}

// 初始模型
func initModels() {
	if models != nil {
		return
	}
	models = make(map[string]interface{})
}

// 初始数据
func initDb() {
	// 判断连接是否存在
	if db != nil {
		err := db.Ping()
		if err == nil {
			return
		}
		Log.Out(logger.ERR_LOG, err.Error())
		db.Close()
	}

	// 建立连接
	d, err := sql.Open(conn.Name, conn.Info)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return
	}
	db = d
}

// 转字符串处理
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return string(v)
	case[]string:
		return strings.Join(v, ",")
	default:
		return ""
	}
}

// 从数据库返回数组中获取指定col的所有数值，各数值通过逗号连接
func ReToString(re []map[string]string, col... string) map[string]string{
	// 判断传入的数据
	if len(re) == 0 || len(col) == 0 {
		return nil
	}

	// 遍历出数值
	m := make(map[string]string)
	for _, dic := range re{
		for _, c := range col {
			word, ok := dic[c]

			// 没有指定Col，返回空
			if !ok {
				return nil
			}

			// 合并各值
			word = fmt.Sprintf("'%s'", word)
			if len(m[c]) == 0 {
				m[c] = word
				continue
			}
			m[c] = fmt.Sprintf("%s,%s", m[c], word)
		}
	}
	return m
}


// 获取数据表列，逗号连接各列名，args为排除列名
func (m *Model) getColumn(args ...string) (col string) {
	if len(m.columns) == 0 {
		return
	}
	if len(args) == 0 {
		args = m.columns
	}
	col = fmt.Sprintf("`%s`", strings.Join(args, "`, `"))
	return
}

// 获取UpdateOrInsert表列，逗号连接各列名，args为排除列名
func (m *Model) getUpdateColumn(primary string, args ...string) (col string) {
	if len(m.columns) == 0 {
		return
	}

	if len(args) == 0 {
		args = m.columns
	}

	for _, v := range args {
		if v == "first"  || v == primary {
			continue
		}
		// 排除参数中的列
		str := fmt.Sprintf("`%s` = VALUES(`%s`)", v, v)
		if len(col) == 0 {
			col = str
			continue
		}
		col = fmt.Sprintf("%s, %s", col, str)
	}
	return
}

// 执行SQL语言
func (m *Model) SetData(query string, args ...interface{}) []byte {
	if len(query) == 0 {
		m.SetDataWithPrimary()
	}
	return m.Exce(query, args...)
}

// 执行默认生成的SQL语句
func (m *Model) SetDataWithPrimary() []byte {
	var data  interface{}
	d := m.values.Get("_data")
	if len(d) > 0 {
		if err := json.Unmarshal([]byte(d), &data); err != nil {
			Log.Out(logger.ERR_LOG, err)
		}
	}

	query := m.InsertUpdateStr(data)
	return m.Exce(query)
}

// 查询QL语句
func (m *Model) GetData(query string, args ...interface{}) []byte {
	if len(query) == 0 {
		m.GetDataWithPrimary()
	}
	return m.Query(query, args...)
}

// 查询默认生成的SQL语句
func (m *Model) GetDataWithPrimary() []byte {
	if len(m.primary) == 0 {
		return Result(0, "Model DelDate's Primary is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `%s` = ?", m.getColumn(), m.table, m.primary)
	return m.Query(query, m.values.Get(m.primary))
}

// 分页查询
func (m *Model) GetDataWithPage(page, where, order string) []byte {
	p := 0
	if len(page) > 0 {
		p, _ = strconv.Atoi(page)
	}
	if len(where) > 0 {
		where = fmt.Sprintf("WHERE %s", where)
	}
	if len(order) > 0 {
		order = fmt.Sprintf("ORDER BY %s", order)
	}
	limit := fmt.Sprintf("%d, %d", p*paging, paging)
	query := fmt.Sprintf("SELECT %s FROM %s %s %s LIMIT %s", m.getColumn(), m.table, where, order, limit)
	return m.Query(query)
}

// 删除数据
func (m *Model) DelData(query string, args ...interface{}) []byte {
	if len(query) == 0 {
		m.DelDataWithPrimary()
	}
	return m.Exce(query, args...)
}

// 按默认条件删除
func (m *Model) DelDataWithPrimary() []byte {
	if len(m.primary) == 0 {
		return Result(0, "Model DelDate'Primary value is nil", nil)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?",  m.table, m.primary)
	return m.Exce(query, m.values.Get(m.primary))
}

// 插入或更新
func (m *Model) InsertUpdateStr(data interface{}) string {
	if len(m.columns) == 0 {
		return ""
	}

	var (
		values  string // insert的值
		exclude = make([]string, 0) // 用在update的列
	)

	// 无指定数据，从传入表单请求
	if data == nil {
		for _, v := range m.columns {
			_, ok := m.values[v]
			value, state := m.valueWithKey(v, m.values.Get(v), ok)
			if state < 0 {
				return ""
			}
			if state == 0 {
				continue
			}
			exclude = append(exclude, v)
			values = m.ValuesAdd(values, value)
		}
		values = fmt.Sprintf("(%s)", values)
		return m.formatInsertUpdateStr(values, exclude)
	}

	// 传入map
	if dic, ok := data.(map[string]string); ok {
		for _, v := range m.columns {
			val, ok := dic[v]
			value, state := m.valueWithKey(v, val, ok)

			if state < 0 {
				return ""
			}
			if state == 0 || !ok {
				continue
			}
			exclude = append(exclude, v)
			values = m.ValuesAdd(values, value)
		}
		values = fmt.Sprintf("(%s)", values)
		return m.formatInsertUpdateStr(values, exclude)
	}

	// 传入[]map[string]string
	if arr, ok := data.([]interface{}); ok {
		for k, i := range arr {
			dic, o := i.(map[string]interface{})
			if !o {
				return ""
			}

			var s string
			for _, v := range m.columns {
				// 判断是否有操作人员column
				val, ok := dic[v]
				value, state := m.valueWithKey(v, val, ok)
				if state < 0 {
					return ""
				}
				if state == 0 {
					continue
				}
				if k == 0{
					exclude = append(exclude, v)
				}
				s = m.ValuesAdd(s, value)
			}
			if len(s) == 0 {
				continue
			}
			s = fmt.Sprintf("(%s)",s)
			values = m.ValuesAdd(values, s)
		}
		return m.formatInsertUpdateStr(values, exclude)
	}
	return ""
}


func (m *Model) formatInsertUpdateStr(values string, exclude []string) string {
	if len(values) == 0 {
		return ""
	}
	columns := m.getColumn(exclude...)
	updateColumns := m.getUpdateColumn(m.primary, exclude...)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s", m.table, columns, values, updateColumns)
}

// 对各key的值进行处理，-1，无主健退出；0，无效操作不处理，1,有效操作
func (m *Model) valueWithKey(key string, val interface{}, ok bool) (string, int) {
	// 如果没有key,或key为date，则不处理，date由数据库自行修改时间
	if len(key) == 0 || key == "id" || key == "date"{
		return "", 0
	}

	// 如果主键没数据，则无法进行insert update
	if key == m.primary && val == nil {
		return "", -1
	}

	// 如果是first键，则insert当前时间，但不在更新里
	if key == "first" {
		return "NOW()", 1
	}

	if key == "editor" {
		return ToString(Session.Get("user")), 1
	}

	if !ok {
		return "", 0
	}

	// 格式化值:转为字符串，转字符，加单引号
	value := ToString(val)
	value = html.EscapeString(value)
	value = fmt.Sprintf("'%s'", value)
	return value, 1
}

func (m Model)ValuesAdd(values, value string) string{
	if len(value) == 0{
		return values
	}
	if len(values) == 0 {
		return value
	}
	return fmt.Sprintf("%s, %s", values, value)
}

// 查询
func (m *Model) Query(query string, args ... interface{}) []byte {
	d, err := m.Fetch(query, args...)
	if err != nil {
		return Result(0, err.Error(), nil)
	}
	return Result(1, "Query is ok", d)
}

// query, sql语句，col独立获取的列
func (m *Model) Fetch(query string, args ... interface{}) (data []map[string]string, e error) {
	// fmt.Println(query, args)
	if len(query) == 0 {
		e = errors.New("Model.Fetch query is nil")
		return
	}
	initDb()
	rows, err := db.Query(query, args...)
	if err != nil {
		Log.Out(logger.ERR_LOG, fmt.Sprintf("Err:%s, Query:%s, args:%s", err.Error(), query, args))
		e = err
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		Log.Out(logger.ERR_LOG, fmt.Sprintf("Err:%s, Query:%s, args:%s",err.Error(), query, args))
		e = err
		return
	}

	data = make([]map[string]string, 0)
	d := make([]sql.RawBytes, len(columns)) // 接收数据组
	a := make([]interface{}, len(columns))  // 接收数组各成员的地址

	for k := range columns {
		a[k] = &d[k] //使用&d[k],不能直接用 key, val := range d 中的&v
	}

	text := make([]string,0)
	for rows.Next() {
		if err := rows.Scan(a...); err != nil {
			text = append(text, err.Error())
			break
		}

		m := make(map[string]string)
		for k, v := range d {
			s := columns[k]
			value := string(v)
			m[s] = html.UnescapeString(value)
		}
		data = append(data, m)
	}
	// 如果出错
	if len(text) > 0 {
		e = errors.New(strings.Join(text, ","))
	}
	return
}

// 执行
func (m *Model) Exce(query string, args ... interface{}) []byte {
	 fmt.Println(query, args, m.values)
	id, rows, err := m.Set(query, args...)
	if err != nil {
		return Result(0, err.Error(), nil)
	}
	return Result(1, "Exce is ok", map[string]interface{}{"id": id, "rows": rows})
}

func (m *Model) Set(query string, args ... interface{}) (id, rows int64, err error) {
	//fmt.Println(query, args)
	if len(query) == 0 {
		err = errors.New("Model.Set query is nil")
		return
	}
	initDb()
	var (
		re sql.Result
	)
	re, err = db.Exec(query, args...)
	if err != nil {
		Log.Out(logger.ERR_LOG, fmt.Sprintf("Err:%s, Query:%s, args:%s",err.Error(), query, args))
		return
	}
	id, _ = re.LastInsertId()
	rows, _ = re.RowsAffected()
	return
}

/**
以下为全局函数
 */

// 统一返回
func Result(code int, info string, data interface{}) []byte {
	m := make(map[string]interface{})
	m["code"] = code
	m["info"] = info
	m["data"] = data
	b, err := json.Marshal(m)
	if err != nil {
		Log.Out(logger.ERR_LOG, fmt.Sprintf("ResultErr: %s", err.Error()))
		return nil
	}
	if code == 0{
		Log.Out(logger.ERR_LOG, fmt.Sprintf("ResultCode: %s", info))
	}
	return b
}

// 入口方法，根据handle处理请求，有权限认证及对请求过多处理
func Router(values url.Values) []byte {
	// handle
	handle := values.Get("handle")
	if len(handle) == 0 {
		return Result(0, "Request 'handle' is nil", nil)
	}

	// func
	function := values.Get("func")
	if len(function) == 0 {
		return Result(0, "Request 'func' is nil", nil)
	}

	m, ok := models[handle]
	if !ok {
		return Result(0, fmt.Sprintf("Request handle'%s Not defined", handle), nil)
	}
	return m.(Mod).Router(function, values)
}

func Upload(r *http.Request) []byte {
	// 上传文件多个文件
	p := r.Form.Get("path")
	var data = make(map[string]string)
	for k, v := range r.MultipartForm.File {
		if len(v) == 0 {
			continue
		}
		for _, value := range v {
			name, err := file(p, k, value)
			if err != nil {
				Log.Out(logger.RUN_LOG, err.Error())
				continue
			}
			data[k] = name
		}
	}
	// 保存数据
	if len(r.Form.Get("func")) == 0 {
		return Result(1, "Upload is ok", data)
	}
	r.Form.Del("path")
	for k, v := range data {
		r.Form.Set(k, v)
	}
	return Router(r.Form)
}

func file(p, k string, header *multipart.FileHeader) (string, error) {
	if len(p) == 0 || len(k) == 0 || len(header.Filename) == 0 {
		return "", errors.New(fmt.Sprintf("Image args is nil, %s, %s, %s", p, k, header.Filename))
	}

	f, err := header.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	outputFilePath := fmt.Sprintf("./file/%s/%s/%s", p, k, header.Filename)
	pa := path.Dir(outputFilePath)
	//判断目录是否存在
	if fi, err := os.Stat(pa); err != nil || !fi.IsDir() {
		os.MkdirAll(pa, os.ModePerm)
	}

	writer, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	io.Copy(writer, f)
	return outputFilePath, nil
}

// 推送
func pushing(){
	users, _ := GetToken(nil).GetAllToken()
	if len(users) == 0 {
		return
	}
	p := GetPusher(nil);
	pusher, _ := p.GetLatestOne()
	if pusher == nil {
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", pusher["push"])
	if err != nil {
		return
	}

	// 在10秒范围内的发送
	d := t.Sub(time.Now())
	if d.Seconds() < 10 {
		b := push.Pusher.Send(pusher["title"],pusher["body"], pusher["image"], pusher["category"], users)
		p.values.Set("pid", pusher["pid"])
		p.values.Set("status", "0")
		p.SetDataWithPrimary()
		// push 对象错误
		if !b {
			return
		}
	}
	time.Sleep(d - 10 * time.Second)
	pushing()
}

// 加盐md5
func (m *Model) Md5(code, user string, b bool) string {
	if !b {
		return code
	}
	h := md5.New()
	h.Write([]byte(user + "JiangYouHua" + code))
	return hex.EncodeToString(h.Sum(nil))
}


