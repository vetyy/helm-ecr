package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/service/ecr"
	"helm.sh/helm/v3/pkg/chart"
)

var (
	MaxResults int64 = 1000
)

const (
	HelmArtifactV1 = "application/vnd.cncf.helm.config.v1+json"
)

type ImageManifest struct {
	SchemaVersion int
	Config        *ImageManifestLayer
	Layers        []*ImageManifestLayer
}

type ImageManifestLayer struct {
	MediaType string
	Digest    string
	Size      int
}

func PopulateIndex(svc *ecr.ECR, repository *Repository, idx *IndexV3) error {
	repositoryFullName := repository.FullName()
	input := &ecr.DescribeImagesInput{
		RegistryId:     repository.RegistryID,
		RepositoryName: &repositoryFullName,
		MaxResults:     &MaxResults,
	}

	// Currently support listing maximum of 1000 images
	result, err := svc.DescribeImages(input)
	if err != nil {
		return fmt.Errorf("failed to describe images: %v", err)
	}

	// Populate Helm Index (index.yaml)
	for _, r := range result.ImageDetails {
		if r.ArtifactMediaType != nil && *r.ArtifactMediaType == HelmArtifactV1 {
			for _, imageTag := range r.ImageTags {
				chartMetadata := &chart.Metadata{
					Version:    *imageTag,
					APIVersion: chart.APIVersionV2,
					Name:       *repository.Name,
					Type:       "application",
				}

				err := idx.Add(chartMetadata, *imageTag, repository.URI(), *r.ImageDigest)
				if err != nil {
					return fmt.Errorf("failed to add to index: %v", err)
				}
			}
		}
	}

	return nil
}

func FetchChartAndWrite(svc *ecr.ECR, repository *Repository, dst io.Writer) error {
	imageTag := repository.Filename
	imageIdentifiers := []*ecr.ImageIdentifier{{ImageTag: imageTag}}

	repositoryFullName := repository.FullName()
	batchImages, err := svc.BatchGetImage(&ecr.BatchGetImageInput{
		RegistryId:     repository.RegistryID,
		RepositoryName: &repositoryFullName,
		ImageIds:       imageIdentifiers,
	})
	if err != nil {
		return fmt.Errorf("failed to batch get image: %v", err)
	}

	if len(batchImages.Images) == 0 {
		return fmt.Errorf("imageTag '%s' was not found", *imageTag)
	}

	imageManifest := &ImageManifest{}
	err = json.Unmarshal([]byte(*batchImages.Images[0].ImageManifest), imageManifest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal image manifest: %v", err)
	}

	result, err := svc.GetDownloadUrlForLayer(&ecr.GetDownloadUrlForLayerInput{
		RegistryId:     repository.RegistryID,
		RepositoryName: &repositoryFullName,
		LayerDigest:    &imageManifest.Layers[0].Digest,
	})
	if err != nil {
		return fmt.Errorf("failed to get download url for layer: %v", err)
	}

	resp, err := http.Get(*result.DownloadUrl)
	if err != nil {
		return fmt.Errorf("failed to download chart: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write chart content: %v", err)
	}

	return nil
}
