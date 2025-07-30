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
	"fmt"
)

type (
	DeleteProfile struct {
		DeleteProfileSettings DeleteProfileSettings
	}

	DeleteProfileType string

	DeleteProfileSettings interface {
		GetDeteProfileType() DeleteProfileType
		AsProfile() *DeleteProfile

		json.Marshaler
	}

	TabularDeleteProfileHard struct{}

	TabularDeleteProfileSoft struct {
		ExpirationSeconds int32 `json:"expiration-seconds"`
	}
)

const (
	HardDeleteProfileType DeleteProfileType = "hard"
	SoftDeleteProfileType DeleteProfileType = "soft"
)

var (
	_ DeleteProfileSettings = (*TabularDeleteProfileHard)(nil)
	_ DeleteProfileSettings = (*TabularDeleteProfileSoft)(nil)
)

func NewTabularDeleteProfileHard() *TabularDeleteProfileHard {
	return &TabularDeleteProfileHard{}
}

func (*TabularDeleteProfileHard) GetDeteProfileType() DeleteProfileType {
	return HardDeleteProfileType
}

func (d *TabularDeleteProfileHard) AsProfile() *DeleteProfile {
	return &DeleteProfile{DeleteProfileSettings: d}
}

func (d TabularDeleteProfileHard) MarshalJSON() ([]byte, error) {
	aux := struct {
		Type string `json:"type"`
	}{
		Type: string(d.GetDeteProfileType()),
	}
	return json.Marshal(aux)
}

func NewTabularDeleteProfileSoft(expirationSeconds int32) *TabularDeleteProfileSoft {
	return &TabularDeleteProfileSoft{
		ExpirationSeconds: expirationSeconds,
	}
}

func (*TabularDeleteProfileSoft) GetDeteProfileType() DeleteProfileType {
	return SoftDeleteProfileType
}

func (d *TabularDeleteProfileSoft) AsProfile() *DeleteProfile {
	return &DeleteProfile{DeleteProfileSettings: d}
}

func (d TabularDeleteProfileSoft) MarshalJSON() ([]byte, error) {
	type Alias TabularDeleteProfileSoft
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(d.GetDeteProfileType()),
		Alias: Alias(d),
	}
	return json.Marshal(aux)
}

func (sc *DeleteProfile) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch peek.Type {
	case "hard":
		var cfg TabularDeleteProfileHard
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.DeleteProfileSettings = &cfg
	case "soft":
		var cfg TabularDeleteProfileSoft
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.DeleteProfileSettings = &cfg
	default:
		return fmt.Errorf("unsupported delete profile type: %s", peek.Type)
	}
	return nil
}

func (sc DeleteProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.DeleteProfileSettings)
}
