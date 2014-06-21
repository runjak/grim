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

// Creates an empty StRingBuffer with a given size.
func NewStRingBuffer(size int) *StRingBuffer {
	return &StRingBuffer{lines: make([]string, size)}
}

// Generates a string describing a StRingBuffer.
func (s *StRingBuffer) String() string {
	return fmt.Sprintf("StRingBuffer{start: %d, end: %d, lines: %s, length: %d}", s.start, s.end, s.lines, s.length)
}

/*
  Tells, if a StRingBuffer is currently full.
  Adding more elements implies that old ones be forgotten.
*/
func (s *StRingBuffer) Full() bool {
	return s.length == len(s.lines)
}

// Tells, if there are no elements in a StRingBuffer.
func (s *StRingBuffer) Empty() bool {
	return s.length == 0
}

// Gives the number of elements currently in the StRingBuffer.
func (s *StRingBuffer) Length() int {
	return s.length
}

// Gives the number of elements a StRingBuffer can hold.
func (s *StRingBuffer) Capacity() int {
	return len(s.lines)
}

/*
  Append a variable number of strings to a StRingBuffer.
  If the StRingBuffer is full, the start will be moved, and elements will be overwritten.
  The original StRingBuffer is returned for chaining.
*/
func (s *StRingBuffer) Push(lines ...string) *StRingBuffer {
	for _, line := range lines {
		//Adjusting start, iff s.Full():
		if s.Full() {
			s.start = s.mod(s.start + 1)
		} else {
			s.length++
		}
		//Writing values:
		s.lines[s.end] = line
		s.end = s.mod(s.end + 1)
	}
	//Return for chaining:
	return s
}

/*
  Return the last string in a StRingBuffer.
  The string will be removed from the StRingBuffer.
  If the StRingBuffer is empty, "" is returned.
*/
func (s *StRingBuffer) Pop() string {
	if s.Empty() {
		return ""
	} else {
		s.length--
	}
	s.end = s.mod(s.end - 1)
	ret := s.lines[s.end]
	s.lines[s.end] = ""
	return ret
}

/*
  Prepend a variable number of strings to a StRingBuffer.
  If the StRingbuffer is full, the end will be moved, and elements will be overwritten.
  The original StRingBuffer is returned for chaining.
*/
func (s *StRingBuffer) Unshift(lines ...string) *StRingBuffer {
	for _, line := range lines {
		//Adjusting end, iff s.Full():
		if s.Full() {
			s.end = s.mod(s.end - 1)
		} else {
			s.length++
		}
		//Writing values
		s.start = s.mod(s.start - 1)
		s.lines[s.start] = line
	}
	//Return for chaining:
	return s
}

/*
  Return the first string in a StringBuffer.
  The string will be removed from the StRingBuffer.
  If the StRingBuffer is empty, "" is returned.
*/
func (s *StRingBuffer) Shift() string {
	if s.Empty() {
		return ""
	} else {
		s.length--
	}
	ret := s.lines[s.start]
	s.lines[s.start] = ""
	s.start = s.mod(s.start + 1)
	return ret
}

/*
  Execute a function on all strings in a StRingBuffer, from start to finish.
  Strings are replaced with the ones returned by the given function.
  The original StRingBuffer is returned for chaining.
*/
func (s *StRingBuffer) Map(f func(string) string) *StRingBuffer {
	l := s.Length()
	for i := 0; i < l; i++ {
		s.Push(f(s.Shift()))
	}
	return s
}

// Like Map, but in reverse order.
func (s *StRingBuffer) MapR(f func(string) string) *StRingBuffer {
	l := s.Length()
	for i := 0; i < l; i++ {
		s.Unshift(f(s.Pop()))
	}
	return s
}

// Like Map, but without replacing strings.
func (s *StRingBuffer) Each(f func(string)) *StRingBuffer {
	return s.Map(mkId(f))
}

// Like Each, but in reverse order.
func (s *StRingBuffer) EachR(f func(string)) *StRingBuffer {
	return s.MapR(mkId(f))
}

// Returns a slice representing the contents of a StRingBuffer.
func (s *StRingBuffer) Slice() []string {
	slice := make([]string, s.Length())
	if !s.Empty() {
		if s.start < s.end {
			end := s.end + 1
			copy(slice, s.lines[s.start:end])
		} else {
			part1 := s.lines[s.start:]
			part2 := s.lines[0:s.end]
			copy(slice, part1)
			copy(slice[len(part1):], part2)
		}
	}
	return slice
}

/*
  Returns a StRingBuffer representing a given []string.
  The StRingBuffer will be full.
  Shift will hold the first element of the given slice,
  and Pop will hold the last.
*/
func Unslice(lines []string) *StRingBuffer {
	srb := &StRingBuffer{lines: lines, length: len(lines)}
	srb.end = srb.mod(srb.length)
	return srb
}

// Helper function for the two Each functions.
func mkId(f func(string)) func(string) string {
	return func(s string) string {
		f(s)
		return s
	}
}

/*
  Helper function to make sure indices stay inside a StRingBuffer.
  Since a%b may be < 0, we need to add the Capacity of a StRingBuffer, to make sure s.mod(a) >= 0.
*/
func (s *StRingBuffer) mod(x int) int {
	l := len(s.lines)
	return (x + l) % l
}
