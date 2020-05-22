package v124

import (
	"github.com/dehorsley/dbbc3mcast/versions"
)

const Version = "DDC_V,124"

func init() {
	versions.Add(Version, &Dbbc3DdcMulticast{})
}
