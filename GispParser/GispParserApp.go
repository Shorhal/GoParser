package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func main() {

	//connect()
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	emulation.SetUserAgentOverride("WebScrapper 1.0")

	url := "https://gisp.gov.ru/pp719v2/pub/org/"

	var pagesInfo string
	var pagesCount int
	var dataGrid string
	var allDataGrid string

	if err := chromedp.Run(ctx,

		chromedp.Navigate(url),
		//chromedp.Click("dx-page-size dx-selection", chromedp.BySearch),
		chromedp.SetAttributeValue(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, "style", "display: block", chromedp.ByID),
		chromedp.Text(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, &pagesInfo),
		chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table > tbody`, &dataGrid, chromedp.ByID),
	); err != nil {

		log.Fatal(err)
	}

	pagesCount, err := strconv.Atoi(strings.Split(pagesInfo, " ")[3])
	if err == nil {
		fmt.Println("Error")
	}

	for i := 0; i < pagesCount; i++ {
		if err := chromedp.Run(ctx,
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &dataGrid, chromedp.ByID),
		); err != nil {

			log.Fatal(err)
		}
		allDataGrid += dataGrid
	}
	//fmt.Println(allDataGrid)
	fmt.Println(pagesCount)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(allDataGrid))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

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

	fmt.Println(Orgs)

}
