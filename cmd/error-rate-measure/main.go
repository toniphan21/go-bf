package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
	"github.com/toniphan21/go-bf/internal"
	"math/rand"
)

func main() {
	N := 1_000_000
	filledN := 1_000_000
	//cf := bf.WithAccuracy(0.001, uint32(filledN))
	cf := bf.WithCapacity(14350730, 10)
	fmt.Println(cf.Info())
	filter, _ := bf.New(cf)
	for i := 0; i < filledN; i++ {
		item := []byte(internal.RandString(10 + rand.Intn(10)))
		filter.Add(item)
		after := filter.Exists(item)
		if !after {
			panic("Bloom Filter doesn't work")
		}
	}
	count := 0
	for i := 0; i < N; i++ {
		item := []byte(internal.RandString(9))
		if filter.Exists(item) {
			count++
		}
	}
	fmt.Println(fmt.Sprintf("False Positive Count: %v", count))
	fmt.Println(fmt.Sprintf("False Positive Rate: %v", float64(count)/float64(N)))
}
