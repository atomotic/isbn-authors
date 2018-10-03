package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/knakk/sparql"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const oclcAPI = "http://classify.oclc.org/classify2/Classify?isbn="

type Book struct {
	Work     Work     `xml:"work" json:"book"`
	Response Response `xml:"response" json:"-"`
	Authors  []Author `xml:"authors>author" json:"authors"`
}
type Work struct {
	Title string `xml:"title,attr" json:"title"`
}
type Response struct {
	Code int `xml:"code,attr" json:"-"`
}
type Author struct {
	Viaf     string `xml:"viaf,attr" json:"viaf"`
	Wikidata string `json:"wikidata"`
	Text     string `xml:",chardata" json:"name"`
}

func getBookAuthors(isbn string) (Book, error) {
	var book Book

	resp, err := http.Get(fmt.Sprintf("%s%s", oclcAPI, isbn))
	if err != nil {
		return book, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = xml.Unmarshal([]byte(body), &book)
	if err != nil {
		return book, err
	}

	if book.Response.Code == 0 || book.Response.Code == 2 || book.Response.Code == 4 {
		for i, author := range book.Authors {
			wd, err := sparql.NewRepo("https://query.wikidata.org/sparql")
			if err != nil {
				return book, err
			}
			query := fmt.Sprintf("SELECT ?q WHERE { ?q wdt:P214 '%s'.}", author.Viaf)
			res, err := wd.Query(query)
			if err != nil {
				return book, err
			}
			if len(res.Results.Bindings) > 0 {
				Q := strings.Split(res.Results.Bindings[0]["q"].Value, "/")[4]
				book.Authors[i].Wikidata = Q
			}
		}
		return book, nil
	} else {
		return book, errors.New("0 results.")
	}

}

func IsbnAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("books"))

		cached := b.Get([]byte(isbn))

		if cached != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", cached)
		} else {
			book, err := getBookAuthors(isbn)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "<book not found>")
				return nil
			}

			bookJSON, _ := json.Marshal(book)
			err = b.Put([]byte(isbn), []byte(bookJSON))
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", bookJSON)
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	db, err = bolt.Open("isbn-authors.cache", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("books"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	r := mux.NewRouter()
	r.HandleFunc("/isbn/{isbn}", IsbnAPI)
	fmt.Println("server running on http://localhost:8093/isbn/{$ISBN}")

	http.ListenAndServe(":8093", r)

}
