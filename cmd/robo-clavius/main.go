package main

import (
	"flag"
	"fmt"
	"time"
)

const layout = "2006-01-02"

type targets struct {
	dates []time.Time
}

func (t *targets) String() string {
	return fmt.Sprint(t.dates)
}

func (t *targets) Set(value string) error {
	date, err := time.Parse(layout, value)

	if err != nil {
		return err
	}

	t.dates = append(t.dates, date)
	return nil
}

var dates targets

func init() {
	flag.Var(&dates, "target", "date the money needs to be out of the pot on")
}

func main() {
	flag.Parse()
	fmt.Println(dates.String())
}
