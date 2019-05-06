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

## Coding style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Patrons

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) |
|---|
| [Scott Rayapoullé](https://github.com/soulcramer) |

Want to see [your own name here](https://www.patreon.com/eduardurbach)?

## Author

| [![Eduard Urbach on Twitter](https://gravatar.com/avatar/16ed4d41a5f244d1b10de1b791657989?s=70)](https://twitter.com/eduardurbach "Follow @eduardurbach on Twitter") |
|---|
| [Eduard Urbach](https://eduardurbach.com) |

[godoc-image]: https://godoc.org/github.com/blitzprog/home?status.svg
[godoc-url]: https://godoc.org/github.com/blitzprog/home
[report-image]: https://goreportcard.com/badge/github.com/blitzprog/home
[report-url]: https://goreportcard.com/report/github.com/blitzprog/home
[tests-image]: https://cloud.drone.io/api/badges/blitzprog/home/status.svg
[tests-url]: https://cloud.drone.io/blitzprog/home
[coverage-image]: https://codecov.io/gh/blitzprog/home/graph/badge.svg
[coverage-url]: https://codecov.io/gh/blitzprog/home
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
