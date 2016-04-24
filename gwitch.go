package gwitch

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"unicode/utf8"
)

//TwitchChat containing Username, OAuth, Channel, Conn, and RawData necessary for connection
type TwitchChat struct {
	Username string
	OAuth    string
	Channel  string
	Conn     net.Conn
	RawData string
}

//Data containing Method, Username, and Message
type Data struct {
	Method   string
	Username string
	Message  string
}

//New Gwitch instance
func New(user string, oauth string, channel string, conn net.Conn) (UI *TwitchChat) {
	UI = &TwitchChat{user, oauth, channel, conn, ""}
	return
}

//Connect to the Twitch IRC servers
func (chat *TwitchChat) Connect() (err error) {
	chat.Conn, err = net.Dial("tcp", "irc.twitch.tv:6667")
	
	//SendRawData standard user info
	chat.SendRawData("PASS " + chat.OAuth)
	chat.SendRawData("NICK " + chat.Username)
	chat.SendRawData("JOIN " + chat.Channel)

	//SendRawData to request the ability to have IRCV3 and membership updates(JOIN & PART)
	chat.SendRawData("CAP REQ :twitch.tv/membership")
    chat.SendRawData("CAP REQ :twitch.tv/commands")


	return
}

//Close the connection
func (chat *TwitchChat) Close() (err error) {
	err = chat.Conn.Close()

	return
}

//SendRawData send raw data to the Twitch server
func (chat *TwitchChat) SendRawData(data string) (err error) {
	_, err = fmt.Fprintf(chat.Conn, "%s\r\n", data)

	return
}

//SendMessage sends message to chat using SendRawData
func (chat *TwitchChat) SendMessage(data string) (err error) {
	message := fmt.Sprintf(":%s!%s@%s.tmi.twitch.tv PRIVMSG %s :%s", chat.Username, chat.Username, chat.Username, chat.Channel, data)
	err = chat.SendRawData(message)

	return
}

//RawRead reads raw data from conn from TwitchChat
func (chat *TwitchChat) RawRead() (data string) {
	reader := bufio.NewReader(chat.Conn)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()

		if err != nil {
			return ""
		}

		if strings.HasPrefix(line, "PING") {
			chat.SendRawData("PONG " + line[5:])
		} else {
			data = line
			return
		}

	}
}

//ReadData directly from RawData, interprets if is join, part, message, op, and de-op
func (chat *TwitchChat) ReadData() (data *Data) { //Here comes some spaghetti, let's hope someone or myself will make this non-trashy
	var formattedstring string

	for {

		chat.RawData = chat.RawRead()

		if chat.RawData == "" {
			return &Data{"ERROR", "", ""}
		}

		formattedstring = fmt.Sprintf(".tmi.twitch.tv PRIVMSG %s :", chat.Channel)

		switch {

		case strings.Contains(chat.RawData, formattedstring):
			tempstring1 := strings.Split(chat.RawData, formattedstring)
			tempstring2 := strings.Split(tempstring1[0], "@")

			if tempstring2[1] != chat.Username { //A quick fix to ignore messages created by self
				data = &Data{"MESSAGE", strings.ToLower(tempstring2[1]), tempstring1[1]}
				return
			}

		case strings.Contains(chat.RawData, ".tmi.twitch.tv JOIN "+chat.Channel):
			tempstring1 := strings.Split(chat.RawData, ".tmi.twitch.tv JOIN "+chat.Channel)
			tempstring2 := strings.Split(tempstring1[0], "@")

			data = &Data{"JOIN", strings.ToLower(tempstring2[1]), ""}
			return

		case strings.Contains(chat.RawData, ".tmi.twitch.tv PART "+chat.Channel):
			tempstring1 := strings.Split(chat.RawData, ".tmi.twitch.tv PART "+chat.Channel)
			tempstring2 := strings.Split(tempstring1[0], "@")

			data = &Data{"PART", strings.ToLower(tempstring2[1]), ""}
			return

		case strings.HasPrefix(chat.RawData, ":jtv MODE "+chat.Channel+" +o "):
			data = &Data{"GAINOP", chat.RawData[14+utf8.RuneCountInString(chat.Channel):], ""}
			return

		case strings.HasPrefix(chat.RawData, ":jtv MODE "+chat.Channel+" -o "):
			data = &Data{"LOSEOP", chat.RawData[14+utf8.RuneCountInString(chat.Channel):], ""}
			return
		}

	}
}
