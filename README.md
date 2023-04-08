# Spacebox-crawler

![linter](https://github.com/bro-n-bro/spacebox-crawler/actions/workflows/lint.yml/badge.svg)


Spacebox-crawler is a central part of [Spacebox](https://github.com/bro-n-bro/spacebox) indexer. Crawler pull data from the node, parse, and puts info into appropriate topics in the Apache Kafka.

## build

```bash
docker build -t spacebox-crawler:latest .
```

## run

Running crawler standalone is pretty much pointless, so please refer to the main [Spacebox repo](https://github.com/bro-n-bro/spacebox#readme) to find out how to start the whole setup.
