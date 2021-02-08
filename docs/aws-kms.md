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
