package profile

import (
	"encoding/json"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	// ADLSStorageSettings represents the storage settings for a warehouse
	// where data are stored on Azure Data Lake Storage.
	ADLSStorageSettings struct {
		// Name of the azure storage account.
		AccountName string `json:"account-name"`
		// Name of the adls filesystem, in blobstorage also known as container.
		Filesystem string `json:"filesystem"`
		// Allow alternative protocols such as wasbs:// in locations.
		// This is disabled by default. We do not recommend to use this setting
		// except for migration of old tables via the register endpoint.
		AllowAlternativeProtocols *bool `json:"allow-alternative-protocols,omitempty"`
		// The authority host to use for authentication.
		// Default: https://login.microsoftonline.com.
		AuthorityHost *string `json:"authority-host,omitempty"`
		// The host to use for the storage account. Default: dfs.core.windows.net.
		Host *string `json:"host,omitempty"`
		// Subpath in the filesystem to use.
		KeyPrefix *string `json:"key-prefix,omitempty"`
		// The validity of the sas token in seconds. Default: 3600.
		SASTokenValiditySeconds *int64 `json:"sas-token-validity-seconds,omitempty"`
	}

	ADLSStorageSettingsOptions func(*ADLSStorageSettings)
)

func (sp *ADLSStorageSettings) GetStorageFamily() StorageFamily {
	return StorageFamilyADLS
}

// NewADLSStorageSettings creates a new ADLS storage profile considering
// the options given.
func NewADLSStorageSettings(accountName, fs string, opts ...ADLSStorageSettingsOptions) *ADLSStorageSettings {
	// Default configuration
	profile := ADLSStorageSettings{
		AccountName:             accountName,
		Filesystem:              fs,
		AuthorityHost:           core.Ptr("https://login.microsoftonline.com"),
		Host:                    core.Ptr("dfs.core.windows.net"),
		SASTokenValiditySeconds: core.Ptr(int64(3600)),
	}

	// Apply options
	for _, v := range opts {
		v(&profile)
	}

	return &profile
}

func WithADLSAlternativeProtocols() ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) {
		sp.AllowAlternativeProtocols = core.Ptr(true)
	}
}

func WithAuthorityHost(host string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) {
		sp.AuthorityHost = &host
	}
}

func WithADLSKeyPrefix(prefix string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) {
		sp.KeyPrefix = &prefix
	}
}

func WithSASTokenValiditySeconds(seconds int64) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) {
		sp.SASTokenValiditySeconds = &seconds
	}
}

func WithHost(host string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) {
		sp.Host = &host
	}
}

func (s *ADLSStorageSettings) AsProfile() StorageProfile {
	return StorageProfile{s}
}

func (s ADLSStorageSettings) MarshalJSON() ([]byte, error) {
	type Alias ADLSStorageSettings
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyADLS),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
