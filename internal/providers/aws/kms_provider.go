package aws

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/rustrial/k8s-gitops-secrets/apis/secrets/v1beta1"
	secretsv1beta1 "github.com/rustrial/k8s-gitops-secrets/apis/secrets/v1beta1"
	"github.com/rustrial/k8s-gitops-secrets/internal/providers"
)

// KmsProviderFactory implementation.
type KmsProviderFactory struct {
	config  aws.Config
	pattern *regexp.Regexp
}

func pattern() *regexp.Regexp {
	return regexp.MustCompile(`^arn:aws:kms:([^:]+):[^:]+:(key|alias)/.+$`)
}

// NewKmsProviderFactory creates new KmsProviderFactory instance.
func NewKmsProviderFactory(config aws.Config) *KmsProviderFactory {
	pattern := pattern()
	return &KmsProviderFactory{
		config,
		pattern,
	}
}

// NewProvider creates new DecryptionProvider instance.
func (p *KmsProviderFactory) NewProvider(ctx context.Context, provider *secretsv1beta1.Provider) (providers.DecryptionProvider, error) {
	if provider.AwsKms != nil && p.pattern.MatchString(provider.KeyEncryptionKeyID) {
		return NewKmsProvider(p.config.Copy()), nil
	}
	return nil, nil
}

// KmsProvider providing AWS KMS implementation.
//
// Versions 1:
//
// Uses AES-256 bit symmtectric encryption in GCM mode, which provides AEAD.
type KmsProvider struct {
	config  aws.Config
	pattern *regexp.Regexp
}

// NewKmsProvider instance.
func NewKmsProvider(config aws.Config) *KmsProvider {
	pattern := pattern()
	return &KmsProvider{
		config,
		pattern,
	}
}

func (p *KmsProvider) adaptConfig(ctx context.Context, arn string) (*aws.Config, error) {
	match := p.pattern.FindStringSubmatch(arn)
	if match == nil {
		return nil, fmt.Errorf("Invalid AWS KSM Key ARN '%s'", arn)
	}
	cfg := p.config.Copy()
	cfg.Region = match[1]
	return &cfg, nil
}

// Encrypt plaintext into Envelope.
func (p *KmsProvider) Encrypt(ctx context.Context, plainText []byte, arn string) (*secretsv1beta1.Envelope, error) {
	cfg, err := p.adaptConfig(ctx, arn)
	if err != nil {
		return nil, err
	}
	svc := kms.NewFromConfig(*cfg)
	return p.encryptV1(ctx, svc, plainText, arn)
}

// DO NOT modify this function, otherwise we break the contract for Additional Authenticated Data (AAD)
// of AEAD and decrypt of existing SealedSecrets with provider.Version `1` will no longer work.
func v1Aad(dataKey []byte, kek string, provider *v1beta1.AwsKmsProvider) []byte {
	version := []byte{
		byte(provider.Version & 0xFF),
		byte((provider.Version >> 8) & 0xFF),
		byte((provider.Version >> 16) & 0xFF),
		byte((provider.Version >> 24) & 0xFF),
	}
	arn := []byte(kek)
	alg := []byte(provider.EncryptionAlgorithm)
	return append(append(append(append(version, arn...), dataKey...), alg...), provider.Nonce...)
}

func (p *KmsProvider) kmsGenerateDataKey(ctx context.Context, svc *kms.Client, arn string) (*kms.GenerateDataKeyOutput, error) {
	if arn == "arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000" {
		// This hack enables testing without depending on real AWS KMS service
		return &kms.GenerateDataKeyOutput{
			CiphertextBlob: testkey,
			Plaintext:      testkey,
			KeyId:          &arn,
		}, nil
	}
	return svc.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
		KeyId:   &arn,
		KeySpec: types.DataKeySpecAes256,
	})
}

