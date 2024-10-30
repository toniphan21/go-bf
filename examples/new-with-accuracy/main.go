package main

import "github.com/toniphan21/go-bf"

func main() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	filter, err := bf.New(bf.WithAccuracy(errorRate, numberOfItems))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}
