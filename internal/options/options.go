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

type Job struct {
	Date   time.Time
	Amount float64
	Event  string
}

type Options struct {
	Key  string
	Jobs []Job
}

var config Options
var dryrun bool

func init() {
	flag.BoolVar(&dryrun, "dryrun", false, "print what you would have done ratehr than doing it")
	flag.Var(&config, "config", "the path to the yaml config file")
}

func exit(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
	}

	flag.Usage()
	os.Exit(1)
}

func GetOptions() (Options, bool) {
	flag.Parse()

	if config.Key == "" {
		exit(errors.New("you must specify a key in your config file"))
	}

	if len(config.Jobs) == 0 {
		exit(errors.New("you must schedule some jobs in your config file"))
	}

	return config, dryrun
}

func (o *Options) String() string {
	return fmt.Sprintf("Options { %v, %v }", o.Key, o.Jobs)
}

func (o *Options) Set(value string) error {
	if o.Key != "" {
		return errors.New("config flag already set")
	}

	body, err := ioutil.ReadFile(value)

	if err != nil {
		return err
	}

	var file map[string]interface{}

	if err := yaml.Unmarshal(body, &file); err != nil {
		return err
	}

	o.Key = file["key"].(string)

	jobs := file["jobs"].([]interface{})
	o.Jobs = make([]Job, len(jobs))

	for _, job := range jobs {
		job := job.(map[string]interface{})

		d := job["date"].(time.Time)
		d = dateutil.TruncateToMidnight(&d)
		o.Jobs = append(o.Jobs, Job{
			Date:   d,
			Amount: job["amount"].(float64),
			Event:  job["event"].(string),
		})
	}

	return nil
}

func (j *Job) String() string {
	return fmt.Sprintf("Job { %v, %v, %v }", j.Date.Format(dateutil.ISO8601), j.Amount, j.Event)
}
