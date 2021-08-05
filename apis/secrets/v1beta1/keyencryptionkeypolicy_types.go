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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeyEncryptionKeyPolicySpec defines the desired state of KeyEncryptionKeyPolicy
type KeyEncryptionKeyPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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

	// White-list of namespaces, which are entitled to use this KEK
	// to decrypt DataEncrpytionKeys.
	//
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`
}

// KeyEncryptionKeyPolicyStatus defines the observed state of KeyEncryptionKeyPolicy
type KeyEncryptionKeyPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:categories=all

// KeyEncryptionKeyPolicy is the Schema for the keyencryptionkeypolicies API
type KeyEncryptionKeyPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyEncryptionKeyPolicySpec   `json:"spec,omitempty"`
	Status KeyEncryptionKeyPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeyEncryptionKeyPolicyList contains a list of KeyEncryptionKeyPolicies
type KeyEncryptionKeyPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeyEncryptionKeyPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeyEncryptionKeyPolicy{}, &KeyEncryptionKeyPolicyList{})
}
