package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/quipo/statsd"
	"log"
	"net"
	"os"
	"time"
)

var (
	addr = flag.String("addr", "127.0.0.1:30002", "\":port\" or \"ip:port\" of beast server to connect to")
)

func main() {
	flag.Parse()

	// StatsD
	prefix := "go1090."
	statsdclient := statsd.NewStatsdClient("172.16.107.114:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 1 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	fmt.Println("Connecting to server...")
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	messages := startAVRClient(conn)

	for {
		go parseAVRMessage(<-messages, stats)
		stats.Incr("messages.total", 1)
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

func parseAVRMessage(message []byte, stats *statsd.StatsdBuffer) {

	// http://wiki.modesbeast.com/Mode-S_Beast:Data_Output_Formats

	var content []byte

	isMlat := message[0] == '@'
	if isMlat {
		//		fmt.Println("MLAT message")
		timestamp := message[1:13]
		content = message[13 : len(message)-1]

		// Pack timestamp
		dst := make([]byte, hex.DecodedLen(len(timestamp)))
		_, err := hex.Decode(dst, timestamp)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Standard message")
		content = message[1 : len(message)-1]
	}

	// Pack message
	packedMessage := make([]byte, hex.DecodedLen(len(content)))
	_, err := hex.Decode(packedMessage, content)
	if err != nil {
		log.Fatal(err)
	}

	switch len(packedMessage) {
	case 2:
		//		fmt.Println("ModeAC")
		stats.Incr("messages.modeac", 1)

	case 7:
		//		fmt.Println("ModeS56")
		stats.Incr("messages.modes56", 1)
		stats.Incr("messages.modes", 1)

		//parseModeS(packedMessage)

	case 14:
		stats.Incr("messages.modes112", 1)
		stats.Incr("messages.modes", 1)

		//		fmt.Println("ModeS112")
		//parseModeS(packedMessage)
	}
}

func parseTime(timebytes []byte) time.Time {
	// http://wiki.modesbeast.com/Radarcape:Firmware_Versions#The_GPS_timestamp

	// SecondsOfDay are using the upper 18 bits of the timestamp
	upper := []byte{timebytes[0]<<2 + timebytes[1]>>6, timebytes[1]<<2 + timebytes[2]>>6, timebytes[3] >> 6}
	// Nanoseconds are using the lower 30 bits.
	lower := []byte{timebytes[2] & 0x3F, timebytes[3], timebytes[4], timebytes[5]}

	// the 48bit timestamp is 18bit day seconds | 30bit nanoseconds
	daySeconds := int(binary.BigEndian.Uint16(upper))
	nanoSeconds := int(binary.BigEndian.Uint32(lower))

	hour := int(daySeconds / 3600)
	minutes := int(daySeconds / 60 % 60)
	seconds := int(daySeconds % 60)

	// The date itself is not part of the timestamp
	utcDate := time.Now().UTC()

	return time.Date(
		utcDate.Year(), utcDate.Month(), utcDate.Day(),
		hour, minutes, seconds, nanoSeconds, time.UTC)
}
