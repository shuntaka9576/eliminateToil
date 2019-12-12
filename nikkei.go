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

	// initialize chrome driver
	driver, page, err := InitChromeDriver(cfg)
	if err != nil {
		return xerrors.Errorf("faild to initialize chrome driver:%v", err)
	}

	// start collect voucher
	err = RunCollectVoucher(driver, page, cacheDir, cfg.Nikkei)
	if err != nil {
		return xerrors.Errorf("faild to collect Voucher:%v", err)
	}

	// concat images
	err = ConcatImage(cacheDir)
	if err != nil {
		return xerrors.Errorf("faild to concat images:%v", err)
	}
	return nil
}

func InitChromeDriver(cfg config) (driver *agouti.WebDriver, page *agouti.Page, err error) {
	driver = agouti.ChromeDriver(cfg.ChromeDriver.DriverPath)
	if err = driver.Start(); err != nil {
		return driver, page, xerrors.Errorf("chrome driver start error:%v", err)
	}
	page, err = driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return driver, page, xerrors.Errorf("initialize chrome browser error:%v", err)
	}

	return driver, page, nil
}

func RunCollectVoucher(driver *agouti.WebDriver, page *agouti.Page, cacheDir string, nikkeiCfg nikkei) error {
	defer driver.Stop()
	for i := nikkeiCfg.Times; i > 0; i-- {
		err := ScreenshotNikkeiVoucher(page, cacheDir, nikkeiCfg, i)
		if err != nil {
			return xerrors.Errorf("take secrennshot nikkei voucher error:%v", err)
		}
		first, err := time.Parse("200601", nikkeiCfg.Start)
		if err != nil {
			return xerrors.Errorf("invalid config start value.time parse error:%v", err)
		}
		next := first.AddDate(0, 1, 0)
		nikkeiCfg.Start = next.Format("200601")
	}
	return nil
}

func ScreenshotNikkeiVoucher(page *agouti.Page, cacheDir string, nikkeiCfg nikkei, count int) error {
	if err := page.Navigate("https://id.nikkei.com/charge/payment/history/search/" + nikkeiCfg.Start); err != nil {
		return xerrors.Errorf("navigate voucher page error:%v", err)
	}

	if count == nikkeiCfg.Times {
		page.FindByID("LA7110Form01:LA7110Email").Fill(nikkeiCfg.Email)
		page.FindByID("LA7110Form01:LA7110Password").Fill(nikkeiCfg.Password)
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
	page.Screenshot(filepath.Join(cacheDir, nikkeiCfg.Start+".png"))
	page.NextWindow()
	return nil
}

func cmdNikkeiClean(c *cli.Context) error {
	var cfg config
	err := cfg.load()
	if err != nil {
		return xerrors.Errorf("execute nikkei clean config parse error:%v", err)
	}
	cacheDir := filepath.Join(cfg.Base.WorkDir, "nikkei", "cache")
	if err := os.RemoveAll(cacheDir); err != nil {
		return xerrors.Errorf("faild to remove directory %v: %v", cacheDir, err)
	}
	return nil
}
