package yum

import (
	"fmt"
	"net/url"
)

type Package struct {
	Name         string
	Architecture string
	Epoch        int32
	Version      string
	Release      string
	Size         int64
	Repository   string
	Summary      string
	URL          *url.URL
	License      string
	Description  string
}

func (p *Package) String() string {
	return fmt.Sprintf("%v-%v:%v-%v.%v", p.Name, p.Epoch, p.Version, p.Release, p.Architecture)
}
