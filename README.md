# Docker Metrics Processor

Just a small tool to parse docker logs and insert them to a storage like timeserie db or mysql for visualising

## Build
```bash
go build
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