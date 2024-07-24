package utils

import "golang.org/x/crypto/ssh"

func Init(u *User, channel ssh.Channel){
    if u == nil{
        panic("User was nil")
    }
    if checkPermissions(u){
        channel.Write([]byte("you are in\n"))
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
