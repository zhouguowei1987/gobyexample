package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	EditDoc88EnableHttpProxy = false
	EditDoc88HttpProxyUrl    = "111.225.152.186:8089"
)

func EditDoc88SetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(EditDoc88HttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type QueryEditDoc88ListFormData struct {
	MenuIndex  int
	ClassifyId string
	FolderId   int
	Sort       int
	Keyword    string
	ShowIndex  int
}

type EditDoc88ResponseData struct {
	Result     string `json:"result"`
	EditTitle  string `json:"edit_title"`
	Class      string `json:"class"`
	UpdateInfo string `json:"updateinfo"`
	State      string `json:"state"`
	SaveFile   string `json:"savefile"`
	Other      string `json:"other"`
}

type EditDoc88FormData struct {
	DocCode        string
	Title          string
	Intro          string
	PCid           string
	Keyword        string
	ShareToDoc     string
	Download       string
	PPrice         string
	PDefaultPoints string
	PPageCount     string
	PDocFormat     string
	Act            string
	GroupList      string
	GroupFreeList  string
}

// ychEduSpider 编辑道客巴巴文档
// @Title 编辑道客巴巴文档
// @Description https://www.doc88.com/，编辑道客巴巴文档
func main() {
	curPage := 1
	for {
		pageListUrl := fmt.Sprintf("https://www.doc88.com/uc/doc_manager.php?act=ajax_doc_list&curpage=%d", curPage)
		fmt.Println(pageListUrl)
		queryEditDoc88ListFormData := QueryEditDoc88ListFormData{
			MenuIndex:  2,
			ClassifyId: "all",
			FolderId:   0,
			Sort:       1,
			Keyword:    "",
			ShowIndex:  1,
		}
		pageListDoc, err := QueryEditDoc88List(pageListUrl, queryEditDoc88ListFormData)
		if err != nil {
			fmt.Println(err)
			break
		}
		liNodes := htmlquery.Find(pageListDoc, `//div[@id="detailed"]/ul[@class="bookshow3"]/li`)
		if len(liNodes) <= 0 {
			break
		}
		for _, liNode := range liNodes {
			PId := htmlquery.SelectAttr(liNode, "id")
			PId = PId[5:len(PId)]

			detailUrl := "https://www.doc88.com/uc/usr_doc_manager.php?act=getDocInfo"
			detailDoc, err := QueryEditDoc88Detail(detailUrl, PId)
			if err != nil {
				fmt.Println(err)
				break
			}

			editUrl := "https://www.doc88.com/uc/index.php"

			DocCodeNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/input`)
			DocCode := htmlquery.SelectAttr(DocCodeNode, "value")

			TitleNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[1]/td[2]/input`)
			Title := htmlquery.SelectAttr(TitleNode, "value")
			fmt.Println(Title)

			IntroNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[2]/td[2]/textarea`)
			Intro := htmlquery.InnerText(IntroNode)

			PCidNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[1]/div[@class="booksedit"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/div[@class="layers"]/input`)
			PCid := htmlquery.SelectAttr(PCidNode, "value")

			PPageCountNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[2]`)
			PPageCount := htmlquery.SelectAttr(PPageCountNode, "value")

			PDocFormatNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[3]/td[2]/input[3]`)
			PDocFormat := htmlquery.SelectAttr(PDocFormatNode, "value")

			PPriceNode := htmlquery.FindOne(detailDoc, `//dl[@class="editlayout"]/form/dd[2]/div[@class="booksedit booksedit-bdr"]/table[@class="edit-table"]/tbody/tr[2]/td[2]/span/input`)
			PPrice := htmlquery.SelectAttr(PPriceNode, "value")

			filePageNum, _ := strconv.Atoi(PPageCount)
			PPriceNew := ""
			// 根据页数设置价格
			if filePageNum > 0 && filePageNum <= 8 {
				PPriceNew = "288"
			} else if filePageNum > 8 && filePageNum <= 18 {
				PPriceNew = "388"
			} else if filePageNum > 18 && filePageNum <= 28 {
				PPriceNew = "488"
			} else if filePageNum > 28 && filePageNum <= 38 {
				PPriceNew = "588"
			} else if filePageNum > 38 && filePageNum <= 48 {
				PPriceNew = "688"
			} else if filePageNum > 48 && filePageNum <= 58 {
				PPriceNew = "788"
			} else {
				PPriceNew = "888"
			}

			// 新旧价格一样，则跳过
			fmt.Println(PPrice, PPriceNew)
			if PPrice == PPriceNew {
				continue
			}
			fmt.Println("===========开始修改价格=============")

			editDoc88FormData := EditDoc88FormData{
				DocCode:        DocCode,
				Title:          Title,
				Intro:          Intro,
				PCid:           PCid,
				Keyword:        "",
				ShareToDoc:     "1",
				Download:       "2",
				PPrice:         PPriceNew,
				PDefaultPoints: "3",
				PPageCount:     PPageCount,
				PDocFormat:     PDocFormat,
				Act:            "save_info",
				GroupList:      "",
				GroupFreeList:  "",
			}
			editDoc88ResponseData, err := EditDoc88(editUrl, editDoc88FormData)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(editDoc88ResponseData)
			time.Sleep(time.Second * 15)
		}
		curPage++
	}
}

