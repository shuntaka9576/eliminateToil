package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

const (
	configFile = "settings.toml"
)

type base struct {
	WorkDir   string `toml:"base"`
	Tomlepath string `toml:"tomlPath"`
}

type chromeDriver struct {
	Version    string `toml:"version"`
	DriverPath string `toml:"driverPath"`
}

type config struct {
	Base         base         `toml:"base"`
	Nikkei       nikkei       `toml:"nikkei"`
	ChromeDriver chromeDriver `toml:"chromeDriver"`
}

func (cfg *config) load() error {
	// Create application folder variable.
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "eliminateToil")
		} else {
			dir = filepath.Join(dir, "eliminateToil")
		}
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "eliminateToil")
	}

	/*
		Search config file in current and app directory.
		If it does not exist, create it as the current.
	*/
	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			return xerrors.Errorf("read config file error: %v", err)
		}

		_, err = toml.Decode(string(stripBOM(config)), cfg)
		if err != nil {
			return xerrors.Errorf("toml config parse error:%v", err)
		}
	} else if _, err := os.Stat(filepath.Join(dir, configFile)); err == nil {
		if err == nil {
			config, err := ioutil.ReadFile(filepath.Join(dir, configFile))
			if err != nil {
				return xerrors.Errorf("read config file error: %v", err)
			}

			_, err = toml.Decode(string(stripBOM(config)), cfg)
			if err != nil {
				return xerrors.Errorf("toml config parse error:%v", err)
			}
		}
	} else {
		f, err := os.Create(configFile)
		if err != nil {
			return xerrors.Errorf("create config file error:%v", err)
		}
		cfg.Nikkei.Email = "hogehoge@gmail.com"
		cfg.Nikkei.Password = "hogehoge"
		cfg.Nikkei.Start = "201911"
		cfg.Nikkei.Times = 6
		return toml.NewEncoder(f).Encode(cfg)
	}

	// Settings config.
	cfg.Base.WorkDir = dir
	if err := os.MkdirAll(dir, 0700); err != nil {
		return xerrors.Errorf("cannot create %v directory:%v", dir, err)
	}
	cfg.ChromeDriver.DriverPath = filepath.Join(dir, "chromedriver", cfg.ChromeDriver.Version)

	//if the chromedriver does not exist, install it.
	if _, err := os.Stat(cfg.ChromeDriver.DriverPath); err != nil {
		driverURL := "https://chromedriver.storage.googleapis.com"
		driverURL += "/" + cfg.ChromeDriver.Version + "/"
		if err := os.MkdirAll(cfg.ChromeDriver.DriverPath, 0700); err != nil {
			return xerrors.Errorf("cannot create %v directory: %v", cfg.ChromeDriver.DriverPath, err)
		}
		switch runtime.GOOS {
		case "windows":
			driverURL += "chromedriver_win32.zip"
		case "darwin":
			driverURL += "chromedriver_mac64.zip"
		case "linux":
			driverURL += "chromedriver_linux64.zip"
		}

		resp, err := http.Get(driverURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(filepath.Join(cfg.ChromeDriver.DriverPath, "chromedriver.zip"))
		if err != nil {
			return err
		}
		defer out.Close()

		io.Copy(out, resp.Body)
		Unzip(filepath.Join(cfg.ChromeDriver.DriverPath, "chromedriver.zip"), cfg.ChromeDriver.DriverPath)
	}
	return nil
}
