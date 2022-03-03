# mongo-go-prometheus

<p>
  <img src="https://img.shields.io/github/workflow/status/globocom/mongo-go-prometheus/Go?style=flat-square">
  <a href="https://github.com/globocom/mongo-go-prometheus/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/globocom/mongo-go-prometheus?color=blue&style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/globocom/mongo-go-prometheus?style=flat-square">
  <a href="https://pkg.go.dev/github.com/globocom/mongo-go-prometheus">
    <img src="https://img.shields.io/badge/Go-reference-blue?style=flat-square">
  </a>
</p>

Monitors that export Prometheus metrics for the MongoDB Go driver

## Installation

	go get github.com/globocom/mongo-go-prometheus

## Usage

```golang
package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/globocom/mongo-go-prometheus"
)

func main() {
	monitor := mongoprom.NewCommandMonitor(
		mongoprom.WithInstanceName("database"),
		mongoprom.WithNamespace("my_namespace"),
		mongoprom.WithDurationBuckets([]float64{.001, .005, .01}),
	)

	poolMonitor := mongoprom.NewPoolMonitor(
		mongoprom.PoolWithInstanceName("database"),
		mongoprom.PoolWithNamespace("my_namespace"),
	)
	
	opts := options.Client().
		ApplyURI("mongodb://localhost:27019").
		SetMonitor(monitor).SetPoolMonitor(poolMonitor)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// run MongoDB commands...
}
```

## Exported metrics

The command monitor exports the following metrics:

- Commands:
    - Histogram of commands: `mongo_commands{instance="db", command="insert"}`
    - Counter of errors: `mongo_command_errors{instance="db", command="update"}`
- Pool:
  - Max number of connections allowed in Connection Pool: `mongodb_connection_pool_max{instance="db"}`
  - Min number of connections allowed in Connection Pool: `mongodb_connection_pool_min{instance="db"}`
  - Actual connections in usage: `mongodb_connection_pool_usage{instance="db"}` 

## API stability

The API is unstable at this point and it might change before `v1.0.0` is released.
