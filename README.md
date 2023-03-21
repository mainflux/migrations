# Mainflux Migration Tool

Tool that is used to migrate from one version of [mainflux][mainfluxLink] to another.

## Installation

```bash
git clone https://github.com:mainflux/migrations.git
cd migrations
make migrate
```

## Usage

```bash
./build/mainflux-migrate --help
Tool for migrating from one version of mainflux to another.It migrates things, channels and thier connections.
                                Complete documentation is available at https://docs.mainflux.io

Usage:
  migrations [flags]

Flags:
  -f, --fromversion string   mainflux version you want to migrate from (default "0.13.0")
  -h, --help                 help for migrations
  -o, --operation string     export dataor import data to a new mainflux deployment (default "export")
  -t, --toversion string     mainflux version you want to migrate to (default "0.14.0")
```

## Example

### 1. Export From Version 0.13.0

Make sure you have started mainflux deployment with version 0.13.0

```bash
./build/mainflux-migrate -f 0.13.0 -o export

2023/03/21 18:07:17 Loaded configuration
{"level":"info","message":"starting export from version 0.13.0","ts":"2023-03-21T15:07:17.315215382Z"}
{"level":"debug","message":"connected to users database","ts":"2023-03-21T15:07:17.320802332Z"}
{"level":"debug","message":"connected to things database","ts":"2023-03-21T15:07:17.325035945Z"}
{"level":"debug","message":"retrieved users from database","ts":"2023-03-21T15:07:17.32749626Z"}
{"level":"debug","message":"written users to csv file","ts":"2023-03-21T15:07:17.328086728Z"}
{"level":"debug","message":"retrieved things from database","ts":"2023-03-21T15:07:17.345744676Z"}
{"level":"debug","message":"written things to csv file","ts":"2023-03-21T15:07:17.353103031Z"}
{"level":"debug","message":"retrieved channels from database","ts":"2023-03-21T15:07:17.364124658Z"}
{"level":"debug","message":"written channels to csv file","ts":"2023-03-21T15:07:17.371581915Z"}
{"level":"debug","message":"retrieved connections from database","ts":"2023-03-21T15:07:17.421364996Z"}
{"level":"debug","message":"written connections to csv file","ts":"2023-03-21T15:07:17.44838793Z"}
{"level":"info","message":"finished exporting from version 0.13.0","ts":"2023-03-21T15:07:17.448434447Z"}
```

### 2. Import To Version 0.14.0

Make sure you have started mainflux deployment with version 0.14.0

```bash
./build/mainflux-migrate -t 0.14.0 -o import

{"level":"info","message":"starting importing to version 0.14.0","ts":"2023-03-21T15:24:09.439589085Z"}
{"level":"debug","message":"created user token","ts":"2023-03-21T15:24:09.508898562Z"}
{"level":"debug","message":"created things","ts":"2023-03-21T15:24:12.748550379Z"}
{"level":"debug","message":"created channels","ts":"2023-03-21T15:24:15.504673694Z"}
{"level":"debug","message":"created connections","ts":"2023-03-21T15:26:41.718172029Z"}
{"level":"info","message":"finished importing to version 0.14.0","ts":"2023-03-21T15:26:41.719890531Z"} 
```

## Testing

If you want to seed the database from version 0.13.0 you can run

```bash
max=10
for i in $(bash -c "echo {1..${max}}"); do ./provision -u testa$i@example.com -p 12345678 --num 50 --prefix seed; done
```

This will create things and channels and connect them

[mainfluxLink]: https://github.com/mainflux/mainflux
