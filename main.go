package main

import (
	"fmt"
	"get-toman/configuration"
	"get-toman/model"
	"math"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

///get papers
func main() {
	db := configuration.GetDB()
	defer db.Close()

	tgjuRate := getEuroTgju()
	bazar360Rate := getEuroBazar360()

	var tomanRate model.Tomanrate

	tomanRate.Date = time.Now().Format("2006-01-02 15:04:05")
	tomanRate.Currency = 1
	tomanRate.Rate = rateValidation(tgjuRate, bazar360Rate)

	model.InsertTomanrate(db, tomanRate)

	fmt.Println(tomanRate)
	fmt.Println("-------Get Toman @25Mordad")
}

type sett struct {
	margin float64
}

func Setting() sett {
	return sett{margin: configuration.GetMargin()}
}

func rateValidation(tgjuRate, bazar360Rate float64) float64 {
	finalRate := 0.0

	if tgjuRate != 0 && bazar360Rate != 0 {
		finalRate = (tgjuRate + bazar360Rate) / 2
	} else {
		finalRate = math.Max(tgjuRate, bazar360Rate)
	}
	setting := Setting()
	finalRate = (setting.margin * finalRate / 100) + finalRate

	return finalRate
}

///get EUR tgju.org
func getEuroTgju() float64 {

	db := configuration.GetDB()
	defer db.Close()

	c := colly.NewCollector()
	thisRate := 0.0
	////////////////// http://www.tgju.org/chart/price_eur
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if e.Attr("itemprop") == "price" {
			thisRate, _ = strconv.ParseFloat(strings.Replace(strings.Trim(e.Text, " "), ",", "", -1), 64)
			fmt.Println("tgju.org: ", thisRate/10)
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {

	})

	c.Visit("http://www.tgju.org/chart/price_eur")
	return thisRate / 10
}

///get EUR bazar360
func getEuroBazar360() float64 {

	db := configuration.GetDB()
	defer db.Close()

	// Instantiate default collector
	c := colly.NewCollector()
	thisRate := 0.0
	////////////////// https://bazar360.com/en/Currencies/
	i := 0
	c.OnHTML("td.success", func(e *colly.HTMLElement) {
		if strings.Trim(e.Text, " ") != "" {
			i = i + 1
			if i == 6 {
				strArray := strings.Fields(strings.Trim(e.Text, " "))
				thisRate, _ = strconv.ParseFloat(strings.Replace(strArray[0], ",", "", 1), 64)
				fmt.Println("bazar360: ", thisRate)

			}

		}

	})

	c.OnError(func(r *colly.Response, err error) {

	})

	c.Visit("https://bazar360.com/en/Currencies/")
	return thisRate
}
