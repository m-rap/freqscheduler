package main

import (
	"fmt"
	"time"
)

const CounterToSleepMax = 50
const SleepMs = 1

var counterToSleep int = CounterToSleepMax

func countToSleep() {
	counterToSleep--
	if counterToSleep <= 0 {
		time.Sleep(SleepMs * time.Millisecond)
		counterToSleep = CounterToSleepMax
	}
}

type Scheduler struct {
	StartMs     int64
	PrevFrameMs int64
	Interval    int64
	NextFrameMs int64
}

type ISchedule interface {
	Schedule()
}

func (s *Scheduler) Schedule() {
	nowMs := time.Now().UnixMilli()
	nowRelativeMs := nowMs - s.StartMs
	if nowRelativeMs >= s.NextFrameMs {
		overtimeMs := (nowRelativeMs - s.PrevFrameMs) % s.Interval
		currFrameMs := nowRelativeMs - overtimeMs
		fmt.Printf("at %v (overtime %v)\n", currFrameMs, overtimeMs)
		s.PrevFrameMs = currFrameMs
		s.NextFrameMs = s.PrevFrameMs + s.Interval
	}
}

func main() {
	fmt.Println("freqscheduler")
	s := Scheduler{
		StartMs:     time.Now().UnixMilli(),
		PrevFrameMs: int64(0),
		Interval:    int64(500),
		NextFrameMs: int64(0),
	}
	for {
		s.Schedule()
		countToSleep()
	}
}
