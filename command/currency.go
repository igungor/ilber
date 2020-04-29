package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

const (
	yahooFinanceURL = "https://query1.finance.yahoo.com/v8/finance/chart/"
)

func defaultCurrencies() []query {
	return []query{
		{
			from: "USD",
			to:   "TRY",
		},
		{
			from: "EUR",
			to:   "TRY",
		},
	}
}

func runCurrency(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	queries := parseMessage(msg.Args())

	var errs []error
	for _, fn := range []queryfunc{queryYahooFinance} {
		s, err := fn(queries)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		_, _ = b.SendMessage(msg.Chat.ID, s, telegram.WithParseMode(telegram.ModeMarkdown))
		break
	}

	if len(errs) > 0 {
		b.Logger.Printf("Error parsing query: %q\n", errs)
		_, _ = b.SendMessage(msg.Chat.ID, "birtakım hatalar sözkonusu")
	}
}

func maybeNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 32)
	return err == nil
}

func parseMessage(terms []string) []query {
	if terms == nil {
		return defaultCurrencies()
	}

	// parse the calculation. eg: "4 USD in TRY"
	//
	// 4 term queries are calculations, such as:
	// "3 USD in TRY" or
	// "5 BTC to TRY"
	if len(terms) == 4 && maybeNumber(terms[0]) {
		amount, _ := strconv.ParseFloat(terms[0], 32)
		from := terms[1]
		to := terms[3]

		return []query{
			{
				amount: amount,
				from:   from,
				to:     to,
				isCalc: true,
			},
		}
	}

	var currencies []string
	for _, term := range terms {
		if len(term) != 3 {
			continue
		}
		currencies = append(currencies, term)
	}
	var queries []query
	for _, currency := range currencies {
		query := query{
			amount: 1,
			from:   currency,
			to:     "TRY",
		}
		queries = append(queries, query)
	}
	return queries
}

type query struct {
	amount float64
	from   string
	to     string
	isCalc bool
}

type queryfunc func([]query) (string, error)

func queryYahooFinance(queries []query) (string, error) {
	if len(queries) == 0 {
		return "", fmt.Errorf("yahoo: no query found")
	}

	request := func(q query) (float64, error) {
		u, _ := url.Parse(yahooFinanceURL)
		u.Path += fmt.Sprintf("%v%v=%v", q.from, q.to, "X")
		params := u.Query()
		params.Set("range", "1d")
		u.RawQuery = params.Encode()

		resp, err := httpclient.Get(u.String())
		if err != nil {
			return 0, fmt.Errorf("yahoo: could not fetch response: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("yahoo: unexpected status code %v", resp.StatusCode)
		}

		var response yahooFinanceResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return 0, fmt.Errorf("yahoo: could not parse json: %v", err)
		}

		result := response.Chart.Result
		quote := result[len(result)-1].Indicators.Quote
		close := quote[len(quote)-1].Close
		prevClose := result[len(result)-1].Meta.PreviousClose

		if len(close) == 0 {
			// some currencies are not available in 'close data', such as BGN.
			// dont let people down.
			if prevClose != 0 {
				return prevClose, nil
			}
			return 0, fmt.Errorf("yahoo: no value found for %v", q)
		}

		var rates []float64
		for _, v := range close {
			rate, ok := v.(float64)
			// skip unrecognized values to a list for later use
			if !ok {
				continue
			}
			rates = append(rates, rate)
		}
		return rates[len(rates)-1], nil
	}

	if len(queries) == 1 && queries[0].isCalc {
		q := queries[0]
		rate, err := request(q)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"%4.2f %v = %4.4f %v",
			q.amount,
			normalize(q.from),
			q.amount*rate,
			normalize(q.to),
		), nil
	}

	var rates []float64
	for _, q := range queries {
		rate, err := request(q)
		if err != nil {
			return "", err
		}
		rates = append(rates, rate)
	}

	var buf bytes.Buffer
	for i, rate := range rates {
		fmt.Fprintf(&buf, "%v = %4.4f %v\n",
			normalize(queries[i].from),
			rate,
			normalize(queries[i].to),
		)
	}
	return buf.String(), nil
}

type yahooFinanceResponse struct {
	Chart struct {
		Error  interface{} `json:"error"`
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close  []interface{} `json:"close"`
					High   []interface{} `json:"high"`
					Low    []interface{} `json:"low"`
					Open   []interface{} `json:"open"`
					Volume []interface{} `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
			Meta struct {
				PreviousClose float64 `json:"previousClose"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
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
