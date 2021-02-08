<h1>GitOps Secrets API reference</h1>
<p>Packages:</p>
<ul class="simple">
<li>
<a href="#secrets.rustrial.org%2fv1beta1">secrets.rustrial.org/v1beta1</a>
</li>
</ul>
<h2 id="secrets.rustrial.org/v1beta1">secrets.rustrial.org/v1beta1</h2>
<p>Package v1beta1 contains API Schema definitions for the SealedSecrets v1beta1 API group</p>
Resource Types:
<ul class="simple"></ul>
<h3 id="secrets.rustrial.org/v1beta1.Authorization">Authorization
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecretStatus">SealedSecretStatus</a>)
</p>
<p>Authorization source object</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>kind</code><br>
<em>
string
</em>
</td>
<td>
<p>The kind of the authorization source.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
<p>The name of the authorization source.</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The namespace of the authorization source.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.AwsKmsProvider">AwsKmsProvider
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.Provider">Provider</a>)
</p>
<p>AwsKmsProvider is the AWS KMS provider.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>version</code><br>
<em>
uint32
</em>
</td>
<td>
<p>Version of the provider spec</p>
</td>
</tr>
<tr>
<td>
<code>encryptionAlgorithm</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionAlgorithm specifies the encryption algorithm that will be used to decrypt.
<a href="https://docs.aws.amazon.com/kms/latest/APIReference/API_Decrypt.html">https://docs.aws.amazon.com/kms/latest/APIReference/API_Decrypt.html</a></p>
</td>
</tr>
<tr>
<td>
<code>nonce</code><br>
<em>
[]byte
</em>
</td>
<td>
<em>(Optional)</em>
<p>Nonce of stream cipher</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.Envelope">Envelope
</h3>
<p>Envelope contains encrypted payload as well as key material and metadata needed to decrypt that payload.</p>
<p>As of now, we banned any public key cryptograhy (PKI) support, as all available PKI schemes are
not quantum-safe.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Provider</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.Provider">
Provider
</a>
</em>
</td>
<td>
<p>
(Members of <code>Provider</code> are embedded into this type.)
</p>
<p>Provider holds the encryption proider specific data.</p>
</td>
</tr>
<tr>
<td>
<code>dataEncryptionKey</code><br>
<em>
[]byte
</em>
</td>
<td>
<em>(Optional)</em>
<p>DataEncryptionKey holds the encrypted symmetric data-key, needed to decrypt the CipherText.
The data-key is either encrypted using envelope encryption provided by the Provider.</p>
</td>
</tr>
<tr>
<td>
<code>cipherText</code><br>
<em>
[]byte
</em>
</td>
<td>
<p>CipherText holds the encrypted payload, encrypted by the symmetric DataKey.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.Envelopes">Envelopes
(<code>[]./apis/secrets/v1beta1.Envelope</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecretSpec">SealedSecretSpec</a>)
</p>
<p>Envelopes array of Envelopes, which allows support for encrypting
a single entry with multiple providers. During decryption, the controller will stop on the first
successful provider and use the decrypted value it returns.</p>
<h3 id="secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicy">KeyEncryptionKeyPolicy
</h3>
<p>KeyEncryptionKeyPolicy is the Schema for the keyencryptionkeypolicies API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicySpec">
KeyEncryptionKeyPolicySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>keyEncryptionKeyId</code><br>
<em>
string
</em>
</td>
<td>
<p>KeyEncryptionKeyId is the provider specific unique ID of
the Key Encryption Key (KEK) use to encrypt/decrypt the
Data Encryption Key (DEK).</p>
<p>This ID must uniquely identify the KEK and provider and
is used in authorization rules to decide which namespaces
can access which KEKs.</p>
</td>
</tr>
<tr>
<td>
<code>namespaces</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>White-list of namespaces, which are entitled to use this KEK
to decrypt DataEncrpytionKeys.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicyStatus">
KeyEncryptionKeyPolicyStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicySpec">KeyEncryptionKeyPolicySpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicy">KeyEncryptionKeyPolicy</a>)
</p>
<p>KeyEncryptionKeyPolicySpec defines the desired state of KeyEncryptionKeyPolicy</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>keyEncryptionKeyId</code><br>
<em>
string
</em>
</td>
<td>
<p>KeyEncryptionKeyId is the provider specific unique ID of
the Key Encryption Key (KEK) use to encrypt/decrypt the
Data Encryption Key (DEK).</p>
<p>This ID must uniquely identify the KEK and provider and
is used in authorization rules to decide which namespaces
can access which KEKs.</p>
</td>
</tr>
<tr>
<td>
<code>namespaces</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>White-list of namespaces, which are entitled to use this KEK
to decrypt DataEncrpytionKeys.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicyStatus">KeyEncryptionKeyPolicyStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.KeyEncryptionKeyPolicy">KeyEncryptionKeyPolicy</a>)
</p>
<p>KeyEncryptionKeyPolicyStatus defines the observed state of KeyEncryptionKeyPolicy</p>
<h3 id="secrets.rustrial.org/v1beta1.Provider">Provider
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.Envelope">Envelope</a>)
</p>
<p>Provider only one can be defined.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>keyEncryptionKeyId</code><br>
<em>
string
</em>
</td>
<td>
<p>KeyEncryptionKeyId is the provider specific unique ID of
the Key Encryption Key (KEK) use to encrypt/decrypt the
Data Encryption Key (DEK).</p>
<p>This ID must uniquely identify the KEK and provider and
is used in authorization rules to decide which namespaces
can access which KEKs.</p>
</td>
</tr>
<tr>
<td>
<code>awsKms</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.AwsKmsProvider">
AwsKmsProvider
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AwsKms provider.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.SealedSecret">SealedSecret
</h3>
<p>SealedSecret is the Schema for the sealedsecrets API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecretSpec">
SealedSecretSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Metadata to use to create the Secret. For example to create it in with
a different name or to add labes and annotations.</p>
<p>Note, any namespace defined here will be ignored, an the corresponding
secret will always be created in the namespace of the SealedSecret.</p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>immutable</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Immutable, if set to true, ensures that data stored in the Secret cannot
be updated (only object metadata can be modified).
If not set to true, the field can be modified at any time.
Defaulted to nil.
This is a beta field enabled by ImmutableEphemeralVolumes feature gate.</p>
</td>
</tr>
<tr>
<td>
<code>data</code><br>
<em>
map[string][]byte
</em>
</td>
<td>
<em>(Optional)</em>
<p>Data contains the secret data. Each key must consist of alphanumeric
characters, &lsquo;-&rsquo;, &lsquo;_&rsquo; or &lsquo;.&rsquo;. The serialized form of the secret data is a
base64 encoded string, representing the arbitrary (possibly non-string)
data value here. Described in <a href="https://tools.ietf.org/html/rfc4648#section-4">https://tools.ietf.org/html/rfc4648#section-4</a></p>
</td>
</tr>
<tr>
<td>
<code>encryptedData</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.Envelopes">
map[string]./apis/secrets/v1beta1.Envelopes
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptedData contains the secret envelope encrypted secret data. Each key must consist of alphanumeric
characters, &lsquo;-&rsquo;, &lsquo;_&rsquo; or &lsquo;.&rsquo;. The values are arrays of Envelopes, which allows support for encrypting
a single entry with multiple providers. During decryption, the controller will stop on the first
successful provider and use the decrypted value it returns.</p>
</td>
</tr>
<tr>
<td>
<code>stringData</code><br>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>StringData allows specifying non-binary secret data in string form.
It is provided as a write-only convenience method.
All keys and values are merged into the data field on write, overwriting any existing values.
It is never output when reading from the API.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#secrettype-v1-core">
Kubernetes core/v1.SecretType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Used to facilitate programmatic handling of secret data.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecretStatus">
SealedSecretStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.SealedSecretSpec">SealedSecretSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecret">SealedSecret</a>)
</p>
<p>SealedSecretSpec defines the desired state of SealedSecret</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Metadata to use to create the Secret. For example to create it in with
a different name or to add labes and annotations.</p>
<p>Note, any namespace defined here will be ignored, an the corresponding
secret will always be created in the namespace of the SealedSecret.</p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>immutable</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Immutable, if set to true, ensures that data stored in the Secret cannot
be updated (only object metadata can be modified).
If not set to true, the field can be modified at any time.
Defaulted to nil.
This is a beta field enabled by ImmutableEphemeralVolumes feature gate.</p>
</td>
</tr>
<tr>
<td>
<code>data</code><br>
<em>
map[string][]byte
</em>
</td>
<td>
<em>(Optional)</em>
<p>Data contains the secret data. Each key must consist of alphanumeric
characters, &lsquo;-&rsquo;, &lsquo;_&rsquo; or &lsquo;.&rsquo;. The serialized form of the secret data is a
base64 encoded string, representing the arbitrary (possibly non-string)
data value here. Described in <a href="https://tools.ietf.org/html/rfc4648#section-4">https://tools.ietf.org/html/rfc4648#section-4</a></p>
</td>
</tr>
<tr>
<td>
<code>encryptedData</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.Envelopes">
map[string]./apis/secrets/v1beta1.Envelopes
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptedData contains the secret envelope encrypted secret data. Each key must consist of alphanumeric
characters, &lsquo;-&rsquo;, &lsquo;_&rsquo; or &lsquo;.&rsquo;. The values are arrays of Envelopes, which allows support for encrypting
a single entry with multiple providers. During decryption, the controller will stop on the first
successful provider and use the decrypted value it returns.</p>
</td>
</tr>
<tr>
<td>
<code>stringData</code><br>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>StringData allows specifying non-binary secret data in string form.
It is provided as a write-only convenience method.
All keys and values are merged into the data field on write, overwriting any existing values.
It is never output when reading from the API.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#secrettype-v1-core">
Kubernetes core/v1.SecretType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Used to facilitate programmatic handling of secret data.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="secrets.rustrial.org/v1beta1.SealedSecretStatus">SealedSecretStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#secrets.rustrial.org/v1beta1.SealedSecret">SealedSecret</a>)
</p>
<p>SealedSecretStatus defines the observed state of SealedSecret</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<p>Conditions</p>
</td>
</tr>
<tr>
<td>
<code>authorizations</code><br>
<em>
<a href="#secrets.rustrial.org/v1beta1.Authorization">
map[string]./apis/secrets/v1beta1.Authorization
</a>
</em>
</td>
<td>
<p>Authorization sources of the Key Encryption Keys references in this object.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<div class="admonition note">
<p class="last">This page was automatically generated with <code>gen-crd-api-reference-docs</code></p>
</div>
