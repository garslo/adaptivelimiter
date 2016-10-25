package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/garslo/adaptivelimiter"
)

func main() {
	var (
		p          float64
		i          float64
		d          float64
		history    int
		resolution time.Duration
		target     time.Duration
		minRate    int64
		maxRate    int64
		urlToGet   string
		minDelta   float64
		maxDelta   float64
		gain       float64
	)
	flag.Float64Var(&p, "p", .01, "p")
	flag.Float64Var(&i, "i", .01, "i")
	flag.Float64Var(&d, "d", .01, "d")
	flag.Float64Var(&minDelta, "mindelta", 1, "min delta")
	flag.Float64Var(&maxDelta, "maxdelta", 1, "max delta")
	flag.Float64Var(&gain, "gain", 100, "gain")
	flag.DurationVar(&resolution, "r", 1*time.Second, "axis resolution")
	flag.DurationVar(&target, "t", 500*time.Millisecond, "target download time")
	flag.IntVar(&history, "history", 5, "history")
	flag.Int64Var(&minRate, "minrate", 1, "min rate")
	flag.Int64Var(&maxRate, "maxrate", 10, "max rate")
	flag.StringVar(&urlToGet, "url", "http://localhost:8090", "what to download")
	flag.Parse()

	axis := adaptivelimiter.NewPartitionedAxis(float64(resolution.Nanoseconds()))
	controller := adaptivelimiter.NewPIDController(axis, p, i, d, history)
	smoothers := []adaptivelimiter.SmoothingFunction{
		adaptivelimiter.LogSmoother(),
		adaptivelimiter.ConstantGainSmoother(gain),
		adaptivelimiter.MinSmoother(minDelta),
		adaptivelimiter.MaxSmoother(maxDelta),
	}
	limiter := adaptivelimiter.NewLimiter(minRate, maxRate, axis, controller, target, resolution, smoothers)
	for {
		limiter.Spawn(func() {
			resp, err := http.Get(urlToGet)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		})
	}
	select {}
}
