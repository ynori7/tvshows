# TvShows
The tvshows application scans the list of recent new premieres and for each one 
checks its IMDB rating and other details and then informs you.such as notifactions 
about interesting new releases.

## Features

### Premieres
The premieres command gathers information about new premieres from Metacritic based on
the configured genres, and for each one it looks up additional data from IMDB. It can then generate 
 an HTML report which is sent by email. 

**Additional details**
- Shows with a score (a rating weighted by the rating count) which is too low (<40 for 
   returning series, <20 for new series) will be filtered out.
-  Emails are sent using Mailjet.
 

**Usage:**

```
go run cmd/premieres/main.go --config config.yaml \
    --last-processed-path /path/to/put/file --output out
```

- `--config` This flag is required and is the path to the configuration YAML.
- `--last-processed-path` This tells the application where it should save the last date
which it processed so that it doesn't miss things or send duplicates
- `--output` This is an optional flag to indicate where html files should be saved. 
By default it's `./out`

Be sure to first copy `config.yaml.dist` to `config.yaml` and fill in the missing blanks

**Set up cronjob:**

First, build the binary:
```
go build -o premieresemailer cmd/premieres/main.go
```

Then set up the cronjob:
```
0 14 * * 1 /path/to/goprojects/src/github.com/ynori7/tvshows/premieresemailer --config /path/to/goprojects/src/github.com/ynori7/tvshows/config.yaml --output /path/to/goprojects/src/github.com/ynori7/tvshows/out --last-processed-path /path/to/put/file
```

Note that the new premieres page gets updated at irregular intervals. That's why it's necessary
to save the last processed date. 

## Project Structure

Commands are located in `cmd` and are the main entry points.

The `premieres` command gathers configuration from the `config` package, then sets up
a `application/` which orchestrates fetching data from `tvshow` and then filtering 
using the `enrich` worker pool.
