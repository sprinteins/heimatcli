# Heimat CLI

The Heimat CLI is a CLI/REPL client for [SprintEins'](https://www.sprinteins.com/) [Intranet](https://heimat.sprinteins.com/).
It uses the heimat's REST API. 
It helps you quickly track your times in heimat with the help of autocomplete.

## Install

You can install the heimat CLI with homebrew

```sh
brew install sprinteins/tap/heimatcli
```

in terminal:

```sh
heimat
```

## Features

- track times
- show tracked times per day
- show tracked time per month
- copy tracked times from a day to another one
- show profile and stats of logged in user


### Date Handling

By default the commands operate on the current day, but if you want to change the date of the commands, the heimat cli understands relative and absolute dates.

To change date relative to today just add `+X` or `-X`.  
E.g.: `time show day +1`

You can switch to a date by adding the date of the day to commands.  
E.g.: `time show day 4`

> ⚠ Attention! You can only set the day of the current month

## Commands

- `time show day [+/-]X`: Shows the tracked times of a given date  
  For Example:
  - `time show day`: shows today's tracked times
  - `time show day +1`: shows tomorrow's tracked times
  - `time show day -1`: shows yesterday's tracked times
  - `time show day 4` shows the 4th day of the current month's tracked times
- `time add [+/-x]`: Go into the time-add mode on the given date. You can recognize if you are on a different date than today on the command prefix. It will show the absolute date.  
    For example: `heimat > time add 4` => `heimat > time add (04.12) > `
- `time copy [+/-]X [[+/-]Y]`: this command accepts two arguments.  
    The first one (`X`) is the source date; it will copy from this date the tracked times.   
    The second (`Y`) is the target; it will copy from this date tracked times. It is optional and falls back to today.  
    For example: `heimat > time copy -1 +1`: copy the tracked times from yesterday to tomorrow.
- `time show month`: shows the tracked times of the current month
- `profile`: show the profile and statistics about the logged in user