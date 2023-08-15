package gopomobeat

import (
	"fmt"
	"time"
)

type Event int8

const (
	WorkTime Event = iota
	ShortBreak
	LongBreak
)

type Pomobeat struct {
	time   time.Time
	events []Event
	config map[Event]int64
}

func NewPomobeat() *Pomobeat {
	return &Pomobeat{
		time: time.Now(),
		events: []Event{
			WorkTime,
			ShortBreak,
			WorkTime,
			ShortBreak,
			WorkTime,
			ShortBreak,
			WorkTime,
			LongBreak,
		},
		config: map[Event]int64{
			WorkTime:   540,
			ShortBreak: 60,
			LongBreak:  360,
		},
	}
}

func (p Pomobeat) GetCalender() string {
	return fmt.Sprintf("#%d-%d", p.GetSets(), p.GetEvent())
}

func (p Pomobeat) GetEvent() int {
	t := p.time.Unix() % p.GetSetSeconds()

	seconds := int64(0)
	var i int
	var event Event
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
