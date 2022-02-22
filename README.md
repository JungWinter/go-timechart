# go-timechart

## Install

```shell
$ go get github.com/JungWinter/go-timechart
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/JungWinter/go-timechart"
)

func main() {
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

	now := time.Date(2022, 2, 1, 10, 0, 0, 0, time.UTC)

	got = f.FormatWithTime(ss, now)
	fmt.Println(got)

	// Output: ├─━┿━━┿━━┿━━┿━━┿━─┼──┼──┼──┼──╂──┼──┤├──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┼──┤
}

```
