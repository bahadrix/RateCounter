# RateCounter

RateCounter Calculates hit frequency within a given time window and precision. and runs synchronously and works seamlessly with the main thread without blocking.

## Window and Bucket Size

You can control memory and CPU usage against the precision with these two parameters.  Rate is calculated by `total_hits_between_window/window_size`, hence the `window` parameter represents the duration of the last hits accumulated. So the  hits older than the window will be discarded.

Time window is divided into buckets, the duration of which is determined by `bucketSize` parameter. This parameter also identifies the precision of calculation. You can set this parameter according to how often you will receive the rate information.

# Space Complexity
Counter array length is 
```go
n = windowSize / bucketSize
```

## Temporal Complexity

- For each Hit() call
  - O(1) When the hit goes into the same bucket
  - O(n) When the hit goes into the next bucket
- For each Get() call 
  - O(n)

## Usage

```go
import "github.com/bahadrix/RateCounter"

counter, _ := RateCounter.NewCounter(5 * time.Second, time.Second)

const hitPerSecond = 5

for i := 0; i < c.hitPerSecond; i++ {
    counter.Hit()
    time.Sleep(time.Duration(1000/c.hitPerSecond) * time.Millisecond)
}

fmt.Println(int(math.Round(counter.Get())))

// Prints 5



```
