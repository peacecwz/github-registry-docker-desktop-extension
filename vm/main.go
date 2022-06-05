package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/peacecwz/github-registry-docker-desktop-extension/vm/internal"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// TODO (peacecwz): Refactor each endpoints

func getOrganizations(ctx echo.Context) error {
	jsonMap := make(map[string]interface{})
	json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)

	token := jsonMap["token"].(string)

	organizationResult := internal.GetOrganizations(token)

	return ctx.JSON(http.StatusOK, organizationResult)
}

func getPackages(ctx echo.Context) error {
	jsonMap := make(map[string]interface{})
	json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)

	token := jsonMap["token"].(string)
	organizationId := jsonMap["organizationId"].(string)

	packagesResult := internal.GetPackages(token, organizationId)

	return ctx.JSON(http.StatusOK, packagesResult)
}

func listenOnUnixSocket(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

func main() {
	isBrowser := os.Getenv("BROWSER")

	router := echo.New()
	router.HideBanner = true
	startURL := ""

	if isBrowser == "true" {
		router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:8081"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
		startURL = ":8090"
	} else {
		var socketPath string
		flag.StringVar(&socketPath, "socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
		flag.Parse()
		os.RemoveAll(socketPath)
		logrus.New().Infof("Starting listening on %s\n", socketPath)

		ln, err := listenOnUnixSocket(socketPath)
		if err != nil {
			log.Fatal(err)
		}
		router.Listener = ln
	}

	router.POST("/organizations", getOrganizations)
	router.POST("/packages", getPackages)

	log.Fatal(router.Start(startURL))
}
