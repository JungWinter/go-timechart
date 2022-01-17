package timechart

import "strings"

type charType uint8

const (
	hour charType = iota + 1
	edge
	slot
)

type Char interface {
	Start() Char
	End() Char
	Fill() Char
	Now() Char

	Hour() Char
	Edge() Char
	Slot() Char

	String() string
}

type UnicodeChar struct {
	t charType

	start bool
	end   bool

	now bool
	in  bool
}

var _ Char = (*UnicodeChar)(nil)

func NewUnicodeChar() Char {
	return UnicodeChar{}
}

func (c UnicodeChar) Start() Char {
	c.start = true
	return c
}

func (c UnicodeChar) End() Char {
	c.end = true
	return c
}

func (c UnicodeChar) Fill() Char {
	c.in = true
	return c
}

func (c UnicodeChar) Now() Char {
	c.now = true
	return c
}

func (c UnicodeChar) Hour() Char {
	c.t = hour
	return c
}

func (c UnicodeChar) Edge() Char {
	c.t = edge
	return c
}

func (c UnicodeChar) Slot() Char {
	c.t = slot
	return c
}

func (c UnicodeChar) String() string {
	switch c.t {
	case hour:
		return c.hour()
	case edge:
		return c.edge()
	case slot:
		return c.slot()
	default:
		return ""
	}
}

func (c UnicodeChar) hour() string {
	switch {
	case c.now && c.in && c.start && !c.end:
		return "╊"
	case c.now && c.in && c.end && !c.start:
		return "╉"
	case c.now && c.in:
		return "╋"
	case c.now:
		return "╂"
	case c.in && c.start && !c.end:
		return "┾"
	case c.in && c.end && !c.start:
		return "┽"
	case c.in:
		return "┿"
	default:
		return "┼"
	}
}

func (c UnicodeChar) edge() string {
	switch {
	case c.start && c.now && c.in:
		return "┣"
	case c.start && c.now:
		return "┠"
	case c.start && c.in:
		return "┝"
	case c.start:
		return "├"
	case c.now && c.in:
		return "┫"
	case c.now:
		return "┨"
	case c.in:
		return "┥"
	default:
		return "┤"
	}
}

func (c UnicodeChar) slot() string {
	switch {
	case c.in:
		return "━"
	default:
		return "─"
	}
}

type Chars []Char

func (cc Chars) String() string {
	ss := make([]string, len(cc))
	for i, c := range cc {
		ss[i] = c.String()
	}
	return strings.Join(ss, "")
}

func (cc Chars) Repeat(n int) Chars {
	if n <= 0 {
		return nil
	}
	repeated := make(Chars, 0, len(cc)*n)
	for i := 0; i < n; i++ {
		repeated = append(repeated, cc...)
	}
	return repeated
}
