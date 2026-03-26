[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/k8s-gitops-secrets)](https://artifacthub.io/packages/search?k8s-gitops-secrets)

![Publish Charts](https://github.com/rustrial/k8s-gitops-secrets/workflows/publish/badge.svg)

# Kubernetes GitOps Secrets

Secure, easy-to-use and GitOps ready Kubernetes `Secret`s management using envelope
encrypted `SealedSecret` objects, which can be stored in Git repositories.
The solution is based on two software components:

- A [command line tool](#command-line-tool) called `seals` to envelope encrypt secret data.
- A Kubernetes [controller](#controller) deployed to your Kubernetes clusters, which will watch for
  `SealedSecret` and `KeyEncryptionKeyPolicy` objects and manages the corresponding (decrypted) `Secret` objects.

As of today, it support the following encryption providers:

- [AWS KMS](docs/aws-kms.md) which offers symetric AES-256-GCM envelope encryption.

_If you miss your favorite encryption provider in this list please feel free to open
a pull request._

---

# Concepts

`SealedSecret` objects are envelope encrypted secrets. `KeyEncryptionKeyPolicy`
objects define which namespaces are authorized to use a specific _Key Encryption Key_
(_KEK_) to decrypt `SealedSecret` objects into plaintext Kubernetes `Secret` objects.
If a `SealedSecret` references any \_KEK\_, for which its namespace is not authorized,
then the contoller will reject it and not decrypt that `SealedSecret`.

For example the following `KeyEncryptionKeyPolicy` authorizes the namespaces
`kube-system` and `sre-services` to use the _KEK_ with ID
`arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000`.

```yaml
apiVersion: secrets.rustrial.org/v1beta1
kind: KeyEncryptionKeyPolicy
metadata:
  name: sre-kek-policy
  namespace: kube-system
spec:
  keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
  namespaces:
    - "kube-system"
    - "sre-services"
```

The following `KeyEncryptionKeyPolicy` authorizes **all** namespaces to use the _KEK_
with ID `arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000`,
by using the wildcart "`*`". WARNING, only use such global authorization if you really
know what you are doing, because they might introduce security/privacy leaks.

```yaml
apiVersion: secrets.rustrial.org/v1beta1
kind: KeyEncryptionKeyPolicy
metadata:
  name: fully-open-kek-policy
  namespace: kube-system
spec:
  keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
  namespaces:
    - "*"
```

## Audience Narrowing

Audience narrowing allows you to restrict **where** a `SealedSecret` can be decrypted by
binding the encrypted data to a specific context (e.g. namespace, secret name, AWS region,
AWS account, or AWS partition). The audience constraints are cryptographically enforced via
the encryption provider's context mechanism (e.g.
[AWS KMS Encryption Context](https://docs.aws.amazon.com/kms/latest/developerguide/encrypt_context.html)),
so they cannot be bypassed without access to the _KEK_.

> **⚠️ Security Warning — multi-tenant clusters:**
> Using audience narrowing is **highly recommended**, especially in multi-tenant setups.
> Without it, an attacker who has read-only access to a `SealedSecret` (e.g. via a Git
> repository) could simply copy it into their own namespace and have the controller decrypt
> it there. By specifying `--namespace` and/or `--name` at encryption time you ensure that
> the secret can **only** be decrypted into the intended namespace and secret name.

The following audience constraints are available when encrypting with the `seals` CLI:

| CLI flag | Constraint |
|---|---|
| `--namespace` | Kubernetes namespace(s) in which the `SealedSecret` can be decrypted |
| `--name` | Kubernetes Secret name(s) into which the `SealedSecret` can be decrypted |
| `--region` | AWS Region(s) in which the `SealedSecret` can be decrypted |
| `--account` | AWS Account ID(s) in which the `SealedSecret` can be decrypted |
| `--partition` | AWS Partition(s) in which the `SealedSecret` can be decrypted |

Each flag can be specified multiple times to allow multiple values. The controller will
verify at decryption time that the `SealedSecret`'s actual context (namespace, name,
controller's AWS region/account/partition) matches at least one of the allowed values for
every constraint that was set. If any constraint is not satisfied, decryption is refused.

# Usage

## Controller

The controller can be operated in cluster or namespace scoped mode. In cluster scope mode
it will track `SealedSecret` objects in all namespaces and in namespaced mode it will only
track objects in one specific namespace. To enable namespaces mode use the `--namespace`
command line flag (e.g. `--namespace my-namespace`).

No matter whether the controller is running in cluster or namespaced mode, it will only
track `KeyEncryptionKeyPolicy` objects in the namespace in which the controller itself
is running. `KeyEncryptionKeyPolicy` in other namespaces will be ignored.

---

## Command Line Tool

**Encrypt secret data**:

_Envelope encrypting is provider specific and for demonstration purpose we use
the AWS KMS provider. Please consult the provider specific documentation
(see above) for instructions for your encryption provider._

Use the `aws-kms` subcommand of the `seals` command line tool to encrypt data. It expects one or several
AWD KMS _CMK_ ARNs and will read the secret data from `STDIN` and output the encrypted Envelopes as YAML
array on `STDOUT`.

```shell
cat my-secret.txt | seals aws-kms arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
```

To narrow the audience so the secret can only be decrypted in namespace `default` with
secret name `sealedsecret-sample` (recommended — see [Audience Narrowing](#audience-narrowing)):

```shell
cat my-secret.txt | seals aws-kms \
  --namespace default \
  --name sealedsecret-sample \
  arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
```

**Example output**

```yaml
- keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
  awsKms:
    encryptionAlgorithm: AES_256
    encryptionContext:
      names: sealedsecret-sample
      namespaces: default
    nonce: 1Mjlj3HkLXjH2x/N
    version: 1
  cipherText: O595N8v+MrHeL6G8jOaS1/yh
  dataEncryptionKey: |-
    AQIDAHhyGbrAGHIA3UUTHqS9CfjLPE3hSeN3A8tabfy8f+NyfwFSFLgzFz1Z08XtFLRVkdfHAAA
    AfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMulVVae
    SGSNLWMe5zAgEQgDuwHu37T2zba4wZzTSMnH/o+shkvYqv8jp7ngPIvJI7Un9H/L8TwKKuBGAee
    2FkO9U/xMLak0/XMl++qg==
```

The encrypted output can then be integrated into a `SealedSecret` object, which usally
will be stored in Git and later deployed to Kubernetes clusters.

```yaml
apiVersion: secrets.rustrial.org/v1beta1
kind: SealedSecret
metadata:
  name: sealedsecret-sample
  namespace: default
spec:
  # Optional metadata
  metadata:
    labels:
      mylabel: my-label-value
  type: opaque
  # Plaintext entries
  data:
    some-non-sensitive-data: SSBhbSBub3QgYSBzZWNyZXQgOi0pCg==
  stringData:
    more-non-sensitive-data: "I am not a secret either :-)"
  # Envelope encrypted entries
  encryptedData:
    my-secret:
      # Note, this is an array to enable encrypting the same entry with
      # multple KEKs (providers).
      - keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
        awsKms:
          encryptionAlgorithm: AES_256
          encryptionContext:
            names: sealedsecret-sample
            namespaces: default
          nonce: 1Mjlj3HkLXjH2x/N
          version: 1
        cipherText: O595N8v+MrHeL6G8jOaS1/yh
        dataEncryptionKey: |-
          AQIDAHhyGbrAGHIA3UUTHqS9CfjLPE3hSeN3A8tabfy8f+NyfwFSFLgzFz1Z08XtFLRVkdfHAAA
          AfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMulVVae
          SGSNLWMe5zAgEQgDuwHu37T2zba4wZzTSMnH/o+shkvYqv8jp7ngPIvJI7Un9H/L8TwKKuBGAee
          2FkO9U/xMLak0/XMl++qg==
```

Above `SealedSecret` will result in below `Secret` once it is decrypted by the controller.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sealedsecret-sample
  namespace: default
  labels:
    mylabel: my-label-value
type: opaque
data:
  some-non-sensitive-data: SSBhbSBub3QgYSBzZWNyZXQgOi0pCg==
  my-secret: d29ybGQK
stringData:
  more-non-sensitive-data: "I am not a secret either :-)"
```

## Getting Started

**Add Helm Repository**

_GitOps secret controller_ can be installed via Helm Chart, which by default will use the prebuilt OCI Images for Linux (`amd64` and `arm64`) from [DockerHub](https://hub.docker.com/r/rustrial/k8s-gitops-secrets-controller).

```shell
helm repo add aws-eks-iam-auth-controller https://rustrial.github.io/k8s-gitops-secrets
```

**Install Helm Chart**

```shell
helm install k8s-gitops-secrets-controller k8s-gitops-secrets/rustrial-k8s-gitops-secrets-controller
```

---
