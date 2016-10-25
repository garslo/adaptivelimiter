package adaptivelimiter

import "time"

type AxisGetRequest struct {
	X     float64
	Reply chan float64
}

func NewPartitionedAxis(partitionSize float64) *PartitionedAxis {
	axis := &PartitionedAxis{[]float64{}, partitionSize, make(chan float64),
		make(chan AxisGetRequest), make(chan chan struct{}), 0}
	axis.Init()
	return axis
}

type PartitionedAxis struct {
	Buckets          []float64
	PartitionSize    float64
	IncChan          chan float64
	GetChan          chan AxisGetRequest
	RegisterWaitChan chan chan struct{}
	Zero             float64
}

func average(fs []float64) float64 {
	nFs := len(fs)
	if nFs == 0 {
		return 0
	}
	sum := float64(0)
	for _, f := range fs {
		sum += f
	}
	return sum / float64(nFs)
}

func (me *PartitionedAxis) Init() {
	t := time.Tick(time.Duration(me.PartitionSize) * time.Nanosecond)
	thisBucket := []float64{}
	nBuckets := 0
	me.Zero = float64(time.Now().UnixNano())
	notifications := [](chan struct{}){}
	go func() {
		for {
			select {
			case <-t:
				me.Buckets = append(me.Buckets, average(thisBucket))
				//fmt.Printf("average %f\n", average(thisBucket))
				thisBucket = []float64{}
				nBuckets++
			case waitCh := <-me.RegisterWaitChan:
				notifications = append(notifications, waitCh)
			case count := <-me.IncChan:
				thisBucket = append(thisBucket, count)
			case reply := <-me.GetChan:
				if reply.X <= me.PartitionSize {
					reply.Reply <- float64(0)
					break
				}
				bucketNumber := int(reply.X / me.PartitionSize)
				if bucketNumber == 0 || nBuckets < bucketNumber {
					reply.Reply <- float64(0)
					break
				}
				reply.Reply <- me.Buckets[bucketNumber-1]
			}
		}
	}()
}

func (me *PartitionedAxis) Wait() {
	waitChan := make(chan struct{})
	me.RegisterWaitChan <- waitChan
	<-waitChan
}

func (me *PartitionedAxis) Inc(n float64) {
	me.IncChan <- n
}

func (me *PartitionedAxis) Get(x float64) float64 {
	req := AxisGetRequest{x, make(chan float64)}
	me.GetChan <- req
	return <-req.Reply
}

func (me *PartitionedAxis) Last() float64 {
	last, ok := me.LastN(1)
	if ok {
		return last[0]
	}
	return 0
}

func (me *PartitionedAxis) LastN(n int) ([]float64, bool) {
	if len(me.Buckets) < n {
		return nil, false
	}
	now := float64(time.Now().UnixNano()) - me.Zero
	rets := make([]float64, n)
	for i := 0; i < n; i++ {
		x := now - float64(i)*me.PartitionSize
		rets[i] = me.Get(x)
	}
	return rets, true
}
