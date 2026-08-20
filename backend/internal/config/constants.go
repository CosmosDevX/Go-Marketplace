package config

import "time"

const (
	ProductPageSize = 16
	AccessTokenTTL  = 15 * time.Minute
	DefaultUserRole = "user"
	SellerRole      = "seller"
	AdminRole       = "admin"
	UploadsPath     = "/uploads"
)
