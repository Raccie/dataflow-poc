package helpers

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func sendLineByLine(r io.Reader, host string, msg, hostname chan string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		msg <- scanner.Text()
		hostname <- host
	}
}

// TCPServer starts a simple TCP Server that listens that binds on 'address' and 'port' passed as arguments
// On each message, the Server sends each line separated with \n to the 'out' channel until client closes connection
func TCPServer(msg, hostname chan string, host, port string) {
	defer close(msg)
	defer close(hostname)

	fmt.Println("Serving on", host+":"+port)
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		panic(err.Error())
	}
	defer l.Close()

	for {
		fmt.Println("Accepting new request")
		conn, err := l.Accept()
		if err != nil {
			panic(err.Error)
		}
		go sendLineByLine(conn, conn.RemoteAddr().String(), msg, hostname)
	}
}
