package main

import (
    "code.google.com/p/go.crypto/bcrypt"
    "encoding/base64"
    "fmt"
    "flag"
    "os"
    "crypto/rand"
)

func main() {
    password := flag.String("password", "", "the password")

    flag.Parse()

    if *password == "" || len(*password) < 5 {
        fmt.Println("Bad password")
        os.Exit(1)
    }

    fmt.Println("starting")

    saltRaw := make([]byte, 64)
    rand.Read(saltRaw)
    salt := base64.StdEncoding.EncodeToString(saltRaw)
    fmt.Println("Salt: ", salt)

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(*password + salt), 15)
    if err != nil {
        fmt.Println("Failed to generate password")
        os.Exit(1)
    }
    fmt.Println("Password: ", base64.StdEncoding.EncodeToString(passwordHash))
}