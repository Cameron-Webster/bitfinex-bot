# Bitfinex Trades Ingester

This program will programatically subscribe to the Bitfinex trades data of the trading pairs supplied by the `tickers` variable. All trade data will be stored in a TimescaleDB instance.

Note: A database connection string should be supplied via the ENV variable `DATABASE_OPTS`.
