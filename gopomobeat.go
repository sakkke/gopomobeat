package gopomobeat

import (
	"fmt"
	"time"
)

type EventType int8
type EventListener func(p Pomobeat)
type EventID int

const (
	WorkTime EventType = iota
	ShortBreak
	LongBreak
)

type Pomobeat struct {
	time      time.Time
	events    []EventType
	config    map[EventType]int64
	listeners map[EventType]map[EventID]EventListener
	nextID    EventID
}

func NewPomobeat() *Pomobeat {
	pomobeat := &Pomobeat{
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
		listeners: map[EventType]map[EventID]EventListener{},
		nextID:    0,
	}

	for _, event := range pomobeat.events {
		pomobeat.listeners[event] = map[EventID]EventListener{}
	}

	return pomobeat
}

func (p *Pomobeat) AddEventListener(e EventType, listener EventListener) EventID {
	p.nextID++
	p.listeners[e][p.nextID] = listener
	return p.nextID
}

func (p Pomobeat) GetCalender() string {
	return fmt.Sprintf("#%d-%d", p.GetSets(), p.GetEvent())
}

func (p Pomobeat) GetDurationUntilNextEvent() time.Duration {
	e := p.GetEvent()
	t := p.time.Unix() % p.GetSetSeconds()

	seconds := int64(0)
	if t < p.config[p.events[0]] {
		seconds = p.config[p.events[0]]
	} else {
		for i, event := range p.events {
			if i > e {
				break
			}

			seconds += p.config[event]
		}
	}

	return time.Duration(seconds-int64(t)) * time.Second
}

func (p Pomobeat) GetEvent() int {
	t := p.time.Unix() % p.GetSetSeconds()

	if t < p.config[p.events[0]] {
		return len(p.events)
	}

	seconds := int64(0)
	var i int
	var event EventType
	for i, event = range p.events {
		seconds += p.config[event]

		if seconds > t {
			break
		}
	}

	return i
}

func (p Pomobeat) GetEventType() EventType {
	event := p.GetEvent() - 1
	if event == 0 {
		event = len(p.events)
	}

	return p.events[event-1]
}

func (p Pomobeat) GetEventListeners(e EventType) map[EventID]EventListener {
	return p.listeners[e]
}

func (p Pomobeat) GetNextEventType() EventType {
	event := p.GetEvent()
	if event == len(p.events) {
		event = 0
	}

	return p.events[event]
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

func (p Pomobeat) Listen() {
	for {
		p.WaitForNextEvent()
		eventType := p.GetNextEventType()
		p.Sync()
		p.NotifyEventListeners(eventType)
	}
}

func (p Pomobeat) NotifyEventListeners(e EventType) {
	for _, listener := range p.listeners[e] {
		listener(p)
	}
}

func (p *Pomobeat) RemoveEventListener(e EventType, TargetID EventID) {
	listeners := map[EventID]EventListener{}
	for id, listener := range p.listeners[e] {
		if id == TargetID {
			continue
		}

		listeners[id] = listener
	}

	p.listeners[e] = listeners
}

func (p *Pomobeat) SetTime(t time.Time) {
	p.time = t
}

func (p *Pomobeat) Sync() {
	p.SetTime(time.Now())
}

func (p Pomobeat) WaitForNextEvent() {
	time.Sleep(p.GetDurationUntilNextEvent())
}
