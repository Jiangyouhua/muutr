package push

import (
	"golang.org/x/net/http2"
	"crypto/tls"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"crypto/ecdsa"
	"golang.org/x/crypto/pkcs12"
	"ido/logger"
	"time"
)

type Push struct{
  privateKey	*ecdsa.PrivateKey
  client	*http.Client
}

var (
	Pusher *Push
	Log *logger.Logger
	p12File string
	authFile string
    auth *ecdsa.PrivateKey
	client *http.Client

	debug bool  // 是否为调试状态
	root string // 根目录
	host string // 请求服务器
	alg     = "ES256"
	kid     = "N4AA9X3JTC"  // AuthKey ID
	iss     = "6Q84LRSMLQ"  // App ID
)

func init() {
	// 初始化条件
	authFile = root + "/certificate/AuthKey_N4AA9X3JTC.p8"
	p12File = root + "/certificate/apns-app-key.p12"
	host = "api.push.apple.com:443"
	if debug {
		p12File = root + "/certificate/apns-dev-key.p12"
		host = "api.development.push.apple.com:443"
	}
}

func getAuthKey() (*ecdsa.PrivateKey) {
	// 读私钥
	b, err := ioutil.ReadFile(authFile)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	return key.(*ecdsa.PrivateKey)
}

func getClient() *http.Client {
	cert := tls.Certificate{}
	b, err := ioutil.ReadFile(p12File)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	key, cer, err := pkcs12.Decode(b, "hotkiss521")
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}
	cert.PrivateKey = key
	cert.Certificate = [][]byte{cer.Raw}

	tr := &http.Transport{
		TLSClientConfig:&tls.Config{
			Certificates: []tls.Certificate{cert},
			InsecureSkipVerify:true,
		},
	}

	err = http2.ConfigureTransport(tr)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return nil
	}

	client := &http.Client{
		Transport: tr,
	}
	return client
}

func(p *Push)Send(title, body, image, category string, user []map[string]string) bool{
	defer func() {
		if err := recover(); err != nil{
			Log.Out(logger.ERR_LOG, err)
		}
	}()
	if p.privateKey == nil || p.client == nil{
		Log.Out(logger.RUN_LOG, "push.Send p.privateKey or p.client is nil")
		return false
	}
	if len(title) == 0 || len(body) == 0 {
		Log.Out(logger.RUN_LOG, "push.Send'tilte or body is nil")
		return true
	}
	if len(user) == 0 {
		Log.Out(logger.RUN_LOG, "push.Send'suer is nil")
		return true
	}

	Log.Out(logger.RUN_LOG, "into push.Send")


	if Pusher == nil {
		Pusher = &Push{
			getAuthKey(),
			getClient(),
		}
	}

	// 准备发送
	s := `{"aps":{"alert":{"body":"%s","title":"%s"}}}`
	s = fmt.Sprintf(s, body, title)

	// 制作Token
	token, err := p.getToken()
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return false
	}

	for _, m := range user {
		t, ok := m["token"]
		if !ok{
			continue
		}
		if !p.send(t, s, token){
			break
		}
	}
	return true
}

func(p *Push) getToken() (string, error) {
	// jwt
	claims := &jwt.StandardClaims{
		IssuedAt: time.Now().Unix(),
		Issuer:   iss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = kid
	token.Header["alg"] = alg

	tokenString, err := token.SignedString(auth)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return "", err
	}
	return tokenString, err
}

// 发送
func(p *Push)send(t, s, token string) bool{
	defer func() {
		if err := recover(); err != nil {
			Log.Out(logger.ERR_LOG, err)
		}
	}()
	if len(t) == 0 {
		return false
	}

	// 发送的地址
	url := fmt.Sprintf("https://%s/3/device/%s", host, t)
	r, err := http.NewRequest("POST", url, strings.NewReader(s))
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return false
	}
	if r == nil {
		Log.Out(logger.ERR_LOG, "push.Pushing http.NewRequest request is nil")
		return false
	}
	r.Header.Set("apns-expiration", "0")
	r.Header.Set("apns-priority", "10")
	r.Header.Set("authorization", token)

	resp, err := client.Do(r)
	if err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode > 200 {
		Log.Out(logger.ERR_LOG, fmt.Sprintf("push.Pushing resp.StatusCode is %v, info %s", resp.StatusCode, resp.Status))
		return false
	}
	Log.Out(logger.RUN_LOG, fmt.Sprintf("Push: %s", url))
	return true
}



