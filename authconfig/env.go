package authconfig

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/gobuffalo/envy"
)

// ParseEnvironment loads a sibling `.env` file then looks through all environment
// variables to set global configuration.
func ParseEnvironment() error {
	envy.Load()
	azureEnv, _ := azure.EnvironmentFromName("AzurePublicCloud") // shouldn't fail
	authorizationServerURL = azureEnv.ActiveDirectoryEndpoint

	// AZURE_GROUP_NAME and `authconfig.GroupName()` are deprecated.
	// Use AZURE_BASE_GROUP_NAME and `authconfig.GenerateGroupName()` instead.
	groupName = envy.Get("AZURE_GROUP_NAME", "azure-go-samples")
	baseGroupName = envy.Get("AZURE_BASE_GROUP_NAME", groupName)

	locationDefault = envy.Get("AZURE_LOCATION_DEFAULT", "westus2")

	var err error
	useDeviceFlow, err = strconv.ParseBool(envy.Get("AZURE_USE_DEVICEFLOW", "0"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_USE_DEVICEFLOW, disabling")
		useDeviceFlow = false
	}
	keepResources, err = strconv.ParseBool(envy.Get("AZURE_SAMPLES_KEEP_RESOURCES", "0"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_SAMPLES_KEEP_RESOURCES, discarding")
		keepResources = false
	}

	// these must be provided by environment
	// clientID
	clientID, err = envy.MustGet("AZURE_CLIENT_ID")
	if err != nil {
		return fmt.Errorf("expected env vars not provided: %s", err)
	}

	// clientSecret
	clientSecret, err = envy.MustGet("AZURE_CLIENT_SECRET")
	if err != nil && useDeviceFlow != true { // don't need a secret for device flow
		return fmt.Errorf("expected env vars not provided: %s", err)
	}

	// tenantID (AAD)
	tenantID, err = envy.MustGet("AZURE_TENANT_ID")
	if err != nil {
		return fmt.Errorf("expected env vars not provided: %s", err)
	}

	// subscriptionID (ARM)
	subscriptionID, err = envy.MustGet("AZURE_SUBSCRIPTION_ID")
	if err != nil {
		return fmt.Errorf("expected env vars not provided: %s", err)
	}

	return nil
}
