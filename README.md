# http

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