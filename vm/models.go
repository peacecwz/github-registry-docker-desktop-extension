package main

import "time"

type GithubUser struct {
	Login string `json:"login"`
	Token string `json:"token"`
}

type Organization []struct {
	Id   int64  `json:"id"`
	Name string `json:"login"`
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

type Timestamp struct {
	time.Time
}

type PackageVersion struct {
	Id               int64     `json:"id"`
	Name             string    `json:"name"`
	PackageType      string    `json:"package_type"`
	Url              string    `json:"url"`
	ImageUrl         string    `json:"image_url"`
	OrganizationName string    `json:"organizationName"`
	CreatedAt        Timestamp `json:"createdAt"`
}

type GetPackagesResponses struct {
	Message  string    `json:"message"`
	Packages []Package `json:"packages"`
}
