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
	var npmParam string

	var student = Student{
		Nama: STUDENT_NAME,
		Npm:  STUDENT_NPM,
	}
	greeterName = "Pin"
	validURI = "http://127.0.0.1:3000/greet/2206082650" //set default
	// jika NPM pada URI tidak sama dengan NPM saya

	validURI = req.Uri
	parsedURI, err := url.Parse(validURI)
	fmt.Println(strings.Split(validURI, "/"))
	//jika url tidak valid
	if len(strings.Split(validURI, "/")) < 4 && !strings.EqualFold(validURI, "http://127.0.0.1:3000/") {
		data = ""
		return HttpResponse{
			Version:       req.Version,
			StatusCode:    "404",
			ContentType:   contentType,
			ContentLength: contentLength,
			Data:          data,
		}
	} else {
		splittedURI := strings.Split(validURI, "/")
		npmWithParams := splittedURI[len(splittedURI)-1]
		npmParam = strings.Split(npmWithParams, "?")[0] //ambil param

		// Output hanya angkanya

	}

	if (!strings.HasPrefix(validURI, "http://") && !strings.HasPrefix(validURI, "https://")) || err != nil || (!strings.EqualFold(npmParam, STUDENT_NPM) && npmParam != "") {

		// && strings.EqualFold(parsedURI, "http://127.0.0.1:3000/"
		fmt.Println(strings.EqualFold(npmParam, STUDENT_NPM))
		contentType = ""

		data = ""
		contentLength = len(data)

		return HttpResponse{
			Version:       req.Version,
			StatusCode:    "200",
			ContentType:   contentType,
			ContentLength: contentLength,
			Data:          data,
		}
	}
	//jika uri nya ada parameter --> set GreeterNama jadi param
	paramValue = parsedURI.Query().Get("name")
	if paramValue != "" {

		greeterName = paramValue
	}
	greeter := GreetResponse{
		Student: student,
		Greeter: greeterName,
	}
	endpoints := strings.Split(req.Uri, "/")
	fmt.Print(endpoints)

	if strings.Contains(req.Accept, "application/xml") {
		contentType = "application/xml"
		var dataJson, err = xml.Marshal(greeter)
		if err != nil {
			log.Fatalln("Error parsing to JSON")
		}
		data = string(dataJson)
		contentLength = len(data)
	} else if strings.Contains(req.Accept, "text/html") || string((req.Uri)[len(req.Uri)-1]) == "/" {
		contentType = "text/html"
		data = fmt.Sprintf("<html><body><h1>Halo, dunia! aku %s</h1></body></html>", STUDENT_NAME)
		contentLength = len(data)
	} else if strings.Contains(req.Accept, "application/json") || strings.EqualFold(req.Uri, "http") {
		contentType = "application/json"
		var dataJson, err = json.Marshal(greeter)
		if err != nil {
			log.Fatalln("Error parsing to JSON")
		}
		data = string(dataJson)
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
	return response

}

func RequestDecoder(bytestream []byte) HttpRequest {
	// Put the decoding program for HTTP Request Packet here
	msgSplit := strings.Split(string(bytestream), "\r\n") //[GET djkfhsjdjf HTTP/1.1 Host: 127.0.0.1 Accept: dfhksjdfh ] len:3

	var host, accept string
	host = strings.TrimSpace((strings.TrimPrefix(msgSplit[1], "Host:")))
	accept = strings.TrimSpace(msgSplit[2])
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
	if res.Data == "" {
		responseString := fmt.Sprintf("%s %d\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
			res.Version, 404, res.ContentType, res.ContentLength, res.Data)

		return []byte(responseString)

	}
	responseString := fmt.Sprintf("%s %d\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		res.Version, 200, res.ContentType, res.ContentLength, res.Data)
	return []byte(responseString)
}
