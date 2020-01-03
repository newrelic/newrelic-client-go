package plugins

// Plugin represents information about a New Relic plugin.
type Plugin struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name,omitempty"`
	GUID                string          `json:"guid,omitempty"`
	Publisher           string          `json:"publisher,omitempty"`
	ComponentAgentCount int             `json:"component_agent_count"`
	Details             PluginDetails   `json:"details"`
	SummaryMetrics      []SummaryMetric `json:"summary_metrics"`
}

// PluginDetails represents information about a New Relic plugin.
type PluginDetails struct {
	Description           int    `json:"description"`
	IsPublic              string `json:"is_public"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
	LastPublishedAt       string `json:"last_published_at,omitempty"`
	HasUnpublishedChanges bool   `json:"has_unpublished_changes"`
	BrandingImageURL      string `json:"branding_image_url"`
	UpgradedAt            string `json:"upgraded_at,omitempty"`
	ShortName             string `json:"short_name"`
	PublisherAboutURL     string `json:"publisher_about_url"`
	PublisherSupportURL   string `json:"publisher_support_url"`
	DownloadURL           string `json:"download_url"`
	FirstEditedAt         string `json:"first_edited_at,omitempty"`
	LastEditedAt          string `json:"last_edited_at,omitempty"`
	FirstPublishedAt      string `json:"first_published_at,omitempty"`
	PublishedVersion      string `json:"published_version"`
}

// SummaryMetric represents summary information for a specific metric.
type SummaryMetric struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Metric        string          `json:"metric"`
	ValueFunction string          `json:"value_function"`
	Thresholds    MetricThreshold `json:"thresholds"`
	Values        MetricValue     `json:"values"`
}

// MetricThreshold represents the different thresholds for a metric in an alert.
type MetricThreshold struct {
	Caution  float64 `json:"caution"`
	Critical float64 `json:"critical"`
}

// MetricValue represents the observed value of a metric.
type MetricValue struct {
	Raw       float64 `json:"raw"`
	Formatted string  `json:"formatted"`
}
