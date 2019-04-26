package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/shuntaka9576/agouti"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

type nikkei struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
	Start    string `toml:"start"`
	Times    int    `toml:"times"`
}

func cmdNikkei(c *cli.Context) error {
	var cfg config
	err := cfg.load()
	if err != nil {
		return xerrors.Errorf("excute nikkeicmd config parse error:%v", err)
	}

	// create image cache directory
	cacheDir := filepath.Join(cfg.Base.WorkDir, "nikkei", "cache", time.Now().Format("20060102-1504-05"))
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return xerrors.Errorf("create image cache directory error: %v", err)
	}

	err = RunChromeDriver(cfg, cacheDir)
	if err != nil {
		return xerrors.Errorf("faild to run chromedriver:%v", err)
	}
	return nil
}

func RunChromeDriver(cfg config, cacheDir string) error {
	var driver *agouti.WebDriver

	driver = agouti.ChromeDriver(cfg.ChromeDriver.DriverPath)
	if err := driver.Start(); err != nil {
		return xerrors.Errorf("chrome driver start error:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return xerrors.Errorf("initialize chrome browser error:%v", err)
	}

	for i := cfg.Nikkei.Times; i > 0; i-- {
		err := ScreenshotNikkeiVoucher(page, cacheDir, cfg.Nikkei, i)
		if err != nil {
			return xerrors.Errorf("take secrennshot nikkei voucher error:%v", err)
		}
		first, err := time.Parse("200601", cfg.Nikkei.Start)
		if err != nil {
			return xerrors.Errorf("invalid config start value.time parse error:%v", err)
		}
		next := first.AddDate(0, 1, 0)
		cfg.Nikkei.Start = next.Format("200601")
	}

	ConcatImage(cacheDir)
	return nil
}

func ScreenshotNikkeiVoucher(page *agouti.Page, cacheDir string, nikkeiConfig nikkei, count int) error {
	if err := page.Navigate("https://id.nikkei.com/charge/payment/history/search/" + nikkeiConfig.Start); err != nil {
		return xerrors.Errorf("navigate voucher page error:%v", err)
	}

	if count == nikkeiConfig.Times {
		page.FindByID("LA7110Form01:LA7110Email").Fill(nikkeiConfig.Email)
		page.FindByID("LA7110Form01:LA7110Password").Fill(nikkeiConfig.Password)
		page.Find("#LA7110Form01 > div > div.form__area > div > div.btn-box.btn-box--flex > input").Click()

		text, err := page.Find("#LA7110Form01 > div.form > div.form__area > div > dl:nth-child(1) > dd").Text()
		if text == "メールアドレスまたはパスワードが間違っています。" && err == nil {
			return xerrors.Errorf("ether email address or password is invalid")
		}
	}

	text, err := page.FindByClass("card_number").Text()
	if text == "お支払い履歴情報がありません。" && err == nil {
		return xerrors.Errorf("not found voucher page")
	}

	page.Find("#contentInner > div.table-history > table > tbody > tr > td.service > p.receipt_link > a").Click()
	page.Find("#contentInner > div.table-history > p > a").Click()
	page.NextWindow()
	page.Screenshot(filepath.Join(cacheDir, nikkeiConfig.Start+".png"))
	page.NextWindow()
	return nil
}

func cmdNikkeiClean(c *cli.Context) error {
	var cfg config
	err := cfg.load()
	if err != nil {
		return xerrors.Errorf("excute nikkeiCleancmd config parse error:%v", err)
	}
	cacheDir := filepath.Join(cfg.Base.WorkDir, "nikkei", "cache")
	if err := os.RemoveAll(cacheDir); err != nil {
		return xerrors.Errorf("remove faild directory %v: %v", cacheDir, err)
	}
	return nil
}
