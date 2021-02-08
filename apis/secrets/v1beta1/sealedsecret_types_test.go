/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"encoding/json"
	"testing"
)

func Test_v1beta1_unmarshal_AwsKmsProvider(t *testing.T) {
	var p AwsKmsProvider
	err := json.Unmarshal([]byte(`{"encryptionAlgorithm": "AES_256"}`), &p)
	if err != nil {
		t.Error(err)
	}
	if p.EncryptionAlgorithm != "AES_256" {
		t.Logf("Invalid encryptionAlgorithm: epected 'AES_256' got %s", p.EncryptionAlgorithm)
		t.Fail()
	}
}

func Test_v1beta1_unmarshal_Provider(t *testing.T) {
	var p Provider
	err := json.Unmarshal([]byte(`{"keyEncryptionKeyId": "x", "awsKms": {"encryptionAlgorithm": "AES_256"}}`), &p)
	if err != nil {
		t.Error(err)
	}
	if p.KeyEncryptionKeyID != "x" {
		t.Logf("Invalid keyEncryptionKeyId: epected 'x' got %s", p.KeyEncryptionKeyID)
		t.Fail()
	}
	if p.AwsKms.EncryptionAlgorithm != "AES_256" {
		t.Logf("Invalid encryptionAlgorithm: epected 'AES_256' got %s", p.AwsKms.EncryptionAlgorithm)
		t.Fail()
	}
}

func Test_v1beta1_unmarshal_Envelope(t *testing.T) {
	var p Envelope
	err := json.Unmarshal([]byte(`{"keyEncryptionKeyId": "x", "awsKms": {"encryptionAlgorithm": "AES_256"}, "dataEncryptionKey": "c2VjCg==", "cipherText": "c2VjCg=="}`), &p)
	if err != nil {
		t.Error(err)
	}
	if p.KeyEncryptionKeyID != "x" {
		t.Logf("Invalid keyEncryptionKeyId: epected 'x' got %s", p.KeyEncryptionKeyID)
		t.Fail()
	}
	if p.Provider.AwsKms.EncryptionAlgorithm != "AES_256" {
		t.Logf("Invalid encryptionAlgorithm: epected 'AES_256' got %s", p.Provider.AwsKms.EncryptionAlgorithm)
		t.Fail()
	}
}

func Test_v1beta1_unmarshal_Envelopes(t *testing.T) {
	var p Envelopes
	err := json.Unmarshal([]byte(`[{"keyEncryptionKeyId": "x", "awsKms": {"encryptionAlgorithm": "AES_256"}, "dataEncryptionKey": "c2VjCg==", "cipherText": "c2VjCg=="}]`), &p)
	if err != nil {
		t.Error(err)
	}
	if p[0].KeyEncryptionKeyID != "x" {
		t.Logf("Invalid keyEncryptionKeyId: epected 'x' got %s", p[0].KeyEncryptionKeyID)
		t.Fail()
	}
	if p[0].Provider.AwsKms.EncryptionAlgorithm != "AES_256" {
		t.Logf("Invalid encryptionAlgorithm: epected 'AES_256' got %s", p[0].Provider.AwsKms.EncryptionAlgorithm)
		t.Fail()
	}
}

func Test_v1beta1_unmarshal_EnvelopesMap(t *testing.T) {
	var p map[string]Envelopes
	err := json.Unmarshal([]byte(`{"a": [{"keyEncryptionKeyId": "x", "awsKms": {"encryptionAlgorithm": "AES_256"}, "dataEncryptionKey": "c2VjCg==", "cipherText": "c2VjCg=="}]}`), &p)
	if err != nil {
		t.Error(err)
	}
	if p["a"][0].KeyEncryptionKeyID != "x" {
		t.Logf("Invalid keyEncryptionKeyId: epected 'x' got %s", p["a"][0].KeyEncryptionKeyID)
		t.Fail()
	}
	if p["a"][0].Provider.AwsKms.EncryptionAlgorithm != "AES_256" {
		t.Logf("Invalid encryptionAlgorithm: epected 'AES_256' got %s", p["a"][0].Provider.AwsKms.EncryptionAlgorithm)
		t.Fail()
	}
}
