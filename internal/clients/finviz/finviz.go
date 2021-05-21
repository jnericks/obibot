package finviz

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jnericks/obibot/internal/log"
)

var DefaultExecAllocatorOptions = []chromedp.ExecAllocatorOption{
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	// chromedp.Headless,

	// After Puppeteer's default behavior.
	chromedp.Flag("disable-background-networking", true),
	chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
	chromedp.Flag("disable-background-timer-throttling", true),
	chromedp.Flag("disable-backgrounding-occluded-windows", true),
	chromedp.Flag("disable-breakpad", true),
	chromedp.Flag("disable-client-side-phishing-detection", true),
	chromedp.Flag("disable-default-apps", true),
	chromedp.Flag("disable-dev-shm-usage", true),
	chromedp.Flag("disable-extensions", true),
	chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
	chromedp.Flag("disable-hang-monitor", true),
	chromedp.Flag("disable-ipc-flooding-protection", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("disable-prompt-on-repost", true),
	chromedp.Flag("disable-renderer-backgrounding", true),
	chromedp.Flag("disable-sync", true),
	chromedp.Flag("force-color-profile", "srgb"),
	chromedp.Flag("metrics-recording-only", true),
	chromedp.Flag("safebrowsing-disable-auto-update", true),
	chromedp.Flag("enable-automation", true),
	chromedp.Flag("password-store", "basic"),
	chromedp.Flag("use-mock-keychain", true),
}

type Response struct {
	URL string
}

func GetHeatMap(ctx context.Context) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(ctx, DefaultExecAllocatorOptions...)
	defer cancel()

	ch := make(chan string)
	errs := make([]error, 3)
	for i := 0; i < 3; i++ {
		go func(ctx context.Context, i int) {
			r, err := oneTab(ctx)
			if err != nil {
				errs[i] = err
			} else {
				select {
				case ch <- r:
				case <-ctx.Done():
				}
			}
		}(ctx, i)
	}

	return &Response{URL: <-ch}, nil
}

func oneTab(ctx context.Context) (string, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://finviz.com/map.ashx`),
		chromedp.WaitVisible(`#share-map`, chromedp.ByID),
		chromedp.Click(`#share-map`),
		chromedp.WaitVisible(`#static`, chromedp.ByID),
		chromedp.Value(`#static`, &result),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func actionFuncLog(msg string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		log.Info(ctx, msg)
		return nil
	}
}

func line() string {
	_, _, l, _ := runtime.Caller(1)
	return strconv.Itoa(l)
}
