# What this is
[Telegram Channel](https://t.me/covidphtesttracker) with data sourced from https://doh.gov.ph/covid19tracker.

Inspired by the good work being done over at [PH Coronavirus Updates](https://t.me/phcoronavirus). 

# Requirements
1. go 1.0+
2. `telegram.json` file
``` json
{
    "chat_id":"",
    "bot_id":"",
    "url": "https://api.telegram.org"
}
```

Instructions [here](https://core.telegram.org/bots#3-how-do-i-create-a-bot) for setting up a telegram bot.

# Building
Just invoke `go build`

# Running
`./phcovidtracker -ta file.csv -d 2021-02-11 -l http://bit.ly/3rJLY6Y -tc telegram.json`
## Flags
1. `-ta`: Testing Aggregates CSV file downloaded from https://doh.gov.ph/covid19tracker
2. `-d`: Date to check
3. `-l`: Link to the file provided in `-ta`
4. `-tc`: Your telegram configuration
