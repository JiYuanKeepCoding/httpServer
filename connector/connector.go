package connector

import (
	"context"
	"fmt"
	"httpServer/converter"
	"httpServer/router"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Server struct {
	activeConns sync.Map

	connCounts *int32

	Ctx context.Context

	Cancel context.CancelFunc

	Signal chan os.Signal
}

func NewSever() *Server {
	server := new(Server)
	ctx, cancel := context.WithCancel(context.Background())
	server.Ctx = ctx
	server.Cancel = cancel
	var i int32
	server.connCounts = &i
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	server.Signal = c
	return server
}

func (s *Server) Run(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	s.addTermHandler(ln)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				if strings.Contains(err.Error(), "closed") {
					break
				}
				continue
			}
			s.activeConns.Store(conn, 0)
			atomic.AddInt32(s.connCounts, 1)
			go s.handleConnection(&conn)
		}
	}()

	<-s.Ctx.Done()
}

func (s *Server) addTermHandler(ln net.Listener) {
	go func() {
		oscall := <-s.Signal
		log.Printf("system call:%+v", oscall)

		//graceful shutdown http server wait for all transaction complete util timeout
		ln.Close()
		shutdownContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
		for *s.connCounts > 0 {
			select {
			case <-shutdownContext.Done():
				s.activeConns.Range(func(key, value interface{}) bool {
					key.(net.Conn).Close()
					return true
				})
				break
			case <-time.After(1 * time.Second):
			}
		}
		s.Cancel()
	}()
}

func (s *Server) handleConnection(conn *net.Conn) {
	defer func() {
		s.activeConns.Delete(conn)
		atomic.AddInt32(s.connCounts, -1)
		(*conn).Close()
	}()

	size := 1024
	b := make([]byte, size)
	var parseContext converter.ParseContext
	addr := (*conn).RemoteAddr().(*net.TCPAddr)
	parseContext.Ip = addr.IP.String()
	parseContext.HeaderMap = make(map[string]string)
	payloadSize := 0
	for {
		n, err := (*conn).Read(b)
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
	converter.WriteResponse(*conn, response)
}
