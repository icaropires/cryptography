package main

import (
	//"encoding/base64"

	"encoding/binary"
	"fmt"
	"io/ioutil"
)

func main() {
	key := make([]byte, 32)
	out := stream_salsa20(10000, []byte{1, 1, 4, 2, 4, 2, 4, 2}, key)
	a := binary.BigEndian.Uint64(out)
	fmt.Println(out, a)

	ioutil.WriteFile("out", out, 0644)
}
