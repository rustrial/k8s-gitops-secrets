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
	apiCoreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AwsKmsProvider is the AWS KMS provider.
type AwsKmsProvider struct {
	// Version of the provider spec
	// +required
	Version uint32 `json:"version,omitempty"`

	// EncryptionAlgorithm specifies the encryption algorithm that will be used to decrypt.
	// https://docs.aws.amazon.com/kms/latest/APIReference/API_Decrypt.html
	//
	// +optional
	EncryptionAlgorithm string `json:"encryptionAlgorithm,omitempty"`

	// Nonce of stream cipher
	// +optional
	Nonce []byte `json:"nonce,omitempty"`
}

// Provider only one can be defined.
type Provider struct {
	// KeyEncryptionKeyId is the provider specific unique ID of
	// the Key Encryption Key (KEK) use to encrypt/decrypt the
	// Data Encryption Key (DEK).
	//
	// This ID must uniquely identify the KEK and provider and
	// is used in authorization rules to decide which namespaces
	// can access which KEKs.
	//
	// +required
	KeyEncryptionKeyID string `json:"keyEncryptionKeyId,omitempty"`

	// AwsKms provider.
	//
	// +optional
	AwsKms *AwsKmsProvider `json:"awsKms,omitempty"`
}

// Envelope contains encrypted payload as well as key material and metadata needed to decrypt that payload.
//
// As of now, we banned any public key cryptograhy (PKI) support, as all available PKI schemes are
// not quantum-safe.
type Envelope struct {
	// Provider holds the encryption proider specific data.
	Provider `json:",inline"`

	// DataEncryptionKey holds the encrypted symmetric data-key, needed to decrypt the CipherText.
	// The data-key is either encrypted using envelope encryption provided by the Provider.
	// +optional
	DataEncryptionKey []byte `json:"dataEncryptionKey,omitempty" yaml:"dataEncryptionKey,omitempty"`

	// CipherText holds the encrypted payload, encrypted by the symmetric DataKey.
	// +required
	CipherText []byte `json:"cipherText,omitempty" yaml:"cipherText"`
}

// Envelopes array of Envelopes, which allows support for encrypting
// a single entry with multiple providers. During decryption, the controller will stop on the first
// successful provider and use the decrypted value it returns.
type Envelopes []Envelope

// SealedSecretSpec defines the desired state of SealedSecret
type SealedSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Metadata to use to create the Secret. For example to create it in with
	// a different name or to add labes and annotations.
	//
	// Note, any namespace defined here will be ignored, an the corresponding
	// secret will always be created in the namespace of the SealedSecret.
	//
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Metadata metav1.ObjectMeta `json:"metadata,omitempty"`

	// Immutable, if set to true, ensures that data stored in the Secret cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// This is a beta field enabled by ImmutableEphemeralVolumes feature gate.
	// +optional
	Immutable *bool `json:"immutable,omitempty"`

	// Data contains the secret data. Each key must consist of alphanumeric
	// characters, '-', '_' or '.'. The serialized form of the secret data is a
	// base64 encoded string, representing the arbitrary (possibly non-string)
	// data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
	// +optional
	Data map[string][]byte `json:"data,omitempty"`

	// EncryptedData contains the secret envelope encrypted secret data. Each key must consist of alphanumeric
	// characters, '-', '_' or '.'. The values are arrays of Envelopes, which allows support for encrypting
	// a single entry with multiple providers. During decryption, the controller will stop on the first
	// successful provider and use the decrypted value it returns.
	// +optional
	EncryptedData map[string]Envelopes `json:"encryptedData,omitempty"`

	// StringData allows specifying non-binary secret data in string form.
	// It is provided as a write-only convenience method.
	// All keys and values are merged into the data field on write, overwriting any existing values.
	// It is never output when reading from the API.
	// +k8s:conversion-gen=false
	// +optional
	StringData map[string]string `json:"stringData,omitempty"`

	// Used to facilitate programmatic handling of secret data.
	// +optional
	Type apiCoreV1.SecretType `json:"type,omitempty"`
}

const (
	// Failure indicates reconcile failed.
	Failure = "Failure"
	// Success indicates reconcile succeeded.
	Success = "Success"
)

// Authorization source object
type Authorization struct {
	// The kind of the authorization source.
	// +required
	Kind string `json:"kind,omitempty"`
	// The name of the authorization source.
	// +required
	Name string `json:"name,omitempty"`
	// The namespace of the authorization source.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SealedSecretStatus defines the observed state of SealedSecret
type SealedSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Authorization sources of the Key Encryption Keys references in this object.
	Authorizations map[string]Authorization `json:"authorizations,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Authorized",type="string",JSONPath=".status.conditions[?(@.type==\"Authorized\")].status",description="Whether all Key Encryption Keys references are authorized for this namespace"
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description="Ready"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description="Status message"

// SealedSecret is the Schema for the sealedsecrets API
type SealedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SealedSecretSpec   `json:"spec,omitempty"`
	Status SealedSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SealedSecretList contains a list of SealedSecret
type SealedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SealedSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SealedSecret{}, &SealedSecretList{})
}
