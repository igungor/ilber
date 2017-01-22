package bot

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var (
	currencyBucket = []byte("currency")
)

type Store struct {
	path string
	db   *bolt.DB
}

func NewStore(path string) *Store {
	if path == "" {
		usr, _ := user.Current()
		path = filepath.Join(usr.HomeDir, ".ilber")
	}

	_ = os.MkdirAll(path, 0755)
	path = filepath.Join(path, "ilber.db")
	return &Store{path: path}
}

func (s *Store) Open() error {
	db, err := bolt.Open(s.path, 0644, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	s.db = db

	err = s.db.Update(func(tx *bolt.Tx) error {
		currencyBkt, err := tx.CreateBucketIfNotExists(currencyBucket)
		if err != nil {
			return err
		}

		buckets := [][]byte{
			[]byte("dollar"),
			[]byte("euro"),
		}

		for _, bucket := range buckets {
			_, err = currencyBkt.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Path() string { return s.path }

func (s *Store) SaveCurrency(currency string, value float64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(currencyBucket).Bucket([]byte(currency))

		k := time.Now().UTC().Format(time.RFC3339)
		v := fmt.Sprintf("%.4f", value)

		return bucket.Put([]byte(k), []byte(v))
	})
}

type TimeCurrencyPair struct {
	T time.Time
	V float64
}

func (s *Store) Values(currency string) ([]TimeCurrencyPair, error) {
	var r []TimeCurrencyPair
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(currencyBucket).Bucket([]byte(currency))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			t, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}

			v, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				return err
			}

			pair := TimeCurrencyPair{T: t, V: v}
			r = append(r, pair)
		}

		return nil
	})

	return r, err
}
