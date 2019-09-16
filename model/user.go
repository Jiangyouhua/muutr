package model

import (
	"fmt"
	"net/url"
	"ido/alidayu"
	"time"
	"regexp"
	"strings"
	"math/rand"
	"os/exec"
	"strconv"
	//"gopkg.in/gomail.v2"
	//"io"
)

type User struct {
	Model
}

var ali *alidayu.Alidayu

func init() {
	initModels()
	m := Model{
		table:   "jyh_user",
		columns: []string{"account", "code", "username", "first", "date", "alert", "claim", "cellular", "role", "uuid", "status"},
		primary: "account",
	}
	u := &User{m}
	models["User"] = u

	// 阿里
	ali = &alidayu.Alidayu{
		"目图Muutr",
		"SMS_130910637",
		"LTAI8N6EAp9r2SCW",
		"htTkdHbpWGPqVQLb55StNeTqlgkWhe",
	}
}

func GetUser(v url.Values) *User {
	m := models["User"].(*User)
	m.values = v
	return m
}

// 路由
func (u *User) Router(k string, v url.Values) []byte {
	if len(k) == 0 {
		return Result(0, "User.Routing'key arg is nil", nil)
	}
	u.values = v
	m := map[string]func() []byte{
		"SendCode":    u.SendCode,
		"GetWithUser": u.GetWithUser,
		"SetCode":     u.SetCode,
		"Login":       u.Login,
		"GetAdmin":    u.GetAdmin,
		"SetData":     u.SetDataWithPrimary,
		"GetData":     u.GetDataWithPrimary,
		"DelData":     u.DelDataWithPrimary,
	}
	f, ok := m[k]
	if !ok {
		return Result(0, fmt.Sprintf("user is not this func (%s)", k), nil)
	}
	return f()
}

