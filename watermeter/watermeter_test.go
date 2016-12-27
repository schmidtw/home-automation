package watermeter

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func setNow(w *Watermeter, min int) {
	w.now = func() time.Time {
		return time.Date(2016, time.December, 25, 1, min, 0, 0, time.UTC)
	}
}

func setUsage(assert *assert.Assertions, w *Watermeter, wg *sync.WaitGroup,
	expected_gallons uint64, expected_flow float64) {

	wg.Add(1)
	w.Usage = func(actual_gallons uint64, actual_flow float64) {
		defer wg.Done()
		assert.Equal(expected_gallons, actual_gallons, "should be equal")
		assert.Equal(expected_flow, actual_flow)
	}
}

func TestWatermeterSimple(t *testing.T) {
	assert := assert.New(t)

	timeout, _ := time.ParseDuration("4s")
	wm := Watermeter{Initial: 0, Timeout: timeout}
	wm.Init()

	assert.Equal(uint64(0), wm.Gallons())
}

func TestWatermeterString(t *testing.T) {
	assert := assert.New(t)

	timeout, _ := time.ParseDuration("4s")
	wm := Watermeter{
		Initial: 0,
		Timeout: timeout,
		now: func() time.Time {
			return time.Date(2016, time.December, 25, 1, 0, 0, 0, time.UTC)
		},
	}
	wm.Init()

	wm.now = nil
	assert.Equal("{\n\tInitial: 0,\n\tTimeout: 4s,\n\tUsage: <nil>,\n\tnow: <nil>,\n\tlast_gallon{ time: 2016-12-25 01:00:00 +0000 UTC, total: 0 },\n\ttotal: 0,\n\tevents { \n\t\t{ time: 2016-12-25 01:00:00 +0000 UTC, total: 0 }\n\t}\n}", wm.String())
}

func TestWatermeterDeep(t *testing.T) {
	var wg sync.WaitGroup

	assert := assert.New(t)

	timeout, _ := time.ParseDuration("4m")
	wm := Watermeter{
		Initial: 500,
		Timeout: timeout,
		now: func() time.Time {
			return time.Date(2016, time.December, 25, 1, 0, 0, 0, time.UTC)
		},
	}
	wm.Init()

	assert.Equal(uint64(0), wm.Gallons())
	setNow(&wm, 1)
	setUsage(assert, &wm, &wg, 1, 0.5)
	wm.Update(500)

	setNow(&wm, 2)
	wm.Update(250)

	setNow(&wm, 3)
	wm.Update(250)

	setNow(&wm, 4)
	wm.Update(250)

	setNow(&wm, 5)
	setUsage(assert, &wm, &wg, 2, 0.25)
	wm.Update(250)

	duration, _ := time.ParseDuration("2s")
	assert.Equal(0.0, wm.GetFlow(duration))

	duration, _ = time.ParseDuration("2m")
	assert.Equal(0.250, wm.GetFlow(duration))

	setNow(&wm, 15)
	wm.Update(550)
	duration, _ = time.ParseDuration("11m")
	assert.Equal(0.050, wm.GetFlow(duration))

	wg.Wait()
}
