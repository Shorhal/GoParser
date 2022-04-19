package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*func parseToOrg(toParse string) []Org {
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
		if err != nil {
			fmt.Println("Error")
		}
		ogrn, err := strconv.Atoi(tempOrg[2])
		if err != nil {
			fmt.Println("Error")
		}
		organization := Org{
			tempOrg[0],
			tempOrg[0],
			inn,
			ogrn,
			kpp,
			rating,
			tempOrg[3],
			tempOrg[3],
			tempOrg[3],
			tempOrg[3],
			tempOrg[3],
			Index,
			tempOrg[3],
			tempOrg[3],
			tempOrg[3],
		}
		Orgs = append(Orgs, organization)

	})
	return Orgs
}*/

func parseToProd(toParse string) []Prod {
	doc := getHtmlDocumentReader(toParse)

	var Prods []Prod
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		fmt.Println("----------tr----------")
		tempOrg := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			text := strings.TrimSpace(td.Text())
			return text
		})
		fmt.Println(tempOrg[0])

		OKPD2, err := strconv.Atoi(tempOrg[2])
		if err != nil {
			fmt.Println("Error")
		}
		TNVED, err := strconv.Atoi(tempOrg[3])
		if err != nil {
			fmt.Println("Error")
		}
		point, err := strconv.Atoi(tempOrg[5])
		if err != nil {
			fmt.Println("Error")
		}
		product := Prod{
			tempOrg[0],
			tempOrg[1],
			OKPD2,
			TNVED,
			tempOrg[4],
			point,
		}
		Prods = append(Prods, product)
	})
	return Prods
}

// Функция получения ссылок на карточку предприятия а также их продукции
func getURLs(toParse string) []URL {
	doc := getHtmlDocumentReader(toParse)

	//var stri []string
	var URLs []URL
	doc.Find("td").Each(func(i int, td *goquery.Selection) {
		var stri []string
		td.Find("a").Each(func(i int, a *goquery.Selection) {
			str, ok := a.Attr("href")
			if ok {
				fmt.Println("")
			}
			stri = append(stri, str)
		})

		if stri != nil {
			URLs = append(URLs, URL{
				stri[0],
				stri[1],
			})
		}

	})

	return URLs
}

func getHtmlDocumentReader(toRead string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(toRead))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}
	return doc
}
