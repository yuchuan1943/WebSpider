package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	xmlpath "gopkg.in/xmlpath.v1"
)

// DownloadFile 下载文件
func DownloadFile(url string, path string, refener string) bool {
	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper("get"), url, strings.NewReader(""))
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")
	req.Header.Set("Referer", refener)

	res, err := client.Do(req)
	if err != nil {
		return false
	}

	defer res.Body.Close()

	dirPath := pathRegexp.FindString(path)

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return false
	}

	f, err := os.Create(path)
	if err != nil {
		return false
	}
	io.Copy(f, res.Body)
	return true
}

// GetHTML 获取HTML
func GetHTML(method string, url string) (html string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(""))

	if err != nil {
		return "", nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body[:]), nil
}

// FindImageURL 发现网址中的图片链接
func FindImageURL(url string, deep int) {
	execed := false
	for i := 0; i < len(urls); i++ {
		if urls[i] == url {
			execed = true
			break
		}
	}
	deep--
	if validHref.MatchString(url) && !execed && deep > 0 && !hrefFilter.MatchString(url) {
		defer func() {

		}()
		urls = append(urls, url)
		html, err := GetHTML("get", url)
		if err != nil {
			return
		}

		path, err := xmlpath.Compile("//@href")
		if err != nil {
			return
		}
		root, err := xmlpath.ParseHTML(strings.NewReader(html))
		if err != nil {
			return
		}
		iter := path.Iter(root)
		for iter.Next() {
			href := iter.Node().String()
			FindImageURL(href, deep)
		}

		pathImg, err := xmlpath.Compile("//@src")
		if err != nil {
			return
		}
		iterImg := pathImg.Iter(root)
		for iterImg.Next() {
			src := iterImg.Node().String()
			if imageRegexp.MatchString(src) {
				downloaded := false
				for i := 0; i < len(images); i++ {
					if images[i] == src {
						downloaded = true
						break
					}
				}
				if !downloaded {
					images = append(images, src)
					fmt.Println(src)
					go DownloadFile(src, strings.Replace(hostRegexp.ReplaceAllString(src, imagePathRoot), "/", "\\", -1), url)
				}
			}
		}
	}
}

// GetCurrentPath 获取当前路径
func GetCurrentPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return path
}

var urls []string
var images []string
var hrefFilter = regexp.MustCompile("(.*\\.ico)|(.*\\.css)|(.*\\.js)|(.*\\.png)")
var imageRegexp = regexp.MustCompile("http(s?)://.*\\.jpg")
var validHref = regexp.MustCompile("http(s?)://.*")
var hostRegexp = regexp.MustCompile("http(s?)://.*?/")
var pathRegexp = regexp.MustCompile(".*\\\\")
var imagePathRoot string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	imagePathRoot = GetCurrentPath() + "\\Images\\"
	FindImageURL("http://test.com", 10)
	fmt.Scanln()
	// MainWindow{
	// 	Title:   "SCREAMO",
	// 	MinSize: Size{800, 640},
	// 	Layout:  VBox{},
	// 	Children: []Widget{
	// 		HSplitter{
	// 			Children: []Widget{
	// 				TextEdit{AssignTo: &outTE, ReadOnly: true},
	// 			},
	// 		},
	// 		PushButton{
	// 			Text: "开始下载",
	// 			OnClicked: func() {
	// 				go FindImageURL("http://test.com", 5)
	// 				outTE.AppendText("任务完成\n")
	// 			},
	// 		},
	// 	},
	// }.Run()
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		outTE.AppendText("报错啦")
	// 	}
	// }()
}
