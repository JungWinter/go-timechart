package timechart

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHalfHourIncrementFormatter_Format(t *testing.T) {
	cases := []struct {
		name      string
		schedules []Schedule
		expected  string
	}{
		{
			name:      "empty",
			schedules: nil,
			expected:  "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "half hour",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
			},
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "multi schedules",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "hours",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
			},
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┾━━┿━━┿━━┽──┼──┼──┼──┼──┤",
		},
		{
			name: "all day",
			schedules: []Schedule{
				{newTime(0, 0), newTime(24, 0)},
			},
			expected: "┝━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┥┝━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┿━━┥",
		},
		{
			name: "patch",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
				{newTime(12, 0), newTime(16, 0)},
			},
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤┝━━┿━━┿━━┿━━┿━━┿━━┿━━┽──┼──┼──┼──┼──┤",
		},
		{
			name: "overlapped",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
				{newTime(12, 0), newTime(17, 0)},
			},
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤┝━━┿━━┿━━┿━━┿━━┿━━┿━━┽──┼──┼──┼──┼──┤",
		},
		{
			name: "both sides of noon",
			schedules: []Schedule{
				{newTime(11, 30), newTime(12, 30)},
			},
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼─━┥┝━─┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f := NewHalfHourIncrementFormatter(NewUnicodeChar)

			got := f.Format(tc.schedules)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestHalfHourIncrementFormatter_FormatWithTime(t *testing.T) {
	cases := []struct {
		name      string
		schedules []Schedule
		t         time.Time
		expected  string
	}{
		{
			name: "t is hour",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 5, 0, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━╋━─┼──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "truncated if t is half",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 5, 30, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─╂──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "00:00",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "┠─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "09:30",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 9, 30, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──╂──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "12:00",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┨├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "12:01",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 12, 1, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤┠─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "13:00",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 13, 0, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├─━╋━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤",
		},
		{
			name: "16:00",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
			},
			t:        time.Date(1970, 1, 1, 16, 0, 0, 0, time.UTC),
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──╊━━┿━━┿━━┽──┼──┼──┼──┼──┤",
		},
		{
			name: "19:00",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
			},
			t:        time.Date(1970, 1, 1, 19, 0, 0, 0, time.UTC),
			expected: "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┾━━┿━━┿━━╉──┼──┼──┼──┼──┤",
		},
		{
			name: "23:59",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 1, 23, 59, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┨",
		},
		{
			name: "24:00",
			schedules: []Schedule{
				{newTime(0, 30), newTime(5, 30)},
				{newTime(12, 30), newTime(17, 30)},
			},
			t:        time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┨",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f := NewHalfHourIncrementFormatter(NewUnicodeChar)

			got := f.FormatWithTime(tc.schedules, tc.t)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestHalfHourIncrementFormatter_timeToIndex(t *testing.T) {
	cases := []struct {
		t           time.Time
		expectedIdx int
		expected    string
	}{
		{
			t:           newTime(0, 0),
			expectedIdx: 0,
			expected:    "",
		},
		{
			t:           newTime(0, 30),
			expectedIdx: 2,
			expected:    "├─",
		},
		{
			t:           newTime(1, 0),
			expectedIdx: 4,
			expected:    "├──┼",
		},
		{
			t:           newTime(1, 30),
			expectedIdx: 5,
			expected:    "├──┼─",
		},
		{
			t:           newTime(2, 0),
			expectedIdx: 7,
			expected:    "├──┼──┼",
		},
		{
			t:           newTime(5, 30),
			expectedIdx: 17,
			expected:    "├──┼──┼──┼──┼──┼─",
		},
		{
			t:           newTime(12, 00),
			expectedIdx: 37,
			expected:    "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
		},
		{
			t:           newTime(12, 30),
			expectedIdx: 39,
			expected:    "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├─",
		},
		{
			t:           newTime(24, 0),
			expectedIdx: 74,
			expected:    "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.t.Format("15:04"), func(t *testing.T) {
			f := NewHalfHourIncrementFormatter(NewUnicodeChar)

			got := f.timeToIndex(tc.t)
			assert.Equal(t, tc.expectedIdx, got)
			assert.Equal(t, tc.expected, f.empty()[:got].String())
		})
	}
}

func TestHalfHourIncrementFormatter_pickRange(t *testing.T) {
	cases := []struct {
		s        string
		start    time.Time
		end      time.Time
		expected string
	}{
		{
			s:        "├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(0, 30),
			end:      newTime(5, 30),
			expected: "━┿━━┿━━┿━━┿━━┿━",
		},
		{
			s:        "┝━─┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(0, 0),
			end:      newTime(0, 30),
			expected: "┝━",
		},
		{
			s:        "├─━┽──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(0, 30),
			end:      newTime(1, 0),
			expected: "━┽",
		},
		{
			s:        "├──┼──┾━━┿━━┽──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(2, 0),
			end:      newTime(4, 0),
			expected: "┾━━┿━━┽",
		},
		{
			s:        "├──┼──┾━─┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(2, 0),
			end:      newTime(2, 30),
			expected: "┾━",
		},
		{
			s:        "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┾━─┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(13, 0),
			end:      newTime(13, 30),
			expected: "┾━",
		},
		{
			s:        "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼─━┥┝━─┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(11, 30),
			end:      newTime(12, 30),
			expected: "━┥┝━",
		},
		{
			s:        "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤┝━━┽──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤",
			start:    newTime(12, 0),
			end:      newTime(13, 0),
			expected: "┝━━┽",
		},
		{
			s:        "├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼─━┥",
			start:    newTime(23, 30),
			end:      newTime(24, 0),
			expected: "━┥",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%s (%s-%s)", tc.expected, tc.start.Format("15:04"), tc.end.Format("15:04")), func(t *testing.T) {
			f := NewHalfHourIncrementFormatter(NewUnicodeChar)

			start, end := f.pickRange(tc.start, tc.end)
			got := string([]rune(tc.s)[start:end])
			assert.Equal(t, tc.expected, got)
		})
	}
}
