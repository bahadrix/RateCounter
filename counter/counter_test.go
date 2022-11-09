package counter

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestRateCounter_Basic(t *testing.T) {

	cases := []struct {
		window       time.Duration
		tickLength   time.Duration
		hitPerSecond int
		delay        time.Duration
	}{
		{window: 5 * time.Second, tickLength: time.Second, hitPerSecond: 10},
		{window: 13 * time.Second, tickLength: 7 * time.Second, hitPerSecond: 3},
	}

	for i, c := range cases {
		t.Logf("Runing case %d/%d", i+1, len(cases))
		counter, err := NewRateCounter(c.window, c.tickLength)

		assert.Nil(t, err)

		for i := 0; i < c.hitPerSecond; i++ {
			counter.Hit()
			time.Sleep(time.Duration(1000/c.hitPerSecond) * time.Millisecond)
		}

		assert.Equal(t, c.hitPerSecond, int(math.Round(counter.Get())), "Rate calculation is wrong")
	}

}

func TestRateCounter_DelayArray(t *testing.T) {

	cases := []struct {
		window     time.Duration
		tickLength time.Duration
		periods    []float64
		want       int
	}{
		{window: 2 * time.Second, tickLength: 1 * time.Second, periods: []float64{1, 1, 1, 1, 1}, want: 1},
		{window: 2 * time.Second, tickLength: 1 * time.Second, periods: []float64{0.1, 0.1, 0.1, 3.0, 1}, want: 1},
		{window: 5 * time.Second, tickLength: 1 * time.Second, periods: []float64{1, 1, 1, 1}, want: 1},
		{window: 5 * time.Second, tickLength: 1 * time.Second, periods: []float64{0.5, 0.5, 0.5}, want: 2},
		{window: 5 * time.Second, tickLength: 1 * time.Second, periods: []float64{0.5, 0.3, 0.5}, want: 2},
		{window: 5 * time.Second, tickLength: 1 * time.Second, periods: []float64{1, 0.3, 0.3, 0.5, 1, 0.5, 1}, want: 2},
		{window: 2 * time.Second, tickLength: 1 * time.Second, periods: []float64{1, 3}, want: 0},
	}

	for i, c := range cases {
		t.Logf("Runing case %d/%d", i+1, len(cases))
		counter, err := NewRateCounter(c.window, c.tickLength)

		assert.Nil(t, err)

		for _, d := range c.periods {
			counter.Hit()
			time.Sleep(time.Duration(1000*d) * time.Millisecond)
		}
		rate := counter.Get()
		assert.Equalf(t, c.want, int(math.Round(rate)), "Rate calculation is wrong, actual: %f", rate)
	}

}
