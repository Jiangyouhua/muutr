package session

/**
* Session类，Session集合类
* 2017.09.18， 单一Session
 */

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
	"strconv"
	"ido/logger"
)

var (
	Log *logger.Logger
	// Sessions 全局单例SessionSet
	sessions *SessionSet
	// Validity session有效时间
	Validity int
)

// Session 网页Session类
type Session struct {
	ID      string                 //当前Session id
	Values  map[string]interface{} //Session 数据集
	Time    time.Time              //最新时间
	Sync    *sync.RWMutex          //多线程操作锁
}

// SessionSet 页面Session类集合
type SessionSet struct {
	Values map[string]*list.Element // Session数据集合
	List   *list.List          // Session集合的有效顺序链，按有效时效最早到最晚
	Sync   *sync.RWMutex       //多线程操作锁
}

// 启用一个全局变量Session集合
func init() {
	sessions = &SessionSet{
		Values: make(map[string]*list.Element),
		List:   new(list.List),
		Sync:   new(sync.RWMutex),
	}
	Validity = 1440 // 流览器默认的过期时间

}

// Start 启用Session
func Start(w http.ResponseWriter, r *http.Request) *Session {
	name := "JiangYouHua"
	key := ""
	//从请求获取cookie key
	cookie, err := r.Cookie(name)
	if err == nil {
		key = cookie.Value
	} else {
		//新生成cookie key
		key = sessions.ID()
		c := &http.Cookie{
			Name:     name,
			Value:    key,
			Path:     "/",
			MaxAge:   Validity,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
	s := sessions.Get(key)
	if s == nil{
		s = sessions.Set(key)
	}
	return s
}

func Update(d time.Duration)  {
	for {
		time.Sleep(d)
		sessions.Update()
	}
}

// ID 生成SessionSet全局id
func (ss *SessionSet) ID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		Log.Out(logger.ERR_LOG, err.Error())
		t := time.Now()
		n := t.UnixNano()
		s := strconv.FormatInt(n,10)
		b = []byte(s)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// Get 从SessionSet集合中获取session，
// 每一次页面请求均调用Start再调用Get，所以只需要在GET里更新一下Session的时间与位置就可以了
func (ss *SessionSet) Get(key string) *Session {
	if len(key) == 0{
		return nil
	}
	ss.Sync.RLock()
	defer ss.Sync.RUnlock()

	element, ok := ss.Values[key]
	if !ok || element == nil{
		return nil
	}
	ss.List.MoveToBack(element)

	s, ok := element.Value.(*Session)
	if !ok {
		return nil
	}
	s.Time = time.Now()
	return s
}

// Set 向SessionSet集合添加新的Session
func (ss *SessionSet) Set(key string)  *Session{
	if len(key) == 0{
		return nil
	}
	ss.Sync.Lock()
	defer ss.Sync.Unlock()

	s := &Session{key, make(map[string]interface{}), time.Now(), new(sync.RWMutex)}
	element := ss.List.PushBack(s)
	ss.Values[key] = element
	return s
}

// Update 更新集合中的有效性，有效更新
func (ss *SessionSet)Update() {
	ss.Sync.Lock()
	defer ss.Sync.Unlock()
	for{
		element := ss.List.Front()
		if element == nil {
			break
		}
		i := ss.List.Remove(element)
		s, ok := i.(*Session)
		if !ok {
			continue
		}
		if !s.Time.IsZero() && int(time.Now().Sub(s.Time).Seconds()) < Validity {
			break
		}
		delete(ss.Values, s.ID)
	}
}

// Get 获取保存在Session的内容
func (s *Session) Get(key string) interface{} {
	s.Sync.RLock()
	defer s.Sync.RUnlock()

	if v, ok := s.Values[key]; ok {
		return v
	}
	return nil
}

// Set 保存内容在Session
func (s *Session) Set(key string, value interface{}) {
	s.Sync.Lock()
	defer s.Sync.Unlock()
	s.Values[key] = value
}

// Del 清除Session
func (s *Session) Del(key string) {
	s.Sync.Lock()
	defer s.Sync.Unlock()
	delete(s.Values, key)
}

// clear 清空
func (s *Session)Clear()  {
	s.Sync.Lock()
	defer s.Sync.Unlock()
	s.Values = make(map[string]interface{})
}
