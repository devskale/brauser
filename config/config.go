package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSConfig represents the JavaScript compatibility configuration
type JSConfig struct {
	JavaScriptCompatibility struct {
		Enabled                bool `json:"enabled"`
		TimeoutSeconds         int  `json:"timeout_seconds"`
		MaxExecutionTimeSeconds int  `json:"max_execution_time_seconds"`
		Categories             struct {
			Console struct {
				Enabled bool     `json:"enabled"`
				Methods []string `json:"methods"`
			} `json:"console"`
			DOM struct {
				Enabled bool `json:"enabled"`
				Methods struct {
					Document   []string `json:"document"`
					Element    []string `json:"element"`
					ClassList  []string `json:"classList"`
					ParentNode []string `json:"parentNode"`
				} `json:"methods"`
			} `json:"dom"`
			Browser struct {
				Enabled bool `json:"enabled"`
				Window  struct {
					Location struct {
						Protocol string `json:"protocol"`
						Host     string `json:"host"`
						Pathname string `json:"pathname"`
					} `json:"location"`
					Methods []string `json:"methods"`
				} `json:"window"`
				Navigator struct {
					UserAgent string `json:"userAgent"`
				} `json:"navigator"`
			} `json:"browser"`
			Storage struct {
				Enabled        bool `json:"enabled"`
				LocalStorage   struct {
						Methods []string `json:"methods"`
					} `json:"localStorage"`
				SessionStorage struct {
					Methods []string `json:"methods"`
				} `json:"sessionStorage"`
			} `json:"storage"`
			WebAPI struct {
				Enabled    bool `json:"enabled"`
				MatchMedia struct {
					Enabled    bool     `json:"enabled"`
					Properties []string `json:"properties"`
					Methods    []string `json:"methods"`
				} `json:"matchMedia"`
				CustomEvent struct {
					Enabled    bool     `json:"enabled"`
					Properties []string `json:"properties"`
				} `json:"CustomEvent"`
				URLSearchParams struct {
					Enabled bool     `json:"enabled"`
					Methods []string `json:"methods"`
				} `json:"URLSearchParams"`
			} `json:"webapi"`
			Frameworks struct {
				Enabled bool `json:"enabled"`
				JQuery  struct {
					Enabled    bool     `json:"enabled"`
					Methods    []string `json:"methods"`
					Properties []string `json:"properties"`
				} `json:"jquery"`
			} `json:"frameworks"`
			SiteSpecific struct {
				Enabled bool `json:"enabled"`
				Globals map[string]struct {
					Enabled     bool   `json:"enabled"`
					Description string `json:"description"`
				} `json:"globals"`
			} `json:"site_specific"`
		} `json:"categories"`
	} `json:"javascript_compatibility"`
}

// LoadJSConfig loads JavaScript configuration from file
func LoadJSConfig(configPath string) (*JSConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	var config JSConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}
	
	return &config, nil
}

