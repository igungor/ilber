package command

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/igungor/tlbot"
	"golang.org/x/net/context"
)

func init() {
	register(cmdCurrency)
}

var cmdCurrency = &Command{
	Name:      "kur",
	ShortLine: "kurlar ne alemde?",
	Run:       runCurrency,
}

var defaultCurrencies = []string{"USD", "EUR"}

const financeURL = "http://finance.yahoo.com/d/quotes.csv?e=.csv&f=c4l1"

func runCurrency(ctx context.Context, b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()

	var currencies []string
	if args != nil {
		currencies = args
	} else {
		currencies = defaultCurrencies
	}

	u, err := url.Parse(financeURL)
	if err != nil {
		log.Printf("Error while parsing url '%v'. Err: %v", financeURL, err)
		return
	}

	var qs []string
	for i, currency := range currencies {
		currency = normalize(currency)
		currencies[i] = currency
		qs = append(qs, currency+"TRY=X")
	}

	// query string be like: USDTRY=X,EURTRY=X
	params := u.Query()
	params.Set("s", strings.Join(qs, ","))
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("Error while fetching currency information. Err: %v", err)
		return
	}
	defer resp.Body.Close()

	cr := csv.NewReader(resp.Body)
	records, err := cr.ReadAll()
	if err != nil {
		log.Printf("Error while parsing currency information. Err: %v", err)
		return
	}

	if len(records) != len(currencies) {
		err := b.SendMessage(msg.Chat.ID, "verdiğin kurlardan biri ya da birkaçı hatalı", tlbot.ModeNone, false, nil)
		if err != nil {
			log.Printf("Error while sending message. Err: %v\n", err)
		}
		return
	}

	var buf bytes.Buffer
	for i, record := range records {
		buf.WriteString(fmt.Sprintf("%v = %v ₺\n", currencies[i], record[1]))
	}

	err = b.SendMessage(msg.Chat.ID, buf.String(), tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
	}
}

// popular currencies
var m = map[string]string{
	"₺":     "TRY",
	"lira":  "TRY",
	"liras": "TRY",

	"$":       "USD",
	"dolar":   "USD",
	"dollar":  "USD",
	"dollars": "USD",

	"€":     "EUR",
	"euro":  "EUR",
	"euros": "EUR",
	"yuro":  "EUR",
	"avro":  "EUR",

	"£":        "GBP",
	"pound":    "GBP",
	"pounds":   "GBP",
	"sterlin":  "GBP",
	"sterlins": "GBP",
}

func normalize(s string) string {
	sl := strings.ToLower(s)
	currency, ok := m[sl]
	if !ok {
		return strings.ToUpper(s)
	}
	return currency
}
