package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetUser() (*GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", accessToken)},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user *GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}
	user.Token = accessToken
	return user, nil
}

func Auth() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://github.com/login/device/code?client_id=%s&scope=%s", GithubClientId, Scope), nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	deviceCode = q.Get("device_code")

	return q.Get("user_code"), nil
}

func AuthComplete() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&device_code=%s&grant_type=%s", GithubClientId, deviceCode, GrantType), nil)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	accessToken = q.Get("access_token")
	return accessToken, nil
}

func GetOrganizations(token string) (Organization, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/orgs", nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var organizations Organization
	err = json.Unmarshal(body, &organizations)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return organizations, err
}

func DeletePackage(token, organizationId, packageName string, packageVersionId int64) error {
	client, ctx := initializeGithubClient(token)
	packageType := "container"
	isOrg := username != organizationId
	
	var resp *github.Response
	var err error
	if isOrg {
		resp, err = client.Organizations.PackageDeleteVersion(ctx, organizationId, packageType, url.QueryEscape(packageName), packageVersionId)
	} else {
		resp, err = client.Users.PackageDeleteVersion(ctx, organizationId, packageType, url.QueryEscape(packageName), packageVersionId)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New("cannot delete it")
	}

	return nil
}

func GetPackages(token, organizationId string) GetPackagesResponses {
	client, ctx := initializeGithubClient(token)

	packageType := "container"
	isOrg := username != organizationId

	var githubPackages []*github.Package
	var err error

	if isOrg {
		githubPackages, _, err = client.Organizations.ListPackages(ctx, organizationId, &github.PackageListOptions{
			PackageType: &packageType,
		})
	} else {
		githubPackages, _, err = client.Users.ListPackages(ctx, organizationId, &github.PackageListOptions{
			PackageType: &packageType,
		})
	}

	if err != nil {
		return GetPackagesResponses{Message: fmt.Sprintf("\n error: %v\n", err)}
	}

	var packages []Package

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {

			packages = append(packages, Package{
				Name:            *githubPackage.Name,
				Id:              *githubPackage.ID,
				Url:             *githubPackage.HTMLURL,
				PackageVersions: getPackageVersions(client, ctx, organizationId, *githubPackage.Owner.Login, *githubPackage.Name, isOrg),
			})
		}
	}

	return GetPackagesResponses{Packages: packages}
}

func getPackageVersions(client *github.Client, ct context.Context, organizationId string, organizationName string, packageName string, isOrg bool) []PackageVersion {
	var githubPackages []*github.PackageVersion
	if isOrg {
		githubPackages, _, _ = client.Organizations.PackageGetAllVersions(ct, organizationId, "container", url.QueryEscape(packageName), nil)
	} else {
		githubPackages, _, _ = client.Users.PackageGetAllVersions(ct, organizationId, "container", url.QueryEscape(packageName), nil)
	}

	var packages []PackageVersion

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {
			packages = append(packages, PackageVersion{
				Name:             githubPackage.Metadata.Container.Tags[0],
				Id:               *githubPackage.ID,
				PackageType:      "container",
				CreatedAt:        Timestamp(*githubPackage.CreatedAt),
				Url:              *githubPackage.HTMLURL,
				OrganizationName: organizationName,
				ImageUrl:         fmt.Sprintf("ghcr.io/%s/%s@%s", organizationName, packageName, *githubPackage.Name),
			})
		}
	}

	return packages
}

func initializeGithubClient(token string) (*github.Client, context.Context) {
	ct := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
			TokenType:   "bearer",
		},
	)
	tc := oauth2.NewClient(ct, ts)

	return github.NewClient(tc), ct
}
