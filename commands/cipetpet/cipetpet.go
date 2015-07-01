package cipetpet

import (
	"math/rand"

	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/cipetpet", cipetpet)
}

var cirling = []string{
	"cipetpetpetpet",
	"cipetpet cipetpetpetpet",
	"cibili cibili sak sak",
	"cibili cibili cibili sak sak sak",
	"bunlar cipetpet ama *lezzetli* cipetpet degil!",
	"hadi yavrum bi daha vur sak sak sak",
	"cirling cirling prrr sak sak sak",
	"bicinik bicinik aniyyaaa",
	"v harfinden kus 4 turlu on yapar: velis velis velis baska, veste veste veste baska, vais vais vais baska, ves ves ves baska.",
	"k harfinden: kis kis kis baska, kah kah baska, kaf kaf baska, kiya kiya baska",
	"cipii cipii cipii sak sak sak vicoooo kis kis kis aniyyaa aniyyaa aniyyaa kiya kiya kiya",
	"cibini picuviiii sak sak vicooo",
	"o onden yaptigi cipii cipiii sak sak sak vicoo'ya baglamiyo, o kaba saksagi daha taneli vuruyo",
	"cipii cipiii sak sak sak sak sak sak. bu da eklemeli kaba saksak",
	"vesle vesle vesle cececece. bozuk, eziyo.",
}

func cipetpet(args ...string) string {
	return cirling[rand.Intn(len(cirling))]
}
