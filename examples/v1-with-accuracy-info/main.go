package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
)

func main() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	config := bf.WithAccuracy(errorRate, numberOfItems)
	fmt.Println(config.Info())
}
