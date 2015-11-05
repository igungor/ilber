package command

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

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
	movieCache      = map[string][]string{}
)

func runMovies(b *tlbot.Bot, msg *tlbot.Message) {
	movies, err := fetchOrCache()
	if err != nil {
		log.Printf("[movies] Error while fetching movies: %v\n", err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString(" 🎬 İstanbul'da vizyon filmleri\n")
	for _, movie := range movies {
		buf.WriteString(fmt.Sprintf("🔸 %v\n", movie))
	}

	err = b.SendMessage(msg.Chat.ID, buf.String(), tlbot.ModeNone, false, nil)
	if err != nil {
		log.Printf("[movies] Error while sending message: %v\n", err)
		return
	}
}

func fetchOrCache() ([]string, error) {
	now := time.Now().UTC()
	year, week := now.ISOWeek()
	// YYYYWW is our cache key. Theaters keep their movies for about a week. We
	// don't need a fresh movie list every day or hour. Using year and ISO week
	// gives us the convenience to avoid cache invalidation. everybody hates
	// cache invalidation. thank you ISO week.
	nowstr := fmt.Sprintf("%v%v", year, week)

	// friday nights and saturday mornings are the times theaters renew their
	// movie list. fetching new list on these days are a waste. just go to
	// cache.
	if now.Weekday() > time.Thursday {
		movies, ok := movieCache[nowstr]
		if !ok {
			return nil, fmt.Errorf("unfortunately today is friday/saturday and the cache is empty. nothing to do")
		}
		return movies, nil
	}

	movies, ok := movieCache[nowstr]
	if ok {
		return movies, nil
	}

	// cache-miss. fetch a fresh list.
	movies = fetchMovies()
	if movies == nil {
		return nil, fmt.Errorf("fetched new movies but the list came empty")
	}

	// put the new list in cache
	movieCache[nowstr] = movies

	return movies, nil
}

func fetchMovies() []string {
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
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		offset := strconv.Itoa(10 * i)
		wg.Add(1)
		go fetch(movieURL+offset, &wg)
	}
	wg.Wait()

	// sort map by its values. map values contain frequency of a movie by
	// theater count. most frequent movie in a theater is most probably
	// screened near the caller's neighborhood.
	vs := newValSorter(movies)
	sort.Sort(vs)

	return vs.Keys
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
