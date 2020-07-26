package main

import (
	"fmt"
	"grpcourse/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
