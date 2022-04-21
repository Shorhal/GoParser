package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseToOrg(toParse string) Org {
	doc := getHtmlDocumentReader(toParse)

	var OrgModel Org

	doc.Find("p").Each(func(i int, p *goquery.Selection) {
		title := p.Find("span").Text()
		value := strings.Split(p.Text(), ":")
		fmt.Println(value[1])
		fmt.Println(p.Text())

		switch title {
		case "Полное наименование предприятия:":
			OrgModel.FullName = value[1]
		case "Сокращенное наименование предприятия:":
			OrgModel.ShortName = value[1]
		case "ОГРН:":
			ogrn, err := strconv.Atoi(value[1])
			if err == nil {
				fmt.Println("Error", value[1])
			}
			OrgModel.OGRN = ogrn
		case "ИНН:":
			inn, err := strconv.Atoi(value[1])
			if err == nil {
				fmt.Println("Error", value[1])
			}
			OrgModel.INN = inn
		case "КПП:":
			kpp, err := strconv.Atoi(value[1])
			if err == nil {
				fmt.Println("Error", value[1])
			}

			OrgModel.KPP = kpp
		case "ОКВЭД 2:":
			OrgModel.INDUSTRY = p.Text()
		case "Страна":
			OrgModel.Country = value[1]
		case "Регион:":
			OrgModel.Region = value[1]
		case "Город:":
			OrgModel.City = value[1]
		case "Адрес:":
			OrgModel.Adress = value[1]
		case "Индес:":
			index, err := strconv.Atoi(value[1])
			if err == nil {
				fmt.Println("Error", value[1])
			}
			OrgModel.Index = index
		case "www:":
			OrgModel.www = value[1]
		}
	})

	rate, err := strconv.Atoi(doc.Find(".value").Text())
	if err == nil {
		fmt.Println("Error", doc.Find(".value").Text())
	}

	OrgModel.Rating = rate

	return OrgModel
}

func parseToProd(toParse string) []Prod {
	doc := getHtmlDocumentReader(toParse)

	var Prods []Prod
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		tempOrg := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			text := strings.TrimSpace(td.Text())
			return text
		})

		//convert to int
		OKPD2, err := strconv.Atoi(tempOrg[2])
		if err != nil {
			fmt.Println("Error -", tempOrg[2])
		}
		TNVED, err := strconv.Atoi(tempOrg[3])
		if err != nil {
			fmt.Println("Error -", tempOrg[3])
		}
		point, err := strconv.Atoi(tempOrg[5])
		if err != nil {
			fmt.Println("Error -", tempOrg[5])
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
			if !ok {
				fmt.Println("Cannot get url")
			}
			stri = append(stri, str)
		})

		if stri != nil {
			URLs = append(URLs, URL{
				strings.TrimSpace(stri[1]),
				strings.TrimSpace(stri[0]),
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
