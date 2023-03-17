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
for i in $(bash -c "echo {2..${max}}"); do ./provision -u testa$i@example.com -p 12345678 --num 50 --prefix seed; done
```

This will create things and channels and connect them

[mainfluxLink]: https://github.com/mainflux/mainflux
