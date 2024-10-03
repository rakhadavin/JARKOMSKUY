package main

import (
	"bufio"
	"encoding/json"
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

	var parsedJson GreetResponse
	var messageMime string
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
	var host string
	var req HttpRequest

	fmt.Printf("[%s] Input the Type:  ", SERVER_TYPE)
	messageMime, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	if len(strings.Split(messageUrl, "/")) < 4 {
		req = HttpRequest{
			Method:  "GET",
			Uri:     strings.TrimSpace(messageUrl),
			Version: "HTTP/1.1",
			Host:    "",
			Accept:  strings.TrimSpace(messageMime),
		}
	} else {

		fmt.Println(strings.Split(messageUrl, " "))
		ipRaw := strings.Split(messageUrl, "/")[2]
		host = ipRaw[:strings.Index(ipRaw, ":")]
		req = HttpRequest{
			Method:  "GET",
			Uri:     strings.TrimSpace(messageUrl),
			Version: "HTTP/1.1",
			Host:    host,
			Accept:  strings.TrimSpace(messageMime),
		}

	}

	response := Fetch(req, socket)
	fmt.Printf("Status Code: %s\nBody: %s\n", response.StatusCode, response.Data)
	if (strings.Contains(response.ContentType, "application/json")) || (strings.Contains(response.ContentType, "application/xml")) {
		var _ = json.Unmarshal([]byte(response.Data), &parsedJson)

		fmt.Println("Parsed: ", parsedJson)
	}

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
	var body string

	response := string(bytestream)
	fragments := strings.Split(response, "\r\n\r\n")
	headerPart := fragments[0]
	if len(fragments) > 1 {
		body = fragments[1]
	}
	var contentType = strings.Split(fragments[0], " ")[2]

	headerLines := strings.Split(headerPart, "\r\n")
	statusLine := headerLines[0]
	headers := map[string]string{}

	statusFr := strings.Split(statusLine, " ")
	statusCode, _ := strconv.Atoi(statusFr[1])

	for _, line := range headerLines[1:] {
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = headerParts[1]
		}
	}

	return HttpResponse{

		Version:       "HTTP/1.1",
		StatusCode:    strconv.Itoa(statusCode),
		ContentType:   contentType,
		ContentLength: int(len(bytestream)),
		Data:          body,
	}
}

func RequestEncoder(req HttpRequest) []byte {
	requestMessage := fmt.Sprintf("%s %s %s\r\nHost: %s\r\nAccept: %s\r\n\r\n",
		req.Method, req.Uri, req.Version, req.Host, req.Accept)

	return []byte(requestMessage)
}
