package sonarr

type SonarrApiIndexerPost struct {
	EnableRss               bool   `json:"enableRss"`
	EnableAutomaticSearch   bool   `json:"enableAutomaticSearch"`
	EnableInteractiveSearch bool   `json:"enableInteractiveSearch"`
	SupportsRss             bool   `json:"supportsRss"`
	SupportsSearch          bool   `json:"supportsSearch"`
	Protocol                string `json:"protocol"`
	Priority                int    `json:"priority"`
	DownloadClientId        int    `json:"downloadClientId"`
	Name                    string `json:"name"`
	Fields                  []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value,omitempty"`
	} `json:"fields"`
	ImplementationName string `json:"implementationName"`
	Implementation     string `json:"implementation"`
	ConfigContract     string `json:"configContract"`
	InfoLink           string `json:"infoLink"`
	Tags               []int  `json:"tags"`
}
