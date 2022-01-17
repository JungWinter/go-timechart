package timechart

import (
	"fmt"
	"time"
)

type Schedule struct {
	Start time.Time
	End   time.Time
}

func (s Schedule) String() string {
	return fmt.Sprintf("%s-%s", s.Start.Format("15:04"), s.End.Format("15:04"))
}

func (s Schedule) IsOverlapped(other Schedule) bool {
	switch {
	case s.Contains(other.Start) && timeLTE(s.End, other.End):
		return true
	case s.Contains(other.End) && timeGTE(s.Start, other.Start):
		return true
	case other.Contains(s.Start) && other.Contains(s.End):
		return true
	case s.Contains(other.Start) && s.Contains(other.End):
		return true
	default:
		return false
	}
}

func (s Schedule) Overlap(other Schedule) Schedule {
	//   ├──────┤     s
	// +     ├──────┤ other
	// = ├───━━━━───┤
	if s.Contains(other.Start) && timeLTE(s.End, other.End) {
		s.End = other.End
		return s
	}
	//       ├──────┤ s
	// + ├──────┤     other
	// = ├───━━━━───┤
	if s.Contains(other.End) && timeGTE(s.Start, other.Start) {
		s.Start = other.Start
		return s
	}
	//     ├───┤   s
	// + ├───────┤ other
	// = ├─━━━━━─┤
	if other.Contains(s.Start) && other.Contains(s.End) {
		return other
	}
	//   ├───────┤ s
	// +   ├───┤   other
	// = ├─━━━━━─┤
	if s.Contains(other.Start) && s.Contains(other.End) {
		return s
	}
	return s
}

// Contains returns true if t is in s.
// e.g. if ├──────────┤ s
//           t1 (true)  t2 (false)
func (s Schedule) Contains(t time.Time) bool {
	return timeGTE(t, s.Start) && timeLTE(t, s.End)
}

// OverlapSchedules returns new overlapped schedules from given ss.
// e.g. if [14:00-16:00, 13:00-15:00] are given, returns [13:00-16:00]
func OverlapSchedules(ss []Schedule) []Schedule {
	if len(ss) == 0 {
		return nil
	}

	exists := ss[:1]
	not := make([]Schedule, 0, len(ss))
	for _, schedule := range ss[1:] {
		found := false
		for i, exist := range exists {
			if exist.IsOverlapped(schedule) {
				exists[i] = exist.Overlap(schedule)
				found = true
				break
			}
		}
		if !found {
			not = append(not, schedule)
		}
	}
	return append(exists, not...)
}

// timeGT returns t1 is greater than t2.
func timeGT(t1, t2 time.Time) bool {
	return t1.After(t2)
}

// timeGTE returns t1 is greater than or equal to t2.
func timeGTE(t1, t2 time.Time) bool {
	return timeGT(t1, t2) || t1.Equal(t2)
}

// timeLT returns t1 is less than t2.
func timeLT(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

// timeLTE returns t1 is less than or equal to t2.
func timeLTE(t1, t2 time.Time) bool {
	return timeLT(t1, t2) || t1.Equal(t2)
}
