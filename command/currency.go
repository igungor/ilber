package command

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
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

func runCurrency(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	opts := &telegram.SendOptions{}
	s, err := parseQuery(msg.Args())
	if err != nil {
		log.Printf("Error parsing query: %v\n", err)
		_, _ = b.SendMessage(msg.Chat.ID, "birtakım hatalar sözkonusu", opts)
		return
	}

	_, err = b.SendMessage(msg.Chat.ID, s, opts)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
	}
}

func parseQuery(terms []string) (string, error) {
	if terms == nil {
		terms = defaultCurrencies
	}

	u, _ := url.Parse(financeURL)
	params := u.Query()

	var isQuestion bool
	var f float64

	// check if the query is a calculation statement (contains 4 tokens):
	// x EUR in TRY
	// x dollars in pounds
	// x dolar kaç lira
	if len(terms) == 4 {
		var err error
		f, err = strconv.ParseFloat(terms[0], 32)
		if err == nil {
			isQuestion = true
		}
	}

	currencies := make([]string, len(terms))
	if isQuestion {
		currencies[0], currencies[1] = normalize(terms[1]), normalize(terms[3])
		params.Set("s", fmt.Sprintf("%v%v=X", currencies[0], currencies[1]))
	} else {
		// query string be like: USDTRY=X,EURTRY=X
		var qs []string
		for i, cur := range terms {
			cur = normalize(cur)
			currencies[i] = cur
			qs = append(qs, cur+"TRY=X")
		}
		params.Set("s", strings.Join(qs, ","))
	}

	u.RawQuery = params.Encode()
	resp, err := httpclient.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("error fetching currencies: %v", err)
	}
	defer resp.Body.Close()

	cr := csv.NewReader(resp.Body)
	records, err := cr.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error reading csv: %v", err)
	}

	var buf bytes.Buffer
	if isQuestion {
		rate, err := strconv.ParseFloat(records[0][1], 32)
		if err != nil {
			return "", fmt.Errorf("error reading record as float: %v", err)
		} else {
			r := f * rate
			buf.WriteString(fmt.Sprintf("%.2f %v = %.2f %v\n", f, currencies[0], r, currencies[1]))
		}
		return buf.String(), nil
	}

	for i, record := range records {
		buf.WriteString(fmt.Sprintf("%v = %v ₺\n", currencies[i], record[1]))
	}

	return buf.String(), nil
}

// popular currencies
var m = map[string]string{
	"₺":     "TRY",
	"tl":    "TRY",
	"lira":  "TRY",
	"liras": "TRY",

	"$":       "USD",
	"dolar":   "USD",
	"dollar":  "USD",
	"dollars": "USD",
	"dolares": "USD",

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
