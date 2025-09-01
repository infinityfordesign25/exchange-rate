# exchange-rate

A Go command-line tool to convert currencies using [ExchangeRate-API](https://www.exchangerate-api.com/docs/free).  
This project uses the **Open API** tier:

- No API key required  
- Updates once per day  
- Rate limited  

## Features

- Fetches daily exchange rates from [ExchangeRate-API](https://www.exchangerate-api.com).
- Converts between any two ISO 4217 currency codes (e.g. USD, EUR, JPY).
- Prints local and UTC timestamps for the last updated time.
- No external dependencies beyond the Go standard library.

## Installation

### Local Build

```bash
git clone https://github.com/jsubroto/exchange-rate.git
cd exchange-rate
go build -o bin/exchange-rate
```

Run with:

```bash
./bin/exchange-rate -source USD -target EUR -amount 100
```

### Global install

```bash
go install github.com/jsubroto/exchange-rate@latest
```

Run with:

```bash
exchange-rate -source USD -target EUR -amount 100
```

Example output:

```text
100.00 USD = 85.61 EUR
Last updated on:
Local: 2025-08-31T12:34:56-04:00
UTC:   2025-08-31T16:34:56Z
```

## Flags

- `-source` : source currency (3-letter code, required)
- `-target` : target currency (3-letter code, required)
- `-amount` : amount to convert (default: 1)

## Attribution
<a href="https://www.exchangerate-api.com">Rates By Exchange Rate API</a>

## License
[MIT](LICENSE)
