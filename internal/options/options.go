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
	Date    time.Time `yaml:"date"`
	Amount  float64   `yaml:"amount"`
	Event   string    `yaml:"event"`
	Monthly bool      `yaml:"monthly,omitempty"`
}

type Options struct {
	Key  string `yaml:"key"`
	Jobs []Job  `yaml:"jobs"`
}

var config Options
var dryrun bool

func init() {
	flag.BoolVar(&dryrun, "dryrun", false, "print what you would have done rather than doing it")
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

	if err := yaml.Unmarshal(body, o); err != nil {
		return err
	}

	for i, job := range o.Jobs {
		o.Jobs[i].Date = dateutil.TruncateToMidnight(&job.Date)
	}

	return nil
}

func (j *Job) String() string {
	return fmt.Sprintf("Job { %v, %v, %v }", j.Date.Format(dateutil.ISO8601), j.Amount, j.Event)
}
