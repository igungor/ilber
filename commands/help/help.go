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

yo - yigit ozgur seysi
img - resim filan ara
vizyon - sinema felan
hava - nem fena nem
bugunkandilmi - is it candle?
caps - incicaps gibi degil gibi
imdb - ayemdiibii
tatil - ne zaman
echo - cok cahilsin
`
}
