package command

// TODO(ig): implement a cache mechanism.
// Currently we do 3 HTTP calls concurrently for each `/vizyon` call. As a
// side note, theaters refresh their movie list on every Friday night/Saturday
// morning.  So it is better to invalidate caches before Saturday noon and
// fetch a fresh movie list.

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/igungor/tlbot"
)

func init() {
	register(cmdMovies)
}

var cmdMovies = &Command{
	Name:      "vizyon",
	ShortLine: "sinema filan",
	Run:       runMovies,
}

var (
	movieURL        = "http://www.google.com/movies?near=Kadikoy,Istanbul&start="
	chromeUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36"
)

func runMovies(b *tlbot.Bot, msg *tlbot.Message) {
	var wg sync.WaitGroup

	// mu guards movies map access
	var mu sync.Mutex
	movies := make(map[string]int)

	// fetch is a closure to fetch the movies in given movieurl and save the
	// result in movies map.
	fetch := func(movieurl string, wg *sync.WaitGroup) {
		defer wg.Done()

		req, _ := http.NewRequest("GET", movieURL, nil)
		req.Header.Set("User-Agent", chromeUserAgent)

		r, err := httpclient.Do(req)
		if err != nil {
			log.Printf("[movies] Error while fetching URL '%v'. Error: %v\n", movieurl, err)
			return
		}
		defer r.Body.Close()

		doc, err := goquery.NewDocumentFromResponse(r)
		if err != nil {
			log.Printf("[movies] Error while fetching DOM from url '%v': %v\n", movieurl, err)
			return
		}

		doc.Find(".theater .desc .name a").Each(func(_ int, s *goquery.Selection) {
			s.Closest(".theater").Find(".showtimes .name").Each(func(_ int, sel *goquery.Selection) {
				mu.Lock()
				movies[sel.Text()]++
				mu.Unlock()
			})
		})
	}

	// fetch 3 pages of theaters
	for i := 0; i < 3; i++ {
		offset := strconv.Itoa(10 * i)
		wg.Add(1)
		go fetch(movieURL+offset, &wg)
	}
	wg.Wait()

	// sort by map values. map values contain frequency of a movie by
	// theater count. most frequent movie in a theater is most probably
	// screened near the caller's home.
	vs := newValSorter(movies)
	sort.Sort(vs)

	var buf bytes.Buffer
	buf.WriteString(" 🎬 İstanbul'da vizyon filmleri\n")
	for _, movie := range vs.Keys {
		buf.WriteString(fmt.Sprintf("🔸 %v\n", movie))
	}

	err := b.SendMessage(msg.Chat, buf.String(), tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("[movies] Error while sending message. Err: %v\n", err)
		return
	}
}

// valsorter is used for sorting the map by value
type valsorter struct {
	Keys []string
	Vals []int
}

func (v *valsorter) Len() int           { return len(v.Vals) }
func (v *valsorter) Less(i, j int) bool { return v.Vals[i] > v.Vals[j] }
func (v *valsorter) Swap(i, j int) {
	v.Vals[i], v.Vals[j] = v.Vals[j], v.Vals[i]
	v.Keys[i], v.Keys[j] = v.Keys[j], v.Keys[i]
}

func newValSorter(m map[string]int) *valsorter {
	vs := &valsorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}
