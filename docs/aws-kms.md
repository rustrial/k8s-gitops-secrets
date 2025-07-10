# AWS KMS Provider

The [AWS KMS](https://aws.amazon.com/kms/) provider uses
[envelope encryption](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#enveloping),
where the _dataKey_ is encrypted using
[Customer Master Key](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#master_keys)
(_CMK_) stored in AWS KMS. *CMK*s cannot be extracted from KMS and therefore offer a very high level of
security as long as only authorized clients are entitled to use them for decrypting.

## Role & Operation Model:

- _Cloud Operator_ creates _CMK_ and entitles _DevOps Engineers_ to use it for encryption.
- _DevOps Engineers_ encrypt their secrets using the _CMK_ to envelop encrypt a _Data Key_, which
  is used to encrypt the secret's plaintext value. The encrypted objects can then be stored
  as SealedSecrets in any (public) Git repository, from where it might be deployed into
  (multiple) Kubernetes clusters.
- _Sealed Secret Controller_ (a Kubernetes controller) is deployed to the
  Kubernetes clusters and watches for _SealedSecret_ objects, decrypting them and creating resp.
  updating the corresponding _Secret_ objects.

## Usage

**Encrypt secret data**:

Use the `aws-kms` subcommand of the `seals` command line tool to encrypt data. It expects one of several
AWD KMS _CMK_ ARNs and will read the secret data from `STDIN` and output the encrypted Envelops as YAML
array on `STDOUT`.

```shell
echo 'mysecret' | seals aws-kms arn:aws:kms:eu-central-1:000000000000:key/6a06295d-f3c1-4462-9fba-67f13120963d
```

**Example output**

```yaml
- keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
  awsKms:
    encryptionAlgorithm: AES_256
    nonce: 1Mjlj3HkLXjH2x/N
  cipherText: O595N8v+MrHeL6G8jOaS1/yh
  dataEncryptionKey: |-
    AQIDAHhyGbrAGHIA3UUTHqS9CfjLPE3hSeN3A8tabfy8f+NyfwFSFLgzFz1Z08XtFLRVkdfHAAA
    AfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMulVVae
    SGSNLWMe5zAgEQgDuwHu37T2zba4wZzTSMnH/o+shkvYqv8jp7ngPIvJI7Un9H/L8TwKKuBGAee
    2FkO9U/xMLak0/XMl++qg==
```

### Encryption Context - Audience

Optionally, encrypted values can be associated with an
[Encryption Context](https://docs.aws.amazon.com/kms/latest/developerguide/encrypt_context.html),
which allows further controll on the audience resp. where such values can be decrypted. When
encrypting values the audience can be provided using the below command-line (CLI) flags, each
CLI flag can be provided multiple times to add multiple audiences.

|CLI flag|Audience|
|---|---|
|`--namespace`|Kubernetes [Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) in which the *SealedSecret* can be decryped|
|`--name`|Name of the resulting Kubernetes [Secret](https://kubernetes.io/docs/concepts/configuration/secret/) into which the *SealedSecret* can be decryped|
|`--region`|AWS [Region](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html) in which the *SealedSecret* can be decryped|
|`--account`|AWS [Account ID](https://docs.aws.amazon.com/accounts/latest/reference/manage-acct-identifiers.html#FindAccountId) in which the *SealedSecret* can be decryped|
|`--partition`|AWS [Partition](https://docs.aws.amazon.com/whitepapers/latest/aws-fault-isolation-boundaries/partitions.html) in which the *SealedSecret* can be decryped|

If no audience constraining encryption context is provided, then the encrypted content can be
decrypted by all audiences that can use the corresponding KMS Key for decryption operations.
If audience constraining encryption context is provided, then the Kubernetes [controller](../README.md#controller) will enforce it by creating an Encryption Context for decryption
according to the following rules:

- `namespace`: the value of .metadata.namespace of the SealedSecret instance that it is processing, which is also the namespace of the Secret that it will create, must match any of the provided namespace.
- `name`: the value of .metadata.name of the SealedSecret instance that it is processing, which is also the name of the Secret that it will create, must match any of the provided names.
- `region`: the AWS Region of the controller must match any of the provided AWS Regions.
- `account`: the AWS Account of the controller's IAM Principal (Role) must match any of the provided AWS Accounts.
- `region`: the AWS Partition of the controller's IAM Principal (Role) must match any of the provided AWS Partitions.

For example to constrain the audience to Kubernetes Namespace `my-namespace` we can encrypted values
like this:

```shell
echo 'mysecret' | seals aws-kms --namespace my-namespace arn:aws:kms:eu-central-1:000000000000:key/6a06295d-f3c1-4462-9fba-67f13120963d
```

Or to constrain the audience to Kubernetes Namespaces `my-namespace` and `your-namespace`
in AWS Regions `eu-central-1` and `eu-west-1` we can encrypted values like this:

```shell
echo 'mysecret' | seals aws-kms --namespace my-namespace --namespace your-namespace --region eu-central-1 --region eu-west-1 arn:aws:kms:eu-central-1:000000000000:key/6a06295d-f3c1-4462-9fba-67f13120963d
```
