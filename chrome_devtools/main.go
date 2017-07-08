package main

import (
	"context"
	"flag"
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
	"io/ioutil"
	"log"
	"os"
	"time"

	cdp "github.com/knq/chromedp"
	cdptypes "github.com/knq/chromedp/cdp"
	"github.com/knq/chromedp/client"
)

var (
	driverName = flag.String("d", "headless", "web driver")
	imagePath  = flag.String("path", "/tmp/webdrive/images/", "captured image path")
	testURL    = flag.String("u", "https://stackoverflow.com/", "test URL")
)

var usage = `Usage: %s [options...]
Options:
  -d     web driver.
  -path  captured image path.
e.g.:
  gochromedev -d chrome -path /tmp/gochromedev/
`

func init() {
	lg.InitializeLog(lg.DebugStatus, lg.LogOff, 99, "[GoChromeDevToolsTest]", "/var/log/go/chrome_devtools.log")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}

	flag.Parse()

	if *driverName == "" {
		flag.Usage()

		os.Exit(1)
		return
	}
}

func main() {
	defer tm.Track(time.Now(), "main()")

	switch *driverName {
	case "chrome":
		chrome()
		//fallthrough
	case "headless":
		headless()
		//fallthrough
	case "localhost":
		localhost()
		//fallthrough
	default:
		headless()
	}
}

func localhost() {
	lg.Debug("============ localhost ============")
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome
	cli := client.New(client.URL(fmt.Sprintf("http://chrome-headless:%d/json", 9222)))
	//client.New("http://chrome-headless:9222/json")
	c, err := cdp.New(ctxt, cdp.WithTargets(cli.WatchPageTargets(ctxt)), cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, checkLocal(&site, &res))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("saved screenshot of #testimonials from search result listing `%s` (%s)", res, site)
}

func headless() {
	lg.Debug("============ headless ============")
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome
	cli := client.New(client.URL(fmt.Sprintf("http://chrome-headless:%d/json", 9222)))
	c, err := cdp.New(ctxt, cdp.WithTargets(cli.WatchPageTargets(ctxt)), cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, googleSearch("site:brank.as", "Easy Money Management", &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("saved screenshot of #testimonials from search result listing `%s` (%s)", res, site)
}

func chrome() {
	lg.Debug("============ chrome ============")

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, googleSearch("site:brank.as", "Easy Money Management", &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("saved screenshot of #testimonials from search result listing `%s` (%s)", res, site)
}

func checkLocal(site, res *string) cdp.Tasks {
	var buf []byte
	//sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	return cdp.Tasks{
		//1.TOP
		cdp.Navigate(`http://web`),
		cdp.Sleep(10 * time.Second),
		cdp.WaitVisible(`#body`, cdp.ByID),

		//2.Accounts
		cdp.Click(`#btn1`, cdp.ByQuery),
		//cdp.WaitVisible(`#ext-comp-1060_header-title-textEl`, cdp.ByQuery),
		cdp.Screenshot(`#body`, &buf, cdp.ByID),
		cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
			return ioutil.WriteFile("index.png", buf, 0644)
		}),
		//3.
		//cdp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, cdp.ByQuery),
		//cdp.Location(site),
		//cdp.Screenshot(`#testimonials`, &buf, cdp.ByID),
		//cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
		//	return ioutil.WriteFile("testimonials.png", buf, 0644)
		//}),
	}
}

func googleSearch(q, text string, site, res *string) cdp.Tasks {
	var buf []byte
	sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	return cdp.Tasks{
		cdp.Navigate(`https://www.google.com`),
		cdp.Sleep(2 * time.Second),
		cdp.WaitVisible(`#hplogo`, cdp.ByID),
		cdp.SendKeys(`#lst-ib`, q+"\n", cdp.ByID),
		cdp.WaitVisible(`#res`, cdp.ByID),
		cdp.Text(sel, res),
		cdp.Click(sel),
		cdp.Sleep(2 * time.Second),
		cdp.WaitVisible(`#footer`, cdp.ByQuery),
		cdp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, cdp.ByQuery),
		cdp.Location(site),
		cdp.Screenshot(`#testimonials`, &buf, cdp.ByID),
		cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
			return ioutil.WriteFile("testimonials.png", buf, 0644)
		}),
	}
}
