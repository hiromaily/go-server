package main

import (
	"context"
	"flag"
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
	"log"
	"os"
	"time"

	cdp "github.com/knq/chromedp"
	//cdptypes "github.com/knq/chromedp/cdp"
	"github.com/knq/chromedp/client"

	pg "github.com/hiromaily/go-server/chrome_devtools/pages"
)

var (
	driverName   = flag.String("d", "headless", "web driver name like `chrome`, `headless`, `edge`")
	headlessHost = flag.String("h", "chrome-headless", "when using on doker-compose, headless service name is used")
	//headlessPort = flag.Int("p", 9222, "port for Chrome DevTools Protocol")

	imagePath = flag.String("path", "/tmp/webdrive/images/", "captured image path")
	testName  = flag.String("n", "localhost", "test name")
)

var usage = `Usage: %s [options...]
Options:
  -d     web driver name like 'chrome', 'headless', 'edge'.
  -h     when using on doker-compose, headless service name is used
  //-p     port for Chrome DevTools Protocol
  -path  captured image path.
  -n     test name
e.g.:
  devtool -d chrome -path /tmp/gochromedev/
`

func init() {
	lg.InitializeLog(lg.DebugStatus, lg.LogOff, 99, "[DevToolsTest]", "/var/log/go/chrome_devtools.log")

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
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		tm.Track(time.Now(), "main()")
	}()

	lg.Infof("driver: %s / testName: %s / headlessHosst: %s ======================", *driverName, *testName, *headlessHost)

	c, err := createCDP(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	var site, res string

	//switch by testName
	switch *testName {
	case "google":
		// run task list
		err = c.Run(ctxt, pg.GoogleSearch("site:brank.as", "Easy Money Management", &site, &res))
		//fallthrough
	case "localhost":
		err = c.Run(ctxt, pg.CheckLocal(&site, &res))
		//fallthrough
	default:
		log.Fatal(fmt.Errorf("testName is invalid. (%s)", *testName))
	}
	if err != nil {
		log.Fatal(err)
	}

	//shutdown if chrome is run
	if *driverName == "chrome" {
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
	}

	lg.Infof("res `%s` site(%s)", res, site)
}

func createCDP(ctx context.Context) (*cdp.CDP, error) {
	//chrome
	if *driverName == "chrome" {
		return cdp.New(ctx, cdp.WithLog(log.Printf))
	}

	//headless
	cli := client.New(client.URL(fmt.Sprintf("http://%s:%d/json", *headlessHost, 9222)))
	return cdp.New(ctx, cdp.WithTargets(cli.WatchPageTargets(ctx)), cdp.WithLog(log.Printf))
}
