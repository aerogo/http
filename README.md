# http

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Patreon][patreon-image]][patreon-url]

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

### Response body as bytes

```go
response.Bytes()
```

### Response body as a string

```go
response.String()
```

### Status code

```go
response.StatusCode()
```

### Deserialize response body into an object (JSON)

```go
response.Unmarshal(&obj)
```

### Response body as bytes (without unzipping gzip contents)

```go
response.Raw()
```

### Response body as a string (without unzipping gzip contents)

```go
response.RawString()
```

### Response length (without unzipping gzip contents)

```go
response.RawLength()
```

## Style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Sponsors

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) | [![Eduard Urbach](https://avatars2.githubusercontent.com/u/438936?s=70&v=4)](https://twitter.com/eduardurbach) |
| --- | --- |
| [Scott Rayapoullé](https://github.com/soulcramer) | [Eduard Urbach](https://eduardurbach.com) |

Want to see [your own name here?](https://www.patreon.com/eduardurbach)

[godoc-image]: https://godoc.org/github.com/aerogo/http?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/http
[report-image]: https://goreportcard.com/badge/github.com/aerogo/http
[report-url]: https://goreportcard.com/report/github.com/aerogo/http
[tests-image]: https://cloud.drone.io/api/badges/aerogo/http/status.svg
[tests-url]: https://cloud.drone.io/aerogo/http
[coverage-image]: https://codecov.io/gh/aerogo/http/graph/badge.svg
[coverage-url]: https://codecov.io/gh/aerogo/http
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
