package main

import "github.com/toniphan21/go-bf"

type YourConfig struct {
}

func (y *YourConfig) Info() string {
	return "info about your config"
}

func (y *YourConfig) NumberOfHashFunctions() byte {
	return 5
}

func (y *YourConfig) StorageCapacity() uint32 {
	return 1_000_000
}

func main() {
	config := &YourConfig{}
	filter, err := bf.New(config)
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}
