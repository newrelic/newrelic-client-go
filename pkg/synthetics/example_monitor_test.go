package synthetics

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func Example_monitor() {
	// Initialize the client configuration. A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Create a simple Synthetics monitor.
	simpleMonitor := Monitor{
		Name:         "Example monitor",
		Type:         MonitorTypes.Ping,
		Status:       MonitorStatus.Enabled,
		SLAThreshold: 2.0,
		URI:          "https://www.example.com",
		Frequency:    5,
		Locations:    []string{"AWS_US_EAST_1"},
	}

	created, err := client.CreateMonitor(simpleMonitor)
	if err != nil {
		log.Fatal("error creating Synthetics monitor: ", err)
	}

	// Create a scripted browser monitor.
	scriptedBrowser := Monitor{
		Name:         "Example scriptied browser monitor",
		Type:         MonitorTypes.ScriptedBrowser,
		Status:       MonitorStatus.Enabled,
		SLAThreshold: 2.0,
		URI:          "https://www.example.com",
		Frequency:    1440,
		Locations:    []string{"AWS_US_EAST_1", "AWS_US_WEST_1"},
	}

	monitorScript := MonitorScript{
		Text: `
var assert = require("assert");
$browser.get("http://www.example.com").then(function(){ 

	// Check the H1 title matches "Example Domain" 
	return $browser.findElement($driver.By.css("h1")).then(function(element){ 
		return element.getText().then(function(text){ 
		assert.equal("Example Domain", text, "Page H1 title did not match"); 
		}); 
	}); 
})`,
	}

	created, err = client.CreateMonitor(scriptedBrowser)
	if err != nil {
		log.Fatal("error creating Synthetics monitor: ", err)
	}

	_, err = client.UpdateMonitorScript(created.ID, monitorScript)
	if err != nil {
		log.Fatal("error updating Synthetics monitor script: ", err)
	}

	// Update an existing Synthetics monitor script.
	created.Locations = []string{"AWS_US_WEST_1"}

	updated, err := client.UpdateMonitor(*created)
	if err != nil {
		log.Fatal("error updating Synthetics monitor: ", err)
	}

	// Delete an existing Synthetics monitor script.
	err = client.DeleteMonitor(updated.ID)
	if err != nil {
		log.Fatal("error deleting Synthetics monitor: ", err)
	}

	// Get all valid Synthetics monitor locations.
	locations, err := client.GetMonitorLocations()
	if err != nil {
		log.Fatal("error retrieving valid Synthetics monitor locations: ", err)
	}

	log.Printf("found %d valid Synthetics monitor locations", len(locations))
}
