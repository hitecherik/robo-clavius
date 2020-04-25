# 🤖 `robo-clavius`

![Go](https://github.com/hitecherik/robo-clavius/workflows/Go/badge.svg?branch=master)

Monzo's savings pots are great, but you can't schedule withdrawals very effectively: what date should you schedule your monthly withdrawal if you want your money available on the day of a monthly bill?

`robo-clavius` understands UK bank holidays and weekends and uses IFTTT webhooks to withdraw your money on the right day.

## Installation

1. Set up a recipe on [IFTTT](https://ifttt.com) that is triggered by a [webhook](https://ifttt.com/maker_webhooks) that moves whatever is in `{{Value1}}` out of your chosen pot. Name your webhook event something like `withdraw_from_pot`, but remember the name for later!
2. Create a file `config.yaml` (example at [configs/config-example.yaml](configs/config-example.yaml)) and populate it with `robo-clavius` configuration information (transfer dates, amounts, and IFTTT events; your IFTTT key and a cache file)
3. Build `robo-clavius` my running `go build -o robo-clavius cmd/robo-clavius/main.go`.
4. Cron `./robo-clavius -config path/to/config.yaml` to run once a day.

## Usage

```
Usage of ./robo-clavius:
  -clean
        remove old jobs from the yaml file on completion
  -config value
        the path to the yaml config file
  -dryrun
        print what you would have done rather than doing it
```

## FAQs

### What's behind the name?

It's named after [Christopher Clavius](https://en.wikipedia.org/wiki/Christopher_Clavius), an astronomer who contributed to the Gregorian Calendar.

### Why use IFTTT rather than the Monzo API?

Using IFTTT webhooks is much easier and means we don't have to use OAuth2 to authenticate with Monzo.

Also, using IFTTT makes it much easier to use this code to trigger other recipes on IFTTT – perhaps even interact with other banks if they become available!

## Licence

Copyright &copy; Alexander Nielsen, 2020. Licenced under the [MIT Licence](LICENCE).
