package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	var err error
	//----------------------------------------------------------------------------------------------------------------------
	//Get and decrypt password
	var keyDir, passwDir string
	var key, passw, plaintext []byte

	fmt.Print("Key file: ")
	if _, err = fmt.Scan(&keyDir); err != nil {
		panic(err)
	}

	fmt.Print("Password file: ")
	if _, err = fmt.Scan(&passwDir); err != nil {
		panic(err)
	}

	key, err = ioutil.ReadFile(keyDir)
	if err != nil {
		panic(err)
	}

	passw, err = ioutil.ReadFile(passwDir)
	if err != nil {
		panic(err)
	}
/*
	plaintext, err = decrypt(passw, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Password decrypted")
*/
	fmt.Println(string(passw))
	plaintext = passw
	fmt.Println("Key '" + string(key) + "' is not using")
//----------------------------------------------------------------------------------------------------------------------
	// Connect to the remote SMTP server.
	var addr, port, sender, recipient, subject, message, full string
/*
	fmt.Print("Enter server adress: ")
	if _, err = fmt.Scan(&addr); err != nil {
		panic(err)
	}
*/
	addr = "smtp.yandex.ru"
/*
	fmt.Print("Enter port: ")
	if _, err = fmt.Scan(&port); err != nil {
		panic(err)
	}
*/
	port = "465"
/*
	fmt.Print("Enter sender email: ")
	if _, err = fmt.Scan(&sender); err != nil {
		panic(err)
	}
*/
	sender = "fumihin@yandex.ru"
/*
	fmt.Print("Enter recipient email: ")
	if _, err = fmt.Scan(&recipient); err != nil {
		panic(err)
	}
*/
	recipient = sender

	in := bufio.NewReader(os.Stdin)

	fmt.Println("Enter `Subject`: ")
	subject, err = in.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println("Enter `Message body`: ")
	message, err = in.ReadString('\n')
	if err != nil {
		panic(err)
	}

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = recipient
	headers["Subject"] = subject

	for k, v := range headers {
		full += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	full += "\r\n" + message

	auth := smtp.PlainAuth("", sender, string(plaintext), addr)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         addr,
	}

	conn, err := tls.Dial("tcp", addr+":"+port, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, addr)
	if err != nil {
		log.Panic(err)
	}

	//Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = client.Mail(sender); err != nil {
		log.Panic(err)
	}

	if err = client.Rcpt(recipient); err != nil {
		log.Panic(err)
	}

	// Data
	writer, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = writer.Write([]byte("Lublu sosat" + full))
	if err != nil {
		log.Panic(err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic(err)
	}
}