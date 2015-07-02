package help

import (
	"github.com/igungor/ilberbot"
)

func init() {
	ilberbot.RegisterCommand("/help", help)
}

func help(args ...string) string {
	return `
sunlar var:

iftar - iftar vakti
sahur - sahur vakti
okundumu - is it read?
bugunkandilmi - is it candle?
vizyon - sinema felan
hava - nem fena nem
yo - yigit ozgur seysi
tatil - ne zaman
benkimim - ilber!
echo - cok cahilsin
`
}
