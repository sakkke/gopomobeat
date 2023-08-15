package gopomobeat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPomobeat(t *testing.T) {
	pomobeat := NewPomobeat()
	assert.WithinDuration(t, time.Now(), pomobeat.time, time.Second)
}

func TestGetCalender(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, "#5013543-6", pomobeat.GetCalender())
}

func TestGetDurationUntilNextEvent(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, 225*time.Second, pomobeat.GetDurationUntilNextEvent())
}

func TestGetEvent(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, 6, pomobeat.GetEvent())
}

func TestGetSetSeconds(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, int64(2700), pomobeat.GetSetSeconds())
}

func TestGetSets(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, int64(5013543), pomobeat.GetSets())
}

func TestGetTime(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, pomobeat.time, time.Unix(1692070515, 0))
}

func TestSetTime(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, time.Unix(1692070515, 0), pomobeat.time)
}

func TestWaitForNextEvent(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070740, 0))

	startTime := time.Now()

	// Wait for 0 seconds
	pomobeat.WaitForNextEvent()

	endTime := time.Now()

	assert.WithinDuration(t, startTime, endTime, 5*time.Millisecond)
}
