package llm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)



type Request struct{
    Model string `json:"model"`
    Prompt string `json:"prompt"`
}

type Res struct{
    Model string `json:"model"`
    Created_at string `json:"created_at"`
    Response string `json:"response"`
    Done bool `json:"done"`
}

func (r Request) Send(channel ssh.Channel){
    payload, err := json.Marshal(r)
    if err != nil {
        panic("Error parsing JSON")
    }

    cmd := exec.Command("curl", "http://localhost:11434/api/generate", "-d", string(payload))
    stdout, err := cmd.StdoutPipe()
    if err != nil{
        panic("Could not get stdout")
    }
    if err := cmd.Start(); err != nil{
        panic(err)
    }
    for{
        var token Res
        test := make([]byte, 256)
        n, _:= stdout.Read(test)
        // fmt.Println(string(test))
        if n < 2 {
            fmt.Println("not reading")
            break
        }
        //Only reading 256 bytes, fix this
        reader := bufio.NewReader(stdout)
        line, _, _ := reader.ReadLine()
        fmt.Println(string(line))
        if err := json.NewDecoder(strings.NewReader(string(line))).Decode(&token); err != nil{
            fmt.Println(err, n)
            break
        }
        channel.Write([]byte(token.Response))
    }
}
