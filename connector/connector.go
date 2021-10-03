package connector

import (
	"fmt"
	"httpServer/converter"
	"httpServer/router"
	"io"
	"net"
)

func Run(port string)  {
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	size := 1024
	b := make([]byte, size)
	var parseContext converter.ParseContext
	addr := conn.RemoteAddr().(*net.TCPAddr)
	parseContext.Ip = addr.IP.String()
	parseContext.HeaderMap = make(map[string]string)
	payloadSize := 0
	for {
		n, err := conn.Read(b)
		if err != nil && err != io.EOF {
			println(err)
			return
		}
		payloadSize += n
		if payloadSize == 0 {
			return
		}
		parseContext.Parse(b[:n])
		if parseContext.Status == converter.BODY {
			break
		}
	}
	response := router.HandleHttpRequest(parseContext.HttpPayload)
	fmt.Printf("Client IP: %s, Response Code: %d \r\n", parseContext.Ip, response.Status)
	converter.WriteResponse(conn, response)
}
