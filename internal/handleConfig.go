// Package internal handles node configuration commands
package internal

import (
	"github.com/iotdomain/iotdomain-go/types"
	"github.com/sirupsen/logrus"
)

// HandleConfigCommand for handling node configuration changes
// Not supported
func (app *IsyApp) HandleConfigCommand(address string, config types.NodeAttrMap) types.NodeAttrMap {
	logrus.Infof("IsyApp.HandleConfigCommand for %s. Ignored as this isn't supported", address)
	return nil
}
