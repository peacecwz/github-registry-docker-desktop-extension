package internal

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type Organization struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetOrganizationsResponses struct {
	Message       string         `json:"message"`
	Organizations []Organization `json:"organizations"`
}

type Package struct {
	Id              int64            `json:"id"`
	Name            string           `json:"name"`
	PackageType     string           `json:"package_type"`
	PackageVersions []PackageVersion `json:"package_versions"`
	RepositoryUrl   string           `json:"repository_url"`
	RepositoryName  string           `json:"repository_name"`
	Url             string           `json:"url"`
}

type PackageVersion struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	PackageType    string `json:"package_type"`
	RepositoryUrl  string `json:"repository_url"`
	RepositoryName string `json:"repository_name"`
	Url            string `json:"url"`
}

type GetPackagesResponses struct {
	Message  string    `json:"message"`
	Packages []Package `json:"packages"`
}

func GetOrganizations(token string) GetOrganizationsResponses {
	client, ctx := initializeGithubClient(token)

	orgs, resp, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		return GetOrganizationsResponses{Message: fmt.Sprintf("\n error: %v\n", err)}
	}

	isExpired, result := checkTokenExpiration(resp)
	if isExpired {
		return *result
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

	return GetOrganizationsResponses{Organizations: organizations}
}

func GetPackages(token string, organizationId string) GetPackagesResponses {
	client, ctx := initializeGithubClient(token)

	packageType := "container"
	githubPackages, _, err := client.Organizations.ListPackages(ctx, organizationId, &github.PackageListOptions{
		PackageType: &packageType,
	})

	if err != nil {
		return GetPackagesResponses{Message: fmt.Sprintf("\n error: %v\n", err)}
	}

	var packages []Package

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {

			packages = append(packages, Package{
				Name:            *githubPackage.Name,
				Id:              *githubPackage.ID,
				Url:             *githubPackage.URL,
				PackageVersions: getPackageVersions(client, ctx, organizationId, *githubPackage.Name),
			})
		}
	}

	return GetPackagesResponses{Packages: packages}
}

func getPackageVersions(client *github.Client, ct context.Context, organizationName string, packageName string) []PackageVersion {
	githubPackages, _, _ := client.Organizations.PackageGetAllVersions(ct, organizationName, "container", url.QueryEscape(packageName), nil)
	var packages []PackageVersion

	for _, githubPackage := range githubPackages {
		if githubPackage.Name != nil && *githubPackage.Name != "" {
			packages = append(packages, PackageVersion{
				Name: githubPackage.Metadata.Container.Tags[0],
				Id:   *githubPackage.ID,
			})
		}
	}

	return packages
}

func checkTokenExpiration(resp *github.Response) (bool, *GetOrganizationsResponses) {
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
		return true, &GetOrganizationsResponses{Message: "Token expiration"}
	}

	return false, nil
}

func initializeGithubClient(token string) (*github.Client, context.Context) {
	ct := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ct, ts)

	return github.NewClient(tc), ct
}