func QueryEditDoc88List(requestUrl string, queryEditDoc88ListFormData QueryEditDoc88ListFormData) (doc *html.Node, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if EditDoc88EnableHttpProxy {
		client = EditDoc88SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("menuIndex", strconv.Itoa(queryEditDoc88ListFormData.MenuIndex))
	postData.Add("classify_id", queryEditDoc88ListFormData.ClassifyId)
	postData.Add("folder_id", strconv.Itoa(queryEditDoc88ListFormData.FolderId))
	postData.Add("sort", strconv.Itoa(queryEditDoc88ListFormData.Sort))
	postData.Add("keyword", queryEditDoc88ListFormData.Keyword)
	postData.Add("show_index", strconv.Itoa(queryEditDoc88ListFormData.ShowIndex))
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_READED_PC_ID=%2C440443443; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[show]=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_RW_ID_1652001782=1; Page_Y_28661208441698=-119.39144736842105; Page_28661208441698=1; cdb_RW_ID_1652003396=1; Page_Y_49216489770034=-107.64802631578948; Page_49216489770034=1; cdb_RW_ID_1441436862=1; Page_Y_33273993902621=-138.1809210526316; Page_33273993902621=1; cdb_back[u]=1; show_index=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_RW_ID_1652003342=1; cdb_RW_ID_1446755800=1; Page_Y_63447623881103=-119.39144736842105; Page_63447623881103=1; Page_Y_18461552600944=-138.1809210526316; Page_18461552600944=1; cdb_RW_ID_1448592928=1; Page_Y_38973996541416=-138.1809210526316; Page_38973996541416=1; cdb_RW_ID_1652003569=1; Page_49216489770843=1; cdb_RW_ID_1448117193=1; Page_70487889772761=1; cdb_RW_ID_1652003552=1; Page_Y_84359814005114=47.36513157894737; Page_84359814005114=1; cdb_RW_ID_1652003401=2; Page_Y_40129460552759=-119.39144736842105; Page_40129460552759=1; cdb_RW_ID_1649285366=1; Page_Y_29239265713022=-56.75986842105264; Page_29239265713022=3; cdb_RW_ID_1649286008=1; cdb_back[pcode]=29299246072117; cdb_back[ajax]=1; cdb_back[tm]=3152; cdb_back[member_id]=104598337; Page_Y_29299246072117=-118.99013157894737; Page_29299246072117=1; cdb_RW_ID_1652297248=1; Page_Y_63547623394305=-88.85855263157896; Page_63547623394305=1; cdb_RW_ID_1652297325=1; Page_Y_28961208836780=-119.39144736842105; Page_28961208836780=1; cdb_RW_ID_1652297342=1; Page_Y_27139237754067=-119.39144736842105; Page_27139237754067=1; cdb_pageType=2; cdb_RW_ID_1652300170=1; Page_Y_63547623188748=-119.39144736842105; Page_63547623188748=1; cdb_RW_ID_1432616696=2; Page_89429720494484=1; Page_Y_69299480216066=-138.1809210526316; Page_Y_63547623188849=-119.39144736842105; cdb_RW_ID_1652300033=1; cdb_RW_ID_1434141363=9; Page_50839606969020=1; cdb_RW_ID_1435002537=1; Page_Y_95329726550621=-135.04934210526318; Page_95329726550621=1; Page_Y_40329460255522=-119.39144736842105; Page_40329460255522=1; cdb_RW_ID_1652300089=1; cdb_back[page]=1; Page_Y_20799230811176=-119.39144736842105; Page_20799230811176=1; cdb_RW_ID_1652300079=2; Page_63547623188849=1; cdb_RW_ID_1432609299=3; Page_69299480216066=1; cdb_RW_ID_1435785193=1; Page_Y_14761570690137=-138.1809210526316; Page_14761570690137=1; Page_Y_50987530144085=167.14802631578948; cdb_back[login]=1; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d234d38b35c1054ae59db368050a5a7a43ad92f350c3f26f274; c_login_name=woyoceo; cdb_logined=1; cdb_back[module_type]=7; cdb_back[image_type]=3; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dmyshare; cdb_back[txtloginname]=15238369929; doc88_lt=wx; cdb_tokenid=4efc50xnmPveubMMOm05CKiTVllMx3eqDsnXwcJY7%2FL37iFZdGAbcN%2FDVkrqzEz73yTCHw8SOzW7RnVxKYz9vh7HmLDKltpLbQ5jRk7H328I9lcYOi89DkaAenClxpoXoA; cdb_back[t]=1; cdb_RW_ID_1647381158=1; Page_47316412056685=1; cdb_back[inout]=all; cdb_back[type]=1; cdb_back[pcid]=8371; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1652300246=2; cdb_back[m]=104598337; Page_50987530144085=1; cdb_back[doc_more_id]=1652565970%2C1652565897%2C1652565888%2C1652565868%2C1652565864%2C1652565849%2C1652565828%2C1652565799%2C1652565776%2C1652565755%2C1652565668%2C1652565664%2C1652565647%2C1652565639%2C1652565620%2C1652565609%2C1652565554%2C1652565515%2C; cdb_RW_ID_1652565809=1; Page_Y_27139237323185=-105.29934210526316; Page_27139237323185=1; cdb_RW_ID_1652566036=2; Page_Y_50987530355415=-109.2138157894737; Page_50987530355415=4; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1652566010=1; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; cdb_back[s]=rel; cdb_RW_ID_1443621487=11; Page_Y_99699448209475=167.9309210526316; Page_99699448209475=2; Page_Y_21673251522838=-119.39144736842105; Page_21673251522838=1; cdb_RW_ID_1652598987=2; Page_Y_49316489835352=-119.39144736842105; Page_49316489835352=1; cdb_back[id]=2; cdb_RW_ID_1420304862=2; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0lXizNXiFNXi2jMW2LE51TPQ1LES0qh9or3c3gJACuvi2i3R1jsS1qk%212qnQ3iXiFotZBK363jll2qBj2Oxj0Tpj1Wv%21HW0S0TFiHWBlBqnR0jvU0qsW3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_Y_20199401814720=354.25986842105266; Page_20199401814720=1; cdb_RW_ID_1652599118=2; cdb_back[doctype]=1; cdb_back[len]=2; Page_Y_49316489833665=-252.45; Page_49316489833665=3; cdb_back[doc_id]=1652599178; cdb_back[curpage]=2; cdb_back[pid]=27539237355940; cdb_RW_ID_1652599173=1; cdb_back[srlid]=9bc9JZLleCc640vJ87oO7titHUq910Z7wgSa9IyXu2jvxPEyuCkZbz+wYWyWaAAs14TWW6r1WnF55bluTjoMTgX8AsFrTi8vcp1KqO%2FnZ7wt; cdb_back[p_name]=%E6%99%BA%E6%85%A7%E5%A4%9A%E5%8A%9F%E8%83%BD%E6%9D%86%E5%BB%BA%E8%AE%BE%E6%8A%80%E6%9C%AF%E5%AF%BC%E5%88%99%28DB3502-T+086-2022%29; cdb_back[rel_p_id]=1652599173; Page_Y_27539237355940=-119.39144736842105; Page_27539237355940=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_change_message=1; cdb_msg_num=0; cdb_back[classify_id]=all; cdb_back[p_id]=1652599191; cdb_back[doccode]=1652599191; cdb_back[title]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[intro]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[p_price]=388; cdb_back[p_default_points]=3; cdb_back[p_pagecount]=16; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1688010341; cdb_back[act]=ajax_doc_list; siftState=1; cdb_back[sort]=2")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=all")
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

func QueryEditDoc88Detail(requestUrl string, PId string) (doc *html.Node, err error) {
	// 初始化客户端
	var client *http.Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	if EditDoc88EnableHttpProxy {
		client = EditDoc88SetHttpProxy()
	}
	postData := url.Values{}
	postData.Add("p_id", PId)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return doc, err
	}

	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_READED_PC_ID=%2C440443443; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[show]=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_RW_ID_1652001782=1; Page_Y_28661208441698=-119.39144736842105; Page_28661208441698=1; cdb_RW_ID_1652003396=1; Page_Y_49216489770034=-107.64802631578948; Page_49216489770034=1; cdb_RW_ID_1441436862=1; Page_Y_33273993902621=-138.1809210526316; Page_33273993902621=1; cdb_back[u]=1; show_index=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_RW_ID_1652003342=1; cdb_RW_ID_1446755800=1; Page_Y_63447623881103=-119.39144736842105; Page_63447623881103=1; Page_Y_18461552600944=-138.1809210526316; Page_18461552600944=1; cdb_RW_ID_1448592928=1; Page_Y_38973996541416=-138.1809210526316; Page_38973996541416=1; cdb_RW_ID_1652003569=1; Page_49216489770843=1; cdb_RW_ID_1448117193=1; Page_70487889772761=1; cdb_RW_ID_1652003552=1; Page_Y_84359814005114=47.36513157894737; Page_84359814005114=1; cdb_RW_ID_1652003401=2; Page_Y_40129460552759=-119.39144736842105; Page_40129460552759=1; cdb_RW_ID_1649285366=1; Page_Y_29239265713022=-56.75986842105264; Page_29239265713022=3; cdb_RW_ID_1649286008=1; cdb_back[pcode]=29299246072117; cdb_back[ajax]=1; cdb_back[tm]=3152; cdb_back[member_id]=104598337; Page_Y_29299246072117=-118.99013157894737; Page_29299246072117=1; cdb_RW_ID_1652297248=1; Page_Y_63547623394305=-88.85855263157896; Page_63547623394305=1; cdb_RW_ID_1652297325=1; Page_Y_28961208836780=-119.39144736842105; Page_28961208836780=1; cdb_RW_ID_1652297342=1; Page_Y_27139237754067=-119.39144736842105; Page_27139237754067=1; cdb_pageType=2; cdb_RW_ID_1652300170=1; Page_Y_63547623188748=-119.39144736842105; Page_63547623188748=1; cdb_RW_ID_1432616696=2; Page_89429720494484=1; Page_Y_69299480216066=-138.1809210526316; Page_Y_63547623188849=-119.39144736842105; cdb_RW_ID_1652300033=1; cdb_RW_ID_1434141363=9; Page_50839606969020=1; cdb_RW_ID_1435002537=1; Page_Y_95329726550621=-135.04934210526318; Page_95329726550621=1; Page_Y_40329460255522=-119.39144736842105; Page_40329460255522=1; cdb_RW_ID_1652300089=1; cdb_back[page]=1; Page_Y_20799230811176=-119.39144736842105; Page_20799230811176=1; cdb_RW_ID_1652300079=2; Page_63547623188849=1; cdb_RW_ID_1432609299=3; Page_69299480216066=1; cdb_RW_ID_1435785193=1; Page_Y_14761570690137=-138.1809210526316; Page_14761570690137=1; Page_Y_50987530144085=167.14802631578948; cdb_back[login]=1; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d234d38b35c1054ae59db368050a5a7a43ad92f350c3f26f274; c_login_name=woyoceo; cdb_logined=1; cdb_back[module_type]=7; cdb_back[image_type]=3; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dmyshare; cdb_back[txtloginname]=15238369929; doc88_lt=wx; cdb_tokenid=4efc50xnmPveubMMOm05CKiTVllMx3eqDsnXwcJY7%2FL37iFZdGAbcN%2FDVkrqzEz73yTCHw8SOzW7RnVxKYz9vh7HmLDKltpLbQ5jRk7H328I9lcYOi89DkaAenClxpoXoA; cdb_back[t]=1; cdb_RW_ID_1647381158=1; Page_47316412056685=1; cdb_back[inout]=all; cdb_back[type]=1; cdb_back[pcid]=8371; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1652300246=2; cdb_back[m]=104598337; Page_50987530144085=1; cdb_back[doc_more_id]=1652565970%2C1652565897%2C1652565888%2C1652565868%2C1652565864%2C1652565849%2C1652565828%2C1652565799%2C1652565776%2C1652565755%2C1652565668%2C1652565664%2C1652565647%2C1652565639%2C1652565620%2C1652565609%2C1652565554%2C1652565515%2C; cdb_RW_ID_1652565809=1; Page_Y_27139237323185=-105.29934210526316; Page_27139237323185=1; cdb_RW_ID_1652566036=2; Page_Y_50987530355415=-109.2138157894737; Page_50987530355415=4; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1652566010=1; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; cdb_back[s]=rel; cdb_RW_ID_1443621487=11; Page_Y_99699448209475=167.9309210526316; Page_99699448209475=2; Page_Y_21673251522838=-119.39144736842105; Page_21673251522838=1; cdb_RW_ID_1652598987=2; Page_Y_49316489835352=-119.39144736842105; Page_49316489835352=1; cdb_back[id]=2; cdb_RW_ID_1420304862=2; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0lXizNXiFNXi2jMW2LE51TPQ1LES0qh9or3c3gJACuvi2i3R1jsS1qk%212qnQ3iXiFotZBK363jll2qBj2Oxj0Tpj1Wv%21HW0S0TFiHWBlBqnR0jvU0qsW3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_Y_20199401814720=354.25986842105266; Page_20199401814720=1; cdb_RW_ID_1652599118=2; cdb_back[doctype]=1; cdb_back[len]=2; Page_Y_49316489833665=-252.45; Page_49316489833665=3; cdb_back[doc_id]=1652599178; cdb_back[curpage]=2; cdb_back[pid]=27539237355940; cdb_RW_ID_1652599173=1; cdb_back[srlid]=9bc9JZLleCc640vJ87oO7titHUq910Z7wgSa9IyXu2jvxPEyuCkZbz+wYWyWaAAs14TWW6r1WnF55bluTjoMTgX8AsFrTi8vcp1KqO%2FnZ7wt; cdb_back[p_name]=%E6%99%BA%E6%85%A7%E5%A4%9A%E5%8A%9F%E8%83%BD%E6%9D%86%E5%BB%BA%E8%AE%BE%E6%8A%80%E6%9C%AF%E5%AF%BC%E5%88%99%28DB3502-T+086-2022%29; cdb_back[rel_p_id]=1652599173; Page_Y_27539237355940=-119.39144736842105; Page_27539237355940=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_change_message=1; cdb_msg_num=0; cdb_back[classify_id]=all; cdb_back[p_id]=1652599191; cdb_back[doccode]=1652599191; cdb_back[title]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[intro]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[p_price]=388; cdb_back[p_default_points]=3; cdb_back[p_pagecount]=16; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1688010341; cdb_back[act]=ajax_doc_list; siftState=1; cdb_back[sort]=2")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=all")
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

