package jobfilter

import (
	"time"

	"github.com/hitecherik/robo-clavius/internal/options"
	"github.com/hitecherik/robo-clavius/pkg/dateutil"
	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

func Filter(jobs []options.Job, checker *ukbankholiday.UkBankHoliday) []*options.Job {
	filtered := []*options.Job{}
	today := time.Now()
	today = dateutil.TruncateToMidnight(&today)

	for _, job := range jobs {
		if job.Monthly && !today.Before(job.Date.Time) {
			date := dateutil.SetYearMonth(&job.Date.Time, today.Year(), today.Month())

			if !today.Before(date) {
				date = date.AddDate(0, 1, 0)
			}

			job.Date = options.JobDate{Time: date}
		}

		if !today.Before(job.Date.Time) {
			continue
		}

		working := job.Date.AddDate(0, 0, -1)

		for dateutil.IsWeekend(&working) && !checker.Check(&working) {
			working = working.AddDate(0, 0, -1)
		}

		working = working.AddDate(0, 0, -1)

		if today.Equal(working) {
			filtered = append(filtered, &job)
		}
	}

	return filtered
}
