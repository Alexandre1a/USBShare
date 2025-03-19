package main

import (
	"fmt"
	"log"
	"net"
	"go.bug.st/serial"
)

func main() {
	// Let user choose USB port
	ports, _ := serial.GetPortsList()
	fmt.Println("Available ports:", ports)
	var portName string
	fmt.Print("Enter USB port (e.g., /dev/ttyACM0): ")
	fmt.Scanln(&portName)

	// Open USB serial port
	port, err := serial.Open(portName, &serial.Mode{BaudRate: 115200})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Forward USB traffic to a local SOCKS proxy
	listener, _ := net.Listen("tcp", "127.0.0.1:1080") // SOCKS port
	for {
		clientConn, _ := listener.Accept()
		go func() {
			defer clientConn.Close()
			buf := make([]byte, 1024)
			for {
				n, _ := clientConn.Read(buf)
				port.Write(buf[:n]) // Send to Windows via USB
				respBuf := make([]byte, 1024)
				m, _ := port.Read(respBuf)
				clientConn.Write(respBuf[:m]) // Send back to app
			}
		}()
	}
}
