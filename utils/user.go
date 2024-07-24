package utils

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type User struct{
    username []byte
    password []byte
    logged bool
}

func (u User) String() string{
    return fmt.Sprintf("user: %s\npass: %s\nlog: %t", u.username, u.password, u.logged)
}

func (u *User) Login(){
    u.logged = true
}

func (u *User) LogOut(){
    u.logged = false 
}

func checkUser(u User, users map[string]User) (bool, *User){
    val, ok := users[string(u.username)]
    if !ok{
        return false, nil
    }
    hash := sha256.Sum256(u.password)
    comp := slices.Compare(hash[:], val.password)
    if comp == 0{
        return true, &val
    }
    return false, nil
}

func loadUsers() (map[string]User, error){
    users := make(map[string]User)
    file, err := os.Open("./files/users.txt")
    if err != nil{
        return nil, err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan(){
        line := scanner.Text()
        log.Println("line scanned" + line)
        tokens := strings.Split(line, ":")
        username := tokens[0]
        password := tokens[1]
        //TODO: handle logged in the file
        users[username] = User{
            username: []byte(username),
            password: []byte(password),
            logged: false,
        }
    }
    // log.Printf("this is the map: %s",users)
    return users, nil
}

func saveUser(u User, users map[string]User) error{
    hash := hashData(u)
    _, found := users[string(u.username)]
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
    newUser.logged = false
    err = saveUser(newUser, users)
    if err != nil{
        channel.Write([]byte(fmt.Sprintf("%s\n", err)))
    }
}

func login(channel ssh.Channel) *User {
    users, err := loadUsers()
    if err != nil{
        panic("Could not load Users")
    }
    username := readInput("Username: ", channel)
    password := readInput("Password: ", channel)
    testUser := User{
        username: username,
        password: password,
        logged: false,
    }
    fmt.Println(testUser.String())
    if res, u := checkUser(testUser, users); res == true{
        if u.logged == true{
            channel.Write([]byte("User already logged in"))
            return u
        }
        fmt.Printf("logging in %s\n", u.username)
        channel.Write([]byte("Login Successfull\n"))
        go handleLogintime(u)
        return u
    }
    channel.Write([]byte("Wrong Username or Password"))
    return nil
}

func handleLogintime(u *User){
    u.logged = true
    timer := time.NewTimer(time.Second * 10).C
    <-timer
    u.logged = false
    fmt.Printf("logging out %s\n", u.username)
}
