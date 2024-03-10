package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"pow_server/pkg"
)

func main() {
	client := pkg.NewWordOfWisdomClient(os.Args[1], zap.NewExample())
	q, err := client.GetWordOfWisdom()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(q)
	}
}
