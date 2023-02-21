# Horizon

## Changelog

All notable changes to this project will be documented in [this file](./changelog.md). This project adheres to Semantic Versioning.

## Startup

- Entry point: cmd/horizon/main.go
- Add environment variables (see example config in the config-example.txt)
- Create Horizon database (make sure `DATABASE_URL` env var contains proper name of DB)
- Apply migrations (run Horizon with `db history migrate up` command-line arguments)
- Build and run Horizon

## Contribution notes

* When rendering the response that contains slice fields, **NEVER** render `null` if the slice is empty.


## Ingest Performance 
Data: catchup from start for first 400 ledgers with 16590 transaction in total

|Version of Horizon| v1| v2|
|:---:|:---:|:----:|
|Batch load of 100 ledgers| 2m13.0366828s | 1m49.039161s|
|Fixed accounts map for v2| -|55.0308887s|