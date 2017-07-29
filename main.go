package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"log"
)

var (
	addr = flag.String("addr", "127.0.0.1:30002", "\":port\" or \"ip:port\" of beast server to connect to")
)

func main() {
	flag.Parse()

	fmt.Println("Connecting to server...")
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	messages := startAVRClient(conn)

	for {
		go parseAVRMessage(<-messages)
	}
}

func startAVRClient(conn net.Conn) chan []byte {

	ch := make(chan []byte)
	reader := bufio.NewReader(conn)

	go func() {
		for {
			currentMessage, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			if currentMessage != nil {
				currentMessage = bytes.TrimRight(currentMessage, "\n")
				ch <- currentMessage
			}
		}
	}()
	return ch
}

func parseAVRMessage(message []byte) {

	// http://wiki.modesbeast.com/Mode-S_Beast:Data_Output_Formats

	isMlat := message[0] == '@'

	if isMlat {
		fmt.Println("MLAT message")
	} else {
		fmt.Println("Standard message")
	}

}
