package main

import (
	"fmt"
	"time"

	"github.com/hitecherik/robo-clavius/internal/options"
	"github.com/hitecherik/robo-clavius/pkg/dateutil"
	"github.com/hitecherik/robo-clavius/pkg/ifttt"
	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

const (
	event = "withdraw_savings"
	key   = "[redacted]"
)

func weekend(t *time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func main() {
	opts := options.GetOptions()

	// TODO: fetch cleverly
	checker, err := ukbankholiday.Fetch()

	if err != nil {
		panic(err)
	}

	sender := ifttt.New(event, key)
	today := time.Now()

	today = dateutil.TruncateToMidnight(&today)

	for _, schedule := range opts.Schedules {
		if today.After(schedule.Date) || today.Equal(schedule.Date) {
			continue
		}

		working := schedule.Date.Add(-dateutil.Day)

		for weekend(&working) && !checker.Check(&working) {
			working = working.Add(-dateutil.Day)
		}

		working = working.Add(-dateutil.Day)

		if today.Equal(working) {
			if opts.Dryrun {
				fmt.Printf("Would have triggered event %v with amount %v for date %v\n", event, schedule.Amount, schedule.Date)
				continue
			}

			if err := sender.Send(&ifttt.Payload{Value1: schedule.Amount}); err != nil {
				panic(err)
			}
		}
	}
}
