# go-sensor-collector

A simple golang backend that act as a gRPC server and REST API, that collect sensor data from sensor gRPC client then serving it as REST API.

## Getting Started

How to setup in local machine

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.19 or higher)
- [Protobuf Compiler](https://grpc.io/docs/protoc-installation/)

### Installation

1. Clone the repository:

```sh
git clone https://github.com/muhammadtaufan/go-sensor-collector.git
```

2. go to the directory:

```sh
cd go-sensor-collector
```

3. setup project:

```sh
make setup
```

4. db migration:

```sh
make migrate-up
```

5. add dummy user:

```sh
make create-user
```

6. create a new migration

```sh
make migrate-create name=create_new_table
```

7. build and run the service:

```sh
make run
```
