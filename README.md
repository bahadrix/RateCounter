# Rate Counter

Calculates hit frequency with given time window and precision.
Runs synchronously and works seamlessly with main thread without blocking.

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
