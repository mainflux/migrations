# Mainflux Migration Tool

A tool that is used to migrate from one version of [mainflux][mainflux-url] to another.

## Installation

```bash
git clone https://github.com:mainflux/migrations.git
cd migrations
make migrate
```

## Usage

```bash
Tool for migrating from one version of mainflux to another.It migrates things, channels and thier connections.
                                Complete documentation is available at https://docs.mainflux.io

Usage:
  migrations [flags]

Flags:
  -f, --fromversion string   mainflux version you want to migrate from (default "0.13.0")
  -h, --help                 help for migrations
  -o, --operation string     export data from an existing mainflux deployment or import data to a new mainflux deployment (default "export")
  -t, --toversion string     mainflux version you want to migrate to (default "0.14.0")
```

## Example

### 1. Export From Version 0.13.0

Make sure you have started mainflux deployment with version 0.13.0

```bash
./build/mainflux-migrate -f 0.13.0 -o export

{"level":"info","message":"starting export from version 0.13.0","ts":"2023-03-30T14:32:53.725849074Z"}
{"level":"debug","message":"connected to users database","ts":"2023-03-30T14:32:53.730192129Z"}
{"level":"debug","message":"connected to things database","ts":"2023-03-30T14:32:53.737027714Z"}

 ✓  Finished Retrieveing Users
 ✓  Finished Retrieveing Things
 ✓  Finished Retrieveing Channels
 ✓  Finished Retrieveing Connection
{"level":"info","message":"finished exporting from version 0.13.0","ts":"2023-03-30T14:32:55.816558105Z"}
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

## Benchmarks

Tools used to benchmark are:

- [hyperfine][hyperfine-url]
- [docker-compose][docker-compose-url]

To run the benchmark tool run:

```bash
bash scripts/benchmark.sh "./build/mainflux-migrate -f 0.13.0 -o export"
```

The above example is to run migration benchmark during exporting

Some of the benchmark results are:

### Exporting

| No. of Users | No. of Things | No. of Channels | No. of Connections | Avg Time(s)          | Max Time(s)          | Description                                 |
| ------------ | ------------- | --------------- | ------------------ | -------------------- | -------------------- | ------------------------------------------- |
| 1            | 10            | 10              | 100                | 0.026298354760000003 | 0.029383215860000002 | Each user has 10 things and 10 channels     |
| 1            | 100           | 100             | 10K                | 0.05340002348        | 0.06367090818        | Each user has 100 things and 100 channels   |
| 1            | 1000          | 1000            | 1M                 | 0.05664883321999999  | 0.06368840492        | Each user has 1000 things and 1000 channels |
| 10           | 100           | 100             | 10K                | 0.03275979294        | 0.03625997544        | Each user has 10 things and 10 channels     |
| 10           | 1K            | 1K              | 1M                 | 0.3047701865         | 0.34431983000000005  | Each user has 100 things and 100 channels   |
| 10           | 10K           | 10K             | 100M               | 0.27441669116        | 0.29463570096        | Each user has 1000 things and 1000 channels |
| 100          | 1K            | 1K              | 1M                 | 0.05666465024        | 0.06114762994        | Each user has 10 things and 10 channels     |
| 100          | 10K           | 10K             | 100M               | 2.43419548986        | 2.90670491576        | Each user has 100 things and 100 channels   |

Example of log output is:

```bash
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 1 users, 10 things and 10 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 1 users, 100 things and 100 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 1 users, 1000 things and 1000 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 10 users, 10 things and 10 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 10 users, 100 things and 100 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 10 users, 1000 things and 1000 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 100 users, 10 things and 10 channels on Mainflux...
Stopping Docker Compose...
Starting Docker Compose...
Provisioning 100 users, 100 things and 100 channels on Mainflux...
Stopping Docker Compose...
```

[mainflux-url]: https://github.com/mainflux/mainflux
[hyperfine-url]: https://github.com/sharkdp/hyperfine/
[docker-compose-url]: https://docs.docker.com/compose/
