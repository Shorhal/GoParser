package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseToOrg(toParse string) []Org {
	doc := getHtmlDocumentReader(toParse)

	var Orgs []Org
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		fmt.Println("----------tr----------")
		tempOrg := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			text := strings.TrimSpace(td.Text())
			return text
		})
		fmt.Println(tempOrg[0])

		inn, err := strconv.Atoi(tempOrg[1])
		if err == nil {
			fmt.Println("Error")
		}
		ogrn, err := strconv.Atoi(tempOrg[2])
		if err == nil {
			fmt.Println("Error")
		}
		organization := Org{
			tempOrg[0],
			inn,
			ogrn,
			tempOrg[3],
		}
		Orgs = append(Orgs, organization)

	})
	return Orgs
}

func parseToProd(toParse string) {
	doc := getHtmlDocumentReader(toParse)

	var Prod []Org
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		fmt.Println("----------tr----------")
		tempOrg := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			text := strings.TrimSpace(td.Text())
			return text
		})
		fmt.Println(tempOrg[0])

		inn, err := strconv.Atoi(tempOrg[1])
		if err == nil {
			fmt.Println("Error")
		}
		ogrn, err := strconv.Atoi(tempOrg[2])
		if err == nil {
			fmt.Println("Error")
		}
		organization := Org{
			tempOrg[0],
			inn,
			ogrn,
			tempOrg[3],
		}
		Prod = append(Prod, organization)
	})
}

func getHtmlDocumentReader(toRead string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(toRead))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}
	return doc
}
