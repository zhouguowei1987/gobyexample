package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	EditDocInEnableHttpProxy = false
	EditDocInHttpProxyUrl    = "111.225.152.186:8089"
)

func EditDocInSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(EditDocInHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

var DocInCookie = "docin_session_id=b686d233-9d6a-4d19-bdbd-55ba049ac297; cookie_id=CA1BC7B9D32000011AEA1CB081C0DA70; time_id=20221227213222; partner_tips=1; __bid_n=18553c8c217b5d5e2e4207; FEID=v10-9bb248a4c21a53f72760ecda6234cbecf70a7381; __xaf_fpstarttimer__=1672147945190; __xaf_thstime__=1672147945346; __xaf_fptokentimer__=1672147945861; last_upload_public243402665=yes; indexnoticeupdatetag31526379=unshow; __root_domain_v=.docin.com; _qddaz=QD.151176617720177; aliyungf_tc=bd3e315da1d7ce31af3fca4a3a2cdd702d14a676aaad8720b6105200a0e499e0; ifShowMsg=true; jumpIn=400; search_option_show=0; _ga_ZYR13KTSXC=deleted; _ga_ZYR13KTSXC=deleted; FPTOKEN=yJ0+QSXweZ4Cqq0iuE4OoxgWMakB+ymq7HPZB8+AFbd3gLJfGXg+uhwC0PoTY2vuB9fF5/2qsuUHnMT2yBvsmKwBeb2Es9MY/cERDMA9Eo0Q4NnpM+qitBzZ5FgbqMdP5jaPkSY3peXMVlupZpjbAVsSoBjCQ5h+OaSsGsZHh5XMvNOM2sUd+BUqUvfHYTZcEf1zoMXQMbPlOFUOa1qcC5h6YlPH6Q3uN6f67bocIhZijom17xgVRl5ISjwvfhBdEIDZWVbcvcDY6+VU8WZfDOkd/p66Bm9Sz/OeodH8SiuMetE/mcgTgF5KiLFLh3yS8JVVgndmKGV3Yppl9eVrOxYrlQPZjf0rOfxGNdetjSbCEVs+HX/Usks8sfKMRey9ZSjQ0XKZdAhRoUoQwCk4BA==|FOHf58AIodtuneILuIhr9eSK/3MZMJ4ih8ikwzMvr5Y=|10|f132030d89062ab6a03a37cb0d2b63f1; visitTopicIds=\"279237,285253,285264\"; _qddab=3-722psv.ll31zmje; pbyCookieKey=1692338365244; userChoose=usertags104_174_171_102_175_176_169_177_178_179_180_181_998_000_; lastLoginType=weixin; firstTimeComDoc=1; recharge_from_type=nav-sub1; refererfunction=https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3Dnwb1nAPJ_i_UD2MO7P3U9Z4-K83LNkrqRXW0vE_zGygvqUPm9OMkjP6rnpT1UJBW%26wd%3D%26eqid%3Dca9b81b10000dab7000000046508091a; isbaiduspider=false; buy_type_4357253219=0; _gid=GA1.2.1484009812.1695543436; mobilefirsttip=tip; today_first_in=1; login_email=phone_25957a50cb; user_password=BnKagWbX8Rc9bU093Rt8k1B1HNSQ7XM%2F33Je3YytS4WYZIxiWbPS%2FYbpBL6ZHAF2H_T1JcST9ZQ6cBZHY6cnK%2Bh1bZX%2FSqMZlbEvCy2w%2FisUofQHFGYX%2F5M2re3LyoROQsXpn%2BcgL6ijLEs%0ANml31NYNlOZ2J2YdK5%2Bg; s_from=direct; uaType=chrome; netfunction=\"/my/upload/myUpload.do\"; JSESSIONID=C8EFBA239F060E8A99F8477A4A23780D-n2; remindClickId=-1; _gat_gtag_UA_3158355_1=1; _ga=GA1.2.43085923.1672147943; _ga_ZYR13KTSXC=GS1.1.1695714096.177.1.1695714745.13.0.0"
var downPrice = 5

// ychEduSpider 编辑豆丁文档
// @Title 编辑豆丁文档
// @Description https://www.docin.com/，编辑豆丁文档
func main() {
	currentPage := 1
	beginId := 0
	for {
		pageListUrl := "https://www.docin.com/my/upload/myUpload.do?onlypPublic=1&totalpublicnum=0"
		referer := "https://www.docin.com/my/upload/myUpload.do?onlypPublic=1&totalpublicnum=0"
		if currentPage > 1 {
			pageListUrl = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
				"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&orderFolder=0&folderId=0"+
				"&myKeyword=&publishCount=&onlypPrivate=&totalprivatenum=0&onlypPublic=1"+
				"&totalpublicnum=0&currentPage=%d&pageType=n&beginId=%d", currentPage, beginId)
		}

		fmt.Println(pageListUrl)
		pageListDoc, err := QueryDocInDoc(pageListUrl, referer)
		if err != nil {
			fmt.Println(err)
			break
		}
		tbodyNodes := htmlquery.Find(pageListDoc, `//div[@class="tableWarp"]/table[@class="my-data"]/tbody`)
		if len(tbodyNodes) <= 0 {
			break
		}
		idsArr := make([]string, 0)
		for _, tbodyNode := range tbodyNodes {
			trNode := htmlquery.FindOne(tbodyNode, `./tr`)
			trId := strings.ReplaceAll(htmlquery.SelectAttr(trNode, "id"), "tr", "")
			idsArr = append(idsArr, trId)

			fileTitleNode := htmlquery.FindOne(tbodyNode, `./tr/td[2]/a`)
			fileTitle := htmlquery.SelectAttr(fileTitleNode, "title")
			fmt.Println(fileTitle)

			filePageNode := htmlquery.FindOne(tbodyNode, `./tr/td[4]`)
			filePage := htmlquery.InnerText(filePageNode)
			filePage = strings.TrimSpace(filePage)
			filePage = strings.ReplaceAll(filePage, "页", "")
			// 根据页数设置价格
			filePageNum, _ := strconv.Atoi(filePage)
			if filePageNum > 0 {
				if filePageNum > 0 && filePageNum <= 5 {
					downPrice = 2
				} else if filePageNum > 5 && filePageNum <= 10 {
					downPrice = 3
				} else if filePageNum > 10 && filePageNum <= 15 {
					downPrice = 4
				} else if filePageNum > 15 && filePageNum <= 20 {
					downPrice = 5
				} else if filePageNum > 20 && filePageNum <= 25 {
					downPrice = 6
				} else if filePageNum > 25 && filePageNum <= 30 {
					downPrice = 7
				} else if filePageNum > 30 && filePageNum <= 35 {
					downPrice = 8
				} else if filePageNum > 35 && filePageNum <= 40 {
					downPrice = 9
				} else if filePageNum > 40 && filePageNum <= 45 {
					downPrice = 10
				} else if filePageNum > 45 && filePageNum <= 50 {
					downPrice = 11
				} else {
					downPrice = 12
				}
			}

			// 查看文档原来价格
			filePriceNode := htmlquery.FindOne(tbodyNode, `./tr/td[5]`)
			filePrice := htmlquery.InnerText(filePriceNode)
			filePrice = strings.TrimSpace(filePrice)
			if filePrice != "免费" {
				floatFilePrice, err := strconv.ParseFloat(filePrice, 64)
				if err != nil {
					continue
				}
				originalPrice := int(floatFilePrice)
				if downPrice == originalPrice {
					continue
				}
			}

			// 开始设置价格
			fmt.Println("-----------------开始设置价格--------------------")
			editUrl := fmt.Sprintf("https://www.docin.com/app/my/docin/batchModifyPrice.do?ids=%s&down_price=%d&price_flag=0", trId, downPrice)
			_, err = QueryDocInDoc(editUrl, referer)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("-----------------开始设置价格完结--------------------")
			time.Sleep(time.Microsecond * 100)
		}
		beginId, _ = strconv.Atoi(idsArr[len(idsArr)-1])
		currentPage++
		fmt.Println(currentPage)
		referer = fmt.Sprintf("https://www.docin.com/my/upload/myUpload.do?styleList=1"+
			"&orderName=0&orderDate=0&orderVisit=0&orderStatus=0&orderFolder=0&folderId=0"+
			"&myKeyword=&publishCount=&onlypPrivate=&totalprivatenum=0&onlypPublic=1"+
			"&totalpublicnum=0&currentPage=%d&pageType=n&beginId=%d", currentPage, beginId)
	}
}

func QueryDocInDoc(requestUrl string, referer string) (doc *html.Node, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 20,
		},
	}
	if EditDocInEnableHttpProxy {
		client = EditDocInSetHttpProxy()
	}
	req, err := http.NewRequest("GET", requestUrl, nil) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", DocInCookie)
	req.Header.Set("Host", "www.docin.com")
	req.Header.Set("Origin", "https://www.docin.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return doc, errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}
	doc, err = htmlquery.Parse(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}
