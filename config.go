package main

import "cloud.google.com/go/asset/apiv1/assetpb"

const (
	serverName            = "Cloud Asset Inventory Server"
	serverVersion         = "1.1.0"
	defaultListPageSize   = 100
	defaultSearchPageSize = 100
	maxListPageSize       = 1000
	maxSearchPageSize     = 500
	maxHistoryAssetNames  = 100
	maxEffectiveIAMNames  = 20
	defaultAssetContent   = assetpb.ContentType_RESOURCE
	defaultHistoryContent = assetpb.ContentType_RESOURCE
	rfc3339Example        = "2026-03-19T10:30:00Z"
)

var contentTypeEnumValues = []string{
	"CONTENT_TYPE_UNSPECIFIED",
	"RESOURCE",
	"IAM_POLICY",
	"ORG_POLICY",
	"ACCESS_POLICY",
	"OS_INVENTORY",
	"RELATIONSHIP",
}
