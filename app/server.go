package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

type Response struct {
	version       string
	statusCode    int
	statusMessage string
	headers       map[string]string
	content       string
}

func NewResponse() *Response {
	return &Response{
		version:    "HTTP/1.1",
		statusCode: 200,
		headers:    make(map[string]string),
	}
}

func (r *Response) AddHeader(key, value string) {
	r.headers[key] = value
}

func (r *Response) SetContent(content string, contentType string) {
	if len(content) == 0 {
		return
	}

	r.AddHeader("Content-Length", strconv.Itoa(len(content)))
	r.AddHeader("Content-Type", contentType)
	r.content = content
}

func (r *Response) SetStatus(statusCode int) {
	r.statusCode = statusCode

	// TODO: Add more status codes
	switch statusCode {
	case 200:
		r.statusMessage = "OK"
	case 404:
		r.statusMessage = "Not Found"
	}
}

func (r *Response) Write(conn net.Conn) {
	sb := strings.Builder{}
	sb.WriteString(r.version)
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(r.statusCode))
	sb.WriteString(" ")
	sb.WriteString(r.statusMessage)
	sb.WriteString("\r\n")

	for key, value := range r.headers {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\r\n")
	}

	if _, ok := r.headers["Content-Length"]; ok {
		sb.WriteString("\r\n")
		sb.WriteString(r.content)
	}

	sb.WriteString("\r\n")

	conn.Write([]byte(sb.String()))

	conn.Close()
}

type Request struct {
	version string
	method  string
	path    string
	headers map[string]string
	content string
}

func NewRequest(conn net.Conn) *Request {
	byteReader := bufio.NewReader(conn)

	startLine, _, _ := byteReader.ReadLine()

	headers := make(map[string]string)
	for {
		header, _, err := byteReader.ReadLine()
		if err != nil || len(header) == 0 {
			break
		}

		headerParts := strings.Split(string(header), ": ")
		if len(headerParts) != 2 {
			continue
		}

		headers[headerParts[0]] = headerParts[1]
	}

	startLineParts := strings.Split(string(startLine), " ")
	if len(startLineParts) != 3 {
		return nil
	}

	// If headers have Content-Length, read the content
	var content string
	if contentLength, ok := headers["Content-Length"]; ok {
		contentLengthInt, err := strconv.Atoi(contentLength)
		if err != nil {
			return nil
		}

		contentBytes := make([]byte, contentLengthInt)
		_, err = byteReader.Read(contentBytes)
		if err != nil {
			return nil
		}

		content = string(contentBytes)
	}

	return &Request{
		method:  startLineParts[0],
		path:    startLineParts[1],
		version: startLineParts[2],
		headers: headers,
		content: content,
	}
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Path() string {
	return r.path
}

func (r *Request) Header(key string) string {
	return r.headers[key]
}

func (r *Request) Content() string {
	return r.content
}

func main() {
	l, err := net.Listen("tcp", "localhost:4221")
	if err != nil {
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	request := NewRequest(conn)

	if request.Method() == "GET" {
		switch {
		case request.Path() == "/":

			response := NewResponse()
			response.SetStatus(200)
			response.Write(conn)

			return
		case strings.HasPrefix(request.path, "/echo/"):
			content := strings.TrimPrefix(request.Path(), "/echo/")

			response := NewResponse()
			response.SetStatus(200)
			response.SetContent(content, "text/plain")
			response.Write(conn)

			return
		case request.Path() == "/user-agent":
			userAgent := request.Header("User-Agent")

			response := NewResponse()
			response.SetStatus(200)
			response.SetContent(userAgent, "text/plain")
			response.Write(conn)

			return
		}
	}

	response := NewResponse()
	response.SetStatus(404)
	response.Write(conn)
}
