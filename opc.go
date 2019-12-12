package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/shuntaka9576/agouti"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

func cmdOpc(c *cli.Context) error {
	var cfg config

	err := cfg.load()
	if err != nil {
		return xerrors.Errorf("excute nikkeicmd config parse error:%v", err)
	}
	// create hash map
	opcMap := make(map[string]opc)
	for _, v := range cfg.Opc {
		opcMap[v.Args] = v
	}

	// merge default opc object and specified to two args object
	var requestOpcSetting opc
	if v, ok := opcMap["default"]; ok {
		requestOpcSetting = v
	} else {
		return xerrors.Errorf("please settings [[opc]] args default setting")
	}

	if args := c.Args().Get(0); args != "" {
		if argsOpcSetting, ok := opcMap[args]; ok {
			if val, ok := checkOpcValue(argsOpcSetting.Url); ok {
				requestOpcSetting.Url = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.UserId); ok {
				requestOpcSetting.UserId = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Password); ok {
				requestOpcSetting.Password = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Subject); ok {
				requestOpcSetting.Subject = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Content); ok {
				requestOpcSetting.Content = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Policy); ok {
				requestOpcSetting.Policy = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.UsrId); ok {
				requestOpcSetting.UsrId = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Date); ok {
				requestOpcSetting.Date = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.AccessStartTime); ok {
				requestOpcSetting.AccessStartTime = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.AccessEndTime); ok {
				requestOpcSetting.AccessEndTime = val
			}
			if val, ok := checkOpcValue(argsOpcSetting.Remark); ok {
				requestOpcSetting.Remark = val
			}
			if argsOpcSetting.Count != 0 {
				requestOpcSetting.Count = argsOpcSetting.Count
			}
		}
	}

	err = RunOpcReservation(cfg, requestOpcSetting)
	if err != nil {
		return xerrors.Errorf("faild to run chromedriver:%v", err)
	}
	return nil
}

func RunOpcReservation(cfg config, opcCfg opc) error {
	// create a date to fill the opc form
	var startTime string
	if t, err := time.Parse("15:04", opcCfg.AccessStartTime); err != nil {
		startTime = ""
	} else {
		startTime = t.Format("15:04")
	}

	var endTime string
	if t, err := time.Parse("15:04", opcCfg.AccessEndTime); err != nil {
		endTime = ""
	} else {
		endTime = t.Format("15:04")
	}

	date, err := time.Parse("2006/01/02", opcCfg.Date)
	if err != nil {
		date = time.Now() // 失敗したら現在時刻
	}

	var wg sync.WaitGroup
	for i := 0; i < opcCfg.Count; i++ {
		wg.Add(1)
		date := date.AddDate(0, 0, i).Format("2006/01/02")
		go func() {
			driver, page, _ := InitChromeDriver(cfg)
			ReserveOpc(driver, page, opcCfg, date, startTime, endTime)
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func ReserveOpc(driver *agouti.WebDriver, page *agouti.Page, opcCfg opc, date, startTime, endTime string) {
	if err := page.Navigate(opcCfg.Url); err != nil {
		fmt.Printf("Reserve opc error: %v", err)
		os.Exit(1)
	}

	// login
	page.FindByID("loginuid").Fill(opcCfg.UserId)
	page.FindByID("loginpwd").Fill(opcCfg.Password)
	page.FindByID("proceed").Click()

	// move to access form page
	page.Find(`#\31 0 > span`).Click()
	page.Find("#menu > ul > li:nth-child(2) > ul > li:nth-child(2)").Click()

	// fill opc form
	page.FindByID("_title_id").Fill(opcCfg.Subject)       // 件名
	page.FindByID("_description_id").Fill(opcCfg.Content) // 内容

	// ポリシー
	page.Find("#add").Click()
	page.Find("#row_" + opcCfg.UserId + "_pol > td:nth-child(8) > div > input").Click()
	page.FindByID("_okButton_id").Click()

	time.Sleep(1 * time.Second) // TODO issue #6

	// 開始日時
	page.Find("#d_startDatetimeScreen_id").Fill(date)
	if startTime != "" {
		page.Find("#t_startDatetimeScreen_id").Fill(startTime)
	}

	// 終了日時
	page.Find("#d_endDatetimeScreen_id").Fill(date)
	if endTime != "" {
		page.Find("#t_endDatetimeScreen_id").Fill(endTime)
	}

	page.Find("#_extra_id").Fill(opcCfg.Remark)
}

func checkOpcValue(s string) (string, bool) {
	switch s {
	case "[empty]":
		return "", true
	case "":
		return s, false // if empty character prefer default
	default:
		return s, true
	}
}
