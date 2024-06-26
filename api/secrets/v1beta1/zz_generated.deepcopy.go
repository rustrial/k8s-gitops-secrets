//go:build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Authorization) DeepCopyInto(out *Authorization) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Authorization.
func (in *Authorization) DeepCopy() *Authorization {
	if in == nil {
		return nil
	}
	out := new(Authorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AwsKmsProvider) DeepCopyInto(out *AwsKmsProvider) {
	*out = *in
	if in.Nonce != nil {
		in, out := &in.Nonce, &out.Nonce
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AwsKmsProvider.
func (in *AwsKmsProvider) DeepCopy() *AwsKmsProvider {
	if in == nil {
		return nil
	}
	out := new(AwsKmsProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Envelope) DeepCopyInto(out *Envelope) {
	*out = *in
	in.Provider.DeepCopyInto(&out.Provider)
	if in.DataEncryptionKey != nil {
		in, out := &in.DataEncryptionKey, &out.DataEncryptionKey
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.CipherText != nil {
		in, out := &in.CipherText, &out.CipherText
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Envelope.
func (in *Envelope) DeepCopy() *Envelope {
	if in == nil {
		return nil
	}
	out := new(Envelope)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Envelopes) DeepCopyInto(out *Envelopes) {
	{
		in := &in
		*out = make(Envelopes, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Envelopes.
func (in Envelopes) DeepCopy() Envelopes {
	if in == nil {
		return nil
	}
	out := new(Envelopes)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyEncryptionKeyPolicy) DeepCopyInto(out *KeyEncryptionKeyPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyEncryptionKeyPolicy.
func (in *KeyEncryptionKeyPolicy) DeepCopy() *KeyEncryptionKeyPolicy {
	if in == nil {
		return nil
	}
	out := new(KeyEncryptionKeyPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeyEncryptionKeyPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyEncryptionKeyPolicyList) DeepCopyInto(out *KeyEncryptionKeyPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KeyEncryptionKeyPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyEncryptionKeyPolicyList.
func (in *KeyEncryptionKeyPolicyList) DeepCopy() *KeyEncryptionKeyPolicyList {
	if in == nil {
		return nil
	}
	out := new(KeyEncryptionKeyPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeyEncryptionKeyPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyEncryptionKeyPolicySpec) DeepCopyInto(out *KeyEncryptionKeyPolicySpec) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyEncryptionKeyPolicySpec.
func (in *KeyEncryptionKeyPolicySpec) DeepCopy() *KeyEncryptionKeyPolicySpec {
	if in == nil {
		return nil
	}
	out := new(KeyEncryptionKeyPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyEncryptionKeyPolicyStatus) DeepCopyInto(out *KeyEncryptionKeyPolicyStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyEncryptionKeyPolicyStatus.
func (in *KeyEncryptionKeyPolicyStatus) DeepCopy() *KeyEncryptionKeyPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(KeyEncryptionKeyPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provider) DeepCopyInto(out *Provider) {
	*out = *in
	if in.AwsKms != nil {
		in, out := &in.AwsKms, &out.AwsKms
		*out = new(AwsKmsProvider)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provider.
func (in *Provider) DeepCopy() *Provider {
	if in == nil {
		return nil
	}
	out := new(Provider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SealedSecret) DeepCopyInto(out *SealedSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SealedSecret.
func (in *SealedSecret) DeepCopy() *SealedSecret {
	if in == nil {
		return nil
	}
	out := new(SealedSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SealedSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SealedSecretList) DeepCopyInto(out *SealedSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SealedSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SealedSecretList.
func (in *SealedSecretList) DeepCopy() *SealedSecretList {
	if in == nil {
		return nil
	}
	out := new(SealedSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SealedSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SealedSecretSpec) DeepCopyInto(out *SealedSecretSpec) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	if in.Immutable != nil {
		in, out := &in.Immutable, &out.Immutable
		*out = new(bool)
		**out = **in
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string][]byte, len(*in))
		for key, val := range *in {
			var outVal []byte
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = make([]byte, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.EncryptedData != nil {
		in, out := &in.EncryptedData, &out.EncryptedData
		*out = make(map[string]Envelopes, len(*in))
		for key, val := range *in {
			var outVal []Envelope
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = make(Envelopes, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	if in.StringData != nil {
		in, out := &in.StringData, &out.StringData
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SealedSecretSpec.
func (in *SealedSecretSpec) DeepCopy() *SealedSecretSpec {
	if in == nil {
		return nil
	}
	out := new(SealedSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SealedSecretStatus) DeepCopyInto(out *SealedSecretStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Authorizations != nil {
		in, out := &in.Authorizations, &out.Authorizations
		*out = make(map[string]Authorization, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SealedSecretStatus.
func (in *SealedSecretStatus) DeepCopy() *SealedSecretStatus {
	if in == nil {
		return nil
	}
	out := new(SealedSecretStatus)
	in.DeepCopyInto(out)
	return out
}
