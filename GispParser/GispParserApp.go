package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB

	dsn := "sqlserver://@localhost:52876?database=Gisp"
	//Соединение с БД
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("filed connection")
	} else {
		fmt.Println("Success connect to MSSQL")
	}

	//Подготовка
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	emulation.SetUserAgentOverride("WebScrapper 1.0")

	url := "https://gisp.gov.ru/pp719v2/pub/org/"

	var pagesInfo string
	var dataGrid string
	//var allDataGrid string
	//Переход на главную страницу
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.SetAttributeValue(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, "style", "display: block", chromedp.ByID),
		chromedp.Text(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, &pagesInfo),
		chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &dataGrid, chromedp.ByID),
	); err != nil {

		log.Fatal(err)
	}

	//pagesCount, err := strconv.Atoi(strings.Split(pagesInfo, " ")[3])
	if err != nil {
		fmt.Println("Error")
	}

	//Сбор данных со страниц
	/*for i := 1; i <= pagesCount; i++ {

		if err := chromedp.Run(ctx,
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &dataGrid, chromedp.ByID),
		); err != nil {

			log.Fatal(err)
		}
		allDataGrid += dataGrid
	}*/

	Org := parseToOrg(dataGrid)
	//Prod := parseToProd(prodDataGrid)

	for _, element := range Org {
		db.Create(element)
	}
}
