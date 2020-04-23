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

	// TODO: factor out into own internal package
	for _, job := range opts.Jobs {
		if job.Monthly && !today.Before(job.Date) {
			date := dateutil.SetYearMonth(&job.Date, today.Year(), today.Month())

			if !today.Before(date) {
				date = date.AddDate(0, 1, 0)
			}

			job.Date = date
		}

		if !today.Before(job.Date) {
			continue
		}

		working := job.Date.AddDate(0, 0, -1)

		for weekend(&working) && !checker.Check(&working) {
			working = working.AddDate(0, 0, -1)
		}

		working = working.AddDate(0, 0, -1)

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
