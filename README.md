# tefas-fon-analysis
a Go program acting as a Telegram bot, fetching financial data for investment funds via the TEFAS API. It provides fund prices and performance metrics for different time periods, interacting with users on the Telegram platform using external packages.
# Telegram Investment Fund Bot

## Introduction

This repository contains a Go program for a Telegram bot that fetches financial data for investment funds using the TEFAS API. The bot provides real-time fund prices and performance metrics for different time periods.

## Features

- Get historical fund prices for the last week
- Fetch general information and performance metrics for a specific investment fund
- Interactive user commands for easy data retrieval
- Telegram bot integration for user-friendly access

## Dependencies

- github.com/PuerkitoBio/goquery
- github.com/go-telegram-bot-api/telegram-bot-api/v5
- github.com/nleeper/goment

## How to Use

1. Obtain a Telegram bot token from the BotFather on Telegram.
2. Clone this repository and navigate to the project directory.
3. Replace the placeholder bot token in the `main` function with your actual token.
4. Build and run the program using `go build` and `./main`.
5. Start interacting with the bot on Telegram by using the `/fon` command followed by a fund code (e.g., `/fon GSP`).

## Example

Sure! Below is the README.md file for the Telegram Investment Fund Bot:

markdown

# Telegram Investment Fund Bot

![Telegram Fund Bot](https://example.com/image.png)

## Introduction

This repository contains a Go program for a Telegram bot that fetches financial data for investment funds using the TEFAS API. The bot provides real-time fund prices and performance metrics for different time periods.

## Features

- Get historical fund prices for the last week
- Fetch general information and performance metrics for a specific investment fund
- Interactive user commands for easy data retrieval
- Telegram bot integration for user-friendly access

## Dependencies

- github.com/PuerkitoBio/goquery
- github.com/go-telegram-bot-api/telegram-bot-api/v5
- github.com/nleeper/goment

## How to Use

1. Obtain a Telegram bot token from the BotFather on Telegram.
2. Clone this repository and navigate to the project directory.
3. Replace the placeholder bot token in the `main` function with your actual token.
4. Build and run the program using `go build` and `./main`.
5. Start interacting with the bot on Telegram by using the `/fon` command followed by a fund code (e.g., `/fon GSP`).

## Example

/fon GSP
<b>GARANTI BIREYSEL SERMAYE PİYASASI YATIRIM FONU</b>

<b>Son 1 aylık getirisi</b> : %1,50
<b>Son 3 aylık getirisi</b> : %3,75
<b>Son 6 aylık getirisi</b> : %4,63
<b>Son 1 yıllık getirisi</b> : %16,70

<b>Son 1 Haftalık veriler</b>

<b>27/07/2023</b> : 2.579000
<b>28/07/2023</b> : 2.602000
<b>29/07/2023</b> : 2.613000
<b>30/07/2023</b> : 2.613000
<b>31/07/2023</b> : 2.635000
<b>01/08/2023</b> : 2.626000



## Contribution

Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
