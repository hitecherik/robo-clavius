package ukbankholiday

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const endpoint string = "https://www.gov.uk/bank-holidays.json"
const region string = "england-and-wales"
const layout string = "2006-01-02"

type UkBankHoliday struct {
	response []byte
	days     []time.Time
}

func fromJson(bytes []byte) (*UkBankHoliday, error) {
	var data map[string]interface{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	division := data[region].(map[string]interface{})
	events := division["events"].([]interface{})
	days := make([]time.Time, len(events))

	for _, day := range events {
		holiday := day.(map[string]interface{})
		time, err := time.Parse(layout, holiday["date"].(string))

		if err != nil {
			return nil, err
		}

		days = append(days, time)
	}

	return &UkBankHoliday{bytes, days}, nil
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

func (u *UkBankHoliday) Check(day time.Time) bool {
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())

	for _, holiday := range u.days {
		if day.Equal(holiday) {
			return true
		}
	}

	return false
}
