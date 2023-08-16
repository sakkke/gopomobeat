package gopomobeat

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPomobeat(t *testing.T) {
	pomobeat := NewPomobeat()
	assert.WithinDuration(t, time.Now(), pomobeat.time, time.Second)
}

func TestAddEventListener(t *testing.T) {
	pomobeat := NewPomobeat()
	listener := EventListener(func(p Pomobeat) {})
	id := pomobeat.AddEventListener(WorkTime, listener)
	assert.Equal(t, EventID(1), id)
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

func TestGetEventListeners(t *testing.T) {
	pomobeat := NewPomobeat()
	event := WorkTime
	pomobeat.AddEventListener(event, EventListener(func(p Pomobeat) {}))
	assert.Contains(t, pomobeat.GetEventListeners(event), EventID(1))
}

func TestGetEventType(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, ShortBreak, pomobeat.GetEventType())
}

func TestGetNextEventType(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070515, 0))
	assert.Equal(t, WorkTime, pomobeat.GetNextEventType())
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

func TestListen(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(1692070739, 0))

	ok := false
	pomobeat.AddEventListener(WorkTime, EventListener(func(p Pomobeat) {
		ok = true
	}))

	done := make(chan bool)
	go func() {
		// Wait for 1 second
		pomobeat.Listen()

		done <- true
	}()

	select {
	case <-done:
		assert.Fail(t, "Should never reach here.")
	case <-time.After(2 * time.Second):
		// Pass
	}

	assert.Equal(t, true, ok)
}

func TestNotifyEventListeners(t *testing.T) {
	pomobeat := NewPomobeat()
	event := WorkTime
	pomobeat.AddEventListener(event, EventListener(func(p Pomobeat) {
		fmt.Print("ok")
	}))

	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pomobeat.NotifyEventListeners(event)

	w.Close()
	_, _ = io.Copy(&buf, r)
	os.Stdout = stdout

	assert.Equal(t, "ok", buf.String())
}

func TestRemoveEventListeners(t *testing.T) {
	pomobeat := NewPomobeat()
	event := WorkTime
	id := pomobeat.AddEventListener(event, EventListener(func(p Pomobeat) {}))
	pomobeat.RemoveEventListener(event, id)
	assert.Equal(t, 0, len(pomobeat.GetEventListeners(event)))
}

func TestSync(t *testing.T) {
	pomobeat := NewPomobeat()
	pomobeat.SetTime(time.Unix(0, 0))

	pomobeat.Sync()
	assert.WithinDuration(t, time.Now(), pomobeat.GetTime(), 5*time.Millisecond)
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
