package main

import "github.com/toniphan21/go-bf"

func main() {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	filter := bf.Must(bf.WithCapacity(capacityInBits, numberOfHashFunctions))

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}
