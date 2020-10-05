package main

import (
	"fmt"
	"regexp"
)

var (
	RepositoryRegexp        = regexp.MustCompile(`ecr://(?P<registry_id>[0-9]+).dkr.ecr.(?P<region>.+).amazonaws.com/(?P<namespace>.+/)?(?P<repository>.+)/(?P<filename>.+)`)
	PartialRepositoryRegexp = regexp.MustCompile(`ecr://(?P<namespace>.+/)?(?P<repository>.+)/(?P<filename>.+)`)
)

type Repository struct {
	Filename   *string
	Name       *string
	Namespace  *string
	Region     *string
	RegistryID *string
}

func (r *Repository) FullName() string {
	return fmt.Sprintf("%s%s", *r.Namespace, *r.Name)
}

func NewRepository(uri string) (repository *Repository, err error) {
	match := RepositoryRegexp.FindStringSubmatch(uri)
	if len(match) == 6 {
		return &Repository{
			RegistryID: &match[1],
			Region:     &match[2],
			Namespace:  &match[3],
			Name:       &match[4],
			Filename:   &match[5],
		}, nil
	}

	match = PartialRepositoryRegexp.FindStringSubmatch(uri)
	if len(match) == 4 {
		return &Repository{
			Namespace: &match[1],
			Name:      &match[2],
			Filename:  &match[3],
		}, nil
	}
	return nil, fmt.Errorf("unknown registry string found")
}
