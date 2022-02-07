# Docker Metrics Processor

Just a small tool to parse docker stats logs and insert them to a storage like timeseries db(coming soon) or mysql for visualising

First you need a  bash/whatever script to generate the stat files periodically.

Example just execute this every 5 seconds
`echo -e "$(date -Iseconds)\n$(docker stats --no-stream)" >> $FILE`

It creates file with this format, for example: 
`2021-11-29.log`
```
2021-11-29T00:00:00+00:00
CONTAINER ID   NAME               CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
64e9282ab459   proxy       0.00%     13.09MiB / 1.904GiB   0.67%     16.8GB / 21.6GB   130MB / 4.1kB     3
a335efdf6d91   web         0.06%     87.12MiB / 1.904GiB   4.47%     16.5GB / 16.6GB   476MB / 1.06GB    6
60487af5ff42   api         0.47%     228.6MiB / 1.904GiB   11.73%    85.7GB / 65.5GB   122MB / 561kB     27
2021-11-29T00:00:08+00:00
CONTAINER ID   NAME               CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
64e9282ab459   proxy       0.00%     13.08MiB / 1.904GiB   0.67%     16.8GB / 21.6GB   130MB / 4.1kB     3
a335efdf6d91   web         0.97%     82.44MiB / 1.904GiB   4.23%     16.5GB / 16.6GB   476MB / 1.06GB    5
60487af5ff42   api         0.03%     233MiB / 1.904GiB     11.95%    85.7GB / 65.5GB   122MB / 565kB     28
```

Then the golang tool that I published in this repo is responsible for parsing such files, and inserting the data into a MySQL database (for now only MYSQL support).
Then you can form graphs to see docker containers performance, memory usage, cpu etc

To see the usage instructions run the command with --help flag

```
kristi.jorgji$ ./docker_metrics_processor --help
Usage of ./docker_metrics_processor:
  -batchSize int
        the number of metrics to insert in storage in one batch (default 1000)
  -deleteOnEnd
        if set all logs will be deleted upon successful parsing
  -inputPath string
        folder path where logs are located (default "input/")
  -maxFilesInParallel int
        the number of files to process in parallel. Ex mysql allows 10 connections in parallel so makes no sense process more then 10 files in parallel (default 10)
```

## Build
```bash
go build -o docker_metrics_processor main.go
```

## Run example on Windows
```bash
.\docker_metrics_processor.exe --inputPath='D:\OneDrive\Documents\metrics\docker' --deleteOnEnd=true
```

## Tests

To run all tests
```bash
go test ./... -v
```