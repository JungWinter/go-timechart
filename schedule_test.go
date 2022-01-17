package timechart

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestOverlapSchedules(t *testing.T) {
	cases := []struct {
		name      string
		schedules []Schedule
		expected  []Schedule
	}{
		{
			name: "identify",
			schedules: []Schedule{
				{newTime(12, 0), newTime(16, 0)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(16, 0)},
			},
		},
		{
			name: "patch",
			schedules: []Schedule{
				{newTime(12, 0), newTime(16, 0)},
				{newTime(16, 0), newTime(19, 0)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(19, 0)},
			},
		},
		{
			name: "order insensitive",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
				{newTime(12, 0), newTime(16, 0)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(19, 0)},
			},
		},
		{
			name: "overlapped",
			schedules: []Schedule{
				{newTime(16, 0), newTime(19, 0)},
				{newTime(12, 0), newTime(17, 0)},
				{newTime(13, 0), newTime(14, 0)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(19, 0)},
			},
		},
		{
			name: "not overlapped",
			schedules: []Schedule{
				{newTime(12, 0), newTime(15, 0)},
				{newTime(16, 0), newTime(19, 0)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(15, 0)},
				{newTime(16, 0), newTime(19, 0)},
			},
		},
		{
			name: "contains",
			schedules: []Schedule{
				{newTime(12, 0), newTime(15, 0)},
				{newTime(13, 0), newTime(14, 0)},
				{newTime(12, 0), newTime(14, 30)},
			},
			expected: []Schedule{
				{newTime(12, 0), newTime(15, 0)},
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := OverlapSchedules(tc.schedules)
			assert.Equal(t, tc.expected, got, cmp.Diff(tc.expected, got))
		})
	}
}

func newTime(h, m int) time.Time {
	return NewTime(h, m, 0)
}
