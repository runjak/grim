/*
  This package implements a ring buffer of strings.
  It's intended usage is to have a log of recently passed lines,
  where its simple to add new lines, and old lines will get overwritten at once.
*/
package stRingBuffer

import "fmt"

type StRingBuffer struct {
	start int
	end   int
	lines []string
	// Mainly so we can distinguish empty and full:
	length int
}

func NewStRingBuffer(size int) *StRingBuffer {
	return &StRingBuffer{lines: make([]string, size)}
}

func (s *StRingBuffer) String() string {
	return fmt.Sprintf("StRingBuffer {start: %d, end: %d, lines: %s, length: %d}", s.start, s.end, s.lines, s.length)
}

func (s *StRingBuffer) Full() bool {
	return s.length == len(s.lines)
}

func (s *StRingBuffer) Empty() bool {
	return s.length == 0
}

func (s *StRingBuffer) Length() int {
	return s.length
}

func (s *StRingBuffer) Push(lines ...string) *StRingBuffer {
	for _, line := range lines {
		//Adjusting start, iff s.Full():
		if s.Full() {
			s.start = s.mod(s.start + 1)
		} else {
			s.length++
		}
		//Writing values:
		s.end = s.mod(s.end + 1)
		s.lines[s.end] = line
	}
	//Return for chaining:
	return s
}

func (s *StRingBuffer) Pop() string {
	if s.Empty() {
		return ""
	} else {
		s.length--
	}
	ret := s.lines[s.end]
	s.lines[s.end] = ""
	s.end = s.mod(s.end - 1)
	return ret
}

func (s *StRingBuffer) Unshift(lines ...string) *StRingBuffer {
	for _, line := range lines {
		//Adjusting end, iff s.Full():
		if s.Full() {
			s.end = s.mod(s.end - 1)
		} else {
			s.length++
		}
		//Writing values
		s.lines[s.start] = line
		s.start = s.mod(s.start - 1)
	}
	//Return for chaining:
	return s
}

func (s *StRingBuffer) Shift() string {
	if s.Empty() {
		return ""
	} else {
		s.length--
	}
	s.start = s.mod(s.start + 1)
	ret := s.lines[s.start]
	s.lines[s.start] = ""
	return ret
}

func (s *StRingBuffer) Map(f func(string) string) *StRingBuffer {
	l := s.Length()
	for i := 0; i < l; i++ {
		s.Push(f(s.Shift()))
	}
	return s
}

func (s *StRingBuffer) MapR(f func(string) string) *StRingBuffer {
	l := s.Length()
	for i := 0; i < l; i++ {
		s.Unshift(f(s.Pop()))
	}
	return s
}

func (s *StRingBuffer) Each(f func(string)) *StRingBuffer {
	return s.Map(mkId(f))
}

func (s *StRingBuffer) EachR(f func(string)) *StRingBuffer {
	return s.MapR(mkId(f))
}

// Helper function for the two Each functions
func mkId(f func(string)) func(string) string {
	return func(s string) string {
		f(s)
		return s
	}
}

// By adding the len to x, we make sure to stay >= 0
func (s *StRingBuffer) mod(x int) int {
	l := len(s.lines)
	return (x + l) % l
}
