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

type Task struct {
	StartMs         int64
	PrevFrameMs     int64
	Interval        int64
	NextFrameMs     int64
	PrevIterMs      int64
	RemainToFrameMs int64
}

func (s *Task) UpdateTiming() {
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

func (s *Task) UpdateTiming2() {
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

func (s *Task) UpdateTiming3() int64 {
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
	return s.NextFrameMs
}

type Worker struct {
	busy bool
}

type Scheduler struct {
	Tasks   []Task
	Workers []Worker
}

const MaxUint64 = ^uint64(0)
const MaxInt64 = int64(MaxUint64 >> 1)

func (s *Scheduler) Loop() {
	for {
		minNextMs := MaxInt64
		for i := range s.Tasks {
			nextMs := s.Tasks[i].UpdateTiming3()
			if nextMs < minNextMs {
				minNextMs = nextMs
			}
		}
		nowMs := time.Now().UnixMilli()
		remainMs := minNextMs - nowMs
		if remainMs > 2 {
			fmt.Printf("sleep %v ms\n", remainMs-2)
			time.Sleep(time.Duration(remainMs-2) * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("freqscheduler")
	s := Scheduler{}
	s.Tasks = append(s.Tasks, Task{
		StartMs:     time.Now().UnixMilli(),
		PrevFrameMs: int64(0),
		Interval:    int64(500),
		NextFrameMs: int64(0),
	})
	s.Tasks = append(s.Tasks, Task{
		StartMs:     time.Now().UnixMilli(),
		PrevFrameMs: int64(0),
		Interval:    int64(1000),
		NextFrameMs: int64(0),
	})
	s.Loop()
}
