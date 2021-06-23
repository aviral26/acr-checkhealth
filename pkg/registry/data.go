package registry

import v1 "github.com/opencontainers/image-spec/specs-go/v1"

// Registry REST routes
const (
	// Ping routes
	routeFrontendPing     = "/v2/"
	routeDataEndpointPing = "/"

	// Blob routes
	routeInitiateBlobUpload = "/v2/%s/blobs/uploads/" // add repo name
	routeBlobPull           = "/v2/%s/blobs/%s"       // add repo name and digest

	// Manifest routes
	routeManifest = "/v2/%s/manifests/%s" // add repo name and digest/tag
)

// Artifacts related constants
const (
	checkHealthAuthor       = "ACR Check Health"
	checkHealthMediaType    = "application/acr.checkhealth.test"
	checkHealthArtifactType = "application/acr.checkhealth.artifact.test"
	checkHealthLayerFmt     = "Test layer authored by " + checkHealthAuthor + " at %s" // add time
	checkHealthRepoPrefix   = "acrcheckhealth"
)

// OCI related data
var (
	ociConfig = v1.Image{
		Author: checkHealthAuthor,
	}
)
