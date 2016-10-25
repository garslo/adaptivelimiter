package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Sleeper struct {
	numconns int64
	sleep    func(int64) time.Duration
}

func (me *Sleeper) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	nconns := atomic.AddInt64(&me.numconns, 1)
	fmt.Printf("nconns %d\n", nconns)
	time.Sleep(me.sleep(nconns))
	atomic.AddInt64(&me.numconns, -1)
}

func main() {
	sleeper := &Sleeper{0, func(n int64) time.Duration {
		return time.Duration(10*n+100) * time.Millisecond
	}}
	http.ListenAndServe(":8090", sleeper)
}
