package main

import "fmt"

func write() bool {
	return true
}

func main() {
	if write() == true {
		fmt.Println("OK")
	}
}
