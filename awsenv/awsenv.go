// Package awsenv extends loading aws.Config with additional environment variables.
package awsenv

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
)

// Session returns additional configs for session.NewSession.
func Session() []*aws.Config {
	hasConfig := false
	awsConfig := &aws.Config{}

	if _, ok := os.LookupEnv("AWS_DISABLE_TLS"); ok {
		awsConfig.DisableSSL = aws.Bool(ok)
		hasConfig = true
	}

	if !hasConfig {
		return nil
	}
	return []*aws.Config{awsConfig}
}

// S3 returns additional configs for s3.New.
func S3() []*aws.Config {
	hasConfig := false
	awsConfig := &aws.Config{}

	if _, ok := os.LookupEnv("AWS_S3_FORCE_PATH_STYLE"); ok {
		awsConfig.S3ForcePathStyle = aws.Bool(ok)
		hasConfig = true
	}
	if endpoint, ok := os.LookupEnv("AWS_S3_ENDPOINT"); ok {
		awsConfig.Endpoint = aws.String(endpoint)
		hasConfig = true
	}

	if !hasConfig {
		return nil
	}
	return []*aws.Config{awsConfig}
}
