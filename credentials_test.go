package docker_credentials

import "fmt"

func ExampleOf() {
	fmt.Println(Of("index.docker.io"))
	// Output: username password <nil>
}
