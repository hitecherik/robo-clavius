package main

import (
	"fmt"
	"os"

	"github.com/hitecherik/robo-clavius/internal/cacher"
	"github.com/hitecherik/robo-clavius/internal/jobfilter"
	"github.com/hitecherik/robo-clavius/internal/options"
	"github.com/hitecherik/robo-clavius/pkg/ifttt"
)

func main() {
	clean, opts, dryrun := options.GetOptions()
	checker, err := cacher.CacheOrFetch(opts.CacheFile, !dryrun)

	if err != nil {
		panic(err)
	}

	sender := ifttt.New(opts.Key)
	jobs := jobfilter.Filter(opts.Jobs, checker)
	errors := 0

	for _, job := range jobs {
		if dryrun {
			fmt.Printf("Would have triggered %v\n", job)
			continue
		}

		payload := ifttt.Payload{Value1: fmt.Sprint(job.Amount)}
		if err := sender.Send(job.Event, &payload); err != nil {
			errors = errors + 1
			defer os.Exit(errors)
			os.Stderr.WriteString(fmt.Sprintln(err))
		}
	}

	if clean && errors == 0 {
		opts.Clean()

		if dryrun {
			fmt.Printf("Would have saved config %v\n", opts)
		} else {
			opts.Save()
		}
	}
}
