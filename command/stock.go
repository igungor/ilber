package command

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
)

func init() {
	register(cmdStock)
}

var cmdStock = &Command{
	Name:      "kur",
	ShortLine: "kurlar ne alemde?",
	Run:       runStock,
}

var (
	defaultCurrencies = []string{"USD", "EUR"}
	financeURL        = "http://finance.yahoo.com/d/quotes.csv?e=.csv&f=c4l1"
)

func runStock(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	var currencies []string
	if args != nil {
		currencies = args
	} else {
		currencies = defaultCurrencies
	}

	u, err := url.Parse(financeURL)
	if err != nil {
		log.Printf("[stock] Error while parsing url '%v'. Err: %v", financeURL, err)
		return
	}

	var qs []string
	for i, currency := range currencies {
		currencies[i] = strings.ToUpper(currency)
		qs = append(qs, currency+"TRY=X,")
	}
	params := u.Query()
	params.Set("s", strings.Join(qs, ""))
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("[stock] Error while fetching stock information. Err: %v", err)
		return
	}
	defer resp.Body.Close()

	cr := csv.NewReader(resp.Body)
	records, err := cr.ReadAll()
	if err != nil {
		log.Printf("[stock] Error while parsing stock information. Err: %v", err)
		return
	}

	if len(records) != len(currencies) {
		b.SendMessage(msg.Chat, "verdigin kurlardan biri ya da birkaci hatali", tlbot.ModeNone, false, nil)
		return
	}

	var buf bytes.Buffer
	for i, record := range records {
		buf.WriteString(fmt.Sprintf("%v = %v ₺\n", currencies[i], record[1]))
	}
	b.SendMessage(msg.Chat, buf.String(), tlbot.ModeNone, false, nil)
}
