package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ghodss/yaml"
	"github.com/rustrial/k8s-gitops-secrets/api/secrets/v1beta1"
	awsProvider "github.com/rustrial/k8s-gitops-secrets/internal/providers/aws"
	"github.com/spf13/cobra"
)

var audience = awsProvider.AwsAudience{}

func init() {
	encrypt.Flags().BoolVarP(&verify, "verify", "v", false, "Verify output by decrypting it, needs Decrypt permission on the KMS CMK.")
	encrypt.Flags().StringArrayVarP(&audience.Namespaces, "namespace", "", []string{}, "Audience: Kubernetes namespace(s)")
	encrypt.Flags().StringArrayVarP(&audience.Names, "name", "", []string{}, "Audience: Kubernetes Secret name(s)")
	encrypt.Flags().StringArrayVarP(&audience.Regions, "region", "", []string{}, "Audience: AWS Region(s)")
	encrypt.Flags().StringArrayVarP(&audience.OrgUnits, "account", "", []string{}, "Audience: AWS Account ID(s)")
	encrypt.Flags().StringArrayVarP(&audience.Partitions, "partition", "", []string{}, "Audience: AWS Partition(s)")
	rootCmd.AddCommand(encrypt)
}

var verify = false

var encrypt = &cobra.Command{
	Use:        "aws-kms arn:aws:kms:eu-central-1:000000000000:key/6a06295d-f3c1-4462-9fba-67f13120963d",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "Seal (envelope encrypt) secret with AWS KMS key.",
	Long:       `Seal (envelope encrypt) secret passed on STDIN with AWS KMS key.`,
	Example:    "cat secret.txt | seals aws-kms arn:aws:kms:eu-central-1:000000000000:key/6a06295d-f3c1-4462-9fba-67f13120963d",
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to load AWS SDK config, %s", err)
			os.Exit(1)
		}
		provider := awsProvider.NewKmsProvider(cfg)
		plainText, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading plain text data from STDIN: %s\n", err)
			os.Exit(1)
		}
		plainText = bytes.Trim(plainText, "\n\r")
		output := make(v1beta1.Envelopes, 0)
		for _, arn := range args {
			envelope, err := provider.Encrypt(ctx, plainText, arn, &audience)
			if verify {
				pt, err := provider.Decrypt(ctx, envelope, &audience)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while decrypting plaintext data: %s\n", err)
					os.Exit(1)
				}
				if !bytes.Equal(pt, plainText) {
					fmt.Fprintf(os.Stderr, "Decrypting plaintext does not match the input plaintext\n")
					os.Exit(1)
				}
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while encrypting plain text data: %s\n", err)
				os.Exit(1)
			}
			output = append(output, *envelope)
		}
		o, err := yaml.Marshal(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding result: %s\n", err)
			os.Exit(1)
		}
		fmt.Print(string(o))
	},
}
