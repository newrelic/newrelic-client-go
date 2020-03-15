package dashboards

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_dashboard() {
	// Initialize the client configuration.  An Admin API key is required to
	// communicate with the backend API.
	cfg := config.Config{
		AdminAPIKey: os.Getenv("NEW_RELIC_ADMIN_API_KEY"),
	}

	// Initialize the client.
	client := New(cfg)

	// Search the dashboards for the current account by title.
	listParams := &ListDashboardsParams{
		Title: "Example dashboard",
	}

	dashboards, err := client.ListDashboards(listParams)
	if err != nil {
		log.Fatal("error listing dashboards:", err)
	}

	// Get dashboard by ID.
	dashboard, err := client.GetDashboard(dashboards[0].ID)
	if err != nil {
		log.Fatal("error getting dashboard:", err)
	}

	// Create a new dashboard.
	applicationName := "Example application"

	dashboard = &Dashboard{
		Title: "Example dashboard",
		Icon:  DashboardIconTypes.BarChart,
	}

	requestsPerMinute := DashboardWidget{
		Visualization: VisualizationTypes.Billboard,
		Data: []DashboardWidgetData{
			{
				NRQL: fmt.Sprintf("FROM Transaction SELECT rate(count(*), 1 minute) WHERE appName = '%s'", applicationName),
			},
		},
		Presentation: DashboardWidgetPresentation{
			Title: "Requests per minute",
		},
		Layout: DashboardWidgetLayout{
			Row:    1,
			Column: 1,
		},
	}

	errorRate := DashboardWidget{
		Visualization: VisualizationTypes.Gauge,
		Data: []DashboardWidgetData{
			{
				NRQL: fmt.Sprintf("FROM Transaction SELECT percentage(count(*), WHERE error IS true) WHERE appName = '%s'", applicationName),
			},
		},
		Presentation: DashboardWidgetPresentation{
			Title: "Error rate",
			Threshold: &DashboardWidgetThreshold{
				Red: 2.5,
			},
		},
		Layout: DashboardWidgetLayout{
			Row:    1,
			Column: 2,
		},
	}

	notes := DashboardWidget{
		Visualization: VisualizationTypes.Markdown,
		Data: []DashboardWidgetData{
			{
				Source: "### Helpful Links\n\n* [New Relic One](https://one.newrelic.com)\n* [Developer Portal](https://developer.newrelic.com)",
			},
		},
		Presentation: DashboardWidgetPresentation{
			Title: "Dashboard note",
		},
		Layout: DashboardWidgetLayout{
			Row:    1,
			Column: 3,
		},
	}

	dashboard.Widgets = []DashboardWidget{
		requestsPerMinute,
		errorRate,
		notes,
	}

	created, err := client.CreateDashboard(*dashboard)
	if err != nil {
		log.Fatal("error creating dashboard:", err)
	}

	// Add a widget to an existing dashboard.
	topApdex := DashboardWidget{
		Visualization: VisualizationTypes.FacetTable,
		Data: []DashboardWidgetData{
			{
				NRQL: fmt.Sprintf("FROM Transaction SELECT rate(count(*), 1 minute) FACET name WHERE appName = '%s'", applicationName),
			},
		},
		Presentation: DashboardWidgetPresentation{
			Title: "Requests per minute, by transaction",
		},
		Layout: DashboardWidgetLayout{
			Row:    1,
			Column: 2,
			Width:  3,
		},
	}

	created.Widgets = append(created.Widgets, topApdex)

	updated, err := client.UpdateDashboard(*created)
	if err != nil {
		log.Fatal("error updating dashboard:", err)
	}

	// Delete a dashaboard.
	_, err = client.DeleteDashboard(updated.ID)
	if err != nil {
		log.Fatal("error deleting dashboard:", err)
	}
}
