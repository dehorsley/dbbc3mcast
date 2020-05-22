package v125

import (
	"github.com/dehorsley/dbbc3mcast/versions"
)

const Version = "DDC_U,125"

func init() {
	versions.Add(Version, &Dbbc3DdcMulticast{})
}
