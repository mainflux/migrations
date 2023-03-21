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

2023/03/20 18:21:45 Loaded configuration
{"level":"info","message":"starting export from version 0.13.0","ts":"2023-03-20T15:21:45.766124907Z"}
{"level":"debug","message":"connected to things database","ts":"2023-03-20T15:21:45.772270971Z"}
{"level":"debug","message":"retrieved things from database","ts":"2023-03-20T15:21:45.78720368Z"}
{"level":"debug","message":"written things to csv file","ts":"2023-03-20T15:21:45.795162923Z"}
{"level":"debug","message":"retrieved channels from database","ts":"2023-03-20T15:21:45.809424493Z"}
{"level":"debug","message":"written channels to csv file","ts":"2023-03-20T15:21:45.816320513Z"}
{"level":"debug","message":"retrieved connections from database","ts":"2023-03-20T15:21:46.196243899Z"}
{"level":"debug","message":"written connections to csv file","ts":"2023-03-20T15:21:46.204984177Z"}
{"level":"info","message":"finished exporting from version 0.13.0","ts":"2023-03-20T15:21:46.205019163Z"}
```

### 2. Import To Version 0.14.0

Make sure you have started mainflux deployment with version 0.14.0

```bash
./build/mainflux-migrate -t 0.14.0 -o import
```

## Testing

If you want to seed the database from version 0.13.0 you can run

```bash
max=10
for i in $(bash -c "echo {0..${max}}"); do ./provision -u testa$i@example.com -p 12345678 --num 50 --prefix seed; done
```

This will create things and channels and connect them

[mainfluxLink]: https://github.com/mainflux/mainflux
