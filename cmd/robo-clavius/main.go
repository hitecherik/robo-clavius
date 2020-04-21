package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/hitecherik/robo-clavius/pkg/dateutil"
	"github.com/hitecherik/robo-clavius/pkg/ifttt"
	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

const event string = "withdraw_savings"
const key string = "[redacted]"
const amount string = "0.01"

type targets struct {
	dates []time.Time
}

func (t *targets) String() string {
	return fmt.Sprint(t.dates)
}

func (t *targets) Set(value string) error {
	date, err := time.ParseInLocation(dateutil.ISO8601, value, time.Local)

	if err != nil {
		return err
	}

	t.dates = append(t.dates, date)
	return nil
}

func weekend(t *time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

var dates targets

func init() {
	flag.Var(&dates, "target", "date the money needs to be out of the pot on")
}

func main() {
	flag.Parse()

	// TODO: fetch cleverly
	checker, err := ukbankholiday.Fetch()

	if err != nil {
		panic(err)
	}

	sender := ifttt.New(event, key)
	payload := ifttt.Payload{Value1: amount}
	today := time.Now()

	today = dateutil.TruncateToMidnight(&today)

	for _, date := range dates.dates {
		if today.After(date) || today.Equal(date) {
			continue
		}

		working := date.Add(-dateutil.Day)

		for weekend(&working) && !checker.Check(&working) {
			working = working.Add(-dateutil.Day)
		}

		working = working.Add(-dateutil.Day)

		if today.Equal(working) {
			if err := sender.Send(&payload); err != nil {
				panic(err)
			}
		}
	}
}
