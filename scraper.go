package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	res, err := http.Get("https://kumparan.com/kumparanoto")
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
	//time.Sleep(10 * time.Second)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	file, err := os.Create("kumparan-urls-oto.csv")
	if err != nil {
		log.Fatalln("Failed making csv : ", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"ID", "Url", "Category"})

	//rows := make([]Article, 0)
	rows := make([]URL, 0)

	maxScrap := 20
	j := 1
	//Viewweb__StyledView-sc-b0snvl-0 bASEUm

	doc.Find(".Viewweb__StyledView-sc-b0snvl-0.emkTMh").Children().Each(func(i int, sel *goquery.Selection) {
		time.Sleep(1 * time.Second)
		if j > maxScrap {
			return
		}
		row := new(URL)
		row.ID = j
		//Viewweb__StyledView-sc-b0snvl-0 BXiYe
		row.Url, _ = sel.Find(".Viewweb__StyledView-sc-b0snvl-0.BXiYe a").Attr("href")
		row.Category = sel.Find(".TextBoxweb__StyledTextBox-sc-1noultz-0 span").Text()
		rows = append(rows, *row)
		if row.Url != "" {
			writer.Write([]string{strconv.Itoa(row.ID), "https://kumparan.com" + row.Url, row.Category})
			j++
		}

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
