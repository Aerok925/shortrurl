package main

import (
	"fmt"
	"hash/crc32"
	"log"
)

func main() {
	str := "Hello, world!"
	hash := crc32.ChecksumIEEE([]byte(str))

	log.Println(fmt.Sprintf("%x", hash))
}
