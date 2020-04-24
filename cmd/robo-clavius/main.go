package main

import (
	"fmt"

	"github.com/hitecherik/robo-clavius/internal/jobfilter"
	"github.com/hitecherik/robo-clavius/internal/options"
	"github.com/hitecherik/robo-clavius/pkg/ifttt"
	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

func main() {
	opts, dryrun := options.GetOptions()

	// TODO: fetch cleverly
	checker, err := ukbankholiday.Fetch()

	if err != nil {
		panic(err)
	}

	sender := ifttt.New(opts.Key)
	jobs := jobfilter.Filter(opts.Jobs, checker)

	for _, job := range jobs {
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
