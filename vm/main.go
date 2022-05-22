package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

// TODO (peacecwz): Refactor each endpoints

type GetOrganizationsResponses struct {
	Message       string
	Organizations []Organization
}

type GetPackagesResponses struct {
	Message  string
	Packages []Package
}

type GetPackageVersionsResponses struct {
	Message         string
	PackageVersions []Package
}

type Organization struct {
	Name string
	Id   int64
}

type Package struct {
	RepositoryUrl  string
	RepositoryName string
	Url            string
	PackageType    string
	Name           string
	Id             int64
}

func getOrganizations(ctx echo.Context) error {
	jsonMap := make(map[string]interface{})
	json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)

	token := jsonMap["token"].(string)

	ct := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ct, ts)

	client := github.NewClient(tc)

	orgs, resp, err := client.Organizations.List(ct, "", nil)
	if err != nil {
		return ctx.JSON(http.StatusOK, GetOrganizationsResponses{Message: fmt.Sprintf("\n error: %v\n", err)})
	}

	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
		return ctx.JSON(http.StatusOK, GetOrganizationsResponses{Message: "Token expiration"})
	}

	var organizations []Organization

	for _, org := range orgs {
		if org.Login != nil && *org.Login != "" {
			organizations = append(organizations, Organization{
				Name: *org.Login,
				Id:   *org.ID,
			})
		}
	}

	return ctx.JSON(http.StatusOK, GetOrganizationsResponses{Organizations: organizations})
}

func getPackages(ctx echo.Context) error {
	jsonMap := make(map[string]interface{})
	json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)

	token := jsonMap["token"].(string)
	organizationName := jsonMap["organizationId"].(string)

	ct := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ct, ts)

	client := github.NewClient(tc)
	packageType := "container"
	githubPackages, resp, err := client.Organizations.ListPackages(ct, organizationName, &github.PackageListOptions{
		PackageType: &packageType,
	})
	if err != nil {
		return ctx.JSON(http.StatusOK, GetPackagesResponses{Message: fmt.Sprintf("\n error: %v\n", err)})
	}

	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
		return ctx.JSON(http.StatusOK, GetPackagesResponses{Message: "Token expiration"})
	}

	var packages []Package

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {
			packages = append(packages, Package{
				Name: *githubPackage.Name,
				Id:   *githubPackage.ID,
				Url:  *githubPackage.URL,
			})
		}
	}

	return ctx.JSON(http.StatusOK, GetPackagesResponses{Packages: packages})
}

func getPackageVersions(ctx echo.Context) error {
	jsonMap := make(map[string]interface{})
	json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)

	token := jsonMap["token"].(string)
	organizationName := jsonMap["organizationId"].(string)
	packageName := jsonMap["packageName"].(string)

	ct := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ct, ts)

	client := github.NewClient(tc)
	packageType := "container"
	githubPackages, resp, err := client.Organizations.PackageGetAllVersions(ct, organizationName, packageType, url.QueryEscape(packageName), nil)
	if err != nil {
		return ctx.JSON(http.StatusOK, GetPackageVersionsResponses{Message: fmt.Sprintf("\n error: %v\n", err)})
	}

	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
		return ctx.JSON(http.StatusOK, GetPackageVersionsResponses{Message: "Token expiration"})
	}

	var packages []Package

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {
			packages = append(packages, Package{
				Name: githubPackage.Metadata.Container.Tags[0],
				Id:   *githubPackage.ID,
			})
		}
	}

	return ctx.JSON(http.StatusOK, GetPackageVersionsResponses{PackageVersions: packages})
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
		startURL = ":8080"
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
	router.POST("/package-versions", getPackageVersions)

	log.Fatal(router.Start(startURL))
}
