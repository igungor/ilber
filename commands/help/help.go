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

bugunkandilmi - is it candle?
vizyon - sinema felan
hava - nem fena nem
yo - yigit ozgur seysi
img - resim filan ara
caps - incicaps gibi degil gibi
imdb - ayemdiibii
tatil - ne zaman
benkimim - ilber!
echo - cok cahilsin
`
}
