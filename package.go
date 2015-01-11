package yum

import (
  "fmt"
	"net/url"
)

type Package struct {
	Name         string
	Architecture string
	Epoch        string
	Version      string
	Release      string
	Size         int64
	Repository   string
	Summary      string
	URL          *url.URL
	License      string
	Description  string
}

func (p *Package) String() string{
  return fmt.Sprintf("%v-%v:%v-%v.%v", p.Name, p.Epoch, p.Version, p.Release, p.Architecture)
}

/*
func (pkg *Package) Info() error {
  pkgs, err := Info(fmt.Sprintf("%v-%v:%v-%v.%v", pkg.Name, pkg.Epoch, pkg.Version, pkg.Release, plg.Architecture))
  if err != nil {
    return err
  }

  if len(pkgs) != 1 {
    return nil
  }

  pkg = pkgs[0]
  return nil
}
*/
