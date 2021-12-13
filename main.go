package main

import (
	"dfe/pkg/blocks"
	"dfe/pkg/helpers"
	"net"
	"os"

	"github.com/trustmaster/goflow"
)

func NewApp() *goflow.Graph {
	n := goflow.NewGraph()
	n.Add("osfw", blocks.NewOSFWParser())

	n.Add("printer", blocks.NewPrinter(os.Stdout))
	conn, _ := net.Dial("tcp", "localhost:5000")
	n.Add("printer2", blocks.NewPrinter(conn))
	n.Connect("osfw", "out", "printer", "message")
	n.Connect("osfw", "out", "printer2", "message")

	n.MapInPort("In1", "osfw", "message")
	n.MapInPort("In2", "osfw", "hostname")
	return n
}

func main() {
	n := NewApp()

	msg_in := make(chan string)
	hostname_in := make(chan string)
	n.SetInPort("In1", msg_in)
	n.SetInPort("In2", hostname_in)

	wait := goflow.Run(n)

	go helpers.TCPServer(msg_in, hostname_in, "localhost", "8080")

	<-wait
}
