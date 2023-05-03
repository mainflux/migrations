# Benchmarks

## Prerequisites

Tools used to benchmark are:

- [hyperfine][hyperfine-url]

## Running benchmarks

1. Exporting from version 0.10.0

    ```bash
    cd ../mainflux
    git checkout v0.10.0
    make dockers
    source .env
    cd ../migrations
    bash scripts/benchmark.sh "./build/mainflux-migrate -f 0.10.0 -o export"
    ```

    | No. of Users | No. of Things | No. of Channels | No. of Connections | Avg Time(s)          | Max Time(s)          | Description                                 |
    | ------------ | ------------- | --------------- | ------------------ | -------------------- | -------------------- | ------------------------------------------- |
    | 1            | 10            | 10              | 100                | 0.034566485720000005 | 0.09884963782000002  | Each user has 10 things and 10 channels     |
    | 1            | 100           | 100             | 10K                | 0.0326820298         | 0.0751839759         | Each user has 100 things and 100 channels   |
    | 1            | 1000          | 1000            | 1M                 | 0.0331291773         | 0.06577641320000001  | Each user has 1000 things and 1000 channels |
    | 10           | 100           | 100             | 10K                | 0.03344045772000001  | 0.08359852092        | Each user has 10 things and 10 channels     |
    | 10           | 1K            | 1K              | 1M                 | 0.03190134054        | 0.06909023374        | Each user has 100 things and 100 channels   |
    | 10           | 10K           | 10K             | 100M               | 0.03263579348        | 0.06982383468        | Each user has 1000 things and 1000 channels |
    | 100          | 1K            | 1K              | 1M                 | 0.03779480262        | 0.10620639272        | Each user has 10 things and 10 channels     |
    | 100          | 10K           | 10K             | 100M               | 0.03808557706        | 0.07619243666        | Each user has 100 things and 100 channels   |

2. Exporting from version 0.11.0

    ```bash
    cd ../mainflux
    git checkout v0.11.0
    make dockers
    source .env
    cd ../migrations
    bash scripts/benchmark.sh "./build/mainflux-migrate -f 0.11.0 -o export"
    ```

    | No. of Users | No. of Things | No. of Channels | No. of Connections | Avg Time(s)          | Max Time(s)          | Description                                 |
    | ------------ | ------------- | --------------- | ------------------ | -------------------- | -------------------- | ------------------------------------------- |
    | 1            | 10            | 10              | 100                | 0.035205581740000005 | 0.07446700714        | Each user has 10 things and 10 channels     |
    | 1            | 100           | 100             | 10K                | 0.06272539274000001  | 0.10125019274000001  | Each user has 100 things and 100 channels   |
    | 1            | 1000          | 1000            | 1M                 | 0.06907739071999999  | 0.10694224082        | Each user has 1000 things and 1000 channels |
    | 10           | 100           | 100             | 10K                | 0.03874325420000001  | 0.0803173315         | Each user has 10 things and 10 channels     |
    | 10           | 1K            | 1K              | 1M                 | 0.33314956046        | 0.37253599406000004  | Each user has 100 things and 100 channels   |
    | 10           | 10K           | 10K             | 100M               | 0.02578143656536191  | 0.20272449028000003  | Each user has 1000 things and 1000 channels |
    | 100          | 1K            | 1K              | 1M                 | 0.06146018526000001  | 0.10082907076000001  | Each user has 10 things and 10 channels     |
    | 100          | 10K           | 10K             | 100M               | 3.3249321970400003   | 5.22643773184        | Each user has 100 things and 100 channels   |

3. Exporting from version 0.12.0

    ```bash
    cd ../mainflux
    git checkout v0.12.0
    make dockers
    source docker/.env
    cd ../migrations
    bash scripts/benchmark.sh "./build/mainflux-migrate -f 0.12.0 -o export"
    ```

    | No. of Users | No. of Things | No. of Channels | No. of Connections | Avg Time(s)          | Max Time(s)          | Description                                 |
    | ------------ | ------------- | --------------- | ------------------ | -------------------- | -------------------- | ------------------------------------------- |
    | 1            | 10            | 10              | 100                | 0.031055812600000003 | 0.04967035410000001  | Each user has 10 things and 10 channels     |
    | 1            | 100           | 100             | 10K                | 0.028933612000000004 | 0.0326543716         | Each user has 100 things and 100 channels   |
    | 1            | 1000          | 1000            | 1M                 | 0.028357845040000003 | 0.030214766540000004 | Each user has 1000 things and 1000 channels |
    | 10           | 100           | 100             | 10K                | 0.03073674362        | 0.041998349020000006 | Each user has 10 things and 10 channels     |
    | 10           | 1K            | 1K              | 1M                 | 0.02917237412        | 0.03199942892        | Each user has 100 things and 100 channels   |
    | 10           | 10K           | 10K             | 100M               | 0.02759141196        | 0.029017398660000002 | Each user has 1000 things and 1000 channels |
    | 100          | 1K            | 1K              | 1M                 | 0.028190954500000004 | 0.0310052765         | Each user has 10 things and 10 channels     |
    | 100          | 10K           | 10K             | 100M               | 0.029723761519999996 | 0.031571882620000005 | Each user has 100 things and 100 channels   |

4. Exporting from version 0.13.0

    ```bash
    cd ../mainflux
    git checkout v0.13.0
    make dockers
    source docker/.env
    cd ../migrations
    bash scripts/benchmark.sh "./build/mainflux-migrate -f 0.13.0 -o export"
    ```

    | No. of Users | No. of Things | No. of Channels | No. of Connections | Avg Time(s)          | Max Time(s)          | Description                                 |
    | ------------ | ------------- | --------------- | ------------------ | -------------------- | -------------------- | ------------------------------------------- |
    | 1            | 10            | 10              | 100                | 0.03508890932000001  | 0.05541861202        | Each user has 10 things and 10 channels     |
    | 1            | 100           | 100             | 10K                | 0.03217278026        | 0.03507820016        | Each user has 100 things and 100 channels   |
    | 1            | 1000          | 1000            | 1M                 | 0.03139354810000001  | 0.034810199400000005 | Each user has 1000 things and 1000 channels |
    | 10           | 100           | 100             | 10K                | 0.03047032764        | 0.03227705704        | Each user has 10 things and 10 channels     |
    | 10           | 1K            | 1K              | 1M                 | 0.03190617944        | 0.033973302840000004 | Each user has 100 things and 100 channels   |
    | 10           | 10K           | 10K             | 100M               | 0.032045252680000004 | 0.037247937880000005 | Each user has 1000 things and 1000 channels |
    | 100          | 1K            | 1K              | 1M                 | 0.030650485220000002 | 0.03501621582        | Each user has 10 things and 10 channels     |
    | 100          | 10K           | 10K             | 100M               | 0.032022149439999995 | 0.03474878554        | Each user has 100 things and 100 channels   |

[hyperfine-url]: https://github.com/sharkdp/hyperfine/
