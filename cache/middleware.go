package cache

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"gitlab.com/tokend/go/signcontrol"

	"github.com/pkg/errors"

	"github.com/go-chi/chi/middleware"

	"github.com/golang/groupcache/lru"
)

type MiddlewareCache struct {
	lock             sync.Mutex
	cache            *lru.Cache
	expirationPeriod time.Duration
}

func NewMiddlewareCache(maxEntries int, expPer time.Duration) *MiddlewareCache {
	return &MiddlewareCache{
		cache:            lru.New(maxEntries),
		expirationPeriod: expPer,
	}
}

func (c *MiddlewareCache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || signcontrol.IsSigned(r) {
			next.ServeHTTP(w, r)
			return
		}

		url := r.URL.String()

		c.lock.Lock()
		rawValue, ok := c.cache.Get(url)
		c.lock.Unlock()

		if ok {
			value, ok := rawValue.(value)
			if !ok {
				panic("failed to cast cached value")
			}

			if value.expiration.After(time.Now()) {
				for k, v := range value.header {
					for _, elem := range v {
						w.Header().Add(k, elem)
					}
				}

				w.WriteHeader(value.status)

				_, err := w.Write(value.response)
				if err != nil {
					panic(errors.Wrap(err, "failed to wrote cached response"))
				}
				return
			}

		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		var raw bytes.Buffer
		ww.Tee(&raw)

		next.ServeHTTP(ww, r)

		if (ww.Status() >= 300) || (ww.Status() < 200) {
			return
		}

		value := value{
			expiration: time.Now().Add(c.expirationPeriod),
			response:   raw.Bytes(),
			header:     ww.Header(),
			status:     ww.Status(),
		}

		c.lock.Lock()
		c.cache.Add(url, value)
		c.lock.Unlock()
	})
}

type value struct {
	expiration time.Time
	response   []byte
	header     http.Header
	status     int
}
