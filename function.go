package ilber

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/ilber/command"
	"github.com/igungor/telegram"
)

var ilber *bot.Bot

func init() {
	logger := log.New(os.Stdout, "ilber: ", log.LstdFlags|log.Lshortfile)
	var err error
	ilber, err = bot.New(logger)
	if err != nil {
		logger.Fatalf("Could not initialize the bot: %v\n", err)
	}
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(http.StatusOK)

	var u telegram.Update
	_ = json.NewDecoder(r.Body).Decode(&u)

	msg := &u.Message

	if msg.IsService() {
		ilber.Logger.Printf("incoming service message: %v", msg)
		return
	}

	cmdname := msg.Command()
	if cmdname == "" {
		ilber.Logger.Printf("no command found from message: %v", msg)
		return
	}

	// is the command even registered?
	cmd := command.Lookup(cmdname)
	if cmd == nil {
		ilber.Logger.Printf("unregistered command %v: %v", cmdname, msg)
		return
	}

	cmd.Run(r.Context(), ilber, msg)
}
