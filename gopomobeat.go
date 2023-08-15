package gopomobeat

import (
	"fmt"
	"time"
)

type EventType int8

const (
	WorkTime EventType = iota
	ShortBreak
	LongBreak
)

type Pomobeat struct {
	time   time.Time
	events []EventType
	config map[EventType]int64
}

func NewPomobeat() *Pomobeat {
	return &Pomobeat{
		time: time.Now(),
		events: []EventType{
			WorkTime,
			ShortBreak,
			WorkTime,
			ShortBreak,
			WorkTime,
			ShortBreak,
			WorkTime,
			LongBreak,
		},
		config: map[EventType]int64{
			WorkTime:   540,
			ShortBreak: 60,
			LongBreak:  360,
		},
	}
}

func (p Pomobeat) GetCalender() string {
	return fmt.Sprintf("#%d-%d", p.GetSets(), p.GetEvent())
}

func (p Pomobeat) GetDurationUntilNextEvent() time.Duration {
	e := p.GetEvent()
	t := p.time.Unix() % p.GetSetSeconds()

	seconds := int64(0)
	for i, event := range p.events {
		if i > e {
			break
		}

		seconds += p.config[event]
	}

	return time.Duration(seconds-int64(t)) * time.Second
}

func (p Pomobeat) GetEvent() int {
	t := p.time.Unix() % p.GetSetSeconds()

	seconds := int64(0)
	var i int
	var event EventType
	for i, event = range p.events {
		seconds += p.config[event]

		if seconds >= t {
			break
		}
	}

	return i
}

func (p Pomobeat) GetSetSeconds() int64 {
	seconds := int64(0)
	for _, event := range p.events {
		seconds += p.config[event]
	}

	return seconds
}

func (p Pomobeat) GetSets() int64 {
	return p.time.Unix()/p.GetSetSeconds()*int64(len(p.events)) + int64(p.GetEvent()) + 1
}

func (p Pomobeat) GetTime() time.Time {
	return p.time
}

func (p *Pomobeat) SetTime(t time.Time) {
	p.time = t
}

func (p Pomobeat) WaitForNextEvent() {
	time.Sleep(p.GetDurationUntilNextEvent())
}
