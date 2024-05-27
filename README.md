# Golang setup

Welcome to the Go Project Template! This project provides a robust starting point for building Go applications with essential  routes. It is designed to help you quickly get up and running with writing business logic for your specific use case.

<img align="right" width="180px" src="https://raw.githubusercontent.com/kaynetik/dotfiles/master/svg-resources/ashleymcnamara%40GOPHER_STAR_WARS.png">

## Table of Contents

<!-- TOC -->
* [Features](#features)
  * [Identity Provider and Task Management API](#identity-provider-and-task-management-api)
* [Local Build](#local-build)
  * [Requirements](#requirements)
  * [Initial Setup - Dockerized](#initial-setup---dockerized)
  * [Initial Setup - Native](#initial-setup---native)
* [Additional Tools Used](#additional-tools-used)
* [Environment Variables](#environment-variables)
* [Database Migration](#database-migration)
<!-- TOC -->

## Features

### Identity Provider and Task Management API
The Task Management API grants admins and users the ability to log in and manage tasks on the platform. 

By relying on our IDP for authentication, the API securely permits any user with a valid token to access its resources.

## Local Build

### Requirements
  - Go >1.21 (see `go.mod`)
  - MySQL 8 (ideally, latest stable version)
  - Redis >6

There are a few different ways to run this project. Locally it can be run using Docker Compose, or it can be run by
using the relevant native binaries.

After cloning the repo, the basic setup commands are:

```bash
cd golang-setup
cp .env.example .env
```

### Initial Setup - Dockerized

To run this project locally using Docker Compose, you will need to have Docker and Docker Compose installed. You can
find instructions for installing Docker [here](https://docs.docker.com/install/), and instructions for installing Docker
Compose [here](https://docs.docker.com/compose/install/).

Once every requirement is satisfied, you can run the project using the following command:

```bash
docker-compose build && docker-compose up
```

### Initial Setup - Native

To run this project locally using the native binaries, you must have Go installed. You can find instructions for
installing Go [here](https://golang.org/doc/install). Once you have Go installed, you will need to install the
dependencies for the project. You can do this by running the following command:

```bash
go get -u ./...
go run main.go
```

## Additional Tools Used

On top of everything mentioned so far, for smooth development, we strongly advise that you install the following tools:

1. `golangci-lint` - Linter for Go. You can find instructions for installing it [here](https://golangci-lint.run). This
   is used to statically check / lint the code and enforce a certain level of quality.
2. `gofumpt` - Formatter for Go. You can find instructions for installing it [here](https://github.com/mvdan/gofumpt).
   This is used to format the code consistently.

## Environment Variables
  See [.env.example](.env.example).

## Database Migration
 Database files are automatically executed when service starts. For more details how to use migration please read a [README.md](/tools/storage/mysql/README.md) file.

