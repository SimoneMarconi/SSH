package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
    "GptSSH/utils"

	ssh "golang.org/x/crypto/ssh"
)

func main() {
    KeysMap := readKeys("./files/authorized_keys")
    //config is the config struct to handle SSH for a tcp server
    config := &ssh.ServerConfig{
        NoClientAuth: true,
        PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey)(*ssh.Permissions, error){
            if KeysMap[string(key.Marshal())]{
                return &ssh.Permissions{
                    Extensions: map[string]string{
                        "pubkey-fp": ssh.FingerprintSHA256(key),
                    },
                }, nil
            }
            return nil, fmt.Errorf("Unknown public key for %s", c.User())
        },
    }

    privateKeyBytes, err := os.ReadFile("./files/id_rsa")
    if err != nil {
        panic("Could not read rsa")
    }
    privateKey, err := ssh.ParsePrivateKeyWithPassphrase(privateKeyBytes, []byte("test"))
    if err != nil {
        log.Panicf("Could not parse the private key, err: %s", err)
    }
    config.AddHostKey(privateKey)
    //the configuration to handle SSH is finished, now we have to configure the tcp endpoint listening for connections

    log.Println("Starting to listen")
    listener, err := net.Listen("tcp", "127.0.0.1:2000")
    if err != nil {
        panic("Could not initialize the tcp listener")
    }
    defer listener.Close()
    //listening loop
    for {
        conn, err := listener.Accept()
        if err != nil{
            log.Println("Could not accept the connections")
        }
        go handleConnection(conn, config)
    }

}

func readKeys(file string) map[string]bool{
	authorizedKeysBytes, err := os.ReadFile(file)
    if err != nil{
        panic("Keys file not found")
    }

    authorizedKeysMap := map[string]bool{}

    for len(authorizedKeysBytes) > 0 {
        publicKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
        if err != nil{
            log.Panicf("Error parsing Public Keys: %s", err)
        }

        authorizedKeysMap[string(publicKey.Marshal())] = true
        authorizedKeysBytes = rest
    }
    return authorizedKeysMap
}


func handleConnection(conn net.Conn, config *ssh.ServerConfig){
    //now we have to perform the SSH handshake
    _, ch, reqch, err := ssh.NewServerConn(conn, config)
    if err != nil{
        panic("Error during the Handshake")
    }

    //Handling sync with the different server connections
    var wg sync.WaitGroup
    //with this the function will not exit focus untill all the request have been fulfilled
    defer wg.Wait()

    wg.Add(1)
    go func(){
        //TODO: Handle this 
        ssh.DiscardRequests(reqch)//we are discarding all the requests
        wg.Done()
    }()

    for newChannel := range ch{
        //we only handle shell sessions, these are signed with the type session
        if newChannel.ChannelType() != "session"{
            newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
            continue
        }
        channel, req, err := newChannel.Accept()
        if err != nil{
            panic("Could not accept Channel")
        }
        wg.Add(1)
        //handling requests
        go func(in <-chan *ssh.Request){
            for r := range in{
                r.Reply(r.Type == "shell", nil)
            }
            wg.Done()
        }(req)

        go utils.HandleInput(channel)

    }
}
