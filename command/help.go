package command

import (
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdHelp)
}

var cmdHelp = &Command{
	Name: "help",
	Run:  runHelp,
}

func runHelp(b *tlbot.Bot, msg *tlbot.Message) {
	b.SendMessage(msg.From, helpmsg, tlbot.ModeMarkdown, false, nil)
}

var helpmsg = `sunlar var:

*yo* - yigit ozgur seysi
*img* - resim filan ara
*vizyon* - sinema felan
*hava* - nem fena nem
*bugunkandilmi* - is it candle?
*imdb* - ayemdiibii
*tatil* - ne zaman
*echo* - cok cahilsin
`
