package converter

import (
	"fmt"
	"net"
	"strings"
)

type Response struct {

	Protocol string

	Status int

	HeaderMap map[string]string

	Body string
}


var StatusMapping = map[int]string {
	200: "OK",
	404: "Not Found",
}

func WriteResponse(conn net.Conn, response Response) {
	var sb strings.Builder
	//HTTP/1.1 200 OK
	sb.WriteString(fmt.Sprintf("%s %d %s \r\n", response.Protocol, response.Status, StatusMapping[response.Status]))

	for k, v := range response.HeaderMap {
		if k != "Content-Length" {
			sb.WriteString(fmt.Sprintf("%s:%s\r\n", k, v))
		}

	}
	sb.WriteString(fmt.Sprintf("Content-Length:%d\r\n", len(response.Body)))
	sb.WriteString("\r\n")
	sb.WriteString(response.Body)
	rawResponse := sb.String()
	conn.Write([]byte(rawResponse))
}
