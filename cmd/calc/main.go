package main

import (
	"fmt"
	"math"
)

func main() {
	errorRate := 0.02
	numberOfItems := float64(10000000)
	log2 := math.Abs(math.Log2(errorRate))
	k := log2
	bitPerItem := 1.44 * log2
	capacityInBits := numberOfItems * bitPerItem
	fmt.Printf("log2 = %v\n", log2)
	fmt.Printf("k = %v\n", k)
	fmt.Printf("m/n = %v\n", bitPerItem)
	fmt.Printf("n = %v\n", numberOfItems)
	fmt.Printf("capacity (bits) = %v\n", capacityInBits)
	fmt.Printf("capacity (bytes) = %v\n", capacityInBits/8)
	fmt.Println("------------------------------")
	fmt.Println("Normalized data")
	fmt.Println("------------------------------")
	nK := int(math.Ceil(k))
	nCapacity := int(math.Ceil(capacityInBits))
	fmt.Printf("k = %v\n", nK)
	fmt.Printf("capacity = %v\n", nCapacity)
	nLog2K := int(math.Ceil(math.Log2(float64(nCapacity))))
	fmt.Printf("nLog2K = %v\n", nLog2K)
	fmt.Println("------------------------------")
	fmt.Println("Recalculate error rate by number of items")
	fmt.Println("------------------------------")
	rate := 1000
	for i := 1000; i <= 10000000; i += rate {
		epsilon := math.Pow(1-math.Pow(math.E, (0-float64(nK)*float64(i))/float64(nCapacity)), float64(nK))
		fmt.Printf("n = %v -> epsilon = %v\n", i, epsilon)
		if i > 10000 {
			rate = 5000
		}
		if i > 100000 {
			rate = 100000
		}
		if i > 1000000 {
			rate = 1000000
		}
	}
}
