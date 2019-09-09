package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

// SlackListener data type
type SlackListener struct {
	api       *slack.Client
	botID     string
	channelID string
}

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

// "ListenAndResponse method respond to event changes on slack channel"
func (s *SlackListener) ListenAndResponse(dbConnection *dbAccess) {
	rtm := s.api.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev, dbConnection); err != nil {
				fmt.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent, dbConnection *dbAccess) error {
	// make sure the BOT ID is mentioned
	fmt.Println(dbConnection.commandList)
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 {
		return fmt.Errorf("invalid message")
	}

	switch strings.ToLower(m[0]) {
	case "hello":
		greetings := getWelcomeMessage(dbConnection)
		// response post message to the slack
		if _, _, err := s.api.PostMessage(ev.Channel, slack.MsgOptionText(greetings, true)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	case "translate":
		// get the actual work to translate
		text := m[1:]
		if len(text) == 0 {
			return fmt.Errorf("No word or text to translate")
		}
		// response post message to the slack
		if _, _, err := s.api.PostMessage(ev.Channel, slack.MsgOptionText(translateToAmharic(text), true)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	default:
		var word string
		switch len(dbConnection.commandList) {
		case 0:
			word = dbConnection.getAmharicWord(strings.ToLower(m[0]))
		default:
			word = getAmharicWord(dbConnection, strings.ToLower(m[0]))
		}
		// response post message to the slack
		if _, _, err := s.api.PostMessage(ev.Channel, slack.MsgOptionText(word, true)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	}
	return nil
}

func getWelcomeMessage(dbCon *dbAccess) string {
	var stringBuilder strings.Builder
	botCommands := dbCon.getListofCommands()
	stringBuilder.WriteString("Please use the following commands to learn some amharic:\n")
	for _, c := range botCommands {
		stringBuilder.WriteString("\t" + c.Command + "\n")
	}
	stringBuilder.WriteString("You can replace the '[@bot_user]' with your app name\n")
	stringBuilder.WriteString("use '@bot_user translate any word' to translate to amharic")
	return stringBuilder.String()
}

func getAmharicWord(dbCon *dbAccess, command string) string {
	commandMap := map[string]string{}
	for _, c := range dbCon.commandList {
		commandMap[strings.Split(c.Command, " ")[1]] = strconv.Itoa(c.ID)
	}
	if v, ok := commandMap[command]; ok {
		return dbCon.getAmharicWordByID(v)
	}
	return "Invalid command"
}
