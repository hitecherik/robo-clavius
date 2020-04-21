package options

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hitecherik/robo-clavius/pkg/dateutil"
	"gopkg.in/yaml.v3"
)

type Schedule struct {
	Date   time.Time
	Amount string
}

type Options struct {
	Schedules []Schedule
}

var config Options
var date string
var amount string

func init() {
	flag.StringVar(&date, "date", "", "date (yyyy-mm-dd) the amount needs to be out of the pot on")
	flag.StringVar(&amount, "amount", "", "the amount to transfer on the date")
	flag.Var(&config, "config", "the path to the yaml config file")
}

func exit(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
	}

	flag.Usage()
	os.Exit(1)
}

func GetOptions() Options {
	flag.Parse()

	if date == "" && amount == "" {
		return config
	}

	if date != "" || amount != "" {
		exit(errors.New("either date or amount not provided"))
	}

	date, err := time.ParseInLocation(dateutil.ISO8601, date, time.Local)

	if err != nil {
		exit(err)
	}

	return Options{[]Schedule{{date, amount}}}
}

func (o *Options) String() string {
	return fmt.Sprint(o.Schedules)
}

func (o *Options) Set(value string) error {
	if len(o.Schedules) > 0 {
		return errors.New("config flag already set")
	}

	body, err := ioutil.ReadFile(value)

	if err != nil {
		return err
	}

	records := []interface{}{}

	if err := yaml.Unmarshal(body, &records); err != nil {
		return err
	}

	o.Schedules = make([]Schedule, len(records))

	for _, record := range records {
		record := record.(map[string]interface{})

		d := record["date"].(time.Time)
		d = dateutil.TruncateToMidnight(&d)
		a := fmt.Sprint(record["amount"].(float64))
		o.Schedules = append(o.Schedules, Schedule{d, a})
	}

	return nil
}
