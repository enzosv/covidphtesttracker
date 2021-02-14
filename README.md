# What this is

Source code for bot powering [Telegram Channel](https://t.me/covidphtesttracker).

Data sourced from https://doh.gov.ph/covid19tracker.

Inspired by the good work being done over at [PH Coronavirus Updates](https://t.me/phcoronavirus).

# What it does

1. Download Testing Aggregates csv from provided Google Drive link
2. Iterate through csv and add `daily_output_unique_individuals` and `daily_output_positive_individuals` grouped by provided date
3. Send message to telegram

# Requirements

1. go 1.0+
2. `config.json` file

- Sample:

  ```json
  {
    "telegram": {
      "chat_id": "",
      "bot_id": "",
      "url": "https://api.telegram.org"
    },
    "gdrive": {
      "api_key": "",
      "url": "https://www.googleapis.com/drive/v2/files",
      "filename_substring": "Testing Aggregates.csv"
    }
  }
  ```

- Instructions [here](https://core.telegram.org/bots#3-how-do-i-create-a-bot) for setting up a telegram bot.
- Instructions [here](https://developers.google.com/drive/api/v2/enable-drive-api) for enabling Google Drive API

# Building

Just invoke `go build`

# Running

```
./covidphtesttracker \
-d 2021-02-12 \
-l https://drive.google.com/drive/folders/1x-zy7hTT19cJ9Hin1B0WYrJksrTI6EU4 \
-c config.json
```

## Flags

1. `-c`: Path to config.json file
2. `-d`: Date to check
3. `-l`: Link to the Google Drive folder containing the DOH files
