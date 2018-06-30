package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// TODO: parameters, body, header, ...
func callOriginalAPI(c *gin.Context, host string, dir string) error {
	ctx := appengine.NewContext(c.Request)

	originalURL, err := url.Parse(host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}
	originalURL.Path = path.Join(originalURL.Path, dir)

	resp, err := makeRequest(&ctx, originalURL)
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
	log.Debugf(ctx, fmt.Sprintf("[CallOriginalAPI] url: %s - status: %s - %s\n", originalURL, resp.Status, body))

	var responseJSON map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	c.JSON(resp.StatusCode, responseJSON)

	return nil
}

func makeRequest(ctx *context.Context, url fmt.Stringer) (*http.Response, error) {
	client := urlfetch.Client(*ctx)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MINING-GATEWAY-KEY", "ok")

	return client.Do(req)
}
