package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type coinhiveWithdrawalResponse struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Error   string `json:"error"`
}

func withdraw(c *gin.Context, secret string, name string, amount int) error {
	ctx := appengine.NewContext(c.Request)
	log.Debugf(ctx, fmt.Sprintf("[Withdrawal] start: %s - %d\n", name, amount))

	resp, err := makeWithdrawalRequest(&ctx, secret, name, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}
	body := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": resp.StatusCode,
			"message":     body,
		})
		return err
	}

	log.Debugf(ctx, fmt.Sprintf("[Withdrawal] Coinhive response: %s - %s\n", resp.Status, body))

	// Coinhive API returns 200 even if the request was failed
	respBodyJSON := coinhiveWithdrawalResponse{}
	err = json.Unmarshal(bodyBytes, &respBodyJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	if !respBodyJSON.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"message":     respBodyJSON.Error,
		})
		return errors.New("failed to call withdrawal API")
	}

	log.Debugf(ctx, fmt.Sprintf("[Withdrawal] success: name(%s),amount(%d)\n", name, amount))

	return nil
}

func makeWithdrawalRequest(ctx *context.Context, secret string, name string, amount int) (*http.Response, error) {
	client := urlfetch.Client(*ctx)

	data := url.Values{}
	data.Set("secret", secret)
	data.Set("name", name)
	data.Set("amount", fmt.Sprintf("%d", amount))

	req, err := http.NewRequest("POST", "https://api.coinhive.com/user/withdraw", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.Do(req)
}
