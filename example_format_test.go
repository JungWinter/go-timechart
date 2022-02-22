package timechart_test

import (
	"fmt"
	"time"

	"github.com/JungWinter/go-timechart"
)

func ExampleHalfHourIncrementFormatter_Format() {
	f := timechart.NewHalfHourIncrementFormatter(timechart.NewUnicodeChar)
	ss := []timechart.Schedule{
		{
			timechart.NewTime(0, 30, 0),
			timechart.NewTime(5, 30, 0),
		},
	}
	got := f.Format(ss)
	fmt.Println(got)

	// Output: ├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤
}

func ExampleHalfHourIncrementFormatter_FormatWithTime() {
	f := timechart.NewHalfHourIncrementFormatter(timechart.NewUnicodeChar)
	ss := []timechart.Schedule{
		{
			timechart.NewTime(0, 30, 0),
			timechart.NewTime(5, 30, 0),
		},
	}
	now := time.Date(2022, 2, 1, 10, 0, 0, 0, time.UTC)

	got := f.FormatWithTime(ss, now)
	fmt.Println(got)

	// Output: ├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──╂──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤
}

func ExampleMultipleTimeChart() {
	f := timechart.NewHalfHourIncrementFormatter(timechart.NewUnicodeChar)
	sss := [][]timechart.Schedule{
		{
			{
				timechart.NewTime(0, 30, 0),
				timechart.NewTime(5, 30, 0),
			},
		},
		{
			{
				timechart.NewTime(3, 0, 0),
				timechart.NewTime(10, 30, 0),
			},
			{
				timechart.NewTime(12, 30, 0),
				timechart.NewTime(17, 30, 0),
			},
		},
		{
			{
				timechart.NewTime(16, 0, 0),
				timechart.NewTime(19, 0, 0),
			},
		},
	}
	now := time.Date(2022, 2, 1, 10, 0, 0, 0, time.UTC)

	for _, ss := range sss {
		got := f.FormatWithTime(ss, now)
		fmt.Println(got)
	}

	// Output:
	// ├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──╂──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤
	// ├──┼──┼──┾━━┿━━┿━━┿━━┿━━┿━━┿━━╋━─┼──┤├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──┼──┼──┤
	// ├──┼──┼──┼──┼──┼──┼──┼──┼──┼──╂──┼──┤├──┼──┼──┼──┾━━┿━━┿━━┽──┼──┼──┼──┼──┤
}
