package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"

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


func hashData(username, password []byte) []byte{
    hash := sha256.Sum256(password)
    storing := append(username, ':')
    storing = append(storing, hash[:]...)
    return storing
}
