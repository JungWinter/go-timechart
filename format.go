package timechart

import (
	"time"
)

type Formatter interface {
	Format([]Schedule) string
	FormatNow([]Schedule) string
	FormatWithTime([]Schedule, time.Time) string
}

type HalfHourIncrementFormatter struct {
	fn func() Char
}

func NewHalfHourIncrementFormatter(charset func() Char) HalfHourIncrementFormatter {
	return HalfHourIncrementFormatter{fn: charset}
}

func (f HalfHourIncrementFormatter) Format(ss []Schedule) string {
	return f.fill(ss).String()
}

func (f HalfHourIncrementFormatter) FormatNow(ss []Schedule) string {
	return f.FormatWithTime(ss, time.Now())
}

func (f HalfHourIncrementFormatter) FormatWithTime(ss []Schedule, t time.Time) string {
	base := f.fill(ss)
	i := f.timeToIndex(t.Round(time.Hour))
	if i > 0 {
		i -= 1
	}
	if t.Hour() == 12 && t.Minute() > 0 {
		i += 1
	}
	base[i] = base[i].Now()
	return base.String()
}

func (f HalfHourIncrementFormatter) fill(ss []Schedule) Chars {
	base := f.empty()
	if len(ss) == 0 {
		return base
	}

	for _, schedule := range OverlapSchedules(ss) {
		s, e := f.pickRange(schedule.Start, schedule.End)
		for i, c := range base[s:e] {
			if i == 0 {
				c = c.Start()
			}
			if i == len(base[s:e])-1 {
				c = c.End()
			}
			c = c.Fill()
			base[s+i] = c
		}
	}
	return base
}

func (f HalfHourIncrementFormatter) empty() Chars {
	start := f.fn().Start().Edge()
	slot := f.fn().Slot()
	hour := f.fn().Hour()
	end := f.fn().End().Edge()

	cc := append(
		Chars{start},
		Chars{slot, slot, hour}.Repeat(11)..., // 11 hours
	)
	cc = append(cc, Chars{slot, slot, end}...) // 12 o'clock
	return cc.Repeat(2)                        // am + pm
}

func (f HalfHourIncrementFormatter) pickRange(start, end time.Time) (int, int) {
	startIdx, endIdx := f.timeToIndex(start), f.timeToIndex(end)
	if start.Hour()%12 != 0 && start.Minute() == 0 {
		startIdx -= 1
	}

	return startIdx, endIdx
}

func (f HalfHourIncrementFormatter) timeToIndex(t time.Time) int {
	if t.Day() > 1 {
		return len(f.empty())
	}

	d := t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
	if d == 0 {
		return 0
	}
	iter := int(d / (30 * time.Minute)) // number of slots
	iter += iter / 2                    // number of hours
	iter += 1                           // 00:00s
	if d > 12*time.Hour {
		iter += 1 // 12:00
	}
	return iter
}
