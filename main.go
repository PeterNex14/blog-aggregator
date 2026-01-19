package main

import (
	"fmt"

	"github.com/PeterNex14/blog_aggregator/internal/config"
)

func main() {
	data, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := data.SetUser("Peter"); err != nil {
		fmt.Printf("error setting user: %v\n", err)
		return
	}

	data, err = config.Read()

	fmt.Printf("%+v", data)
}