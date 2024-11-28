# Ecommerce Analysis Project
==========================

## Overview
------------

This project is designed to analyze customer revenue data from an ecommerce platform. It uses a MySQL database to store the data and a Go program to perform the analysis.

## Getting Started
---------------

To get started with this project, you will need to have the following installed on your system:

* Go (version 1.23.3 or later)
* Docker

### Running the Project

To run the project, follow these steps:

1. Clone the repository to your local machine.
2. Start the Docker container for the MySQL database by running the command `docker-compose up`.
3. run the Go program by running the command `go run cmd/analyzer/main.go`.

## Configuration
-------------

The project uses environment variables to configure the database connection and other settings. The following variables are used:

* `DB_HOST`: the hostname of the MySQL database server
* `DB_PORT`: the port number of the MySQL database server
* `DB_USER`: the username to use when connecting to the MySQL database
* `DB_PASSWORD`: the password to use when connecting to the MySQL database
* `DB_NAME`: the name of the MySQL database to use

These variables can be set in the `.env` file or as environment variables on your system.

## Analysis
----------

The Go program performs the following analysis:

* Loads purchase events from the database
* Calculates the revenue for each customer
* Sorts the customers by revenue in descending order
* Calculates quantile statistics for the top 2.5% of customers by revenue
* Exports the results to a new table in the database

The analysis is performed using the following steps:

1. Load purchase events from the database
2. Calculate revenue for each customer
3. Sort customers by revenue in descending order
4. Calculate quantile statistics for top 2.5% of customers
5. Export results to new table in database

## Quantile Statistics
-------------------

The project calculates quantile statistics for the top 2.5% of customers by revenue. The quantile statistics are calculated using the following formula:

* `NumberOfCustomers`: the number of customers in the quantile
* `MaxRevenue`: the maximum revenue of the customers in the quantile
* `MinRevenue`: the minimum revenue of the customers in the quantile

The quantile statistics are calculated for each quantile (i.e. 0-2.5%, 2.5-5%, etc.).

## Exporting Results
------------------

The project exports the results of the analysis to a new table in the database. The table has the following columns:

* `CustomerID`: the ID of the customer
* `Email`: the email address of the customer
* `Revenue`: the revenue of the customer

The results are exported using the following steps:

1. Create a new table in the database to store the results
2. Insert the results into the new table

## Logging
---------

The project uses a custom logging system to log events and errors. The logging system uses the following log levels:

* `INFO`: informational messages
* `WARNING`: warning messages
* `ERROR`: error messages

The log messages are written to the console and to a log file.

## Dependencies
------------

The project uses the following dependencies:

* `github.com/go-sql-driver/mysql`: MySQL driver for Go
* `github.com/joho/godotenv`: environment variable loader
* `filippo.io/edwards25519`: Edwards curve implementation (indirect dependency)
