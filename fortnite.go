package eapi

import (
	"time"
)

// FortniteService handle all Fortnite-specific routes.
type FortniteService service

// CloudstorageSystemFile defines a cloud storage file.
type CloudstorageSystemFile struct {
	UniqueFilename string
	Filename       string
	Hash           string
	Hash256        string
	Length         int
	ContentType    string
	Uploaded       time.Time
	StorageType    string
	DoNotCache     bool
}

// GetCloudstorageSystem gets all system files.
func (s *FortniteService) GetCloudstorageSystem() (res []CloudstorageSystemFile, err error) {
	// Create request
	req, err := s.client.newReq("GET", fortniteURL+"/cloudstorage/system", nil)
	if err != nil {
		return
	}

	// Make request
	res = make([]CloudstorageSystemFile, 0)
	err = s.client.do(req, &res)

	return
}
