package alidayu

import (
	"crypto/hmac"
	"encoding/base64"
	"crypto/sha1"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
	"errors"
)

const (
	ali = "http://dysmsapi.aliyuncs.com"
)

type Alidayu struct{
	SignName string
	TemplateCode string
	AccessKeyId string
	AccessKeySecret string
}

type Result struct{
	BizId string
	Message string
	RequestId string
	Code string
}

func (a *Alidayu)signHMAC(params url.Values, appSecret string) (signature string){
	keys := make([]string,0)
	for k := range params {
		keys = append(keys, k)
	}
	str := ""
	sort.Strings(keys)
	for _, k := range keys {
		str += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(params.Get(k))
	}
	signString := "GET&%2F&" + url.QueryEscape(str[1:])
	mac := hmac.New(sha1.New, []byte(appSecret+"&"))
	mac.Write([]byte(signString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (a *Alidayu)SendCode(user, code string) error{

	if len(user) == 0 || len(code) == 0{
		return errors.New("Alidaye.SendCode user or code is nil")
	}

	// 1. 系统参数
	params := url.Values{}
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureNonce", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("AccessKeyId", a.AccessKeyId)
	params.Set("SignatureVersion", "1.0")
	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Add("Format", "JSON")

	// 2. 业务API参数
	params.Set("Action", "SendSms")
	params.Set("Version", "2017-05-25")
	params.Set("RegionId", "cn-hangzhou")
	params.Set("UserNumbers", user)
	params.Set("SignName", a.SignName)
	params.Set("TemplateParam", "{\"code\":\""+code+"\"}")
	params.Set("TemplateCode", a.TemplateCode)
	params.Set("OutId", "123")

	// 3. 去除签名关键字Key
	if signature := params.Get("Signature"); len(signature) > 0 {
		params.Del("Signature")
	}

	signString := a.signHMAC(params, a.AccessKeySecret)
	params.Set("Signature", signString)
	req, err := http.NewRequest(http.MethodGet, ali+"/?"+params.Encode(), nil)

	//req.Header.Set("x-sdk-client", "Java/2.0.0")
	//req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", "Java/1.6.0_45")

	c := new(http.Client)
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Result
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return err
	}
	if result.Code == "OK" {
		return nil
	}
	return errors.New(result.Message)
}

