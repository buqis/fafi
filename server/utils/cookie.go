package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Jar struct {
	mutex sync.RWMutex
	// map[url]map[cookie.Name]*http.Cookie
	kw   map[string]map[string]*http.Cookie
	once sync.Once
}

func (j *Jar) lazy() {
	j.once.Do(func() {
		j.kw = make(map[string]map[string]*http.Cookie)
	})
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.lazy()

	j.mutex.Lock()
	defer j.mutex.Unlock()

	for _, cookie := range cookies {
		cookie := cookie

		if cookie.MaxAge < 0 {
			// 删除
			delete(j.kw[getDomain(u.Hostname())], cookie.Name)
			continue
		} else if cookie.MaxAge > 0 {
			cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
		}

		if _, ok := j.kw[getDomain(u.Hostname())]; !ok {
			j.kw[getDomain(u.Hostname())] = make(map[string]*http.Cookie)
		}

		j.kw[getDomain(u.Hostname())][cookie.Name] = cookie
	}
}

func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	j.lazy()

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	cookies := make([]*http.Cookie, 0, len(j.kw[getDomain(u.Hostname())]))
	for _, v := range j.kw[getDomain(u.Hostname())] {
		v := v
		// domain
		if len(v.Domain) > 0 && !strings.HasSuffix(u.Hostname(), v.Domain) {
			continue
		}
		// path
		if len(v.Path) > 0 && v.Path != "/" && !strings.HasPrefix(u.Path, v.Path) {
			continue
		}
		cookies = append(cookies, v)
	}
	return cookies
}

func (j *Jar) Export() string {
	j.lazy()

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	data, _ := json.Marshal(j.kw)
	return string(data)
}

func (j *Jar) Import(str string) error {
	j.lazy()

	if len(str) == 0 {
		return nil
	}

	kw := map[string]map[string]*http.Cookie{}
	if err := json.Unmarshal([]byte(str), &kw); err != nil {
		return err
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	now := time.Now()
	for k, v := range kw {
		for name, cookie := range v {
			if cookie.Expires.Sub(now) <= 0 {
				// 过期
				continue
			}

			if _, ok := j.kw[k]; !ok {
				j.kw[k] = map[string]*http.Cookie{}
			}
			j.kw[k][name] = cookie
		}
	}
	return nil
}

func getDomain(host string) string {
	if strings.Count(host, ".") < 2 {
		return host
	}
	return host[strings.Index(host, "."):]
}