// DO NOT modify this function, otherwise we break the contract for Additional Authenticated Data (AAD)
// of AEAD and decrypt of existing SealedSecrets with provider.Version `1` will no longer work.
func (p *KmsProvider) encryptV1(ctx context.Context, svc *kms.Client, plainText []byte, arn string) (*secretsv1beta1.Envelope, error) {
	kp, err := p.kmsGenerateDataKey(ctx, svc, arn)
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS KMS DataKey: %w", err)
	}
	// Input ARN might be a KMS Key alias, but we have to store a specific
	// (immutable) ARN, so let's use the ARN (KeyId) as returned by GenerateDataKey.
	if kp.KeyId != nil && len(*kp.KeyId) > 0 {
		arn = *kp.KeyId
	}
	block, err := aes.NewCipher(kp.Plaintext)
	if err != nil {
		return nil, fmt.Errorf("Error creating block Cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("Error creating GCM stream Cipher: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("Error while generating random nounce: %w", err)
	}
	provider := &v1beta1.AwsKmsProvider{
		Version:             1,
		EncryptionAlgorithm: string(types.DataKeySpecAes256),
		Nonce:               nonce,
	}
	// We use AEAD with Additional Authenticated Data (AAD) to detect whether the
	// secret has been manipulated.
	aad := v1Aad(kp.CiphertextBlob, arn, provider)
	ciphertext := gcm.Seal(make([]byte, 0), nonce, plainText, aad)
	return &v1beta1.Envelope{
		Provider: v1beta1.Provider{
			KeyEncryptionKeyID: arn,
			AwsKms:             provider,
		},
		DataEncryptionKey: kp.CiphertextBlob,
		CipherText:        ciphertext,
	}, nil
}

// Decrypt envelope encrypted secret value.
func (p *KmsProvider) Decrypt(ctx context.Context, envelope *secretsv1beta1.Envelope) ([]byte, error) {
	if envelope.AwsKms == nil {
		return nil, fmt.Errorf("Envelope has no AWS KMS provider information")
	}
	cfg, err := p.adaptConfig(ctx, envelope.KeyEncryptionKeyID)
	if err != nil {
		return nil, err
	}
	kms := kms.NewFromConfig(*cfg)
	switch envelope.AwsKms.Version {
	case 1:
		return p.decryptV1(ctx, kms, envelope)
	}
	return nil, fmt.Errorf("Unsupported encrypted Envelope version %d", envelope.AwsKms.Version)
}

var testkey []byte = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

func (p *KmsProvider) kmsDecrypt(ctx context.Context, svc *kms.Client, envelope *secretsv1beta1.Envelope) ([]byte, error) {
	if envelope.KeyEncryptionKeyID == "arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000" {
		// This hack enables testing without depending on real AWS KMS service
		return envelope.DataEncryptionKey, nil
	}
	dataKey, err := svc.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      envelope.DataEncryptionKey,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecSymmetricDefault,
		KeyId:               &envelope.KeyEncryptionKeyID,
	})
	if err != nil {
		return nil, err
	}
	return dataKey.Plaintext, nil
}

func (p *KmsProvider) decryptV1(ctx context.Context, svc *kms.Client, envelope *secretsv1beta1.Envelope) ([]byte, error) {
	dataKey, err := p.kmsDecrypt(ctx, svc, envelope)
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting dataKey using KMS key %s: %w", envelope.KeyEncryptionKeyID, err)
	}
	block, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, fmt.Errorf("Error creating block Cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("Error creating GCM stream Cipher: %w", err)
	}
	// We use AEAD with Additional Authenticated Data (AAD) to detect whether the
	// secret has been manipulated. This will make sure we do not update Secrets
	// with invalid values, which might cause failure or services that depend on that Secret.
	aad := v1Aad(envelope.DataEncryptionKey, envelope.KeyEncryptionKeyID, envelope.AwsKms)
	plainText, err := gcm.Open(make([]byte, 0), envelope.AwsKms.Nonce, envelope.CipherText, aad)
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting and authenticate secret: %w", err)
	}
	return plainText, nil
}
