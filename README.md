[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/k8s-gitops-secrets)](https://artifacthub.io/packages/search?k8s-gitops-secrets)

![Publish Charts](https://github.com/rustrial/k8s-gitops-secrets/workflows/publish/badge.svg)

# Kubernetes GitOps Secrets

Secure, easy-to-use and GitOps ready Kubernetes `Secret`s management using evelope
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

# Usage

## Controller

The controller can be operated in cluster or namespace scoped mode. In cluster scope mode
it willtrack `SealedSecret` objects in all namespaces and in namespaced mode it will only
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

Use the `aws-kms` subcommand of the `seals` command line tool to encrypt data. It expects one of several
AWD KMS _CMK_ ARNs and will read the secret data from `STDIN` and output the encrypted Envelopes as YAML
array on `STDOUT`.

```shell
cat my-secret.txt | seals arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
```

**Example output**

```yaml
- keyEncryptionKeyId: arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000
  awsKms:
    encryptionAlgorithm: AES_256
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
spec:
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
