package secret

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/argsutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/surveyutils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/spf13/cobra"
)

func awsCmd(opts *options.Options) *cobra.Command {
	input := opts.Create.InputSecret.AwsSecret
	cmd := &cobra.Command{
		Use:   "aws",
		Short: `Create an AWS secret with the given name`,
		Long:  `Create an AWS secret with the given name`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := argsutils.MetadataArgsParse(opts, args); err != nil {
				return err
			}
			if opts.Top.Interactive {
				// and gather any missing args that are available through interactive mode
				if err := AwsSecretArgsInteractive(&opts.Metadata, &input); err != nil {
					return err
				}
			}
			// create the secret
			if err := createAwsSecret(opts.Top.Ctx, opts.Metadata, input); err != nil {
				return err
			}
			fmt.Printf("Created AWS secret [%v] in namespace [%v]\n", opts.Metadata.Name, opts.Metadata.Namespace)
			return nil
		},
	}

	flags := cmd.Flags()
	flagutils.AddMetadataFlags(flags, &opts.Metadata)
	flags.StringVar(&input.AccessKey, "access-key", "", "aws access key")
	flags.StringVar(&input.SecretKey, "secret-key", "", "aws secret key")

	return cmd
}

func AwsSecretArgsInteractive(meta *core.Metadata, input *options.AwsSecret) error {
	if err := surveyutils.InteractiveNamespace(&meta.Namespace); err != nil {
		return err
	}

	if err := cliutil.GetStringInput("name of secret", &meta.Name); err != nil {
		return err
	}

	if err := cliutil.GetStringInput("Enter AWS Access Key ID (leave empty to read credentials from "+
		"~/.aws/credentials): ", &input.AccessKey); err != nil {
		return err
	}
	if err := cliutil.GetStringInput("Enter AWS Secret Key (leave empty to read credentials from "+
		"~/.aws/credentials): ", &input.SecretKey); err != nil {
		return err
	}

	return nil
}

func createAwsSecret(ctx context.Context, meta core.Metadata, input options.AwsSecret) error {
	if input.AccessKey == "" || input.SecretKey == "" {
		fmt.Printf("access key or secret key not provided, reading credentials from ~/.aws/credentials")
		creds := credentials.NewSharedCredentials("", "")
		val, err := creds.Get()
		if err != nil {
			return err
		}
		input.SecretKey = val.SecretAccessKey
		input.AccessKey = val.AccessKeyID
	}
	secret := &gloov1.Secret{
		Metadata: meta,
		Kind: &gloov1.Secret_Aws{
			Aws: &gloov1.AwsSecret{
				AccessKey: input.AccessKey,
				SecretKey: input.SecretKey,
			},
		},
	}

	secretClient := helpers.MustSecretClient()
	if _, err := secretClient.Write(secret, clients.WriteOpts{Ctx: ctx}); err != nil {
		return err
	}
	return nil
}
