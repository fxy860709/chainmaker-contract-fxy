package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// 生成 10 个 [0, 128) 范围的真随机数。
	//for i := 0; i < 10; i++ {
	//	result, _ := rand.Int(rand.Reader, big.NewInt(10000000000))
	//	fmt.Println(result)
	//}

	result, _ := rand.Int(rand.Reader, big.NewInt(10000000000))
	fmt.Println(result)
}
