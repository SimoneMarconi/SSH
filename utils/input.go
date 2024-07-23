package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type User struct{
    username []byte
    password []byte
}

func HandleInput(channel ssh.Channel){
        for {
            channel.Write([]byte("> "))
            data := make([]byte, 256)
            n, err := channel.Read(data)
            if err != nil {
                if err == io.EOF{
                    break
                }
                panic("Error reading channel data")
            }
        input := string(data[:n-1])
        log.Println(input)
        switch (input){
        case "addUser":
            addUser(channel) 
            break
        case "exit":
            channel.Close()
            break
        }

    }
}

func addUser(channel ssh.Channel){
    password := readInput("Insert SuperPassword: ", channel)
    //Password check
    if string(password) != "test"{
        channel.Write([]byte("Wrong SuperPassword\n"))
        return
    }
    var newUser User
    newUser.username = readInput("Username: ", channel)
    newUser.password = readInput("Password: ", channel)
    saveUser(newUser.username, newUser.password)
}

func readInput(prompt string, channel ssh.Channel) []byte{
    channel.Write([]byte(prompt))
    data := make([]byte, 100)
    for{
        n, err := channel.Read(data)
        if err != nil{
            channel.Write([]byte(fmt.Sprint(err)))
        }
        if data[n-1] == '\n'{
            data = data[:n-1]
            break
        }
    }
    return data
}

func saveUser(username, password []byte){
    hash := hashData(username, password)
    file, err := os.OpenFile("./files/users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Panicf("Error with users file, err: %s", err)
    }
    defer file.Close()
    _, err = file.Write(hash)
    if err != nil{
        log.Panicf("Error writing in the file, err: %s", err)
    }
}

func hashData(username, password []byte) []byte{
    hash := sha256.Sum256(password)
    storing := append(username, ':')
    storing = append(storing, hash[:]...)
    return storing
}

func checkUser(username, password []byte){

}
