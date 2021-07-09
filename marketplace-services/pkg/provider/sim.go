package provider

import (
	"math/rand"
	"sync"
	"time"
)

type ChannelNotFound struct {
	msg string
}

func (e ChannelNotFound) Error() string {
	return e.msg
}

type sensorSimulator struct {
	min       int
	max       int
	frequency time.Duration
	timeout   time.Duration
	generator *rand.Rand
	channels  map[chan<- int]struct{}
	done      chan struct{}
	sync.RWMutex
}

func NewSensorSimulator(min int, max int, frequency time.Duration, timeout time.Duration) *sensorSimulator {
	source := rand.NewSource(time.Now().UnixNano())
	return &sensorSimulator{
		min:       min,
		max:       max,
		frequency: frequency,
		timeout:   timeout,
		generator: rand.New(source),
		channels:  make(map[chan<- int]struct{}),
		done:      make(chan struct{}),
	}
}

func (s *sensorSimulator) Simulate() {
	ticker := time.NewTicker(s.frequency)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			measurement := s.generator.Int63n(int64(s.max)-int64(s.min+1)) + int64(s.min)
			if s.timeout != 0 {
				s.NotifyWithTimeout(int(measurement), s.timeout)
			} else {
				s.Notify(int(measurement))
			}
		}
	}
}

func (s *sensorSimulator) Close() {
	close(s.done)
}

func (s *sensorSimulator) Attach(channel chan<- int) {
	s.Lock()
	defer s.Unlock()

	s.channels[channel] = struct{}{}
}

func (s *sensorSimulator) Detach(channel chan<- int) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.channels[channel]
	if !ok {
		return &ChannelNotFound{}
	}

	delete(s.channels, channel)

	return nil
}

func (s *sensorSimulator) Notify(data int) {
	s.RLock()
	defer s.RUnlock()

	for channel := range s.channels {
		channel <- data
	}
}

func (s *sensorSimulator) NotifyWithTimeout(data int, timeout time.Duration) {
	s.RLock()
	defer s.RUnlock()

	for channel := range s.channels {
		select {
		case channel <- data:
		case <-time.After(timeout):
		}
	}
}
