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
	//"gorm.io/driver/sqlserver"
	//"gorm.io/gorm"
)

func main() {
	/*var db *gorm.DB

	dsn := "sqlserver://@localhost:52876?database=Gisp"
	//Соединение с БД
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("filed connection")
	} else {
		fmt.Println("Success connect to MSSQL")
	}*/

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

	// Полученние ссылок на карточки предсприятий и их продукцию
	var pagesInfo string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.SetAttributeValue(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, "style", "display: block", chromedp.ByID),
		chromedp.Text(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, &pagesInfo),
	); err != nil {

		log.Fatal(err)
	}

	//Получение количества страниц
	var tableData string
	pagesCount, err := strconv.Atoi(strings.Split(pagesInfo, " ")[3])
	if err != nil {
		fmt.Println("Error")
	}

	//Получение разметки всех страниц
	for i := 1; i <= pagesCount; i++ {
		var tempContainer string
		xpath := `//*[@aria-label="Page ` + strconv.Itoa(i) + `"]`
		if err := chromedp.Run(ctx,
			chromedp.Click(xpath),
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &tempContainer, chromedp.ByID),
		); err != nil {
			log.Fatal(err)
		}
		tableData += tempContainer
	}

	//Получение ссылок на карточки предприятий и список их продукции
	UrlData := getURLs(tableData)

	var OrgString string
	var ProdsString string
	URLProds := "https://gisp.gov.ru" + UrlData[0].Prods
	if err := chromedp.Run(ctx,
		chromedp.Navigate(UrlData[0].Org),
		chromedp.OuterHTML(`body > main > div > div.content__inner > div > div:nth-child(2) > div > div > div`, &OrgString),
		chromedp.Navigate(URLProds),
		chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &ProdsString),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println(ProdsString)

	//Org := parseToOrg(dataGrid)
	//Prod := parseToProd(prodDataGrid)

	/*for _, element := range Org {
		db.Create(element)
	}
	*/
}

//Получение разметки-Информация об организаии и ее продукции
