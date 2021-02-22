# MrAndreID / Go RabbitMQ

[![Go Reference](https://pkg.go.dev/badge/github.com/MrAndreID/gorabbitmq.svg)](https://pkg.go.dev/github.com/MrAndreID/gorabbitmq) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The `MrAndreID/GoRabbitMQ` package is a collection of functions in the go language for RabbitMQ.

---

## Table of Contents

* [Install](#install)
* [Usage](#usage)
* [Versioning](#versioning)
* [Authors](#authors)
* [License](#license)
* [Official Documentation for Go Language](#official-documentation-for-go-language)
* [More](#more)

---

## Install

To use The `MrAndreID/GoRabbitMQ` package, you must follow the steps below:

```sh
go get -u github.com/MrAndreID/gorabbitmq
```

## Usage

### Client

```go
import "github.com/MrAndreID/gohelpers"

err := gorabbitmq.Client("Andrea Adam", gorabbitmq.Connection{"127.0.0.1", "5672", "account", "account", "account"}, gorabbitmq.QueueSetting{"account", true, false, false, false, nil}, gorabbitmq.OtherSetting{"account", "60000", false, false, 18})

if err != nil {
    gohelpers.HandleResponse(response, 400, "looks like something went wrong", err)

    return
}
```

### Server

```go
import "services/routes"

gorabbitmq.Server(gorabbitmq.Connection{"127.0.0.1", "5672", "account", "account", "account"}, gorabbitmq.QueueSetting{"account", true, false, false, false, nil}, gorabbitmq.ConsumeSetting{"", false, false, false, false, nil}, gorabbitmq.OtherSetting{"account", "60000", false, false, 18}, routes.HandleRequest)
```

### RPC Client

```go
import "github.com/MrAndreID/gohelpers"

result, err := gorabbitmq.RPCClient("Andrea Adam", gorabbitmq.Connection{"127.0.0.1", "5672", "account", "account", "account"}, gorabbitmq.QueueSetting{"account", true, false, false, false, nil}, gorabbitmq.ConsumeSetting{"", true, false, false, false, nil}, gorabbitmq.OtherSetting{"account", "60000", false, false, 18})

if err != nil {
    gohelpers.HandleResponse(response, 400, "looks like something went wrong", err)

    return
}
```

### RPC Server

```go
import "services/routes"

gorabbitmq.RPCServer(gorabbitmq.Connection{"127.0.0.1", "5672", "account", "account", "account"}, gorabbitmq.QueueSetting{"account", true, false, false, false, nil}, gorabbitmq.QosSetting{1, 0, false}, gorabbitmq.ConsumeSetting{"", false, false, false, false, nil}, gorabbitmq.OtherSetting{"account", "60000", false, false, 18}, routes.HandleRequest)
```

## Versioning

I use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository. 

## Authors

**Andrea Adam** - [MrAndreID](https://github.com/MrAndreID/)

## License

MIT licensed. See the LICENSE file for details.

## Official Documentation for Go Language

Documentation for Go Language can be found on the [Go Language website](https://golang.org/doc/).

## More

Documentation can be found [on https://go.dev/](https://pkg.go.dev/github.com/MrAndreID/gorabbitmq).
