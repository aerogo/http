# http

[![Reference][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][codecov-image]][codecov-url]
[![License][license-image]][license-url]

Provides HTTP utilities. Currently it offers a fast and easy-to-use HTTP client.

## Installation

```shell
go get github.com/aerogo/http/client
```

## Request

### Basic GET request

```go
response, err := client.Get("https://example.com").End()
```

### Basic POST request

```go
response, err := client.Post("https://example.com").End()
```

### Sending request headers

```go
response, err := client.Get("https://example.com").Header("Accept", "text/html").End()
```

## Response

### Status code

```go
response.StatusCode()
```

### Response body as a string

```go
response.String()
```

### Response body as bytes

```go
response.Bytes()
```

### Deserialize response body into an object (JSON)

```go
response.Unmarshal(&obj)
```

### Response body as a string without unzipping gzip contents

```go
response.RawString()
```

### Response body as bytes without unzipping gzip contents

```go
response.RawBytes()
```

[godoc-image]: https://godoc.org/github.com/aerogo/http?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/http
[report-image]: https://goreportcard.com/badge/github.com/aerogo/http
[report-url]: https://goreportcard.com/report/github.com/aerogo/http
[tests-image]: https://cloud.drone.io/api/badges/aerogo/http/status.svg
[tests-url]: https://cloud.drone.io/aerogo/http
[codecov-image]: https://codecov.io/gh/aerogo/http/graph/badge.svg
[codecov-url]: https://codecov.io/gh/aerogo/http
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/aerogo/http/blob/master/LICENSE
