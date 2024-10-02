package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 2048
)

type Student struct {
	Nama string
	Npm  string
}

type GreetResponse struct {
	Student Student
	Greeter string
}

type HttpRequest struct {
	Method  string
	Uri     string
	Version string
	Host    string
	Accept  string
}

type HttpResponse struct {
	Version       string
	StatusCode    string
	ContentType   string
	ContentLength int
	Data          string
}

func main() {

	remoteTcpAddress, err := net.ResolveTCPAddr(SERVER_TYPE, net.JoinHostPort("127.0.0.1", "3000"))
	if err != nil {
		log.Fatalln(err)
	}
	socket, err := net.DialTCP(SERVER_TYPE, nil, remoteTcpAddress)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("TCP Client Socket Program Example in Go\n")
	fmt.Printf("[%s] Dialling from %s to %s\n", SERVER_TYPE, socket.LocalAddr(), socket.RemoteAddr())

	defer socket.Close()

	fmt.Printf("[%s] Creating receive buffer of size %d\n", SERVER_TYPE, BUFFER_SIZE)

	fmt.Printf("[%s] Input the url:  ", SERVER_TYPE)
	messageUrl, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(messageUrl)

	fmt.Printf("[%s] Input the Type:  ", SERVER_TYPE)
	messageMime, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(strings.TrimSpace(messageMime))

	// msg := strings.TrimSpace(messageUrl) + strings.TrimSpace(messageMime)

	req := HttpRequest{
		Method:  "GET",
		Uri:     strings.TrimSpace(messageUrl),
		Version: "HTTP/1.1",
		Host:    "127.0.0.1",
		Accept:  strings.TrimSpace(messageMime),
	}
	fmt.Println(req)
	response := Fetch(req, socket)
	fmt.Printf("Status Code: %s\nBody: %s\n", response.StatusCode, response.Data)

}

func Fetch(req HttpRequest, connection net.Conn) HttpResponse {

	// This handles the request-making to the server
	requestBytes := RequestEncoder(req)

	_, err := connection.Write(requestBytes)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	receiveBuffer := make([]byte, BUFFER_SIZE)

	receiveLength, err := connection.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}

	receiveMsg := string(receiveBuffer[:receiveLength])

	// 4. Decode response dari byte stream menjadi HttpResponse
	return ResponseDecoder([]byte(receiveMsg))
}

func ResponseDecoder(bytestream []byte) HttpResponse {
	response := string(bytestream)
	fmt.Println(response)
	// Split response menjadi header dan body
	parts := strings.Split(response, "\r\n\r\n")
	headerPart := parts[0]
	var body string
	if len(parts) > 1 {
		body = parts[1]
	}
	fmt.Println(parts)

	// Split headers berdasarkan baris
	headerLines := strings.Split(headerPart, "\r\n")
	statusLine := headerLines[0]
	headers := map[string]string{}

	// Parsing status line (contoh: "HTTP/1.1 200 OK")
	statusParts := strings.Split(statusLine, " ")
	statusCode, _ := strconv.Atoi(statusParts[1])
	// status := strings.Join(statusParts[2:], " ")

	// Parsing headers
	for _, line := range headerLines[1:] {
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = headerParts[1]
		}
	}

	fmt.Println("BODYYY : ", body)
	// Mengembalikan struct HttpResponse
	return HttpResponse{

		Version:       "HTTP/1.1",
		StatusCode:    strconv.Itoa(statusCode),
		ContentType:   "application/json",
		ContentLength: int(len(bytestream)),
		Data:          body,
	}
}

func RequestEncoder(req HttpRequest) []byte {
	requestMessage := fmt.Sprintf("%s %s %s\r\nHost: %s\r\nAccept: %s\r\n\r\n",
		req.Method, req.Uri, req.Version, req.Host, req.Accept)

	return []byte(requestMessage)
}
