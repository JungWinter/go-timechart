package timechart

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHalfHourIncrementFormatter(t *testing.T) {
	t.Run("WithNewUnicodeChar", func(t *testing.T) {
		f := NewHalfHourIncrementFormatter(NewUnicodeChar)

		t.Run("Format", func(t *testing.T) {
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
					got := f.Format(tc.schedules)
					assert.Equal(t, tc.expected, got)
				})
			}
		})
		t.Run("timeToIndex", func(t *testing.T) {
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
					got := f.timeToIndex(tc.t)
					assert.Equal(t, tc.expectedIdx, got)
					assert.Equal(t, tc.expected, f.empty()[:got].String())
				})
			}
		})
		t.Run("pickRange", func(t *testing.T) {
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
					start, end := f.pickRange(tc.start, tc.end)
					got := string([]rune(tc.s)[start:end])
					assert.Equal(t, tc.expected, got)
				})
			}
		})
	})

}

func newTime(h, m int) time.Time {
	return time.Date(1, 1, 1, h, m, 0, 0, time.UTC)
}
