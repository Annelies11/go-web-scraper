package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// initializing a data structure to keep the scraped data
type Article struct {
	Title, Url, Lang string
}

func main() {
	res, err := http.Get("https://blog.saugi.me/")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	} else {
		fmt.Println("Berhasil!")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("Status code error : %d %s", res.StatusCode, res.Status)
		log.Fatalf("Status code error : %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	file, err := os.Create("lists.csv")
	if err != nil {
		log.Fatalln("Failed making csv : ", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"Judul", "Lang", "Url"})
	rows := make([]Article, 0)
	doc.Find(".card").Children().Each(func(i int, sel *goquery.Selection) {
		row := new(Article)
		row.Title = sel.Find(".card-title-sm").Text()
		row.Url, _ = sel.Find(".card-title-sm span a").Attr("href")
		row.Lang = sel.Find(".lang").Text()
		rows = append(rows, *row)
		writer.Write([]string{row.Title, row.Lang, row.Url})
	})
	defer writer.Flush()
	bts, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	log.Println(string(bts))
	fmt.Println(string(bts))

}
