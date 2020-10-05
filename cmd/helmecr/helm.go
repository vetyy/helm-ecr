package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

type IndexV3 struct {
	index *repo.IndexFile
}

func (idx *IndexV3) Has(name, version string) bool {
	return idx.index.Has(name, version)
}

func (idx *IndexV3) MarshalBinary() (data []byte, err error) {
	return yaml.Marshal(idx.index)
}

func (idx *IndexV3) Add(metadata interface{}, filename, baseURL, digest string) error {
	md, ok := metadata.(*chart.Metadata)
	if !ok {
		return fmt.Errorf("metadata is not *chart.Metadata")
	}

	idx.index.Add(md, filename, baseURL, digest)
	return nil
}

func NewIndex() (*IndexV3, error) {
	_, err := IsHelm3()
	if err != nil {
		return nil, fmt.Errorf("helm v2 is not supported: %v", err)
	}
	return &IndexV3{index: repo.NewIndexFile()}, nil
}

func IsHelm3() (bool, error) {
	if os.Getenv("TILLER_HOST") != "" {
		return false, fmt.Errorf("TILLER_HOST found")
	}

	cmd := exec.Command("helm", "version", "--short", "--client")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to execute 'helm' command: %v", err)
	}

	return strings.HasPrefix(string(out), "v3."), nil
}
