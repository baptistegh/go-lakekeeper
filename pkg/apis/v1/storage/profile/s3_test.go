package profile

import (
	"encoding/json"
	"testing"
)

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
