package main

import (
	"bufio"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "localhost:4221")
	if err != nil {
		return
	}

	conn, err := l.Accept()
	if err != nil {
		return
	}

	var buf []byte
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	byteReader := bufio.NewReader(conn)
	startLine, _, err := byteReader.ReadLine()

	var headers []byte
	for {
		header, _, err := byteReader.ReadLine()
		if err != nil || len(header) == 0 {
			break
		}
		headers = append(headers, header...)
		headers = append(headers, "\r\n"...)
	}

	startLineParts := strings.Split(string(startLine), " ")
	if len(startLineParts) != 3 {
		return
	}

	if startLineParts[0] == "GET" && startLineParts[1] == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			return
		}
		return
	}

	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		return
	}
}
