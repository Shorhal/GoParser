package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	fmt.Println(len(parseToOrg(allDataGrid)))

}
