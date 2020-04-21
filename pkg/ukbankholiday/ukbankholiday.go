package ukbankholiday

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hitecherik/robo-clavius/pkg/dateutil"
)

const (
	endpoint = "https://www.gov.uk/bank-holidays.json"
	region   = "england-and-wales"
)

type UkBankHoliday struct {
	response []byte
	dates    []time.Time
}

func fromJson(bytes []byte) (*UkBankHoliday, error) {
	var data map[string]interface{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	division := data[region].(map[string]interface{})
	events := division["events"].([]interface{})
	dates := make([]time.Time, len(events))

	for _, date := range events {
		holiday := date.(map[string]interface{})
		time, err := time.ParseInLocation(dateutil.ISO8601, holiday["date"].(string), time.Now().Location())

		if err != nil {
			return nil, err
		}

		dates = append(dates, time)
	}

	return &UkBankHoliday{bytes, dates}, nil
}

func Fetch() (*UkBankHoliday, error) {
	resp, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return fromJson(body)
}

func Load(path string) (*UkBankHoliday, error) {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return fromJson(body)
}

func (u *UkBankHoliday) Save(path string) error {
	return ioutil.WriteFile(path, u.response, 0644)
}

func (u *UkBankHoliday) Check(date *time.Time) bool {
	truncated := dateutil.TruncateToMidnight(date)

	for _, holiday := range u.dates {
		if truncated.Equal(holiday) {
			return true
		}
	}

	return false
}
