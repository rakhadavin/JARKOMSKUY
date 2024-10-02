package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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

	var student = Student{
		Nama: STUDENT_NAME,
		Npm:  STUDENT_NPM,
	}

	greeter := GreetResponse{
		Student: student,
		Greeter: student.Nama,
	}
	fmt.Printf(greeter.Greeter)

	if strings.Contains(req.Accept, "application/json") {
		contentType = "application/json"
		var dataJson, err = json.Marshal(greeter)
		if err != nil {
			log.Fatalln("Error parsing to JSON")
		}
		data = string(dataJson)
		contentLength = len(data) // Contoh nilai
	} else if strings.Contains(req.Accept, "text/html") {
		contentType = "text/html"
		data = fmt.Sprintf("<html><body><h1>Halo, dunia! aku %s</h1></body></html>", STUDENT_NAME)
		contentLength = len(data)
	} else {
		contentType = "text"
		// data, _ := json.Marshal(greeter)
		contentLength = len(data)
	}

	response := HttpResponse{
		Version:       req.Version,
		StatusCode:    "200",
		ContentType:   contentType,
		ContentLength: contentLength,
		Data:          data,
	}
	fmt.Printf("Responding to clients: %+v\n", response)
	// Kembali ke HttpResponse struct
	return response

}

func RequestDecoder(bytestream []byte) HttpRequest {
	// Put the decoding program for HTTP Request Packet here
	msgSplit := strings.Split(string(bytestream), "\r\n") //[GET djkfhsjdjf HTTP/1.1 Host: 127.0.0.1 Accept: dfhksjdfh ] len:3

	var host, accept string
	host = strings.TrimSpace((strings.TrimPrefix(msgSplit[1], "Host:")))
	accept = strings.TrimSpace(msgSplit[2])

	fmt.Println("ACCEPT:: ", accept)
	// Membuat struct HttpRequest
	return HttpRequest{
		Method:  strings.Split(string(bytestream), "\r\n")[0],
		Uri:     strings.Split(string(bytestream), "\r\n")[1],
		Version: "HTTP/1.1",
		Host:    host,
		Accept:  accept,
	}
}

func ResponseEncoder(res HttpResponse) []byte {
	// Put the encoding program for HTTP Response Struct here
	fmt.Println("TOLOL")
	fmt.Println(res.Data)
	responseString := fmt.Sprintf("%s %d\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		res.Version, 200, res.ContentType, res.ContentLength, res.Data)
	return []byte(responseString)
}
