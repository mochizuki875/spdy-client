# SPDY Client
This is a simple CLI tool that can be used to send requests to a SPDY server.

## Usage

GET Request
```bash
$ spdy-client get --url="https://<Host>:<Port>/<Path>"
```

POST Request
```bash
$ spdy-client post --url="https://<Host>:<Port>/<Path>"
```

## Build
```bash
$ go build -o spdy-client ./main.go
```