package gwitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

type TwitchChat struct {
	Username string
	OAuth    string
	Channel  string
	Conn     net.Conn
}

func New(user string, oauth string, channel string, conn net.Conn) *TwitchChat {
	UI := &TwitchChat{user, oauth, channel, conn}
	return UI
}

func (chat *TwitchChat) Connect() error {
	var err error

	chat.Conn, err = net.Dial("tcp", "irc.twitch.tv:6667")

	chat.SendRawData("PASS " + chat.OAuth)
	chat.SendRawData("NICK " + chat.Username)
	chat.SendRawData("JOIN " + chat.Channel)

	return err
}

func (chat *TwitchChat) Close() error {
	var err error

	err = chat.Conn.Close()

	return err
}

func (chat *TwitchChat) SendRawData(data string) error {
	var err error

	_, err = fmt.Fprintf(chat.Conn, "%s\r\n", data)

	return err
}

func (chat *TwitchChat) SendMessage(data string) error {
	var err error

	message := fmt.Sprintf(":%s!%s@%s.tmi.twitch.tv PRIVMSG %s :%s", chat.Username, chat.Username, chat.Username, chat.Channel, data)
	err = chat.SendRawData(message)
	fmt.Println(message)

	return err
}

func (chat *TwitchChat) RawRead() string {

	reader := bufio.NewReader(chat.Conn)
	tp := textproto.NewReader(reader)
	for {
		line, _ := tp.ReadLine()

		if strings.HasPrefix(line, "PING") {

			chat.SendRawData("PONG " + line[5:])

		} else {

			return line

		}
	}
}

func (chat *TwitchChat) ReadMessage() (string, string) {
	var rawdata string
	var formattedstring string

	for {

		rawdata = chat.RawRead()
		formattedstring = fmt.Sprintf(".tmi.twitch.tv PRIVMSG %s :", chat.Channel)

		if strings.Contains(rawdata, formattedstring) {

			tempstring1 := strings.Split(rawdata, formattedstring)
			tempstring2 := strings.Split(tempstring1[0], "@")

			return tempstring2[1], tempstring1[1]
		}

	}
}
