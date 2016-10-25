package adaptivelimiter

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type intervalLimiter struct {
	rwlock sync.RWMutex
	ticker *time.Ticker
}

func (me *intervalLimiter) Spawn(fn func()) {
	me.rateLimitSpawn(fn)
}

func (me *intervalLimiter) updateRate(ratePerSecond int64) {
	me.rwlock.Lock()
	if me.ticker != nil {
		me.ticker.Stop()
	}
	interval := time.Second / time.Duration(ratePerSecond)
	me.ticker = time.NewTicker(interval)
	me.rwlock.Unlock()
}

func (me *intervalLimiter) rateLimitSpawn(fn func()) {
	me.rwlock.RLock()
	select {
	case _, ok := <-me.ticker.C:
		if ok {
			me.spawnRoutine(fn)
		}
	}
	me.rwlock.RUnlock()
}

func (me *intervalLimiter) spawnRoutine(fn func()) {
	go func() {
		fn()
	}()
}

type SmoothingFunction func(float64) float64

func NewLimiter(minRate, maxRate int64, axis *PartitionedAxis, controller *PIDController, targetDuration, resolution time.Duration, smoothers []SmoothingFunction) *Limiter {
	l := &Limiter{minRate, maxRate, axis, minRate, controller, targetDuration, resolution,
		&intervalLimiter{}, smoothers}
	l.Init()
	return l
}

type Limiter struct {
	MinRate            int64
	MaxRate            int64
	CountAxis          *PartitionedAxis
	CurrentAllowedRate int64
	Controller         *PIDController
	TargetDuration     time.Duration
	Resolution         time.Duration
	Limiter            *intervalLimiter
	Smoothers          []SmoothingFunction
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func (me *Limiter) Init() {
	me.Limiter.updateRate(me.MinRate) // + (me.MaxRate-me.MinRate)/2)
	t := time.Tick(me.Resolution)
	go func() {
		for {
			<-t
			newRate := me.CurrentAllowedRate
			delta, ok := me.Controller.Delta()
			if !ok {
				continue
			}
			//fmt.Printf("before %f\n", delta)
			for _, smoother := range me.Smoothers {
				delta = smoother(delta)
			}
			//fmt.Printf("after %f\n", delta)
			newRate += int64(delta)
			if newRate > me.MaxRate {
				newRate = me.MaxRate
			}
			if newRate < me.MinRate {
				newRate = me.MinRate
			}
			fmt.Printf("%d\n", newRate)
			me.Limiter.updateRate(newRate)
			me.CurrentAllowedRate = newRate
		}
	}()
}

func (me *Limiter) Spawn(fn func()) error {
	me.Limiter.Spawn(func() {
		start := time.Now()
		fn()
		dur := time.Now().Sub(start)
		// milliseconds
		durMs := time.Duration(dur-me.TargetDuration).Seconds() * 1000
		me.Controller.Update(durMs)
	})
	return nil
}
