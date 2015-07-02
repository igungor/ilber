package echo

import (
	"strings"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/echo", echo)
}

func echo(args ...string) string {
	return strings.Join(args, " ")
}
