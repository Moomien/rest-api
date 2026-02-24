package main

import (
	"encoding/json"
	"os"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{
		Name: "Alex",
		Age:  19,
	}
	json.NewEncoder(os.Stdout).Encode(u)
}
