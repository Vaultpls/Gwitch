package main

import (
	"fmt"
	"github.com/Vaultpls/Gwitch"
)

func main() {
	t := gwitch.New("yourbotname", "your oauth here", "#twitchstream", nil)

	err := t.Connect()

	if err != nil {
		panic(err)
	}

	for {
		user, message := t.ReadMessage()
		if message == "!hi" {
			t.SendMessage("Hi " + user + "!")
		}
		fmt.Println(user + ": " + message)
	}
}
