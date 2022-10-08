// ref. mv/models.go#GithubUser
export interface GithubUser {
  readonly login: string;
  readonly token: string;
}

// ref. mv/models.go#Organization
export interface Organization {
  readonly id: number;
  readonly login: string;
}

// ref. mv/models.go#Package
export interface Package {
  readonly id: number;
  readonly name: string;
  readonly package_type: string;
  readonly package_versions: PackageVersion[];
  readonly repository_url: string;
  readonly repository_name: string;
  readonly url: string;
}

// ref. mv/models.go#PackageVersion
export interface PackageVersion {
  readonly id: number;
  readonly name: string;
  readonly package_type: string;
  readonly url: string;
  readonly image_url: string;
  readonly organizationName: string;
  readonly createdAt: string;
}

// ref. mv/models.go#GetPackagesResponses
export interface GetPackagesResponses {
  readonly message: string;
  readonly packages: Package[];
}
