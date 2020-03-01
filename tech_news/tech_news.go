package pholcus_list

// 基础包
import (
	// "github.com/henrylee2cn/pholcus/common/goquery"                          //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需

	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用
	// "github.com/henrylee2cn/pholcus/logs"
	// net包
	// "net/http" //设置http.Header
	// "net/url"
	// 编码包
	// "encoding/xml"
	//"encoding/json"
	// 字符串处理包
	//"regexp"
	// "strconv"
	// "fmt"
	// "math"
	// "time"

	"github.com/henrylee2cn/pholcus/common/goquery"
)

func init() {
	FileTest.Register()
}

var FileTest = &Spider{
	Name:        "亿欧科技新闻",
	Description: "抓取亿欧科技新闻智慧城市板块",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  "https://www.iyiou.com/smartcity/",
				Rule: "智慧城市",
			})
		},

		Trunk: map[string]*Rule{

			"智慧城市": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取分页导航
					navBox := query.Find("#page-nav a")
					navBox.Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "新闻列表",
							})
						}

					})

				},
			},

			"新闻列表": {
				ItemFields: []string{
					"标题",
					"作者",
					"来源",
					"链接",
					"时间",
				},

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取新闻列表
					newList := query.Find(".newestArticleList li")
					newList.Each(func(i int, s *goquery.Selection) {
						//标题
						newsTitle, _ := s.Find("div[class=text fl]>a").Attr("title")
						//作者
						newsAuthor := s.Find(".name").Text()
						//时间
						newsTime := s.Find(".time").Text()
						//链接
						url, _ := s.Find("div[class=text fl]>a").Attr("href")

						//输出格式
						ctx.Output(map[int]interface{}{
							0: newsTitle,
							1: newsAuthor,
							2: "亿欧",
							3: url,
							4: newsTime,
						})
					})
				},
			},
		},
	},
}
