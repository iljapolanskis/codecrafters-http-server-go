package main

import "net"

func main() {
	l, err := net.Listen("tcp", ":4221")
	if err != nil {
		return
	}

	conn, err := l.Accept()
	if err != nil {
		return
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		return
	}
}
