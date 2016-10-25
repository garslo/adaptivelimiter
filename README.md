# adaptivelimiter

Experimental package implementing a PID Controller in Go, and using
that to control a rate limiter. The limiter attempts to keep the jobs
processing at a constant duration, adjusting the rate if the average
job duration is too high or too low. As with all PID controllers, the
`p`, `i`,and `d` parameters need careful tuning.

See `example/` for an example limiter and server for testing.

Start the "sleep server" - will sleep more or less depending on its
load:

```bash
$ go run sleepserver.go
```

```bash
$ go run main.go -p -1 -i 0 -d 1 -minrate 40 -maxrate 1000 -r 10s -t 400ms -url http://localhost:8090 -gain 1 -mindelta 0 -maxdelta 100
```
