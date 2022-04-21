package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

	actx, cancel := chromedp.NewExecAllocator(
		context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...,
	)
	//Подготовка
	mainCtx, cancel := chromedp.NewContext(
		actx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	mainCtx, cancel = context.WithTimeout(mainCtx, 20*time.Second)
	defer cancel()

	url := "https://gisp.gov.ru/pp719v2/pub/org/"

	//получение количества страниц с информацие об организациях
	var pagesInfo string
	if err := chromedp.Run(mainCtx,
		chromedp.Navigate(url),
		chromedp.SetAttributeValue(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, "style", "display: block", chromedp.ByID),
		chromedp.Text(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, &pagesInfo),
	); err != nil {

		log.Fatal(err)
	}

	//Получение количества страниц
	pagesCount, err := strconv.Atoi(strings.Split(pagesInfo, " ")[3])
	if err != nil {
		fmt.Println("Error")
	}

	//Получение разметки всех страниц
	tableData := getHtmlFromOrgTablePages(mainCtx, pagesCount)

	//Получение ссылок на карточки предприятий и список их продукции
	UrlData := getURLs(tableData)
	fmt.Println(len(UrlData))

	//Обход и обработка всех ссылок
	var Prod []Prod
	for _, item := range UrlData {

		URLToProd := "https://gisp.gov.ru" + item.Prods
		Prod = append(Prod, goToProdUrl(mainCtx, URLToProd)...)
	}

	fmt.Println(Prod[0], "------")
}

//Функция получения организации по ссылке
func goToOrgUrl(ctx context.Context, url string, OrgList *[]Org) {

	var pageHtml string

	cloneCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	fmt.Printf("%s is opening in a new tab\n", url)
	chromedp.Run(cloneCtx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(`body > main > div > div.content__inner > div > div:nth-child(2) > div > div > div`, &pageHtml),
	)

}

//Функция получения продукций по ссылке
func goToProdUrl(ctx context.Context, url string) []Prod {

	var htmlData string
	var pagesInfo string

	cloneCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	fmt.Printf("%s is opening in a new tab\n", url)
	//Количество страниц
	if err := chromedp.Run(cloneCtx,
		chromedp.Navigate(url),
		chromedp.SetAttributeValue(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, "style", "display: block", chromedp.ByID),
		chromedp.Text(`#datagrid > div > div.dx-datagrid-pager.dx-pager > div.dx-pages > div.dx-info`, &pagesInfo),
	); err != nil {

		panic(err)
	}

	var pagesCount int

	pagesCount, err := strconv.Atoi(strings.Split(pagesInfo, " ")[3])
	if err != nil {
		fmt.Println("Error")
	}
	//обход страниц
	for i := 1; i <= pagesCount; i++ {
		var tempContainer string
		xpath := `//*[@aria-label="Page ` + strconv.Itoa(i) + `"]`
		if err := chromedp.Run(cloneCtx,
			chromedp.Click(xpath),
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &tempContainer),
		); err != nil {
			log.Fatal(err)
		}
		htmlData += tempContainer
	}

	res := parseToProd(htmlData)
	return res
}

//Функция получения разметки со всех страниц с информацией об организациях
func getHtmlFromOrgTablePages(ctx context.Context, pagesCount int) string {

	var htmlData string
	for i := 1; i <= pagesCount; i++ {
		var tempContainer string
		xpath := `//*[@aria-label="Page ` + strconv.Itoa(i) + `"]`
		if err := chromedp.Run(ctx,
			chromedp.Click(xpath),
			chromedp.WaitVisible(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`),
			chromedp.OuterHTML(`#datagrid > div > div.dx-datagrid-rowsview.dx-scrollable.dx-visibility-change-handler.dx-scrollable-both.dx-scrollable-simulated.dx-scrollable-customizable-scrollbars > div > div > div.dx-scrollable-content > div > table`, &tempContainer, chromedp.ByID),
		); err != nil {
			log.Fatal(err)
		}
		htmlData += tempContainer
	}
	return htmlData
}
