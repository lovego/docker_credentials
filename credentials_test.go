package docker_credentials_test

import (
	"fmt"

	"github.com/lovego/docker_credentials"
)

func ExampleOf() {
	fmt.Println(docker_credentials.Of("index.docker.io"))
	// Output: username password <nil>
}

func ExampleNew() {
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
