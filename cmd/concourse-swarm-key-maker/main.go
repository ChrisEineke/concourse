package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
)

func main() {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(rand.UintN(256))
	}
	enc := bytes.NewBuffer([]byte{})
	hex.NewEncoder(enc).Write(key)
	fmt.Println("/key/swarm/psk/1.0.0/")
	fmt.Println("/base16/")
	fmt.Println(enc)
}
