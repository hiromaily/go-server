package pages

import (
	"context"
	//"fmt"
	cdp "github.com/knq/chromedp"
	cdptypes "github.com/knq/chromedp/cdp"
	"io/ioutil"
	"time"
)

func CheckLocal(site, res *string) cdp.Tasks {
	var buf []byte
	//sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	return cdp.Tasks{
		//1.TOP
		cdp.Navigate(`http://web/global`),
		cdp.Sleep(2 * time.Second),
		//cdp.WaitVisible(`#body`, cdp.ByID),
		cdp.WaitReady(`.outer-nav`, cdp.ByQuery),
		cdp.Screenshot(`div.intro`, &buf, cdp.ByQuery),
		cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
			return ioutil.WriteFile("images/01.png", buf, 0644)
		}),

		//2.Accounts
		//cdp.Click(`li:nth-child(1)`, cdp.ByQuery),
		//cdp.WaitVisible(`div.work`, cdp.ByQuery),
		//cdp.Screenshot(`div.work`, &buf, cdp.ByQuery),
		//cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
		//	return ioutil.WriteFile("02.png", buf, 0644)
		//}),

		//3.
		//cdp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, cdp.ByQuery),
		//cdp.Location(site),
		//cdp.Screenshot(`#testimonials`, &buf, cdp.ByID),
		//cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
		//	return ioutil.WriteFile("testimonials.png", buf, 0644)
		//}),
	}
}
