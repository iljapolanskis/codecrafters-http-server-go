package main

import (
	"bufio"
	"net"
	"strconv"
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

	if startLineParts[0] == "GET" {
		switch {
		case startLineParts[1] == "/":
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			if err != nil {
				return
			}
			return
		case strings.HasPrefix(startLineParts[1], "/echo/"):
			content := strings.TrimPrefix(startLineParts[1], "/echo/")

			//HTTP/1.1 200 OK
			//Content-Type: text/plain
			//Content-Length: 3

			//abc

			var reponse []byte
			addHeaders(&reponse, "HTTP/1.1 200 OK")
			addHeaders(&reponse, "Content-Type: text/plain")
			addHeaders(&reponse, "Content-Length: "+strconv.Itoa(len(content)))
			addContent(&reponse, content)

			_, err = conn.Write(reponse)
			if err != nil {
				return
			}
			return
		}
	}

	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		return
	}
}

func addHeaders(headers *[]byte, header string) {
	*headers = append(*headers, header...)
	*headers = append(*headers, "\r\n"...)
}

func addContent(content *[]byte, data string) {
	*content = append(*content, "\r\n"...)
	*content = append(*content, data...)
}
