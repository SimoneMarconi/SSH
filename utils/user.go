package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"golang.org/x/crypto/ssh"
)

func checkUser(username, password []byte, users map[string][]byte) bool {
    val, ok := users[string(username)]
    if !ok{
        return false
    }
    comp := slices.Compare(val, password)
    if comp == 0{
        return true
    }
    return false
}

func loadUsers() (map[string][]byte, error){
    users := make(map[string][]byte)
    file, err := os.Open("./files/users.txt")
    if err != nil{
        return nil, err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan(){
        line := scanner.Text()
        log.Println("line scanned" + line)
        idx := strings.Index(line, ":")
        if idx == -1{
            break 
        }
        log.Println(idx)
        username := line[:idx]
        password := line[idx + 1:]
        users[username] = []byte(password)
    }
    log.Printf("this is the map: %s",users)
    return users, nil
}

func saveUser(username, password []byte, users map[string][]byte) error{
    hash := hashData(username, password)
    _, found := users[string(username)]
    if found {
        return errors.New("Username already in use") 
    }

    file, err := os.OpenFile("./files/users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Panicf("Error with users file, err: %s", err)
    }
    defer file.Close()
    _, err = file.Write(append(hash, '\n'))
    if err != nil{
        log.Panicf("Error writing in the file, err: %s", err)
    }
    return nil
}

func addUser(channel ssh.Channel){
    password := readInput("Insert SuperPassword: ", channel)
    //Password check
    if string(password) != "test"{
        channel.Write([]byte("Wrong SuperPassword\n"))
        return
    }
    users, err := loadUsers()
    if err != nil {
        panic(err)
    }
    var newUser User
    newUser.username = readInput("Username: ", channel)
    newUser.password = readInput("Password: ", channel)
    err = saveUser(newUser.username, newUser.password, users)
    if err != nil{
        channel.Write([]byte(fmt.Sprintf("%s\n", err)))
    }
}
