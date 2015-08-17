package gwitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"unicode/utf8"
)

type TwitchChat struct {
	Username string
	OAuth    string
	Channel  string
	Conn     net.Conn
}

type Data struct {
	Method   string
	Username string
	Message  string
}

func New(user string, oauth string, channel string, conn net.Conn) (UI *TwitchChat) {
	UI = &TwitchChat{user, oauth, channel, conn}
	return
}

func (chat *TwitchChat) Connect() (err error) {
	chat.Conn, err = net.Dial("tcp", "irc.twitch.tv:6667")

	chat.SendRawData("PASS " + chat.OAuth)
	chat.SendRawData("NICK " + chat.Username)
	chat.SendRawData("JOIN " + chat.Channel)

	return
}

func (chat *TwitchChat) Close() (err error) {
	err = chat.Conn.Close()

	return
}

func (chat *TwitchChat) SendRawData(data string) (err error) {
	_, err = fmt.Fprintf(chat.Conn, "%s\r\n", data)

	return
}

func (chat *TwitchChat) SendMessage(data string) (err error) {
	message := fmt.Sprintf(":%s!%s@%s.tmi.twitch.tv PRIVMSG %s :%s", chat.Username, chat.Username, chat.Username, chat.Channel, data)
	err = chat.SendRawData(message)

	return
}

func (chat *TwitchChat) RawRead() (data string, err error) {
	var line string
	reader := bufio.NewReader(chat.Conn)
	tp := textproto.NewReader(reader)

	for {
		line, err = tp.ReadLine()

		if err != nil {
			return "", err
		}

		if strings.HasPrefix(line, "PING") {
			chat.SendRawData("PONG " + line[5:])
		} else {
			data = line
			return
		}

	}
}

func (chat *TwitchChat) ReadData() (data *Data) { //Here comes some spaghetti, let's hope someone or myself will make this non-trashy
	var rawdata string
	var formattedstring string
	var err error

	for {

		rawdata, err = chat.RawRead()

		if err != nil {
			return &Data{"ERROR", "", ""}
		}

		formattedstring = fmt.Sprintf(".tmi.twitch.tv PRIVMSG %s :", chat.Channel)

		switch {

		case strings.Contains(rawdata, formattedstring):
			tempstring1 := strings.Split(rawdata, formattedstring)
			tempstring2 := strings.Split(tempstring1[0], "@")

			if tempstring2[1] != chat.Username { //A quick fix to ignore messages created by self
				data = &Data{"MESSAGE", tempstring2[1], tempstring1[1]}
				return
			}

		case strings.Contains(rawdata, ".tmi.twitch.tv JOIN %s"+chat.Channel):
			tempstring1 := strings.Split(rawdata, ".tmi.twitch.tv JOIN %s"+chat.Channel)
			tempstring2 := strings.Split(tempstring1[0], "@")

			data = &Data{"JOIN", tempstring2[1], ""}
			return

		case strings.Contains(rawdata, ".tmi.twitch.tv PART %s"+chat.Channel):
			tempstring1 := strings.Split(rawdata, ".tmi.twitch.tv PART %s"+chat.Channel)
			tempstring2 := strings.Split(tempstring1[0], "@")

			data = &Data{"PART", tempstring2[1], ""}
			return

		case strings.HasPrefix(rawdata, ":jtv MODE "+chat.Channel+" +o "):
			data = &Data{"GAINOP", rawdata[14+utf8.RuneCountInString(chat.Channel):], ""}
			return

		case strings.HasPrefix(rawdata, ":jtv MODE "+chat.Channel+" -o "):
			data = &Data{"LOSEOP", rawdata[14+utf8.RuneCountInString(chat.Channel):], ""}
			return
		}

	}
}
