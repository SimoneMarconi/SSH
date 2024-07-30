package llm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"

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
    reader := bufio.NewReader(stdout)
    for{
        var token Res
        line, _, err := reader.ReadLine()
        if err != nil{
            fmt.Println(err)
            break
        }
        if len(line) != 0 {
            err := json.Unmarshal(line, &token)
            if err != nil {
                panic(err)
            }
            channel.Write([]byte(token.Response))
        }
    }
    channel.Write([]byte("\n"))
}
