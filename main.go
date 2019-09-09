package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

// environment variables used are placed under a truct
type envConfig struct {
	Port              string `envconfig:"PORT" default:"5000"`
	BotToken          string `envconfig:"BOT_TOKEN" required:"true"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`
	BotID             string `envconfig:"BOT_ID" required:"true"`
	ChannelID         string `envconfig:"CHANNEL_ID" required:"false"`
	BotDBUserPassword string `envconfig:"BOTDB_USER_PASSWORD" default:"5000"`
	BotDBUser         string `envconfig:"BOT_DB_USER" default:"botUser"`
	DBPort            string `envconfig:"DB_PORT" default:"3306"`
	DBAddress         string `envconfig:"DB_ADDRESS" default:"localhost"`
	DBName            string `envconfig:"DB_NAME" default:"goBotDB"`
}

// added the struct below to verify this app from slack
type slackResponse struct {
	Token     string
	Challenge string
	Type      string
}

func main() {
	//used github.com/kelseyhightower/envconfig library to get all the list of env variables
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		fmt.Printf("[ERROR] Failed to process env var: %s", err)
	}
	// github.com/nlopes/slack library supports most if not all of the api.slack.com REST calls,
	// as well as the Real-Time Messaging protocol over websocket, in a fully managed way.
	api := slack.New(env.BotToken)

	//create the connection string to be used by apps
	dbConnection := dbAccess{
		conString: getConnectionString(env),
	}

	dbConnection.prepareDatabase()

	slackListener := &SlackListener{
		api:       api,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	//start listening to events on the slack channel and pass the dbconnection struct
	// to be used later when we want to make connection to database
	go slackListener.ListenAndResponse(&dbConnection)

	mux := http.NewServeMux()
	mux.HandleFunc("/learnAmharic", learnAmharic)

	http.ListenAndServe(":5000", mux)
}

//get the connectionstring from env variables
func getConnectionString(env envConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env.BotDBUser,
		env.BotDBUserPassword,
		env.DBAddress,
		env.DBPort,
		env.DBName)
}

func learnAmharic(writer http.ResponseWriter, request *http.Request) {

	slackMessage := slackResponse{}

	err := json.NewDecoder(request.Body).Decode(&slackMessage)
	if err != nil {
		fmt.Println("The marshallong has some issue")
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte(slackMessage.Challenge))
}
