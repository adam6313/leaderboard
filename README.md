# Leaderboard

![](https://img.shields.io/badge/golang-1.17.2-blue)
[![Release](https://img.shields.io/badge/release-v1.1.0-blue)](https://github.com/adam6313/leaderboard/releases/tag/v1.1.0)
[![License](https://img.shields.io/github/license/adam6313/leaderboard)](LICENSE)

## Overview
This is a semple leaderboard project.
## Requirements
- The leaderboard has 2 feature
    - One feature is the receive gaming score from a client (with its client ID). 
    - Second to show latest top 10 highest score clients.
- The leaderboard should reset every 10 minutes.
## Tech
- Using Golang Iris framework to build api
- Using redis for database
  - Use the `zset` method to record customer scores to get a sorted data structure.
  - If the requirement is to start timing after the leaderboard has data. the TTL mechanism can be used to complete the reset leaderboard.
  - If the requirement is is global ticker. there is a cronjob designed to complete the reset leaderboard.
- Layered Architecture Design
- Github Actions for CI

## 3rd party lib
- redismock
  - tests use
- mock
  - tests use
- json-iterator
  - decode encode
- cron
  - time tasks
- cobra
  - for CLI interactions
- testify
  - tests use
- fx
  - dependency injection 
- zap
  - logger

## Getting Started
:::info
:bulb: Please check that the computer environment has docker service!
:::
Run the following command.


##### Clone project
  ```shell
git clone https://github.com/adam6313/leaderboard
  ```
##### Create docker network
```shell
docker network create leaderboard
```
##### Run
```shell
docker network create leaderboard
```
After the above steps are completed, you can use the following command to check whether the service starts successfully
```shell
curl --location --request GET '127.0.0.1:8080'
```


## APIs


| URI               |   Method |   Desc   |
| --------          | -------- | -------- |
| /api/v1/score     | POST     | record client score     |
| /api/v1/leaderboard     | GET     | get latest top 10 highest score clients     |
