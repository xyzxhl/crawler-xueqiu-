package main

import (
	"CrawlBase/db"
	"CrawlBase/pj"
	"CrawlBase/xqCookie"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type NoData struct{}

func (NoData) Error() string {
	return "NoData"
}

var (
	client  *http.Client
	cookie  string
	now     = (time.Now().Unix() - 86400) * 1000
	urlBase = "https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol=%v&begin=%v&period=day&type=before&count=-500&indicator=kline,pe,pb,ps,pcf,market_capital,agt,ggt,balance"
	compch  = make(chan struct{})
	limch   = make(chan struct{}, 20)
	errch   = make(chan error)
	nd      NoData
)

func GetChangeData(symbol string) ([]pj.LChange, error) {
	url := fmt.Sprintf(urlBase, symbol, now)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content := new(bytes.Buffer)
	content.ReadFrom(resp.Body)
	changes, err := pj.LProcessJson(content.Bytes())
	if err != nil {
		errch <- err
		return nil, err
	}

	if len(changes) == 0 {
		return nil, nd
	}
	return changes, nil
}

func Crawl(symbols []string) {
	limch <- struct{}{}
	defer func() {
		<-limch
		compch <- struct{}{}
	}()

	data := make(map[string][]pj.LChange)
	dates := make(map[string]time.Time)

	for _, symbol := range symbols {
		changes, err := GetChangeData(symbol)
		if err != nil {
			errch <- err
			return
		}
		data[symbol] = changes
		dates[symbol] = changes[0].T
	}

	db.FINameUpdateDate(dates)
	db.CHRecordInsert(data)

	errch <- nil
}

func main() {
	client = &http.Client{}
	cookie, _ = xqCookie.GetCookie(client)

	db.InitDB()
	symbols, _ := db.FINameGetAllSymbols()

	for i := 0; i <= len(symbols)/20; i++ {
		if i < len(symbols)/20 {
			go Crawl(symbols[i*20 : i*20+20])
		} else {
			go Crawl(symbols[i*20:])
		}
		<-errch
		fmt.Println(i)
	}
	for range symbols {
		<-compch
		break
	}
}
