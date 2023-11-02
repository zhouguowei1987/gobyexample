package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ZZStepEnableHttpProxy = false
	ZZStepHttpProxyUrl    = "111.225.152.186:8089"
)

func ZZStepSetHttpProxy() (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(ZZStepHttpProxyUrl)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

type ZZStepSubject struct {
	name string
	url  string
}

var subjects = []ZZStepSubject{
	{
		name: "试卷",
		url:  "http://www2.zzstep.com/front/paper/index.html",
	},
	{
		name: "中考",
		url:  "http://www2.zzstep.com/front/beikao/index.html",
	},
}

// ychEduSpider 获取中国教育出版网文档
// @Title 获取中国教育出版网文档
// @Description http://www2.zzstep.com/，获取中国教育出版网文档
func main() {
	for _, subject := range subjects {
		current := 1
		isPageListGo := true
		for isPageListGo {
			subjectIndexUrl := subject.url
			if current > 1 {
				subjectIndexUrl += fmt.Sprintf("?studysection=204&subject=29&page=%d", current)
			}
			subjectIndexDoc, err := htmlquery.LoadURL(subjectIndexUrl)
			if err != nil {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			liNodes := htmlquery.Find(subjectIndexDoc, `//div[@class="zy-list fn-mt20"]/ul[@class="reslist"]/li[@class="fn-pt20 fn-pb20"]`)
			if len(liNodes) <= 0 {
				fmt.Println(err)
				current = 1
				isPageListGo = false
				continue
			}
			for _, liNode := range liNodes {
				fmt.Println("============================================================================")
				fmt.Println("主题：", subject.name)
				fmt.Println("=======当前页为：" + strconv.Itoa(current) + "========")

				fileName := htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a`))
				fileName = strings.TrimSpace(fileName)
				fileName = strings.ReplaceAll(fileName, "/", "-")
				fileName = strings.ReplaceAll(fileName, ":", "-")
				fileName = strings.ReplaceAll(fileName, "：", "-")
				fileName = strings.ReplaceAll(fileName, "（", "(")
				fileName = strings.ReplaceAll(fileName, "）", ")")
				fmt.Println(fileName)

				filePath := "../www2.zzstep.com/www2.zzstep.com/" + subject.name + "/" + fileName
				_, errDoc := os.Stat(filePath + ".doc")
				_, errDocx := os.Stat(filePath + ".docx")
				_, errPdf := os.Stat(filePath + ".pdf")
				_, errPpt := os.Stat(filePath + ".ppt")
				_, errPptx := os.Stat(filePath + ".pptx")
				if errDoc != nil && errDocx != nil && errPdf != nil && errPpt != nil && errPptx != nil {

					// 所需智币
					pointsNode := htmlquery.FindOne(liNode, `./div[@class="btn-item fn-left"]/div[@class="money fn-pt10"]`)
					if pointsNode == nil {
						fmt.Println("没有智币div")
						continue
					}
					pointsText := htmlquery.InnerText(pointsNode)
					fmt.Println(pointsText)
					pointsText = strings.ReplaceAll(pointsText, "智币", "")

					points, err := strconv.Atoi(pointsText)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if points > 0 {
						fmt.Println("需要智币下载", points)
						continue
					}

					viewUrl := "http://www2.zzstep.com" + htmlquery.InnerText(htmlquery.FindOne(liNode, `./div[@class="zy-box fn-left"]/div[@class="subject-t"]/a/@href`))
					fmt.Println(viewUrl)

					downLoadUrl := strings.ReplaceAll(viewUrl, "index", "download")
					fmt.Println(downLoadUrl)

					fmt.Println("=======开始下载" + strconv.Itoa(current) + "========")
					err = downloadZZStep(downLoadUrl, viewUrl, filePath)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("=======开始完成========")
					time.Sleep(time.Millisecond * 200)
				}
			}
			current++
			isPageListGo = true
		}
	}
}

func downloadZZStep(attachmentUrl string, referer string, filePath string) error {
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
	if ZZStepEnableHttpProxy {
		client = ZZStepSetHttpProxy()
	}
	req, err := http.NewRequest("GET", attachmentUrl, nil) //建立连接
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www2.zzstep.com")
	req.Header.Set("Referer", referer)
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //拿到返回的内容
	if err != nil {
		return err
	}
	// 检查HTTP响应头中的Content-Disposition字段获取文件名和后缀
	fileName := getZZStepFileNameFromHeader(resp)
	fileExtension := filepath.Ext(fileName) // 获取文件后缀
	fileExtArr := []string{".doc", ".docx", ".pdf", ".ppt", ".pptx"}
	fmt.Println("文件后缀:", fileExtension)
	if !StrInArrayZZStep(fileExtension, fileExtArr) {
		return errors.New("文件后缀：" + fileExtension + "不在下载后缀列表")
	}
	filePath += fileExtension

	defer resp.Body.Close()
	// 如果访问失败，就打印当前状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status :" + strconv.Itoa(resp.StatusCode))
	}

	// 创建一个文件用于保存
	fileDiv := filepath.Dir(filePath)
	if _, err = os.Stat(fileDiv); err != nil {
		if os.MkdirAll(fileDiv, 0777) != nil {
			return err
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// StrInArrayZZStep str in string list
func StrInArrayZZStep(str string, data []string) bool {
	if len(data) > 0 {
		for _, row := range data {
			if str == row {
				return true
			}
		}
	}
	return false
}

// 从HTTP响应头中获取文件名
func getZZStepFileNameFromHeader(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	fileName := ""
	if contentDisposition != "" {
		fileName = parseZZStepFileNameFromContentDisposition(contentDisposition)
	} else {
		fileName = filepath.Base(resp.Request.URL.Path) // 默认使用URL中的文件名作为本地文件名
	}
	return fileName
}

// 从Content-Disposition字段中解析文件名
func parseZZStepFileNameFromContentDisposition(contentDisposition string) string {
	// 参考：https://tools.ietf.org/html/rfc6266#section-4.3
	// 示例：attachment; filename="example.txt" -> example.txt
	fileNameStart := len("attachment; ") + len("filename=") + 2 // 2为引号的长度
	fileNameEnd := len(contentDisposition) - 1 - len("\"")      // 最后一个双引号的位置
	fileName := contentDisposition[fileNameStart:fileNameEnd]   // 提取文件名字符串
	return fileName[1:]                                         // 去掉字符串开头的引号（如果存在）并返回结果
}
