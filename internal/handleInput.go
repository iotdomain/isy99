// Package internal handles node input commands
package internal

import (
	"strings"
	"time"

	"github.com/iotdomain/iotdomain-go/types"
)

// // HandleConfig updates for nodes
// func (adapter *Isy99Adapter) HandleConfig(
// 	device *nodes.Node, inOutput *nodes.NodeInOutput, changes map[string]string, isEncrypted bool) {

// 	for attrName, configValue := range changes {
// 		if inOutput != nil {
// 			inOutput.UpdateConfig(attrName, configValue, false)
// 		} else {
// 			device.UpdateConfig(attrName, configValue, false)
// 		}
// 	}
// }

// SwitchOnOff turns lights or switch on or off. A payload '0', 'off' or 'false' turns off, otherwise it turns on
func (app *IsyApp) SwitchOnOff(input *types.InputDiscoveryMessage, onOffString string) error {

	// any non-zero, false or off value is considered on
	newValue := true
	if onOffString == "0" || strings.ToLower(onOffString) == "off" || strings.ToLower(onOffString) == "false" {
		newValue = false
	}
	prevValue := "n/a"
	prevOutputValue := app.pub.OutputValues.GetOutputValueByAddress(input.Address)
	if prevOutputValue != nil {
		prevValue = prevOutputValue.Value
	}

	app.logger.Infof("IsyApp.SwitchOnOff: Address %s. Previous value=%s, New value=%v", input.Address, prevValue, newValue)

	// input.UpdateValue(onOffString)
	node := app.pub.Nodes.GetNodeByAddress(input.Address)
	err := app.isyAPI.WriteOnOff(node.NodeID, newValue)
	if err != nil {
		app.logger.Errorf("IsyApp.SwitchOnOff: Input %s: error writing ISY: %v", input.Address, err)
	}
	return err
}

// HandleInputCommand for handling input commands
// Currently very basic. Only switches are supported.
func (app *IsyApp) HandleInputCommand(input *types.InputDiscoveryMessage, s *types.SetInputMessage) {
	app.logger.Infof("IsyApp.InputHandler. Input for '%s'", input.Address)

	// payloadStr := string(payload[:])

	// for now only support on/off
	switch input.InputType {
	case types.InputTypeSwitch:
		//adapter.UpdateOutputValue()device.UpdateSensorCommand(sensor, payloadStr)
		_ = app.SwitchOnOff(input, s.Value)
	default:
		app.logger.Warningf("IsyApp.InputHandler. Input '%s' is Not a switch", input.Address)
	}
	// publish the result. give gateway time to update.
	// TODO: get push notification instead
	// Give gateway time to update.
	time.Sleep(300 * time.Millisecond)
	app.Poll(app.pub)
}
