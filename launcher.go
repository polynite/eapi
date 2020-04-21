package eapi

import "fmt"

// LauncherService handles all launcher-related routes.
type LauncherService service

// CatalogResponse defines a catalog response.
type CatalogResponse struct {
	Elements []Catalog
}

// Catalog is a game-specific catalog.
type Catalog struct {
	AppName      string
	LabelName    string
	BuildVersion string
	Hash         string
	Manifests    []Manifest
}

// Manifest contains links to the actual manifest
type Manifest struct {
	URI         string `json:"uri"`
	QueryParams []struct {
		Name  string
		Value string
	}
}

// GetCatalog gets a launcher catalog
func (s *LauncherService) GetCatalog(platform string, namespace string, catalogItem string, app string, label string) (res *CatalogResponse, err error) {
	// Build url
	url := fmt.Sprintf("%s/public/assets/v2/platform/%s/namespace/%s/catalogItem/%s/app/%s/label/%s", launcherURL, platform, namespace, catalogItem, app, label)

	// Create request
	req, err := s.client.newReq("GET", url, nil)
	if err != nil {
		return
	}

	// Make request
	res = &CatalogResponse{}
	err = s.client.do(req, &res)

	return
}
