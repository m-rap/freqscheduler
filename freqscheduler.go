package main

import (
	"fmt"
	"time"
)

const CounterToSleepMax = 5000
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
	StartMs         int64
	PrevFrameMs     int64
	Interval        int64
	NextFrameMs     int64
	PrevIterMs      int64
	RemainToFrameMs int64
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

func (s *Scheduler) Schedule2() {
	nowMs := time.Now().UnixMilli()
	var iterEllapsedMs int64
	if s.PrevIterMs == 0 {
		iterEllapsedMs = 0
	} else {
		iterEllapsedMs = nowMs - s.PrevIterMs
	}
	s.PrevIterMs = nowMs
	s.RemainToFrameMs -= iterEllapsedMs
	if s.RemainToFrameMs <= 0 {
		overtimeMs := s.RemainToFrameMs * -1
		fmt.Printf("at %v (overtime %v)\n", nowMs, overtimeMs)
		overtimeMs %= s.Interval
		s.RemainToFrameMs = s.Interval - overtimeMs
	}
}

func (s *Scheduler) Schedule3() {
	nowMs := time.Now().UnixMilli()
	if s.NextFrameMs == 0 {
		s.NextFrameMs = nowMs
	}
	if nowMs >= s.NextFrameMs {
		overtimeMs := nowMs - s.NextFrameMs
		fmt.Printf("at %v (overtime %v)\n", nowMs, overtimeMs)
		overtimeMs %= s.Interval
		s.NextFrameMs = nowMs + s.Interval - overtimeMs
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
		//s.Schedule()
		s.Schedule3()
		countToSleep()
	}
}
