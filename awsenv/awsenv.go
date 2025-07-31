// Package awsenv extends AWS SDK configuration with additional environment
// variables.
package awsenv

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 returns additional options for s3.New and s3.NewFromConfig.
//
// Environment variables:
//
//	AWS_S3_USE_PATH_STYLE, AWS_S3_FORCE_PATH_STYLE
//	  Sets s3.Options.UsePathStyle to true.
//
//	AWS_S3_ENDPOINT
//	  Overrides AWS S3 endpoint (e.g. http://localhost:9000 for integration
//	  tests with MinIO).
func S3() []func(o *s3.Options) {
	return []func(o *s3.Options){
		func(o *s3.Options) {
			_, v1 := os.LookupEnv("AWS_S3_FORCE_PATH_STYLE")
			_, v2 := os.LookupEnv("AWS_S3_USE_PATH_STYLE")
			o.UsePathStyle = v1 || v2
		},
		func(o *s3.Options) {
			endpoint, ok := os.LookupEnv("AWS_S3_ENDPOINT")
			if !ok {
				return
			}
			o.EndpointResolver = s3.EndpointResolverFromURL(endpoint)
		},
	}
}

// S3Bucket returns S3 bucket name from environment.
//
// Environment variables:
//
//	AWS_S3_BUCKET
//	  S3 bucket name to use.
func S3Bucket() string {
	return os.Getenv("AWS_S3_BUCKET")
}
