package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	//"github.com/jasonlvhit/gocron"
	//"gorm.io/driver/sqlserver"
	//"gorm.io/gorm"
)

func main() {

	//gocron.Every(2).Minutes().Do(UpdateData)
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
	mainCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	mainCtx, cancel = context.WithTimeout(mainCtx, 20*time.Second)
	defer cancel()

	url := "https://gisp.gov.ru/pp719v2/pub/org/"

	// Полученние ссылок на карточки предсприятий и их продукцию
	var pagesInfo string
	if err := chromedp.Run(mainCtx,
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
		if err := chromedp.Run(mainCtx,
			chromedp.Click(xpath),
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &tempContainer, chromedp.ByID),
		); err != nil {
			log.Fatal(err)
		}
		tableData += tempContainer
	}

	//Получение ссылок на карточки предприятий и список их продукции
	UrlData := getURLs(tableData)
	fmt.Println(len(UrlData))

	orgInfoCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	orgInfoCtx, cancel = context.WithTimeout(orgInfoCtx, 50*time.Second)
	defer cancel()

	var Org []Org
	var Prod []Prod
	for _, item := range UrlData {
		var OrgHtmlString string
		var ProdsHtmlString string
		URLProds := "https://gisp.gov.ru" + item.Prods
		if err := chromedp.Run(orgInfoCtx,
			chromedp.Navigate(item.Org),
			chromedp.OuterHTML(`body > main > div > div.content__inner > div > div:nth-child(2) > div > div > div`, &OrgHtmlString),
			chromedp.Navigate(URLProds),
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &ProdsHtmlString),
		); err != nil {
			log.Fatal(err)
		}
		_org := check(OrgHtmlString)
		_prod := parseToProd(ProdsHtmlString)
		Org = append(Org, _org)
		Prod = append(Prod, _prod...)
	}

	/*for _, element := range Org {
		db.Create(element)
	}
	*/
}

//Получение разметки-Информация об организаии и ее продукции
