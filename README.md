# DeepTruth
- [About](#about)
- [Requirements](#requirements)
- [How to run](#how-to-run)

## About
Da ya ne eby blyat idi nahui

## Requirements
### To run you have to install this software
* [go](https://go.dev/)
* [migrate](https://github.com/golang-migrate/migrate)
* [sqlc](https://sqlc.dev/)
* [make](https://www.gnu.org/software/make/)

## How to run
1.
    ```bash
    make compose-up
    ```
1. Wait about 30 seconds while db starts
1.
    ```bash
    make tidy       # Only on first run
    make migrate-up # Only on first run
    make run
    ```
