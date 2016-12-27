package watermeter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWatermeter(t *testing.T) {
	assert := assert.New(t)

	timeout, _ := time.ParseDuration("4s")
	wm := Watermeter{Initial: 500, Timeout: timeout}

	wm.Init()

	assert.Equal(uint64(0), wm.Gallons())
	wm.Usage = func(gallons uint64) { assert.Equal(gallons, uint64(1)) }
	wm.Update(1000)
	time.Sleep(time.Second * 3)
	wm.Usage = func(gallons uint64) { assert.Equal(gallons, uint64(2)) }
	wm.Update(1000)
	wm.Usage = func(gallons uint64) { assert.Equal(gallons, uint64(3)) }
	wm.Update(1000)
	wm.Usage = nil
	wm.Update(1000)
	assert.Equal(uint64(4), wm.Gallons())
	duration, _ := time.ParseDuration("2s")
	assert.Equal(60.0, wm.GetFlow(duration))
	time.Sleep(time.Second * 1)
	wm.Update(1000)
	assert.Equal(90.0, wm.GetFlow(duration))
}
