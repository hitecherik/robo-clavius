package cacher

import (
	"fmt"
	"os"
	"time"

	"github.com/hitecherik/robo-clavius/pkg/ukbankholiday"
)

const cacheExpiryDays = 7

func fetch(path string, save bool) (*ukbankholiday.UkBankHoliday, error) {
	holidays, err := ukbankholiday.Fetch()

	if err != nil {
		return nil, err
	}

	if save {
		_ = holidays.Save(path)
	} else {
		fmt.Printf("Would have saved bank holiday cache to %v\n", path)
	}
	return holidays, nil
}

func CacheOrFetch(path string, save bool) (*ukbankholiday.UkBankHoliday, error) {
	if path == "nil" {
		return fetch(path, save)
	}

	info, err := os.Stat(path)

	if err != nil || info.ModTime().AddDate(0, 0, cacheExpiryDays).Before(time.Now()) {
		return fetch(path, save)
	}

	holidays, err := ukbankholiday.Load(path)

	if err != nil {
		return fetch(path, save)
	}

	return holidays, nil
}
