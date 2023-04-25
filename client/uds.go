package main

import (
	"fmt"
	"net"
	"os"
)

func udsHandler() error {
	_ = os.Remove(socketPath)

	l, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: socketPath})
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func(c *net.UnixConn) {
			defer c.Close()

			buf := make([]byte, 8)
			_, err = c.Read(buf)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				return
			}

			setByte := buf[0]
			dpidMacAddr := buf[1:7]
			switchPort := buf[7]

			if setByte == 0x00 {
				// This is a get request
				value := read(string(dpidMacAddr))
				_, _ = c.Write([]byte(value))
			} else {
				// This is a set request
				write(string(dpidMacAddr), string(switchPort))
				_, _ = c.Write([]byte("OK"))
			}
		}(conn)
	}
}
