// Package main contains the CLI of the secondary disk image generator.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	builder "github.com/GoogleCloudPlatform/ai-on-gke/gke-disk-image-builder"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return "my string representation"
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var containerImages stringSlice
	projectName := flag.String("project-name", "", "name of a gcp project where the script will be run")
	imageName := flag.String("image-name", "", "name of the image that will be generated")
	zone := flag.String("zone", "", "zone where the resources will be used to create the image creator resources")
	gcsPath := flag.String("gcs-path", "", "gcs location to dump the logs")
	diskSizeGb := flag.Int64("disk-size-gb", 10, "size of a disk that will host the unpacked images")
	gcpOAuth := flag.String("gcp-oauth", "", "path to GCP service account credential file")
	imagePullAuth := flag.String("image-pull-auth", "None", "auth mechanism to pull the container image, valid values: [None, ServiceAccountToken].\nNone means that the images are publically available and no authentication is required to pull them.\nServiceAccountToken means the service account oauth token will be used to pull the images.\nFor more information refer to https://cloud.google.com/compute/docs/access/authenticate-workloads#applications")
	timeout := flag.String("timeout", "20m", "Default timout for each step, defaults to 20m")
	flag.Var(&containerImages, "container-image", "container image to include in the disk image. This flag can be specified multiple times")

	flag.Parse()
	ctx := context.Background()

	td, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Panicf("invalid argument, timeout: %v, err: %v", timeout, err)
	}

	var auth builder.ImagePullAuthMechanism
	switch *imagePullAuth {
	case "":
		auth = builder.None
	case "None":
		auth = builder.None
	case "ServiceAccountToken":
		auth = builder.ServiceAccountToken
	default:
		log.Panicf("Please specify a valid value for the flag --image-pull-auth, valid values are [None, ServiceAccountToken]")
	}

	req := builder.Request{
		ImageName:       *imageName,
		ProjectName:     *projectName,
		Zone:            *zone,
		GCSPath:         *gcsPath,
		DiskSizeGB:      *diskSizeGb,
		GCPOAuth:        *gcpOAuth,
		ContainerImages: containerImages,
		Timeout:         td,
		ImagePullAuth:   auth,
	}

	if err = builder.GenerateDiskImage(ctx, req); err != nil {
		log.Panicf("unable to generate disk image: %v", err)
	}
	fmt.Printf("Image has successfully been created at: projects/%s/global/images/%s\n", req.ProjectName, req.ImageName)
}