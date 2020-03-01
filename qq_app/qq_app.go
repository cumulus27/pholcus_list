package pholcus_list

// 基础包
import (
	"strings"
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
	Name:        "应用宝",
	Description: "抓取应用宝应用列表",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  "https://sj.qq.com/myapp/category.htm?orgame=1",
				Rule: "应用列表",
			})
		},

		Trunk: map[string]*Rule{

			"应用列表": {
				ItemFields: []string{
					"名称",
					"大小",
					"下载量",
					"链接",
				},

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取应用列表
					newList := query.Find(".main li")
					newList.Each(func(i int, s *goquery.Selection) {
						//名称
						appTitle := s.Find(".app-info-desc a").Text()
						//大小
						appSize := s.Find(".size").Text()
						//下载量
						longS := s.Find(".download").Text()
						longS = longS[2 : end-1]
						appDownload := strings.Trim(longS, " ")

						//链接
						url, _ := s.Find(".app-info-desc a").Attr("href")

						//输出格式
						ctx.Output(map[int]interface{}{
							0: appTitle,
							1: appSize,
							2: appDownload,
							3: url,
						})
					})
				},
			},
		},
	},
}
