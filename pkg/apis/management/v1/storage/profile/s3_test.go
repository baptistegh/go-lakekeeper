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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3StorageProfile(t *testing.T) {
	profile := NewS3StorageSettings(
		"bucket1",
		"eu-west-1",
		WithSTSEnabled(),
		WithSTSRoleARN("role-arn"),
		WithS3KeyPrefix("prefix"),
		WithEndpoint("http://minio:9000/"),
		WithS3AlternativeProtocols(),
		WithAssumeRoleARN("assume"),
		WithAWSKMSKeyARN("kms"),
		WithFlavor(S3CompatFlavor),
		WithPathStyleAccess(),
		WithPushS3DeleteDisabled(true),
		WithRemoteSigningURLStyle(VirtualHostSigningURLStyle),
		WithSTSTokenValiditySeconds(7200),
	)

	assert.Equal(t, StorageFamilyS3, profile.GetStorageFamily())
	assert.Equal(t, "bucket1", profile.Bucket)
	assert.Equal(t, "eu-west-1", profile.Region)
	assert.True(t, profile.STSEnabled)
	assert.Equal(t, "role-arn", *profile.STSRoleARN)
	assert.Equal(t, "prefix", *profile.KeyPrefix)
	assert.Equal(t, "http://minio:9000/", *profile.Endpoint)
	assert.True(t, *profile.AllowAlternativeProtocols)
	assert.Equal(t, "assume", *profile.AssumeRoleARN)
	assert.Equal(t, "kms", *profile.AWSKMSKeyARN)
	assert.Equal(t, S3CompatFlavor, *profile.Flavor)
	assert.True(t, *profile.PathStyleAccess)
	assert.True(t, *profile.PushS3DeleteDisabled)
	assert.Equal(t, VirtualHostSigningURLStyle, *profile.RemoteSigningURLStyle)
	assert.Equal(t, int64(7200), *profile.STSTokenValiditySeconds)
}

func TestS3StorageProfile_Marshal(t *testing.T) {
	sp := NewS3StorageSettings(
		"bucket1",
		"eu-west-1",
		WithEndpoint("http://minio:9000/"),
		WithPathStyleAccess(),
		WithS3KeyPrefix("warehouse"),
	)

	expected := `{"type":"s3","bucket":"bucket1","region":"eu-west-1","sts-enabled":false,"endpoint":"http://minio:9000/","flavor":"aws","key-prefix":"warehouse","path-style-access":true,"push-s3-delete-disabled":true,"remote-signing-url-style":"auto","sts-token-validity-seconds":3600}`

	b, err := json.Marshal(sp)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}

	// by Config
	b, err = json.Marshal(sp.AsProfile())
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}
}
