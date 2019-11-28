# {name}

{go:header}

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

### Other HTTP methods

```go
client.Post("https://example.com/user/1").End()
client.Put("https://example.com/user/1").End()
client.Delete("https://example.com/user/1").End()
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

{go:footer}
