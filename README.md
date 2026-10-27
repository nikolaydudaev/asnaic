# Convert ASNs and AI Crawlers Data to DuckDB-friendly parquet files

`asnaic` converts JSON and text files into DuckDB-friendly JSON and Parquet files with some transformations of data along the way.

## Data Sources

These two repositories get all credit for the data:

- [AS IP Blocks](https://github.com/ipverse/as-ip-blocks)
- [AI Crawler IP Ranges](https://github.com/tbeachill/ai-crawler-ip-ranges)

I have no relation to those projects and do not collect data myself.

## Description

`asnaic` goes through files and reads them into Go structures, writes them line by line to JSON file for DuckDB to read them.
It then unnests columns with IP ranges, converts them to `INET`, calculates gateway and broadcast addresses for each range and writes it to `parquet` file.
This all is done in DuckDB instead of Go.

## Installation

```sh
git clone https://git.nikolaydudaev.com/ndudaev/asnaic.git
cd asnaic
go build -o ./build/asnaic .
```

## Usage

```sh
./build/asnaic asn # Generate AS data into JSON/Parquet
./build/asnaic aic # Generate AI Crawlers data into JSON/Parquet
```
