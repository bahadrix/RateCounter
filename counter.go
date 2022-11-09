package RateCounter

import (
	"errors"
	"sync"
	"time"
)

type timeBucket struct {
	counter       uint
	firstHitStamp int64
}

type Counter struct {
	window      time.Duration
	tickLength  time.Duration
	buckets     []timeBucket
	mx          sync.Mutex
	windowStart int64
}

func NewCounter(window time.Duration, bucketSize time.Duration) (rc *Counter, err error) {

	if bucketSize > window {
		err = errors.New("bucket size must be smaller than window")
		return
	}
	rc = &Counter{
		window:     window,
		tickLength: bucketSize,
		buckets:    make([]timeBucket, window/bucketSize),
	}
	return
}

func (rc *Counter) resetBuckets() {
	for i := 0; i < len(rc.buckets); i++ {
		rc.buckets[i].counter = 0
		rc.buckets[i].firstHitStamp = 0
	}
}

func (rc *Counter) Hit() {
	stamp := time.Now().UnixNano()
	rc.mx.Lock()
	defer rc.mx.Unlock()

	if rc.windowStart == 0 {
		rc.windowStart = stamp
	}

	bucket := &rc.buckets[rc.index(stamp)]
	bucket.counter++
	if bucket.firstHitStamp == 0 {
		bucket.firstHitStamp = stamp
	}
}

func (rc *Counter) Get() float64 {
	rc.mx.Lock()
	defer rc.mx.Unlock()

	stamp := time.Now().UnixNano()
	rc.index(stamp)
	const nanoSeconds = float64(1e9)
	var sum uint
	for _, bucket := range rc.buckets {
		sum += bucket.counter
	}
	if sum == 0 || rc.windowStart == stamp {
		return 0
	}

	passedSeconds := float64(stamp-rc.windowStart) / nanoSeconds

	//fmt.Printf("State: %v Sum: %d Secs: %f \n", rc.buckets, sum, passedSeconds)
	return float64(sum) / passedSeconds
}

func (rc *Counter) index(stamp int64) int {
	i := int((stamp - rc.windowStart) / rc.tickLength.Nanoseconds())
	if i >= len(rc.buckets) {
		steps := i - len(rc.buckets) + 1
		if steps >= len(rc.buckets) {
			rc.resetBuckets()
			rc.windowStart = stamp
		} else {
			for n := 0; n < steps; n++ {
				rc.shift()
			}
		}
		return len(rc.buckets) - 1
	}
	return i
}

func (rc *Counter) shift() {
	var i int
	for i = 1; i < len(rc.buckets); i++ {
		rc.buckets[i-1] = rc.buckets[i]
	}
	rc.buckets[i-1] = timeBucket{} // i - 1 -> last element
	if rc.buckets[0].firstHitStamp != 0 {
		rc.windowStart = rc.buckets[0].firstHitStamp
	}

}
