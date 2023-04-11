package main

import (
	"CrawlBase/db"
	"CrawlBase/pj"
	"CrawlBase/xqCookie"
	"fmt"
	"io"
	"net/http"
)

func GetFundsByType(t int, cookie string, client *http.Client) []pj.RFund {
	url := "https://stock.xueqiu.com/v5/stock/screener/fund/list.json?type=" + fmt.Sprintf("%d", t) +
		"&parent_type=1&order=desc&order_by=percent&page=1&size=900"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", cookie)
	resp, _ := client.Do(req)
	content, _ := io.ReadAll(resp.Body)
	data, _ := pj.RProcessJson(content)
	return data
}

func main() {
	db.InitDB()
	client := &http.Client{}
	cookie := xqCookie.GetCookie(client)
	for i := 11; i < 21; i++ {
		Funds := GetFundsByType(i, cookie, client)
		for _, v := range Funds {
			db.FINameInsert(v.Symbol, v.Name)
		}
	}
}
