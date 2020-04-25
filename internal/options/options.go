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

type JobDate struct {
	time.Time
}

type Job struct {
	Date    JobDate `yaml:"date"`
	Amount  float64 `yaml:"amount"`
	Event   string  `yaml:"event"`
	Monthly bool    `yaml:"monthly,omitempty"`
}

type Options struct {
	path      string
	Key       string `yaml:"key"`
	CacheFile string `yaml:"cache_file,omitempty"`
	Jobs      []Job  `yaml:"jobs"`
}

var clean bool
var config Options
var dryrun bool

func init() {
	flag.BoolVar(&clean, "clean", false, "remove old jobs from the yaml file on completion")
	flag.Var(&config, "config", "the path to the yaml config file")
	flag.BoolVar(&dryrun, "dryrun", false, "print what you would have done rather than doing it")
}

func exit(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
	}

	flag.Usage()
	os.Exit(-1)
}

func GetOptions() (bool, Options, bool) {
	flag.Parse()

	if config.Key == "" {
		exit(errors.New("you must specify a key in your config file"))
	}

	if len(config.Jobs) == 0 {
		exit(errors.New("you must schedule some jobs in your config file"))
	}

	return clean, config, dryrun
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

	o.path = value

	return nil
}

func (o *Options) Clean() {
	jobs := []Job{}
	today := time.Now()

	for _, job := range o.Jobs {
		if job.Monthly || !today.After(job.Date.Time) {
			jobs = append(jobs, job)
		}
	}

	o.Jobs = jobs
}

func (o *Options) Save() error {
	out, err := yaml.Marshal(o)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(o.path, out, 0644)
}

func (d JobDate) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.TaggedStyle,
		Value: d.Format(dateutil.ISO8601),
	}, nil
}

func (d *JobDate) UnmarshalYAML(value *yaml.Node) error {
	value.Decode(&d.Time)
	*d = JobDate{Time: dateutil.TruncateToMidnight(&d.Time)}
	return nil
}

func (d *JobDate) String() string {
	return d.Format(dateutil.ISO8601)
}
