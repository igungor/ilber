package whoami

import (
	"math/rand"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/benkimim", whoami)
}

var answers = []string{
	"benim lan, ilber :)",
	"ilber",
	"ilbert",
	"ilberto",
	"ilberto garcia de la marquez",
	"dilbert",
	"cahil cahil konusma",
}

func whoami(args ...string) string {
	return answers[rand.Intn(len(answers))]
}
