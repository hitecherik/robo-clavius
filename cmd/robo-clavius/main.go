package main

import (
	"fmt"
	"time"

	"github.com/hitecherik/robo-clavius/internal/options"
	"github.com/hitecherik/robo-clavius/pkg/dateutil"
	"github.com/hitecherik/robo-clavius/pkg/ifttt"
	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

func weekend(t *time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func main() {
	opts, dryrun := options.GetOptions()

	// TODO: fetch cleverly
	checker, err := ukbankholiday.Fetch()

	if err != nil {
		panic(err)
	}

	sender := ifttt.New(opts.Key)
	today := time.Now()

	today = dateutil.TruncateToMidnight(&today)

	for _, job := range opts.Jobs {
		if today.After(job.Date) || today.Equal(job.Date) {
			continue
		}

		working := job.Date.Add(-dateutil.Day)

		for weekend(&working) && !checker.Check(&working) {
			working = working.Add(-dateutil.Day)
		}

		working = working.Add(-dateutil.Day)

		if today.Equal(working) {
			if dryrun {
				fmt.Printf("Would have triggered %v\n", job.String())
				continue
			}

			payload := ifttt.Payload{Value1: fmt.Sprint(job.Amount)}
			if err := sender.Send(job.Event, &payload); err != nil {
				panic(err)
			}
		}
	}
}
