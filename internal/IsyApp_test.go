package internal

import (
	"testing"
	"time"

	"github.com/iotdomain/iotdomain-go/messaging"
	"github.com/iotdomain/iotdomain-go/publisher"
	"github.com/iotdomain/iotdomain-go/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfigFolder = "../test"
const testCacheFolder = "../test/cache"
const testData = "isy99-testdata.xml"
const deckLightsID = "15 2D A 1"

var messengerConfig = &messaging.MessengerConfig{Domain: "test"}
var appConfig = &IsyAppConfig{}

func TestLoadConfig(t *testing.T) {
	_, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)
	assert.Equal(t, "isy99", appConfig.PublisherID)
	assert.Equal(t, "gateway", appConfig.GatewayID)
}

// Read ISY device and check if more than 1 node is returned. A minimum of 1 is expected if the device is online with
// an additional node for each connected node.
func TestReadIsy(t *testing.T) {
	_, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	isyAPI := NewIsyAPI(appConfig.GatewayAddress, appConfig.LoginName, appConfig.Password)
	// use a simulation file
	isyAPI.address = "file://../test/gateway-config.xml"
	isyDevice, err := isyAPI.ReadIsyGateway()
	assert.NoError(t, err)
	assert.NotEmptyf(t, isyDevice.configuration.AppVersion, "Expected an application version")

	// use a simulation file
	isyAPI.address = "file://../test/gateway-nodes.xml"
	isyNodes, err := isyAPI.ReadIsyNodes()
	if assert.NoError(t, err) {
		assert.True(t, len(isyNodes.Nodes) > 5, "Expected 5 ISY nodes. Got fewer.")
	}
}

func TestPollOnce(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	app := NewIsyApp(appConfig, pub)
	pub.Start()
	assert.NoError(t, err)
	app.Poll(pub)
	time.Sleep(3 * time.Second)
	pub.Stop()
}

// This simulates the switch
func TestSwitch(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	app := NewIsyApp(appConfig, pub)
	app.SetupGatewayNode(pub)

	// FIXME: load isy nodes from file
	pub.Start()
	assert.NoError(t, err)
	app.Poll(pub)
	// some time to publish stuff
	time.Sleep(2 * time.Second)

	// throw a switch
	deckSwitch := pub.GetNodeByID(deckLightsID)
	require.NotNilf(t, deckSwitch, "Switch %s not found", deckLightsID)

	switchInput := pub.GetInput(deckSwitch.NodeID, types.InputTypeSwitch, types.DefaultInputInstance)
	require.NotNil(t, switchInput, "Input of switch node not found on address %s", deckSwitch.Address)
	// switchInput := deckSwitch.GetInput(types.InputTypeSwitch)

	logrus.Infof("TestSwitch: --- Switching deck switch %s OFF", deckSwitch.Address)

	pub.PublishSetInput(switchInput.Address, "false")
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	// fetch result
	switchOutput := pub.GetOutput(deckLightsID, types.OutputTypeSwitch, types.DefaultOutputInstance)
	// switchOutput := deckSwitch.GetOutput(types.InputTypeSwitch)
	require.NotNilf(t, switchOutput, "Output switch of node %s not found", deckLightsID)

	outputValue := pub.GetOutputValue(switchOutput.NodeID, types.OutputTypeSwitch, types.DefaultOutputInstance)
	assert.Equal(t, "false", outputValue.Value)

	logrus.Infof("TestSwitch: --- Switching deck switch %s ON", deckSwitch.Address)
	require.NotNil(t, switchInput, "Input switch of node %s not found", deckSwitch.NodeID)
	pub.PublishSetInput(switchInput.Address, "true")

	time.Sleep(3 * time.Second)
	outputValue = pub.GetOutputValue(switchOutput.NodeID, types.OutputTypeSwitch, types.DefaultOutputInstance)
	assert.Equal(t, "true", outputValue.Value)

	// be nice and turn the light back off
	pub.PublishSetInput(switchInput.Address, "false")

	pub.Stop()
}

func TestStartStop(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	// app := NewIsyApp(appConfig, pub)

	pub.Start()
	time.Sleep(time.Second * 10)
	pub.Stop()
}
