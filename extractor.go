package fms

import (
	"errors"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

type cookieExtractor struct {
	Email    string
	Password string
}

func getDriver() (*agouti.WebDriver, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	if err := driver.Start(); err != nil {
		return nil, err
	}
	return driver, nil
}

func getLoginPage(driver *agouti.WebDriver) (*agouti.Page, error) {
	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}
	if err := page.Navigate(`https://account.line.biz/login?scope=line&redirectUri=https%3A%2F%2Fdevelopers.line.biz%2Fflex-simulator%2F`); err != nil {
		return nil, err
	}
	err = page.Find("body > div > div > div:nth-child(3) > div > form > div > input").Submit()
	if err != nil {
		return nil, err
	}
	return page, nil
}

func waitToLoad(page *agouti.Page) error {
	for {
		html, err := page.HTML()
		if err != nil {
			return err
		}
		if strings.Contains(html, "入力内容に誤りがあります。") {
			return errors.New("Invalid email or password")
		}
		title, err := page.Title()
		if err != nil {
			return err
		}
		if title == "Flex Message Simulator" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func getCookie(page *agouti.Page) (string, error) {

	// get all cookies of the page
	cookies, err := page.GetCookies()
	if err != nil {
		return "", err
	}

	// format the cookies into a shape that can be set for HTTP requests
	var reqCookie string
	for _, cookie := range cookies {
		cookieString := cookie.String()
		reqCookie += cookieString[0:strings.Index(cookieString, ";")] + "; "
	}

	return reqCookie, nil
}

// getCookie is a function that gets the cookies of the user
func (extractor *cookieExtractor) getCookie() (string, error) {

	// get chromedriver
	driver, err := getDriver()
	if err != nil {
		return "", err
	}
	defer driver.Stop()

	// get the login page
	page, err := getLoginPage(driver)
	if err != nil {
		return "", err
	}

	// get the elements
	fieldset := page.Find("#app > div > div > div > div.MdBox01 > div > form > fieldset")
	emailInput := fieldset.Find("div:nth-child(2) > input[type=text]")
	passwordInput := fieldset.Find("div:nth-child(3) > input[type=password]")
	loginButton := fieldset.Find("div.mdFormGroup01Btn > button")

	// login
	emailInput.Fill(extractor.Email)
	passwordInput.Fill(extractor.Password)
	err = loginButton.Submit()
	if err != nil {
		return "", err
	}

	// wait to load the page
	err = waitToLoad(page)
	if err != nil {
		return "", err
	}

	// get the cookies of the page
	cookie, err := getCookie(page)
	if err != nil {
		return "", err
	}

	return cookie, nil
}
