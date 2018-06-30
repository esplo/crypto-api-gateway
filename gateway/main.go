package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type apiData struct {
	Endpoints []endpoint `yaml:"endpoints"`
}
type endpoint struct {
	Cost int    `yaml:"cost"`
	Path string `yaml:"path"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	envman := newEnvman()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Mining-Authorization")
	r.Use(cors.New(config))
	r.Use(headerRequired())

	api, err := loadAPIData()
	if err != nil {
		panic("unable to load API data")
	}

	// set routing for original API
	for _, e := range api.Endpoints {
		// to bind "endpoint"
		r.GET(e.Path, func(ep endpoint) func(ctx *gin.Context) {
			return func(c *gin.Context) {
				err := withdraw(c, envman.CoinhiveSecret, c.GetHeader("X-Mining-Authorization"), ep.Cost)
				if err != nil {
					return
				}
				err = callOriginalAPI(c, envman.APIHost, ep.Path)
				if err != nil {
					return
				}
			}
		}(e))
	}

	http.Handle("/", r)
}

func loadAPIData() (d apiData, err error) {
	buf, err := ioutil.ReadFile("./api.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(buf, &d)
	if err != nil {
		return
	}

	return
}

func headerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("X-Mining-Authorization")
		fmt.Println(auth)
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid access"})
		}
		c.Next()
	}
}
