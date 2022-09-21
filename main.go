//package main
//
//import "utils/logging"
//
//func main() {
//	conf := logging.Conf{
//		Path:    "log/test.log",
//		Encoder: "json",
//	}
//	logging.Init(conf)
//	logging.Errorf("bug 2:%s", "bug")
//}

package main

import (
	"fmt"
	"sync"
	"time"
)

// timeout middleware wraps the request context with a timeout
//func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
//	return func(c *gin.Context) {
//
//		// wrap the request context with a timeout
//		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
//
//		defer func() {
//			// check if context timeout was reached
//			if ctx.Err() == context.DeadlineExceeded {
//
//				// write response and abort the request
//				//c.Writer.WriteHeader(http.StatusGatewayTimeout)
//				c.Abort()
//			}
//
//			//cancel to clear resources after finished
//			cancel()
//		}()
//
//		// replace request with context wrapped request
//		c.Request = c.Request.WithContext(ctx)
//		c.Next()
//	}
//}
//
//func timedHandler(duration time.Duration) func(c *gin.Context) {
//	return func(c *gin.Context) {
//
//		// get the underlying request context
//		ctx := c.Request.Context()
//
//		// create the response data type to use as a channel type
//		type responseData struct {
//			status int
//			body   map[string]interface{}
//		}
//
//		// create a done channel to tell the request it's done
//		doneChan := make(chan responseData)
//
//		// here you put the actual work needed for the request
//		// and then send the doneChan with the status and body
//		// to finish the request by writing the response
//		go func() {
//			time.Sleep(duration)
//			doneChan <- responseData{
//				status: 200,
//				body:   gin.H{"hello": "world"},
//			}
//		}()
//
//		// non-blocking select on two channels see if the request
//		// times out or finishes
//		select {
//
//		// if the context is done it timed out or was cancelled
//		// so don't return anything
//		case <-ctx.Done():
//			return
//
//			// if the request finished then finish the request by
//			// writing the response
//		case res := <-doneChan:
//			c.JSON(res.status, res.body)
//		}
//	}
//}
//
//func main() {
//	// create new gin without any middleware
//	engine := gin.New()
//
//	// add timeout middleware with 2 second duration
//	engine.Use(timeoutMiddleware(time.Second * 2))
//
//	// create a handler that will last 1 seconds
//	engine.GET("/short", timedHandler(time.Second))
//
//	// create a route that will last 5 seconds
//	engine.GET("/long", timedHandler(time.Second*5))
//
//	// run the server
//	log.Fatal(engine.Run(":8080"))
//}

// 并发访问同一个user_id/ip的记录需要上锁
var recordMu map[string]*sync.RWMutex

func init() {
	recordMu = make(map[string]*sync.RWMutex)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type TokenBucket struct {
	BucketSize int                // 木桶内的容量：最多可以存放多少个令牌
	TokenRate  time.Duration      // 多长时间生成一个令牌
	records    map[string]*record // 报错user_id/ip的访问记录
}

// 上次访问时的时间戳和令牌数
type record struct {
	last  time.Time
	token int
}

func NewTokenBucket(bucketSize int, tokenRate time.Duration) *TokenBucket {
	return &TokenBucket{
		BucketSize: bucketSize,
		TokenRate:  tokenRate,
		records:    make(map[string]*record),
	}
}

func (t *TokenBucket) getUidOrIp() string {
	// 获取请求用户的user_id或者ip地址
	return "127.0.0.1"
}

// 获取这个user_id/ip上次访问时的时间戳和令牌数
func (t *TokenBucket) getRecord(uidOrIp string) *record {
	if r, ok := t.records[uidOrIp]; ok {
		return r
	}
	return &record{}
}

// 保存user_id/ip最近一次请求时的时间戳和令牌数量
func (t *TokenBucket) storeRecord(uidOrIp string, r *record) {
	t.records[uidOrIp] = r
	//bytes, _ := json.Marshal(r)
	//stringData := string(bytes)
	fmt.Println("record:", r.last, r.token)
}

// 验证是否能获取一个令牌
func (t *TokenBucket) validate(uidOrIp string) bool {
	// 并发修改同一个用户的记录上写锁
	rl, ok := recordMu[uidOrIp]
	if !ok {
		var mu sync.RWMutex
		rl = &mu
		recordMu[uidOrIp] = rl
	}
	rl.Lock()
	defer rl.Unlock()

	r := t.getRecord(uidOrIp)
	now := time.Now()
	if r.last.IsZero() {
		// 第一次访问初始化为最大令牌数
		r.last, r.token = now, t.BucketSize
	} else {
		if r.last.Add(t.TokenRate).Before(now) {
			// 如果与上次请求的间隔超过了token rate
			// 则增加令牌，更新last
			r.token += max(int(now.Sub(r.last)/t.TokenRate), t.BucketSize)
			r.last = now
		}
	}
	var result bool
	if r.token > 0 {
		// 如果令牌数大于1，取走一个令牌，validate结果为true
		r.token--
		result = true
	}

	// 保存最新的record
	t.storeRecord(uidOrIp, r)
	return result
}

// 返回是否被限流
func (t *TokenBucket) IsLimited() bool {
	return !t.validate(t.getUidOrIp())
}

func main() {
	tokenBucket := NewTokenBucket(2, time.Second)
	for i := 0; i < 3; i++ {
		fmt.Println(tokenBucket.IsLimited())
	}
	//time.Sleep(1 * time.Second)
	//fmt.Println(tokenBucket.IsLimited())
}
