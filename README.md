# README

Repo - https://github.com/puneeth8994/scraper

## Solution Overview
The Scraper application provides a utility to scrape content from web (http) URLs and return the value.

The application expects URLs returning either json or html responses.
Any Non-2xx responses are skipped in the final result.

The application scrapes the title field in the web urls and returns the file response as a map.


## Running the app

Once the project setup is done, you can run the service for default URLs in the code or provide a --file flag and provide files as args

example - `go run cmd/main.go --file="./urls.txt"`

### Running with default URLs

To run the project with only default URLs, you can execute the code as -
`go run cmd/main.go`

## Project Setup

### 1. Clone the repository

- `git clone https://github.com/puneeth8994/scraper.git`
- `cd scraper`

### 2. Install Dependencies

- Make sure Go is installed on your system (Go version 1.16+ is recommended). You can check if Go is installed by running:

```
go version
```

- Once the go version is verified. Install any external dependencies using -

```
go mod tidy
```

## Testing:
The project includes unit tests for core functionalities and utils.
You can find test files with `_test.go` suffix to file names in respective directories

## Run the Tests:

You can execute the tests using the Go testing tool. In the project root directory (or the directory containing the test files), run:

```
go test ./... -v
```

This command will find and run all tests in the current package and sub-packages
You can obtain verbose output using -v flag for better visibility of test results.

##### Check Test Results:
The output will indicate whether each test passed or failed. Any failures will be accompanied by error messages.
