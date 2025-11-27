package profile

import (
	"encoding/json"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestADLSStorageSettings_NoOpts(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem")

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","authority-host":"https://login.microsoftonline.com","host":"dfs.core.windows.net","sas-token-validity-seconds":3600}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestADLSStorageSettings_AlternativeProtocols(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem", WithADLSAlternativeProtocols())

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())
	assert.Equal(t, core.Ptr(true), profile.AllowAlternativeProtocols)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","allow-alternative-protocols":true,"authority-host":"https://login.microsoftonline.com","host":"dfs.core.windows.net","sas-token-validity-seconds":3600}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestADLSStorageSettings_AuthorityHost(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem", WithAuthorityHost("authority"))

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())
	assert.Equal(t, core.Ptr("authority"), profile.AuthorityHost)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","authority-host":"authority","host":"dfs.core.windows.net","sas-token-validity-seconds":3600}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestADLSStorageSettings_KeyPrefix(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem", WithADLSKeyPrefix("prefix"))

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())
	assert.Equal(t, core.Ptr("prefix"), profile.KeyPrefix)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","authority-host":"https://login.microsoftonline.com","host":"dfs.core.windows.net","key-prefix":"prefix","sas-token-validity-seconds":3600}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestADLSStorageSettings_SASTokenValiditySeconds(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem", WithSASTokenValiditySeconds(256))

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())
	assert.Equal(t, core.Ptr(int64(256)), profile.SASTokenValiditySeconds)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","authority-host":"https://login.microsoftonline.com","host":"dfs.core.windows.net","sas-token-validity-seconds":256}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestADLSStorageSettings_Host(t *testing.T) {
	profile := NewADLSStorageSettings("account", "filesystem", WithHost("specific-host"))

	assert.Equal(t, StorageFamilyADLS, profile.GetStorageFamily())
	assert.Equal(t, core.Ptr("specific-host"), profile.Host)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"adls","account-name":"account","filesystem":"filesystem","authority-host":"https://login.microsoftonline.com","host":"specific-host","sas-token-validity-seconds":3600}`

	assert.JSONEq(t, jsonStr, string(b))
}
