afi-file-downloader
===================

Utility for downloading file from HTTP URL and stream into tcp connection.

# Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Build](#build)
    - [Console](#build-with-console)
    - [Docker](#build-with-docker)
    - [Werf](#build-with-werf)
- [Communication](#communication)
- [Usage](#usage)
    - [Console](#run-with-console)
    - [Docker](#run-with-docker)
    - [Werf](#run-with-werf)


# Requirements

- GoLang >= 1.15.3
- Docker >= 19.03.13
- Werf >= v1.1.21+fix32


# Installation

```bash
# with go get
$ go get github.com/morozovcookie/afifiledownloader/cmd/afi-file-downloader/...

# with git
$ git clone https://github.com/morozovcookie/afi-file-downloader.git
```


# Build

## Build With Console

```bash
# build binary file
$ make go-build
```

## Build With Docker

```bash
# build docker image
$ make docker-build

# publish
$ make docker-publish DOCKER_REPOSITORY=sample.registry.com
```

## Build With Werf

```bash
# build docker image
$ make werf-build

# publish docker image
$ make werf-publish DOCKER_REPOSITORY=sample.registry.com
```


# Communication

Communication between user and application based on JSON format message which pass through stdin. Response you will get through stdout or stderr if all goes down.

## Request

|Field                  |Type     |Description                                  |Mandatory|Default|
|-----------------------|:-------:|---------------------------------------------|:-------:|:-----:|
|ignore-ssl-certificates|*Boolean*|Allow insecure connection                    |N        |False  |
|follow-redirects       |*Boolean*|Follow redirects                             |N        |False  |
|max-redirects          |*Long*   |Limit redirects                              |N        |5      |
|url                    |*String* |HTTP URL for downloading                     |Y        |       |
|output                 |*String* |TCP *host:port* for streaming downloaded data|N        |       |
|timeout                |*String* |Request timeout                              |N        |1s     |

## Response

|Field           |Type          |Description                 |
|----------------|:------------:|----------------------------|
|success         |*Boolean*     |Request result              |
|http-code       |*Number*      |HTTP response status code   |
|content-length  |*Long*        |HTTP response content length|
|content-type    |*String*      |HTTP response content type  |
|error-message   |*String*      |Error message               |
|redirects       |*List<String>*|List of redirects           |


# Usage

## Run With Console

```bash
$ echo "<json>" | make go-run
```

## Run With Docker

```bash
$ echo "<json>" | make docker-run
```

## Run With Werf

```bash
$ echo "<json>" | make werf-run
```
