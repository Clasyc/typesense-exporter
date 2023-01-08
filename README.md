# Typesense metrics exporter for Prometheus
This is [Prometheus](https://prometheus.io/) metrics exporter for [Typesense](https://typesense.org) 
[metrics](https://typesense.org/docs/0.23.1/api/cluster-operations.html#cluster-metrics).

## Usage

### Docker

```
docker build -t typesense-exporter .
docker run -d -p 9101:9101 -e TYPESENSE_API_KEY=xyz -e TYPESENSE_HOST=http://typesense:8108 typesense-exporter
```

## Local development

### Run

```
go build -o typesense_exporter
./typesense_exporter
```

### Running with docker-compose
    
```
export TYPESENSE_API_KEY=xyz
docker-compose up
```