// 获取单人用户信息
func (u *User) GetWithUser() [] byte {
	user := Session.Get("user").(string)
	if len(user) == 0 {
		return Result(0, "User.GetWithUser'user is nil", nil)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user` = ?", u.getColumn(), u.table);
	return u.GetData(query, user)
}

// 获取全部下级管理者
func (u *User) GetAdmin() []byte {
	where := fmt.Sprintf("`role` > 2 AND role < %v", Session.Get("role"))
	return u.Model.GetDataWithPage(u.values.Get("page"), where, "`date` DESC")
}

// 设置密码
func (u *User) SetCode() []byte {
	account := Session.Get("account").(string)
	if len(account) == 0 {
		return Result(0, "User.SendCode'account is nil", nil)
	}
	code := u.values.Get("code")
	if len(code) == 0 {
		return Result(0, "User.SendCode'code is nil", nil)
	}
	code = u.Md5(code, account, true)
	u.values.Set("code", code)
	u.values.Set("account", account)
	return u.SetDataWithPrimary()
}

// 发送验证码
func (u *User) SendCode() []byte {
	account := u.values.Get("account")
	if len(account) == 0 {
		return Result(0, "user.SendCode'account is nil", nil)
	}

	// 判断是否为手机或邮件
	reg := regexp.MustCompile(`^[0-9]{11}`)
	b := reg.MatchString(account)
	if !b {
		reg = regexp.MustCompile(`[\w-.]+@[\w-]+.[\w-]+(.[\w-]+)+`)
		if !reg.MatchString(account) {
			return Result(0, "user.SendCode'account is not account and email", nil)
		}
	}

	info := make([]string, 0)
	data, err := u.Fetch("SELECT 1 account FROM jyh_user where account = ?", account)
	if err != nil {
		info = append(info, err.Error())
	}

	// 要发验证码
	r := rand.New(rand.NewSource(time.Now().Unix()))
	number := r.Intn(999999)
	c := fmt.Sprintf("%06v", number)
	code := u.Md5(c, account, true)
	u.CreateOrUpdateUser(account, code)

	if b {
		// 阿里
		err = u.SendSMS(account, c)
	} else {
		err = u.SendMail(account, c)
	}
	if err != nil {
		info = append(info, err.Error())
		return Result(0, strings.Join(info, ","), nil)
	}

	info = append(info, "SendSMS is ok")
	return Result(1, strings.Join(info, ","), data)
}

// 用户登录， 登录10次无效， type=0:用户注册或登录，1:自动登录，2:更新Session
func (u *User) Login() []byte {
	// 尝试登录次数, 尝试的时间间隔
	if !u.tryLoginThreshold() {
		return Result(0, "LoginTryMore", nil)
	}
	// 参数
	account := u.values.Get("account")
	if len(account) == 0 {
		return Result(0, "AccountOrCodeErr", nil)
	}
	code := u.values.Get("code")
	if len(code) == 0 {
		return Result(0, "AccountOrCodeErr", nil)
	}
	code = u.Md5(code, account, true)

	// 请求
	query := "SELECT `id`, `username`, `first`, `date`, `alert`, `claim`, `cellular`, `role`, `uuid`, `status` FROM jyh_user WHERE `account` = ? AND `code` = ?"
	d, err := u.Fetch(query, account, code)

	// 请求失败
	if err != nil {
		return Result(0, "SystemErr", err)
	}
	if len(d) == 0 {
		return Result(0, "AccountOrCodeErr", nil)
	}

	// 用户无效
	user := d[0]
	if v, err := strconv.Atoi(user["status"]); err != nil || v < 1 {
		return Result(0, "UserBanned", nil)
	}

	// uuid不同，账号被别的设备登录，用户注册登录时不匹配uuid
	t, _ := strconv.Atoi(u.values.Get("type"))
	uuid := u.values.Get("uuid")
	if t > 0 && user["uuid"] != uuid {
		if len(u.values.Get("sync")) > 0 {
			return Result(0, "OnAnother", nil)
		}
		query := "UPDATE jyh_user SET `uuid` = ? WHERE `id` = ?"
		u.Exce(query, uuid, user["id"] )
	}

	// 登录成功，设置Session
	Session.Set("times", 0)
	Session.Set("account", account)
	Session.Set("user", user["id"])
	Session.Set("role", user["role"])
	Session.Set("username", user["username"])

	// 更新Session
	if t > 1 {
		return Result(1, "session is ok", nil)
	}

	// 返回用户相关数据
	user["code"] = u.values.Get("code")
	user["account"] = account

	// 自动登录并同步数据
	data, err := GetSync(u.values).SetAndGet()
	if data == nil {
		data = make(map[string]interface{})
	}
	data["user"] = user
	return Result(1, err.Error(), data)
}

/**
 * 防止用户过多尝试登录
 */

func (u *User)tryLoginThreshold() bool {
	t := Session.Get("times")
	d := Session.Get("duration")

	times := 0
	if t != nil {
		times = t.(int)
	}
	times ++

	var duration int64
	if d != nil{
		duration = time.Now().Unix() - d.(int64)
	}

	if duration > 180 {
		times = 0
	}
	Session.Set("duration", time.Now().Unix())

	if times > 10 {
		return false
	}

	// 从数据获取用户信息
	Session.Set("times", times)
	return true
}

/**
 * 用户退出
 */
func (u *User) Logout() []byte {
	Session.Clear()
	return Result(1, "logout is ok", nil)
}

/**
 * 创建或修改有效用户
 * @param account
 * @param code
 */
func (u *User) CreateOrUpdateUser(account, code string) {
	if len(account) == 0 || len(code) == 0 {
		return
	}
	username := u.values.Get("username")
	if len(username) == 0 {
		username = strings.Split(account, "@")[0]
	}
	query := "INSERT INTO jyh_user (`account`, `code`, `username`, `status`, `role`, `first`) VALUES ( ?, ?, ?, 1, 2, NOW()) ON DUPLICATE KEY UPDATE `code`=VALUES(`code`), `role`=VALUES(`role`), `status`=VALUES(`status`)"
	u.Exce(query, account, code, username)
}

func (u *User) SendSMS(account, code string) error {
	return ali.SendCode(account, code)
}

func (u *User) SendMail(account, code string) error {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang='en'>
<head>
    <meta charset='UTF-8'>
    <meta name='viewport'
          content='width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no'>
    <title>Muutr</title>
</head>
<body>
<div style='max-width:980px;line-height: 1.8em; background-color: #eee;'>
    <div style='padding:25px; background-color: #333; '>
        <img src='https://muutr.com/html/image/logo.png' height=60>
    </div>
    <div style='padding: 25px;'>
        Hello,<br>
        We have noticed, <br>
        You are trying to use your Muttr account.<br>
        If you receive verification code information,<br>
        Please enter the following.<br><br>
        <h1 style='color:RGB(0,180,0);text-align: center; font-size: 2em;'>%s</h1><br>
        Thank you!<br>
    </div>
    <div style='text-align:right; padding-right:50px'>
        Muutr<br>%s<br><br>
    </div>
</div>
</body>
</html>
	`, code, time.Now().Format("2006-01-02 15:04:06"))


	str := fmt.Sprintf(`echo "%s" | mutt -e "set content_type=text/html" -e "set from=Muutr<no-reply@muutr.com>" -s "Your Verification Code from Muutr" %s`, body, account)
	cmd := exec.Command("/bin/bash", "-c", str)
	return cmd.Run()
}
