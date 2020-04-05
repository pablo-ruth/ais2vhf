package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"

	"github.com/jacobsa/go-serial/serial"
)

func main() {

	// Open serial connection
	options := serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        4800,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
		ParityMode:      serial.PARITY_NONE,
	}
	port, err := serial.Open(options)
	if err != nil {
		fmt.Printf("failed top open serial connection: %s\n", err)
		return
	}
	defer port.Close()

	c, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		fmt.Printf("failed to connect to NMEA server: %s\n", err)
		return
	}
	defer c.Close()

	for {
		// Read NMEA message from connection
		msg, _, err := bufio.NewReader(c).ReadLine()
		if err != nil {
			fmt.Printf("failed to read msg: %s\n", err)
			return
		}

		// Catch only GPS messages
		match, err := regexp.MatchString("\\$..(GLL|GGA|RMC).*", string(msg))
		if err != nil {
			fmt.Printf("failed to match string: %s\n", err)
			continue
		}
		if !match {
			continue
		}

		// Send NMEA message to serial tty
		_, err = port.Write([]byte(fmt.Sprintf("%s\r\n", string(msg))))
		if err != nil {
			fmt.Printf("failed to write message to serial tty: %s\n", err)
			continue
		}
	}
}