// LoadDefaultJSConfig returns a default configuration if file loading fails
func LoadDefaultJSConfig() *JSConfig {
	return &JSConfig{
		JavaScriptCompatibility: struct {
			Enabled                bool `json:"enabled"`
			TimeoutSeconds         int  `json:"timeout_seconds"`
			MaxExecutionTimeSeconds int  `json:"max_execution_time_seconds"`
			Categories             struct {
				Console struct {
					Enabled bool     `json:"enabled"`
					Methods []string `json:"methods"`
				} `json:"console"`
				DOM struct {
					Enabled bool `json:"enabled"`
					Methods struct {
						Document   []string `json:"document"`
						Element    []string `json:"element"`
						ClassList  []string `json:"classList"`
						ParentNode []string `json:"parentNode"`
					} `json:"methods"`
				} `json:"dom"`
				Browser struct {
					Enabled bool `json:"enabled"`
					Window  struct {
						Location struct {
							Protocol string `json:"protocol"`
							Host     string `json:"host"`
							Pathname string `json:"pathname"`
						} `json:"location"`
						Methods []string `json:"methods"`
					} `json:"window"`
					Navigator struct {
						UserAgent string `json:"userAgent"`
					} `json:"navigator"`
				} `json:"browser"`
				Storage struct {
					Enabled        bool `json:"enabled"`
					LocalStorage   struct {
							Methods []string `json:"methods"`
						} `json:"localStorage"`
					SessionStorage struct {
							Methods []string `json:"methods"`
						} `json:"sessionStorage"`
					} `json:"storage"`
				WebAPI struct {
					Enabled    bool `json:"enabled"`
					MatchMedia struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
						Methods    []string `json:"methods"`
					} `json:"matchMedia"`
					CustomEvent struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
					} `json:"CustomEvent"`
					URLSearchParams struct {
						Enabled bool     `json:"enabled"`
						Methods []string `json:"methods"`
					} `json:"URLSearchParams"`
				} `json:"webapi"`
				Frameworks struct {
					Enabled bool `json:"enabled"`
					JQuery  struct {
						Enabled    bool     `json:"enabled"`
						Methods    []string `json:"methods"`
						Properties []string `json:"properties"`
					} `json:"jquery"`
				} `json:"frameworks"`
				SiteSpecific struct {
					Enabled bool `json:"enabled"`
					Globals map[string]struct {
						Enabled     bool   `json:"enabled"`
						Description string `json:"description"`
					} `json:"globals"`
				} `json:"site_specific"`
			} `json:"categories"`
		}{
			Enabled:                true,
			TimeoutSeconds:         2,
			MaxExecutionTimeSeconds: 3,
			Categories: struct {
				Console struct {
					Enabled bool     `json:"enabled"`
					Methods []string `json:"methods"`
				} `json:"console"`
				DOM struct {
					Enabled bool `json:"enabled"`
					Methods struct {
						Document   []string `json:"document"`
						Element    []string `json:"element"`
						ClassList  []string `json:"classList"`
						ParentNode []string `json:"parentNode"`
					} `json:"methods"`
				} `json:"dom"`
				Browser struct {
					Enabled bool `json:"enabled"`
					Window  struct {
						Location struct {
							Protocol string `json:"protocol"`
							Host     string `json:"host"`
							Pathname string `json:"pathname"`
						} `json:"location"`
						Methods []string `json:"methods"`
					} `json:"window"`
					Navigator struct {
						UserAgent string `json:"userAgent"`
					} `json:"navigator"`
				} `json:"browser"`
				Storage struct {
					Enabled        bool `json:"enabled"`
					LocalStorage   struct {
							Methods []string `json:"methods"`
						} `json:"localStorage"`
					SessionStorage struct {
							Methods []string `json:"methods"`
						} `json:"sessionStorage"`
					} `json:"storage"`
				WebAPI struct {
					Enabled    bool `json:"enabled"`
					MatchMedia struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
						Methods    []string `json:"methods"`
					} `json:"matchMedia"`
					CustomEvent struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
					} `json:"CustomEvent"`
					URLSearchParams struct {
						Enabled bool     `json:"enabled"`
						Methods []string `json:"methods"`
					} `json:"URLSearchParams"`
				} `json:"webapi"`
				Frameworks struct {
					Enabled bool `json:"enabled"`
					JQuery  struct {
						Enabled    bool     `json:"enabled"`
						Methods    []string `json:"methods"`
						Properties []string `json:"properties"`
					} `json:"jquery"`
				} `json:"frameworks"`
				SiteSpecific struct {
					Enabled bool `json:"enabled"`
					Globals map[string]struct {
						Enabled     bool   `json:"enabled"`
						Description string `json:"description"`
					} `json:"globals"`
				} `json:"site_specific"`
			}{
				Console: struct {
					Enabled bool     `json:"enabled"`
					Methods []string `json:"methods"`
				}{Enabled: true, Methods: []string{"log"}},
				DOM: struct {
					Enabled bool `json:"enabled"`
					Methods struct {
						Document   []string `json:"document"`
						Element    []string `json:"element"`
						ClassList  []string `json:"classList"`
						ParentNode []string `json:"parentNode"`
					} `json:"methods"`
				}{Enabled: true},
				Browser:      struct {
					Enabled bool `json:"enabled"`
					Window  struct {
						Location struct {
							Protocol string `json:"protocol"`
							Host     string `json:"host"`
							Pathname string `json:"pathname"`
						} `json:"location"`
						Methods []string `json:"methods"`
					} `json:"window"`
					Navigator struct {
						UserAgent string `json:"userAgent"`
					} `json:"navigator"`
				}{Enabled: true},
				Storage:      struct {
					Enabled        bool `json:"enabled"`
					LocalStorage   struct {
							Methods []string `json:"methods"`
						} `json:"localStorage"`
					SessionStorage struct {
						Methods []string `json:"methods"`
					} `json:"sessionStorage"`
				}{Enabled: true},
				WebAPI:       struct {
					Enabled    bool `json:"enabled"`
					MatchMedia struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
						Methods    []string `json:"methods"`
					} `json:"matchMedia"`
					CustomEvent struct {
						Enabled    bool     `json:"enabled"`
						Properties []string `json:"properties"`
					} `json:"CustomEvent"`
					URLSearchParams struct {
						Enabled bool     `json:"enabled"`
						Methods []string `json:"methods"`
					} `json:"URLSearchParams"`
				}{Enabled: true},
				Frameworks:   struct {
					Enabled bool `json:"enabled"`
					JQuery  struct {
						Enabled    bool     `json:"enabled"`
						Methods    []string `json:"methods"`
						Properties []string `json:"properties"`
					} `json:"jquery"`
				}{Enabled: true},
				SiteSpecific: struct {
					Enabled bool `json:"enabled"`
					Globals map[string]struct {
						Enabled     bool   `json:"enabled"`
						Description string `json:"description"`
					} `json:"globals"`
				}{Enabled: true},
			},
		},
	}
}