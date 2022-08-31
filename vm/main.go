package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	GithubClientId = "adbb61eb6d9e0267a7df"
	GrantType      = "urn:ietf:params:oauth:grant-type:device_code"
	Scope          = "user,repo,read:packages,delete:packages,read:org"
)

var deviceCode, accessToken = "", ""

func main() {
	var socketPath string
	flag.StringVar(&socketPath, "socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	flag.Parse()

	os.RemoveAll(socketPath)

	logrus.New().Infof("Starting listening on %s\n", socketPath)
	router := echo.New()
	router.HideBanner = true

	startURL := ""

	ln, err := listen(socketPath)
	if err != nil {
		log.Fatal(err)
	}
	router.Listener = ln

	router.GET("/me", getMe)
	router.GET("/auth", auth)
	router.GET("/auth/complete", authComplete)
	router.GET("/logout", logout)
	router.GET("/organizations", getOrganizations)
	router.GET("/packages", getPackages)
	router.GET("/package-delete", deletePackage)
	log.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

func getOrganizations(ctx echo.Context) error {
	organizationResult, err := GetOrganizations(accessToken)
	if err != nil {
		log.Fatalln(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, organizationResult)
}

func getPackages(ctx echo.Context) error {
	organizationId := ctx.QueryParam("organizationId")

	packagesResult := GetPackages(accessToken, organizationId)

	return ctx.JSON(http.StatusOK, packagesResult)
}

func deletePackage(ctx echo.Context) error {
	organizationName := ctx.QueryParam("organizationName")
	packageName := ctx.QueryParam("packageName")
	packageVersionIdStr := ctx.QueryParam("packageVersionId")
	packageVersionId, err := strconv.ParseInt(packageVersionIdStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	err = DeletePackage(accessToken, organizationName, packageName, int64(packageVersionId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, "")
}

func auth(ctx echo.Context) error {
	userCode, err := Auth()
	if err != nil {
		log.Fatalln(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, userCode)
}

func authComplete(ctx echo.Context) error {
	token, err := AuthComplete()
	if err != nil {
		log.Fatalln(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, token)
}

func getMe(ctx echo.Context) error {
	user, err := GetUser()
	if err != nil {
		log.Fatalln(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, user)
}

func logout(ctx echo.Context) error {
	deviceCode = ""
	accessToken = ""
	return ctx.JSON(http.StatusOK, "")
}
