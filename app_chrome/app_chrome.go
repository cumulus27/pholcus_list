package pholcus_list

// 基础包
import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"log"
	"strings"
	"time"
)

func init() {
	FileTest.Register()
}

var flagDevToolWsUrl = flag.String("devtools-ws-url", "ws://172.105.194.56:9222/devtools/page/2682788AD45F968BF8F8FB2DFAF1DD33", "DevTools WebSsocket URL")

var FileTest = &Spider{
	Name:        "应用宝-动态滚动",
	Description: "抓取应用宝列表和动态滚动标签",
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
					"标志",
				},

				ParseFunc: func(ctx *Context) {

					flag.Parse()
					if *flagDevToolWsUrl == "" {
						log.Fatal("must specify -devtools-ws-url")
					}
					//create allocator context for use with creating a browser context later
					allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), *flagDevToolWsUrl)
					defer cancel()
					//create context
					ctxt, cancel := chromedp.NewContext(allocatorContext)
					defer cancel()

					sel := `.load-more-btn`

					// navigate
					if err := chromedp.Run(ctxt, chromedp.Navigate(`https://sj.qq.com/myapp/category.htm?orgame=1`)); err != nil {
						log.Printf("could not navigate to qq: %v", err)
					}

					// wait visible
					if err := chromedp.Run(ctxt, chromedp.WaitVisible(sel)); err != nil {
						log.Printf("could not get section: %v", err)
					}

					for i := 0; i < 5; i++ {
						if err := chromedp.Run(ctxt, chromedp.ScrollIntoView(sel)); err != nil {
							log.Printf("could not scroll to section: %v", err)
						}
						time.Sleep(time.Second * 2)
						log.Printf("Get a srcoll.")
					}

					var lastText string
					if err := chromedp.Run(ctxt, chromedp.Text(sel, &lastText)); err != nil {
						log.Printf("could not get last text: %v", err)
					}

					log.Printf("Get last mode: %s", lastText)

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
						longS = longS[2 : len(longS)-1]
						appDownload := strings.Trim(longS, " ")

						//链接
						url, _ := s.Find(".app-info-desc a").Attr("href")

						//输出格式
						ctx.Output(map[int]interface{}{
							0: appTitle,
							1: appSize,
							2: appDownload,
							3: url,
							4: lastText,
						})
					})
				},
			},
		},
	},
}
