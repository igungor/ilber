package ray

import "github.com/igungor/ilberbot"

func init() {
	ilberbot.RegisterCommand("/ray", ray)
}

func ray(args ...string) string {
	return "malifalitiko!"
}