func EditDoc88(requestUrl string, editDoc88FormData EditDoc88FormData) (editDoc88ResponseData EditDoc88ResponseData, err error) {
	client := &http.Client{}
	editDoc88ResponseData = EditDoc88ResponseData{}
	postData := url.Values{}
	postData.Add("doccode", editDoc88FormData.DocCode)
	postData.Add("title", editDoc88FormData.Title)
	postData.Add("intro", editDoc88FormData.Intro)
	postData.Add("pcid", editDoc88FormData.PCid)
	postData.Add("keyword", editDoc88FormData.Keyword)
	postData.Add("sharetodoc", editDoc88FormData.ShareToDoc)
	postData.Add("download", editDoc88FormData.Download)
	postData.Add("p_price", editDoc88FormData.PPrice)
	postData.Add("p_default_points", editDoc88FormData.PDefaultPoints)
	postData.Add("p_pagecount", editDoc88FormData.PPageCount)
	postData.Add("p_doc_format", editDoc88FormData.PDocFormat)
	postData.Add("act", "save_info")
	postData.Add("group_list", editDoc88FormData.GroupList)
	postData.Add("group_free_list", editDoc88FormData.GroupFreeList)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(postData.Encode())) //建立连接

	if err != nil {
		return editDoc88ResponseData, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "__root_domain_v=.doc88.com; _qddaz=QD.155181178889683; _qddab=3-gv2ozy.lib1y9mi; PHPSESSID=r1clbe0fu15io3vrsg41mce152; cdb_sys_sid=r1clbe0fu15io3vrsg41mce152; cdb_back[at]=0; cdb_back[n]=6; cdb_back[book_id]=0; cdb_READED_PC_ID=%2C440443443; cdb_back[mid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_back[show]=1; cdb_back[sharetodoc]=1; cdb_back[download]=2; cdb_RW_ID_1652001782=1; Page_Y_28661208441698=-119.39144736842105; Page_28661208441698=1; cdb_RW_ID_1652003396=1; Page_Y_49216489770034=-107.64802631578948; Page_49216489770034=1; cdb_RW_ID_1441436862=1; Page_Y_33273993902621=-138.1809210526316; Page_33273993902621=1; cdb_back[u]=1; show_index=1; cdb_back[folder_id]=0; cdb_back[show_index]=1; cdb_RW_ID_1652003342=1; cdb_RW_ID_1446755800=1; Page_Y_63447623881103=-119.39144736842105; Page_63447623881103=1; Page_Y_18461552600944=-138.1809210526316; Page_18461552600944=1; cdb_RW_ID_1448592928=1; Page_Y_38973996541416=-138.1809210526316; Page_38973996541416=1; cdb_RW_ID_1652003569=1; Page_49216489770843=1; cdb_RW_ID_1448117193=1; Page_70487889772761=1; cdb_RW_ID_1652003552=1; Page_Y_84359814005114=47.36513157894737; Page_84359814005114=1; cdb_RW_ID_1652003401=2; Page_Y_40129460552759=-119.39144736842105; Page_40129460552759=1; cdb_RW_ID_1649285366=1; Page_Y_29239265713022=-56.75986842105264; Page_29239265713022=3; cdb_RW_ID_1649286008=1; cdb_back[pcode]=29299246072117; cdb_back[ajax]=1; cdb_back[tm]=3152; cdb_back[member_id]=104598337; Page_Y_29299246072117=-118.99013157894737; Page_29299246072117=1; cdb_RW_ID_1652297248=1; Page_Y_63547623394305=-88.85855263157896; Page_63547623394305=1; cdb_RW_ID_1652297325=1; Page_Y_28961208836780=-119.39144736842105; Page_28961208836780=1; cdb_RW_ID_1652297342=1; Page_Y_27139237754067=-119.39144736842105; Page_27139237754067=1; cdb_pageType=2; cdb_RW_ID_1652300170=1; Page_Y_63547623188748=-119.39144736842105; Page_63547623188748=1; cdb_RW_ID_1432616696=2; Page_89429720494484=1; Page_Y_69299480216066=-138.1809210526316; Page_Y_63547623188849=-119.39144736842105; cdb_RW_ID_1652300033=1; cdb_RW_ID_1434141363=9; Page_50839606969020=1; cdb_RW_ID_1435002537=1; Page_Y_95329726550621=-135.04934210526318; Page_95329726550621=1; Page_Y_40329460255522=-119.39144736842105; Page_40329460255522=1; cdb_RW_ID_1652300089=1; cdb_back[page]=1; Page_Y_20799230811176=-119.39144736842105; Page_20799230811176=1; cdb_RW_ID_1652300079=2; Page_63547623188849=1; cdb_RW_ID_1432609299=3; Page_69299480216066=1; cdb_RW_ID_1435785193=1; Page_Y_14761570690137=-138.1809210526316; Page_14761570690137=1; Page_Y_50987530144085=167.14802631578948; cdb_back[login]=1; cdb_back[txtPassword]=abcdqq123456; cdb_back[captchaCode]=1; cdb_login_if=1; cdb_uid=104598337; cdb_token=5176691bb4a2b7d6bd67c231efd81e657d782f6cb333928fd33f234c70382d9a89fe0ad0ebba21c3dc7bc12152ab66ccc2f5b04d04e00e86770e2edff24aa4a84def49f043721d234d38b35c1054ae59db368050a5a7a43ad92f350c3f26f274; c_login_name=woyoceo; cdb_logined=1; cdb_back[module_type]=7; cdb_back[image_type]=3; cdb_back[refer]=%2Fuc%2Fdoc_manager.php%3Fact%3Ddoc_list%26state%3Dmyshare; cdb_back[txtloginname]=15238369929; doc88_lt=wx; cdb_tokenid=4efc50xnmPveubMMOm05CKiTVllMx3eqDsnXwcJY7%2FL37iFZdGAbcN%2FDVkrqzEz73yTCHw8SOzW7RnVxKYz9vh7HmLDKltpLbQ5jRk7H328I9lcYOi89DkaAenClxpoXoA; cdb_back[t]=1; cdb_RW_ID_1647381158=1; Page_47316412056685=1; cdb_back[inout]=all; cdb_back[type]=1; cdb_back[pcid]=8371; cdb_back[p_doc_format]=PDF; cdb_RW_ID_1652300246=2; cdb_back[m]=104598337; Page_50987530144085=1; cdb_back[doc_more_id]=1652565970%2C1652565897%2C1652565888%2C1652565868%2C1652565864%2C1652565849%2C1652565828%2C1652565799%2C1652565776%2C1652565755%2C1652565668%2C1652565664%2C1652565647%2C1652565639%2C1652565620%2C1652565609%2C1652565554%2C1652565515%2C; cdb_RW_ID_1652565809=1; Page_Y_27139237323185=-105.29934210526316; Page_27139237323185=1; cdb_RW_ID_1652566036=2; Page_Y_50987530355415=-109.2138157894737; Page_50987530355415=4; cdb_back[uid]=b99ce806c0b55b3bdccae7bc14f8ca3e; cdb_RW_ID_1652566010=1; ShowSkinTip_1=1; cdb_H5R=1; showAnnotateTipIf=1; cdb_back[s]=rel; cdb_RW_ID_1443621487=11; Page_Y_99699448209475=167.9309210526316; Page_99699448209475=2; Page_Y_21673251522838=-119.39144736842105; Page_21673251522838=1; cdb_RW_ID_1652598987=2; Page_Y_49316489835352=-119.39144736842105; Page_49316489835352=1; cdb_back[id]=2; cdb_RW_ID_1420304862=2; cdb_back[data]=GSxkHoph3jfiuQdE3mNE3jZE3gxlDN9kDW1A0lXizNXiFNXi2jMW2LE51TPQ1LES0qh9or3c3gJACuvi2i3R1jsS1qk%212qnQ3iXiFotZBK363jll2qBj2Oxj0Tpj1Wv%21HW0S0TFiHWBlBqnR0jvU0qsW3iXiDutdHmtSoWlk3jfi0qPU1qk%210T0Q3gU%3D; Page_Y_20199401814720=354.25986842105266; Page_20199401814720=1; cdb_RW_ID_1652599118=2; cdb_back[doctype]=1; cdb_back[len]=2; Page_Y_49316489833665=-252.45; Page_49316489833665=3; cdb_back[doc_id]=1652599178; cdb_back[curpage]=2; cdb_back[pid]=27539237355940; cdb_RW_ID_1652599173=1; cdb_back[srlid]=9bc9JZLleCc640vJ87oO7titHUq910Z7wgSa9IyXu2jvxPEyuCkZbz+wYWyWaAAs14TWW6r1WnF55bluTjoMTgX8AsFrTi8vcp1KqO%2FnZ7wt; cdb_back[p_name]=%E6%99%BA%E6%85%A7%E5%A4%9A%E5%8A%9F%E8%83%BD%E6%9D%86%E5%BB%BA%E8%AE%BE%E6%8A%80%E6%9C%AF%E5%AF%BC%E5%88%99%28DB3502-T+086-2022%29; cdb_back[rel_p_id]=1652599173; Page_Y_27539237355940=-119.39144736842105; Page_27539237355940=1; cdb_back[pm_id]=1486396; cdb_back[friend_id]=0; cdb_change_message=1; cdb_msg_num=0; cdb_back[classify_id]=all; cdb_back[p_id]=1652599191; cdb_back[doccode]=1652599191; cdb_back[title]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[intro]=%E6%99%BA%E6%85%A7%E7%A4%BE%E5%8C%BA%E5%81%A5%E8%BA%AB%E4%B8%AD%E5%BF%83%E5%BB%BA%E8%AE%BE%E4%B8%8E%E8%BF%90%E8%90%A5%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%28DB43-T+2282-2022%29; cdb_back[p_price]=388; cdb_back[p_default_points]=3; cdb_back[p_pagecount]=16; cdb_back[state]=all; cdb_back[menuIndex]=2; cdb_msg_time=1688010341; cdb_back[act]=ajax_doc_list; siftState=1; cdb_back[sort]=2")
	req.Header.Set("Host", "www.doc88.com")
	req.Header.Set("Origin", "https://www.doc88.com")
	req.Header.Set("Referer", "https://www.doc88.com/uc/doc_manager.php?act=doc_list&state=all")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return editDoc88ResponseData, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return editDoc88ResponseData, err
	}
	err = json.Unmarshal(respBytes, &editDoc88ResponseData)
	if err != nil {
		return editDoc88ResponseData, err
	}
	return editDoc88ResponseData, nil
}
