# Echo Framework with Mongo Driver

## Overview
My Boilerplate for Rest API using [echo](https://echo.labstack.com/) and [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

## Setup
- Copy config and modify
```
$ cp config/config.json.example config/config.json
```
- Install and update go mod
```
$ go get -u
```
- Seed user
```
$ go run main.go seed
```
- Generate documentation
```
$ swag init
```
- Run
```
$ go run main.go
```

## Notes
- Roles
    - User can only list or get by id
    - Administrator can add user