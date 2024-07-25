package utils

import (
	"GptSSH/llm"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

func Init(u *User, channel ssh.Channel){
    if u == nil{
        panic("User was nil")
    }
    if checkPermissions(u){
        for {
            channel.Write([]byte("$ "))
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
            req := llm.Request{
                Model: "llama3",
                Prompt: input,
            }
            req.Send(channel)
        }
    }else{
        channel.Write([]byte("auth needed\n"))
    }
}

func checkPermissions(u *User) bool{
    if u.logged == true{
        return true
    }
    return false
}
