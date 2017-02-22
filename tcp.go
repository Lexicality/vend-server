package main

import (
	"bufio"
	"io"
	"net"
	"strings"
	"syscall"
)

func tcpServer(addr string) {
	sv, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %s", addr, err)
		return
	}
	defer sv.Close()

	for {
		conn, err := sv.Accept()
		if err == syscall.EINVAL {
			log.Critical("TCP connection go byebye :(")
			break
		} else if err != nil {
			log.Errorf("Unable to accept connection: %s", err)
			continue
		}

		log.Debug("TCP connection opened")
		go handleTCPRead(conn)
		go handleTCPWrite(conn)
	}
}

func handleTCPWrite(conn net.Conn) {
	log.Debug("Hello I'm a writer")
	in := messageSub()
	out := bufio.NewWriter(conn)
	for {
		msg, ok := <-in
		if !ok {
			log.Debug("Channel closed - closing socket")
			out.WriteString("bye\n")
			conn.Close()
			break
		}

		_, err := out.WriteString(msg)
		// Add the delimeter
		if err == nil {
			out.WriteByte('\n')
		}
		// Send it to ye server
		if err == nil {
			err = out.Flush()
		}

		if err == io.EOF {
			log.Debug("Socket closed")
			break
		} else if err != nil {
			log.Infof("Unable to write to socket: %s", err)
			conn.Close()
			break
		}
	}
	log.Debug("RIP - writer")
}

func handleTCPRead(conn net.Conn) {
	log.Debug("Hello I'm a reader")
	in := bufio.NewReader(conn)
	for {
		pkt, err := in.ReadString('\n')
		if err == io.EOF {
			log.Debug("Socket closed")
			break
		} else if err != nil {
			log.Infof("Unable to read from socket: %s", err)
			conn.Close()
			break
		}
		pkt = strings.TrimSpace(pkt)

		if pkt == "ping" {
			_, err = conn.Write([]byte("pong\n"))
			if err != nil {
				log.Infof("Unable to respond to ping: %s", err)
				conn.Close()
				break
			}
			continue
		}

		// ...?
		log.Debugf("Got TCP message: %s", pkt)
	}
	log.Debug("RIP - reader")
}
