// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package profile

import (
	"encoding/json"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

// S3StorageSettings represents the storage settings for a warehouse
// where data are stored on AWS S3 or S3-compatible storage.
type S3StorageSettings struct {
	// Name of the S3 bucket
	Bucket string `json:"bucket"`
	// Region to use for S3 requests.
	Region string `json:"region"`
	// Wether to enable STS.
	// If true, Lakekeeper will try to assume role
	// from AssumeRoleARN or STSRoleARN. You must provide
	// one of them. Default: false
	STSEnabled bool `json:"sts-enabled"`
	// Allow s3a:// and s3n:// in locations.
	// This is disabled by default.
	// We do not recommend to use this setting except for migration of old hadoop-based tables
	// via the register endpoint.
	// Tables with s3a paths are not accessible outside the Java ecosystem
	AllowAlternativeProtocols *bool `json:"allow-alternative-protocols,omitempty"`
	// Optional ARN to assume when accessing the bucket from Lakekeeper.
	AssumeRoleARN *string `json:"assume-role-arn,omitempty"`
	// ARN of the KMS key used to encrypt the S3 bucket, if any.
	AWSKMSKeyARN *string `json:"aws-kms-key-arn,omitempty"`
	// Optional endpoint to use for S3 requests, if not provided the region will be
	// used to determine the endpoint. If both region and endpoint are provided,
	// the endpoint will be used. Example: http://s3-de.my-domain.com:9000
	Endpoint *string `json:"endpoint,omitempty"`
	// S3 flavor to use. Defaults to AWS
	Flavor *S3Flavor `json:"flavor,omitempty"`
	// Subpath in the bucket to use.
	KeyPrefix *string `json:"key-prefix,omitempty"`
	// Path style access for S3 requests.
	// If the underlying S3 supports both, we recommend to not set path_style_access.
	PathStyleAccess *bool `json:"path-style-access,omitempty"`
	// Controls whether the s3.delete-enabled=false flag is sent to clients.
	//
	// In all Iceberg 1.x versions, when Spark executes DROP TABLE xxx PURGE,
	// it directly deletes files from S3, bypassing the catalog's soft-deletion mechanism.
	// Other query engines properly delegate this operation to the catalog.
	// This Spark behavior is expected to change in Iceberg 2.0.
	//
	// Setting this to true pushes the s3.delete-enabled=false flag to clients,
	// which discourages Spark from directly deleting files during DROP TABLE xxx PURGE operations.
	// Note that clients may override this setting, and it affects other Spark operations that
	// require file deletion, such as removing snapshots.
	//
	// For more details, refer to Lakekeeper's Soft-Deletion documentation.
	// This flag has no effect if Soft-Deletion is disabled for the warehouse.
	PushS3DeleteDisabled *bool `json:"push-s3-delete-disabled,omitempty"`
	// S3 URL style detection mode for remote signing. One of auto, path-style, virtual-host.
	// Default: auto.
	// When set to auto, Lakekeeper will first try to parse the URL as virtual-host
	// and then attempt path-style. path assumes the bucket name is the first path
	// segment in the URL.
	// virtual-host assumes the bucket name is the first subdomain if it is preceding .s3 or .s3-.
	RemoteSigningURLStyle *RemoteSigningURLStyle `json:"remote-signing-url-style,omitempty"`
	// Optional role ARN to assume for sts vended-credentials.
	// If not provided, assume_role_arn is used.
	// Either assume_role_arn or sts_role_arn must be provided if sts_enabled is true.
	STSRoleARN *string `json:"sts-role-arn,omitempty"`
	// The validity of the sts tokens in seconds. Default is 3600
	STSTokenValiditySeconds *int64 `json:"sts-token-validity-seconds,omitempty"`
}

type (
	S3Flavor              string
	RemoteSigningURLStyle string
)

const (
	AWSFlavor      S3Flavor = "aws"
	S3CompatFlavor S3Flavor = "s3-compat"

	AutoSigningURLStyle        RemoteSigningURLStyle = "auto"
	PathSigningURLStyle        RemoteSigningURLStyle = "path-style"
	VirtualHostSigningURLStyle RemoteSigningURLStyle = "virtual-host"
)

func (sp *S3StorageSettings) GetStorageFamily() StorageFamily {
	return StorageFamilyS3
}

type S3StorageSettingsOptions func(*S3StorageSettings)

// NewS3StorageSettings creates a new S3 storage profile considering
// the options given.
func NewS3StorageSettings(bucket, region string, opts ...S3StorageSettingsOptions) *S3StorageSettings {
	// Default configuration
	profile := S3StorageSettings{
		Bucket:                  bucket,
		Region:                  region,
		STSEnabled:              false,
		Flavor:                  core.Ptr(AWSFlavor),
		STSTokenValiditySeconds: core.Ptr(int64(3600)),
		PushS3DeleteDisabled:    core.Ptr(true),
		RemoteSigningURLStyle:   core.Ptr(AutoSigningURLStyle),
	}

	// Apply options
	for _, v := range opts {
		v(&profile)
	}

	return &profile
}

func WithSTSEnabled() S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.STSEnabled = true
	}
}

func WithSTSRoleARN(stsRoleARN string) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.STSRoleARN = &stsRoleARN
	}
}

func WithS3KeyPrefix(prefix string) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.KeyPrefix = &prefix
	}
}

func WithEndpoint(endpoint string) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.Endpoint = &endpoint
	}
}

func WithS3AlternativeProtocols() S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.AllowAlternativeProtocols = core.Ptr(true)
	}
}

func WithAssumeRoleARN(arn string) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.AssumeRoleARN = &arn
	}
}

func WithAWSKMSKeyARN(arn string) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.AWSKMSKeyARN = &arn
	}
}

func WithFlavor(flavor S3Flavor) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.Flavor = &flavor
	}
}

func WithPathStyleAccess() S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.PathStyleAccess = core.Ptr(true)
	}
}

func WithPushS3DeleteDisabled(active bool) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.PushS3DeleteDisabled = &active
	}
}

func WithRemoteSigningURLStyle(style RemoteSigningURLStyle) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.RemoteSigningURLStyle = &style
	}
}

func WithSTSTokenValiditySeconds(seconds int64) S3StorageSettingsOptions {
	return func(sp *S3StorageSettings) {
		sp.STSTokenValiditySeconds = &seconds
	}
}

func (s *S3StorageSettings) AsProfile() StorageProfile {
	return StorageProfile{StorageSettings: s}
}

func (s S3StorageSettings) MarshalJSON() ([]byte, error) {
	type Alias S3StorageSettings
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyS3),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
