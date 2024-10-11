package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
	"math/rand"
	"strings"
	"time"
)

func main() {
	N := 1_000_000
	filledN := 1_000_000
	//cf := bf.WithAccuracy(0.001, uint32(filledN))
	cf := bf.WithCapacity(14350730, 10)
	fmt.Println(cf.Info())
	filter, _ := bf.New(cf)
	for i := 0; i < filledN; i++ {
		item := []byte(randString(10 + rand.Intn(10)))
		filter.Add(item)
		after := filter.Exists(item)
		if !after {
			panic("Bloom Filter doesn't work")
		}
	}
	count := 0
	for i := 0; i < N; i++ {
		item := []byte(randString(9))
		if filter.Exists(item) {
			count++
		}
	}
	fmt.Println(fmt.Sprintf("False Positive Count: %v", count))
	fmt.Println(fmt.Sprintf("False Positive Rate: %v", float64(count)/float64(N)))
}

// credit: moorara
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
