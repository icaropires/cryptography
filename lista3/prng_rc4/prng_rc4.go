package main

import (
	"fmt"
	"io/ioutil"
)

var s = make([]byte, 256)
var t = make([]byte, 256)

func initialize() {
	fmt.Println("--------------")

	var k = []byte("lalala")

	for i := range s {
		s[i] = byte(i)
		t[i] = byte(k[i%len(k)])
	}
}

func initialPermutation() {
	j := 0
	for i := range s {
		j = (j + int(s[i]) + int(t[i])) % 256
		s[i], s[j] = s[j], s[i]
	}
}

func streamGeneration() {
	i, j := 0, 0

	data := make([]byte, 256)

	for l := 0; l < 256; l++ {
		i = (i + 1) % 256
		j = (j + int(s[i])) % 256

		s[i], s[j] = s[j], s[i]

		other_t := (int(s[i]) + int(s[j])) % 256
		k := s[other_t]

		data[l] = k
	}

	ioutil.WriteFile("out", data, 0644)
}

func main() {
	initialize()
	initialPermutation()
	streamGeneration()
}
