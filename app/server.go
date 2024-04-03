package main

import "net"

func main() {
	l, err := net.Listen("tcp", ":4221")
	if err != nil {
		panic(err)
	}

	_, err = l.Accept()
	if err != nil {
		panic(err)
	}
}
