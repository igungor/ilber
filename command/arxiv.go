package command

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aaron-lebo/ocd/feeds/atom"
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdArxiv)
}

var cmdArxiv = &Command{
	Name:      "arxiv",
	ShortLine: "arxiv ne arar la bazarda",
	Hidden:    true,
	Run:       runArxiv,
}

var (
	arxivURL = "http://export.arxiv.org/api/query"
)

func runArxiv(b *tlbot.Bot, msg *tlbot.Message) {
	args := msg.Args()
	if len(args) == 0 {
		err := b.SendMessage(msg.Chat.ID, "boş geçmeyelim 💩", tlbot.ModeNone, false, nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	u, err := url.Parse(arxivURL)
	if err != nil {
		log.Printf("Error while parsing url '%v'. Err: %v", arxivURL, err)
		return
	}

	qs := strings.Join(args, " ")
	params := u.Query()
	params.Set("search_query", qs)
	params.Set("max_results", "1")
	u.RawQuery = params.Encode()

	resp, err := httpclient.Get(u.String())
	if err != nil {
		log.Printf("Error while fetching arxiv document. Err: %v", err)
		return
	}
	defer resp.Body.Close()

	var result atom.Feed
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Error while decoding the response: %v", err)
		return
	}

	if len(result.Entries) == 0 {
		err := b.SendMessage(msg.Chat.ID, "sonuç boş geldi 👐", tlbot.ModeNone, false, nil)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	entry := result.Entries[0]
	pdflink := "pdf linki yok"
	for _, link := range entry.Links {
		if link.Title == "pdf" {
			pdflink = link.Href
			break
		}
	}
	var categories []string
	for _, category := range entry.Categories {
		categories = append(categories, category.Term)
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*title:* %v\n", entry.Title))
	buf.WriteString(fmt.Sprintf("*categories:* %v\n", strings.Join(categories, ", ")))
	buf.WriteString(fmt.Sprintf("*published:* %v\n", entry.Published.Format("2006-01-02")))
	buf.WriteString(fmt.Sprint("*authors:*\n"))
	for _, author := range entry.Authors {
		buf.WriteString(fmt.Sprintf(" - %v\n", author.Name))
	}
	buf.WriteString(fmt.Sprintf("*pdf:* %v", pdflink))

	err = b.SendMessage(msg.Chat.ID, buf.String(), tlbot.ModeMarkdown, false, nil)
	if err != nil {
		log.Printf("Error while sending message: %v\n", err)
	}
}
