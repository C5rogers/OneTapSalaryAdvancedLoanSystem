package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
)

// here define cloudinary client and return it to inject it on the server as a dependency

type CloudinaryClient struct {
	cloudinary.Cloudinary
	*config.Config
}

func NewCloudinaryClient(config *config.Config) (*CloudinaryClient, error) {

	cld, err := cloudinary.NewFromParams(
		config.Cloudinary.CloudName,
		config.Cloudinary.ApiKey,
		config.Cloudinary.ApiSecret,
	)
	if err != nil {
		return nil, err
	}
	return &CloudinaryClient{
		Cloudinary: *cld,
		Config:     config,
	}, nil
}
