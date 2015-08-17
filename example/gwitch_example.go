package main

import (
	"github.com/Vaultpls/Gwitch"
)

func main() {
	t := gwitch.New("quanticbot", "your own oauth", "#vaultpls", nil)

	err := t.Connect()

	if err != nil {
		panic(err)
	}

	for {

		data, err := t.ReadData()

		if err != nil{
			panic(err)
		}

		if data.Method == "MESSAGE" && data.Message == "!hi" {
			t.SendMessage("Hi " + data.Username + "!  I am " + t.Username + "!")
		}

		switch {

		case data.Method == "JOIN":
			t.SendMessage("Hey there " + data.Username + "!")

		case data.Method == "PART":
			t.SendMessage("Bye " + data.Username + "!")

		case data.Method == "GAINOP":
			t.SendMessage(data.Username + " gained op!")

		case data.Method == "LOSEOP":
			t.SendMessage(data.Username + " lost op!")

		}
	}
}
