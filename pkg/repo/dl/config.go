package dl

import "fmt"

type Dl struct {
	url *string
}

func NewDl(url *string) *Dl {
	return &Dl{url: url}
}

func (d *Dl) GetDlUrl(name string, version string) string {
	return fmt.Sprintf("%s/packages/%s/%s", *d.url, name, version)
}
