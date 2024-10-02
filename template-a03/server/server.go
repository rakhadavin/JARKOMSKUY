package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

const (
	SERVER_HOST  = "127.0.0.1"
	SERVER_PORT  = "3000"
	SERVER_TYPE  = "tcp"
	BUFFER_SIZE  = 2048
	STUDENT_NAME = "Pin"
	STUDENT_NPM  = "2206082650"
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

	listenAddress, err :=
		net.ResolveTCPAddr(SERVER_TYPE,
			net.JoinHostPort(SERVER_HOST,
				SERVER_PORT))
	if err != nil {
		log.Fatalln(err)
	}

	socket, err :=
		net.ListenTCP(SERVER_TYPE,
			listenAddress)

	if err != nil {
		log.Fatalln(err)
	}

	defer socket.Close()

	fmt.Println("Server ready to listen !")

	for {

		connection, err := socket.AcceptTCP()
		go HandleConnection(connection)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func HandleConnection(connection net.Conn) {
	receiveBuffer := make([]byte, BUFFER_SIZE)
	defer connection.Close()
	receiveLength, err := connection.Read(receiveBuffer)
	if err != nil {
		log.Fatalln("[CAN'T READ BUFFER] : ", err)
	}
	message := string(receiveBuffer[:receiveLength])

	req := RequestDecoder([]byte(message))

	response := HandleRequest(req)

	responseString := ResponseEncoder(response)
	_, err = connection.Write((responseString))
	if err != nil {
		log.Fatalln(err)
	}

}

func HandleRequest(req HttpRequest) HttpResponse {

	var contentType string
	var contentLength int
	var data string
	var paramValue string
	var validURI string
	var greeterName string

	var student = Student{
		Nama: STUDENT_NAME,
		Npm:  STUDENT_NPM,
	}
	greeterName = "Pin"
	validURI = "HTTPS://greet"
	validURI = req.Uri
	fmt.Println("ini urinya awalnya : ", validURI)

	// jika URI mengandung param -- set nama greeter dengan nama pada param
	validURI = req.Uri
	parsedURI, err := url.Parse(validURI)
	if err != nil {
		log.Fatalln("Error : Invalid URL")
	}
	paramValue = parsedURI.Query().Get("name")
	if paramValue != "" {

		fmt.Println("MASOK")
		greeterName = paramValue
	}
	greeter := GreetResponse{
		Student: student,
		Greeter: greeterName,
	}

	if strings.Contains(req.Accept, "application/json") {
		contentType = "application/json"
		var dataJson, err = json.Marshal(greeter)
		if err != nil {
			log.Fatalln("Error parsing to JSON")
		}
		data = string(dataJson)
		contentLength = len(data)
	} else if strings.Contains(req.Accept, "application/xml") {
		contentType = "application/json"
		var dataJson, err = xml.Marshal(greeter)
		if err != nil {
			log.Fatalln("Error parsing to JSON")
		}
		data = string(dataJson)
		contentLength = len(data)
	} else if strings.Contains(req.Accept, "text/html") {
		contentType = "text/html"
		data = fmt.Sprintf("<html><body><h1>Halo, dunia! aku %s</h1></body></html>", STUDENT_NAME)
		contentLength = len(data)
	} else {
		contentType = "text"
		data, _ := json.Marshal(greeter)
		contentLength = len(data)
	}
	fmt.Println("GRETER NAMA : ", greeter.Student.Nama)

	response := HttpResponse{
		Version:       req.Version,
		StatusCode:    "200",
		ContentType:   contentType,
		ContentLength: contentLength,
		Data:          data,
	}
	fmt.Printf("Responding to clients: %+v\n", response)
	return response

}

func RequestDecoder(bytestream []byte) HttpRequest {
	// Put the decoding program for HTTP Request Packet here
	msgSplit := strings.Split(string(bytestream), "\r\n") //[GET djkfhsjdjf HTTP/1.1 Host: 127.0.0.1 Accept: dfhksjdfh ] len:3

	var host, accept string
	host = strings.TrimSpace((strings.TrimPrefix(msgSplit[1], "Host:")))
	accept = strings.TrimSpace(msgSplit[2])
	fmt.Println("lengkap : ", msgSplit)
	fmt.Println("INI : ", strings.Split(msgSplit[0], " ")[0])
	// Membuat struct HttpRequest
	return HttpRequest{
		Method:  strings.Split(msgSplit[0], " ")[0],
		Uri:     strings.Split(msgSplit[0], " ")[1],
		Version: "HTTP/1.1",
		Host:    host,
		Accept:  accept,
	}
}

func ResponseEncoder(res HttpResponse) []byte {
	// Put the encoding program for HTTP Response Struct here
	fmt.Println(res.Data)
	responseString := fmt.Sprintf("%s %d\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		res.Version, 200, res.ContentType, res.ContentLength, res.Data)
	return []byte(responseString)
}
