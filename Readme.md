# docker\_credentials
Get stored username and password from docker.

[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/docker_credentials?1)](https://goreportcard.com/report/github.com/lovego/docker_credentials)
[![GoDoc](https://godoc.org/github.com/lovego/docker_credentials?status.svg)](https://godoc.org/github.com/lovego/docker_credentials)

## usage
```go
func ExampleOf() {
	fmt.Println(docker_credentials.Of("index.docker.io"))
	// Output: username password <nil>
}

func ExampleCredentials_New() {
	c, err := docker_credentials.New([]byte(`{
    "auths": {
      "https://index.docker.io/v1/": { "auth": "dXNlcm5hbWU6cGFzc3dvcmQ=" }
    }
  }`))
	if err == nil {
		fmt.Println(c.Of("index.docker.io"))
	} else {
		fmt.Println(err)
	}
	// Output: username password <nil>
}
```
