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
	assert.Equal(t, true, profile.STSEnabled)
	assert.Equal(t, "role-arn", *profile.STSRoleARN)
	assert.Equal(t, "prefix", *profile.KeyPrefix)
	assert.Equal(t, "http://minio:9000/", *profile.Endpoint)
	assert.Equal(t, true, *profile.AllowAlternativeProtocols)
	assert.Equal(t, "assume", *profile.AssumeRoleARN)
	assert.Equal(t, "kms", *profile.AWSKMSKeyARN)
	assert.Equal(t, S3CompatFlavor, *profile.Flavor)
	assert.Equal(t, true, *profile.PathStyleAccess)
	assert.Equal(t, true, *profile.PushS3DeleteDisabled)
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
