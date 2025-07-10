package aws

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rustrial/k8s-gitops-secrets/api/secrets/v1beta1"
	secretsv1beta1 "github.com/rustrial/k8s-gitops-secrets/api/secrets/v1beta1"
	"github.com/rustrial/k8s-gitops-secrets/internal/providers"
)

// KmsProviderFactory implementation.
type KmsProviderFactory struct {
	config       aws.Config
	pattern      *regexp.Regexp
	awsRegion    string
	awsPartition string
	awsAccountId string
}

func pattern() *regexp.Regexp {
	return regexp.MustCompile(`^arn:[^:]+:kms:([^:]+):[^:]+:(key|alias)/.+$`)
}

// NewKmsProviderFactory creates new KmsProviderFactory instance.
func NewKmsProviderFactory(ctx context.Context, config aws.Config) (*KmsProviderFactory, error) {
	pattern := pattern()
	p := KmsProviderFactory{}
	p.config = config
	p.pattern = pattern
	// Get current region from config
	p.awsRegion = config.Region
	// Lookup current AWS Partition and Account ID
	stsClient := sts.NewFromConfig(config)
	callerIdentity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return &p, err
	} else {
		if callerIdentity.Account != nil {
			// Extract partition from ARN (format: arn:partition:service:region:account-id:resource)
			if callerIdentity.Arn != nil {
				arnParts := strings.Split(*callerIdentity.Arn, ":")
				if len(arnParts) >= 2 {
					p.awsPartition = arnParts[1]
				}
				if len(arnParts) >= 5 {
					p.awsAccountId = arnParts[4]
				}
			}
		}
	}
	return &p, nil
}

// NewProvider creates new DecryptionProvider instance.
func (pf *KmsProviderFactory) NewProvider(ctx context.Context, provider *secretsv1beta1.Provider) (providers.DecryptionProvider, error) {
	if provider.AwsKms != nil && pf.pattern.MatchString(provider.KeyEncryptionKeyID) {
		p := pf.NewKmsProvider()
		return p, nil
	}
	return nil, nil
}

// NewProvider creates new DecryptionProvider instance.
func (pf *KmsProviderFactory) NewKmsProvider() *KmsProvider {
	p := KmsProvider{
		config:       pf.config.Copy(),
		pattern:      pf.pattern,
		awsRegion:    pf.awsRegion,
		awsPartition: pf.awsPartition,
		awsAccountId: pf.awsAccountId,
	}
	return &p
}

// KmsProvider providing AWS KMS implementation.
//
// Versions 1:
//
// Uses AES-256 bit symmtectric encryption in GCM mode, which provides AEAD.
type KmsProvider struct {
	config       aws.Config
	pattern      *regexp.Regexp
	awsRegion    string
	awsPartition string
	awsAccountId string
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
func (p *KmsProvider) Encrypt(ctx context.Context, plainText []byte, arn string, audience providers.Audience) (*secretsv1beta1.Envelope, error) {
	cfg, err := p.adaptConfig(ctx, arn)
	if err != nil {
		return nil, err
	}
	svc := kms.NewFromConfig(*cfg)
	return p.encryptV1(ctx, svc, plainText, arn, audience)
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

func (p *KmsProvider) kmsGenerateDataKey(ctx context.Context, svc *kms.Client, arn string, encryptionContext map[string]string) (*kms.GenerateDataKeyOutput, error) {
	if arn == "arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000" {
		// This hack enables testing without depending on real AWS KMS service
		return &kms.GenerateDataKeyOutput{
			CiphertextBlob: testkey,
			Plaintext:      testkey,
			KeyId:          &arn,
		}, nil
	}
	return svc.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
		KeyId:             &arn,
		KeySpec:           types.DataKeySpecAes256,
		EncryptionContext: encryptionContext,
	})
}

// DO NOT modify this function, otherwise we break the contract for Additional Authenticated Data (AAD)
// of AEAD and decrypt of existing SealedSecrets with provider.Version `1` will no longer work.
func (p *KmsProvider) encryptV1(ctx context.Context, svc *kms.Client, plainText []byte, arn string, audience providers.Audience) (*secretsv1beta1.Envelope, error) {
	encryptionContext := audience.Audience()
	kp, err := p.kmsGenerateDataKey(ctx, svc, arn, encryptionContext)
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
		EncryptionContext:   encryptionContext,
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

func (p *KmsProvider) Audience(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret, envelope *secretsv1beta1.Envelope) (providers.Audience, []string) {
	context := envelope.AwsKms.EncryptionContext
	audience := awsAudienceFromMap(context)
	failures := make([]string, 0)
	if len(audience.Namespaces) > 0 && !slices.Contains(audience.Namespaces, sealedSecret.Namespace) {
		audience.Namespaces = []string{}
		failures = append(failures, fmt.Sprintf("namespace '%s' is not in audience", sealedSecret.Namespace))
	}
	if len(audience.Names) > 0 && !slices.Contains(audience.Names, sealedSecret.Name) {
		audience.Names = []string{}
		failures = append(failures, fmt.Sprintf("name '%s' is not in audience", sealedSecret.Name))
	}
	if len(audience.Partitions) > 0 && !slices.Contains(audience.Partitions, p.awsPartition) {
		audience.Partitions = []string{}
		failures = append(failures, fmt.Sprintf("partition '%s' is not in audience", p.awsPartition))
	}
	if len(audience.Regions) > 0 && !slices.Contains(audience.Regions, p.awsRegion) {
		audience.Regions = []string{}
		failures = append(failures, fmt.Sprintf("region '%s' is not in audience", p.awsRegion))
	}
	if len(audience.OrgUnits) > 0 && !slices.Contains(audience.OrgUnits, p.awsAccountId) {
		audience.OrgUnits = []string{}
		failures = append(failures, fmt.Sprintf("account '%s' is not in audience", p.awsAccountId))
	}
	return audience, failures
}

// Decrypt envelope encrypted secret value.
func (p *KmsProvider) Decrypt(ctx context.Context, envelope *secretsv1beta1.Envelope, audience providers.Audience) ([]byte, error) {
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
		return p.decryptV1(ctx, kms, envelope, audience)
	}
	return nil, fmt.Errorf("Unsupported encrypted Envelope version %d", envelope.AwsKms.Version)
}

var testkey []byte = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

func (p *KmsProvider) kmsDecrypt(ctx context.Context, svc *kms.Client, envelope *secretsv1beta1.Envelope, audience providers.Audience) ([]byte, error) {
	if envelope.KeyEncryptionKeyID == "arn:aws:kms:eu-central-1:000000000000:key/00000000-0000-0000-0000-000000000000" {
		// This hack enables testing without depending on real AWS KMS service
		return envelope.DataEncryptionKey, nil
	}
	dataKey, err := svc.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      envelope.DataEncryptionKey,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecSymmetricDefault,
		KeyId:               &envelope.KeyEncryptionKeyID,
		EncryptionContext:   audience.Audience(),
	})
	if err != nil {
		return nil, err
	}
	return dataKey.Plaintext, nil
}

func (p *KmsProvider) decryptV1(ctx context.Context, svc *kms.Client, envelope *secretsv1beta1.Envelope, audience providers.Audience) ([]byte, error) {
	dataKey, err := p.kmsDecrypt(ctx, svc, envelope, audience)
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
