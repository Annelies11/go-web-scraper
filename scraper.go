package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// initializing a data structure to keep the scraped data
type Article struct {
	Title, Url, Lang string
}

type URL struct {
	ID       int
	Url      string
	Category string
}

func main() {
	res, err := http.Get("https://kumparan.com/channel/bola-sports")
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

	file, err := os.Create("kumparan-urls-bola-sports.csv")
	if err != nil {
		log.Fatalln("Failed making csv : ", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"ID", "Url", "Category"})

	//rows := make([]Article, 0)
	rows := make([]URL, 0)

	maxScrap := 5

	doc.Find(".sc-5mlv5q-0").Children().Each(func(i int, sel *goquery.Selection) {
		// time.Sleep(1 * time.Second)
		if i > maxScrap-1 {
			return
		}
		row := new(URL)
		row.ID = i + 1
		row.Url, _ = sel.Find(".Viewweb__StyledView-sc-b0snvl-0 a").Attr("href")
		row.Category = sel.Find(".Textweb__StyledText-sc-1ed9ao-0").Text()
		rows = append(rows, *row)
		writer.Write([]string{strconv.Itoa(row.ID), "https://kumparan.com" + row.Url, row.Category})

	})
	defer writer.Flush()
	bts, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	log.Println(string(bts))
	//fmt.Println(string(bts))

}
