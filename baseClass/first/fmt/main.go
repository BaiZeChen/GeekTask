package main

import (
	"encoding/hex"
	"fmt"
)

func main() {

	fmt.Printf("%.2f\n", 3.124234)

	// []byte -> 16进制字符串
	strByte := []byte{'k', 'q', 'd'}
	encodedStr := hex.EncodeToString(strByte)
	fmt.Println(encodedStr)

}
