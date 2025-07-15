package main
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	aic "github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
)


// main is the entry point of the Brauser application.
func main() {
	fmt.Println("Brauser: Minimalistic Terminal Web Browser")
	if len(os.Args) < 2 {
		fmt.Println("Usage: brauser <url>")
		return
	}
	url := os.Args[1]
	content, err := fetchPage(url)
	if err != nil {
		log.Fatalf("Failed to fetch page: %v", err)
	}
	renderHTML(content, url)
}

// fetchPage fetches the content of the given URL and returns it as a string.
func fetchPage(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// renderHTML parses the HTML content, displays it in a structured format, and executes embedded JS.
func renderHTML(htmlContent, baseURL string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Extract and print title
	title := doc.Find("title").Text()
	fmt.Println("Title:", title)

	// Extract and print headings
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%s: %s\n", s.Get(0).Data, s.Text())
	})

	// Extract and print paragraphs
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	// Extract and print links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			fmt.Printf("Link: %s (%s)\n", s.Text(), href)
		}
	})

	// For HN-specific: Extract story items
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".titleline > a").Text()
		link := s.Find(".titleline > a").AttrOr("href", "")
		fmt.Printf("Story: %s (%s)\n", title, link)
	})

	// Render images as ASCII art
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			asciiArt, err := renderImageAsASCII(src, baseURL)
			if err != nil {
				log.Printf("Failed to render image %s: %v", src, err)
				return
			}
			fmt.Println("ASCII Image:")
			fmt.Println(asciiArt)
		}
	})

	// Execute embedded JavaScript
	executeJS(doc)

	// TODO: Add more elements like lists, etc.
	// TODO: Implement DOM manipulation for dynamic content updates
}

// renderImageAsASCII fetches an image from the given src (handling relative URLs) and converts it to ASCII art.
func renderImageAsASCII(src, baseURL string) (string, error) {
	// Handle relative URLs
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	if !u.IsAbs() {
		base, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}
		src = base.ResolveReference(u).String()
	}

	resp, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "brauser-img-*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(imgData); err != nil {
		return "", err
	}
	if err := tempFile.Close(); err != nil {
		return "", err
	}

	flags := aic.DefaultFlags()
	flags.Dimensions = []int{80, 40} // Set desired width and height
	flags.Colored = true

	asciiArt, err := aic.Convert(tempFile.Name(), flags)
	if err != nil {
		return "", err
	}

	return asciiArt, nil
}

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

// JSErrorType represents different categories of JavaScript errors
type JSErrorType int

const (
	SyntaxError JSErrorType = iota
	RuntimeError
	CompatibilityError
	TimeoutError
	HarmlessError
)

// JSError represents a categorized JavaScript error
type JSError struct {
	Type        JSErrorType
	Message     string
	Script      string
	Suppressed  bool
	OriginalErr error
}

// ErrorHandler manages JavaScript error categorization and suppression
type ErrorHandler struct {
	SuppressionPatterns []*regexp.Regexp
	HarmlessPatterns    []*regexp.Regexp
	CompatibilityStubs  map[string]bool
}

// NewErrorHandler creates a new error handler with default patterns
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		SuppressionPatterns: []*regexp.Regexp{
			// Common harmless errors
			regexp.MustCompile(`SyntaxError: Unexpected token :`),
			regexp.MustCompile(`ReferenceError: .* is not defined`),
			regexp.MustCompile(`TypeError: Cannot read property .* of undefined`),
			regexp.MustCompile(`TypeError: Cannot read property .* of null`),
			regexp.MustCompile(`TypeError: Object has no member`),
			regexp.MustCompile(`TypeError: .* is not a function`),
		},
		HarmlessPatterns: []*regexp.Regexp{
			// Analytics and tracking errors
			regexp.MustCompile(`gtag|ga\(|_gaq|dataLayer`),
			regexp.MustCompile(`fbq|_fbq|facebook`),
			regexp.MustCompile(`twitter|twttr`),
			regexp.MustCompile(`google-analytics|googletagmanager`),
			// Ad-related errors
			regexp.MustCompile(`adsystem|doubleclick|googlesyndication`),
			// Social media widgets
			regexp.MustCompile(`instagram|linkedin|pinterest`),
		},
		CompatibilityStubs: make(map[string]bool),
	}
}

// CategorizeError determines the type of JavaScript error
func (eh *ErrorHandler) CategorizeError(err error, script string) *JSError {
	errorMsg := err.Error()
	
	jsErr := &JSError{
		Message:     errorMsg,
		Script:      script,
		OriginalErr: err,
		Suppressed:  false,
	}
	
	// Check for timeout errors
	if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "interrupted") {
		jsErr.Type = TimeoutError
		return jsErr
	}
	
	// Check for syntax errors
	if strings.Contains(errorMsg, "SyntaxError") || strings.Contains(errorMsg, "Unexpected token") {
		jsErr.Type = SyntaxError
		// Check if it's a harmless syntax error
		for _, pattern := range eh.SuppressionPatterns {
			if pattern.MatchString(errorMsg) {
				jsErr.Suppressed = true
				break
			}
		}
		return jsErr
	}
	
	// Check for compatibility errors
	if strings.Contains(errorMsg, "ReferenceError") || strings.Contains(errorMsg, "is not defined") {
		jsErr.Type = CompatibilityError
		// Check if it's a known missing API
		for _, pattern := range eh.SuppressionPatterns {
			if pattern.MatchString(errorMsg) {
				jsErr.Suppressed = true
				break
			}
		}
		return jsErr
	}
	
	// Check for harmless errors (analytics, ads, etc.)
	for _, pattern := range eh.HarmlessPatterns {
		if pattern.MatchString(script) || pattern.MatchString(errorMsg) {
			jsErr.Type = HarmlessError
			jsErr.Suppressed = true
			return jsErr
		}
	}
	
	// Default to runtime error
	jsErr.Type = RuntimeError
	return jsErr
}

// ShouldSuppress determines if an error should be suppressed from logging
func (eh *ErrorHandler) ShouldSuppress(jsErr *JSError) bool {
	return jsErr.Suppressed || jsErr.Type == HarmlessError
}

// LogError logs JavaScript errors with appropriate severity
func (eh *ErrorHandler) LogError(jsErr *JSError) {
	if eh.ShouldSuppress(jsErr) {
		// Only log suppressed errors in debug mode
		log.Printf("[DEBUG] Suppressed JS %s: %s", eh.getErrorTypeName(jsErr.Type), jsErr.Message)
		return
	}
	
	switch jsErr.Type {
	case SyntaxError:
		log.Printf("[WARN] JS Syntax Error: %s", jsErr.Message)
	case RuntimeError:
		log.Printf("[ERROR] JS Runtime Error: %s", jsErr.Message)
	case CompatibilityError:
		log.Printf("[INFO] JS Compatibility Issue: %s", jsErr.Message)
	case TimeoutError:
		log.Printf("[WARN] JS Timeout: %s", jsErr.Message)
	default:
		log.Printf("[ERROR] JS Error: %s", jsErr.Message)
	}
}

// getErrorTypeName returns a human-readable error type name
func (eh *ErrorHandler) getErrorTypeName(errorType JSErrorType) string {
	switch errorType {
	case SyntaxError:
		return "Syntax Error"
	case RuntimeError:
		return "Runtime Error"
	case CompatibilityError:
		return "Compatibility Error"
	case TimeoutError:
		return "Timeout Error"
	case HarmlessError:
		return "Harmless Error"
	default:
		return "Unknown Error"
	}
}

// handleCompatibilityError attempts graceful degradation for compatibility errors
func (env *JSEnvironment) handleCompatibilityError(jsErr *JSError) {
	// Extract API name from error message
	if strings.Contains(jsErr.Message, "is not defined") {
		// Try to extract the undefined variable name
		re := regexp.MustCompile(`ReferenceError: (\w+) is not defined`)
		matches := re.FindStringSubmatch(jsErr.Message)
		if len(matches) > 1 {
			apiName := matches[1]
			env.errorHandler.AddMissingAPI(env, apiName)
		}
	}
}

// AddMissingAPI dynamically adds a missing API stub for graceful degradation
func (eh *ErrorHandler) AddMissingAPI(env *JSEnvironment, apiName string) {
	if eh.CompatibilityStubs[apiName] {
		return // Already added
	}
	
	log.Printf("[INFO] Adding graceful degradation stub for missing API: %s", apiName)
	
	// Add common missing APIs
	switch apiName {
	case "dispatchEvent":
		env.runtime.Set("dispatchEvent", func(call goja.FunctionCall) goja.Value {
			log.Println("Stub: dispatchEvent called")
			return goja.Undefined()
		})
	case "requestAnimationFrame":
		env.runtime.Set("requestAnimationFrame", func(call goja.FunctionCall) goja.Value {
			log.Println("Stub: requestAnimationFrame called")
			return env.runtime.ToValue(1) // Return dummy ID
		})
	case "cancelAnimationFrame":
		env.runtime.Set("cancelAnimationFrame", func(call goja.FunctionCall) goja.Value {
			log.Println("Stub: cancelAnimationFrame called")
			return goja.Undefined()
		})
	case "fetch":
		env.runtime.Set("fetch", func(call goja.FunctionCall) goja.Value {
			log.Println("Stub: fetch called")
			// Return a rejected promise stub
			promise := env.runtime.NewObject()
			promise.Set("then", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(promise) })
			promise.Set("catch", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(promise) })
			return env.runtime.ToValue(promise)
		})
	case "IntersectionObserver":
		env.runtime.Set("IntersectionObserver", func(call goja.ConstructorCall) *goja.Object {
			log.Println("Stub: IntersectionObserver constructor called")
			observer := env.runtime.NewObject()
			observer.Set("observe", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
			observer.Set("unobserve", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
			observer.Set("disconnect", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
			return observer
		})
	case "MutationObserver":
		env.runtime.Set("MutationObserver", func(call goja.ConstructorCall) *goja.Object {
			log.Println("Stub: MutationObserver constructor called")
			observer := env.runtime.NewObject()
			observer.Set("observe", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
			observer.Set("disconnect", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
			return observer
		})
	default:
		// Generic object stub for unknown APIs
		env.runtime.Set(apiName, env.runtime.NewObject())
	}
	
	eh.CompatibilityStubs[apiName] = true
}

// JSEnvironment manages the JavaScript runtime and stubs
type JSEnvironment struct {
	runtime      *goja.Runtime
	config       *JSConfig
	window       *goja.Object
	document     *goja.Object
	errorHandler *ErrorHandler
}

// loadJSConfig loads JavaScript configuration from file
func loadJSConfig(configPath string) (*JSConfig, error) {
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

// loadDefaultJSConfig returns a default configuration if file loading fails
func loadDefaultJSConfig() *JSConfig {
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

// NewJSEnvironment creates a new JavaScript environment with all stubs
func NewJSEnvironment() *JSEnvironment {
	config, err := loadJSConfig("js_config.json")
	if err != nil {
		log.Printf("Failed to load JS config, using defaults: %v", err)
		config = loadDefaultJSConfig()
	}
	
	env := &JSEnvironment{
		runtime:      goja.New(),
		config:       config,
		errorHandler: NewErrorHandler(),
	}
	
	if config.JavaScriptCompatibility.Enabled {
		env.setupAllStubs()
	}
	
	return env
}

// setupAllStubs initializes all JavaScript stubs based on configuration
func (env *JSEnvironment) setupAllStubs() {
	if env.config.JavaScriptCompatibility.Categories.Console.Enabled {
		env.setupConsoleStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.DOM.Enabled {
		env.setupDOMStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.Browser.Enabled {
		env.setupBrowserStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.Storage.Enabled {
		env.setupStorageStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.WebAPI.Enabled {
		env.setupWebAPIStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.Frameworks.Enabled {
		env.setupFrameworkStubs()
	}
	if env.config.JavaScriptCompatibility.Categories.SiteSpecific.Enabled {
		env.setupSiteSpecificStubs()
	}
}

// setupConsoleStubs creates console object with logging methods
func (env *JSEnvironment) setupConsoleStubs() {
	console := env.runtime.NewObject()
	console.Set("log", func(call goja.FunctionCall) goja.Value {
		args := make([]interface{}, len(call.Arguments))
		for i, arg := range call.Arguments {
			args[i] = arg.String()
		}
		log.Println(args...)
		return goja.Undefined()
	})
	env.runtime.Set("console", console)
}

// createMockElement creates a comprehensive mock DOM element
// createMockElement creates a comprehensive DOM element stub with realistic property access patterns
func (env *JSEnvironment) createMockElement() *goja.Object {
	mockElem := env.runtime.NewObject()
	
	// Core element properties with dynamic getters/setters
	mockElem.Set("innerHTML", "")
	mockElem.Set("outerHTML", "<div></div>")
	mockElem.Set("textContent", "")
	mockElem.Set("innerText", "")
	mockElem.Set("className", "")
	mockElem.Set("id", "")
	mockElem.Set("tagName", "DIV")
	mockElem.Set("nodeName", "DIV")
	mockElem.Set("nodeType", 1) // ELEMENT_NODE
	mockElem.Set("nodeValue", goja.Null())
	
	// Style object with common CSS properties
	style := env.runtime.NewObject()
	commonStyleProps := []string{"display", "visibility", "opacity", "position", "top", "left", "width", "height", "margin", "padding", "border", "background", "color", "fontSize", "fontFamily", "zIndex"}
	for _, prop := range commonStyleProps {
		style.Set(prop, "")
	}
	style.Set("setProperty", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	style.Set("getPropertyValue", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue("") })
	style.Set("removeProperty", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("style", style)
	
	// Enhanced classList with comprehensive methods
	mockClassList := env.runtime.NewObject()
	mockClassList.Set("add", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockClassList.Set("remove", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockClassList.Set("toggle", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	mockClassList.Set("contains", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	mockClassList.Set("replace", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockClassList.Set("length", 0)
	mockElem.Set("classList", mockClassList)
	
	// Attributes handling
	mockElem.Set("setAttribute", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("getAttribute", func(call goja.FunctionCall) goja.Value { return goja.Null() })
	mockElem.Set("removeAttribute", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("hasAttribute", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	mockElem.Set("getAttributeNames", func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() })
	
	// Dataset for data-* attributes
	dataset := env.runtime.NewObject()
	mockElem.Set("dataset", dataset)
	
	// Enhanced parent node with comprehensive methods
	mockParent := env.createMockParentNode()
	mockElem.Set("parentNode", mockParent)
	mockElem.Set("parentElement", mockParent)
	
	// Child node management
	mockElem.Set("appendChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockElem.Set("insertBefore", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockElem.Set("removeChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockElem.Set("replaceChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 1 {
			return call.Arguments[1]
		}
		return goja.Undefined()
	})
	
	// Child node collections
	childNodes := env.runtime.NewArray()
	childNodes.Set("length", 0)
	mockElem.Set("childNodes", childNodes)
	mockElem.Set("children", env.runtime.NewArray())
	mockElem.Set("firstChild", goja.Null())
	mockElem.Set("lastChild", goja.Null())
	mockElem.Set("firstElementChild", goja.Null())
	mockElem.Set("lastElementChild", goja.Null())
	mockElem.Set("nextSibling", goja.Null())
	mockElem.Set("previousSibling", goja.Null())
	mockElem.Set("nextElementSibling", goja.Null())
	mockElem.Set("previousElementSibling", goja.Null())
	
	// Query methods
	mockElem.Set("querySelector", func(call goja.FunctionCall) goja.Value {
		return env.runtime.ToValue(env.createMockElement())
	})
	mockElem.Set("querySelectorAll", func(call goja.FunctionCall) goja.Value {
		return env.runtime.NewArray()
	})
	mockElem.Set("getElementsByTagName", func(call goja.FunctionCall) goja.Value {
		return env.runtime.NewArray()
	})
	mockElem.Set("getElementsByClassName", func(call goja.FunctionCall) goja.Value {
		return env.runtime.NewArray()
	})
	
	// Event handling with realistic patterns
	mockElem.Set("addEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("removeEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("dispatchEvent", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(true) })
	
	// Common event properties (for framework compatibility)
	mockElem.Set("onclick", goja.Null())
	mockElem.Set("onload", goja.Null())
	mockElem.Set("onerror", goja.Null())
	mockElem.Set("onchange", goja.Null())
	mockElem.Set("onsubmit", goja.Null())
	
	// Geometry and positioning
	mockElem.Set("offsetWidth", 100)
	mockElem.Set("offsetHeight", 100)
	mockElem.Set("offsetTop", 0)
	mockElem.Set("offsetLeft", 0)
	mockElem.Set("offsetParent", mockParent)
	mockElem.Set("clientWidth", 100)
	mockElem.Set("clientHeight", 100)
	mockElem.Set("scrollWidth", 100)
	mockElem.Set("scrollHeight", 100)
	mockElem.Set("scrollTop", 0)
	mockElem.Set("scrollLeft", 0)
	
	// Bounding box method
	mockElem.Set("getBoundingClientRect", func(call goja.FunctionCall) goja.Value {
		rect := env.runtime.NewObject()
		rect.Set("top", 0)
		rect.Set("left", 0)
		rect.Set("bottom", 100)
		rect.Set("right", 100)
		rect.Set("width", 100)
		rect.Set("height", 100)
		rect.Set("x", 0)
		rect.Set("y", 0)
		return env.runtime.ToValue(rect)
	})
	
	// Form-specific properties (for input elements)
	mockElem.Set("value", "")
	mockElem.Set("checked", false)
	mockElem.Set("disabled", false)
	mockElem.Set("readonly", false)
	mockElem.Set("selected", false)
	mockElem.Set("type", "text")
	mockElem.Set("name", "")
	mockElem.Set("form", goja.Null())
	
	// Focus methods
	mockElem.Set("focus", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("blur", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockElem.Set("click", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	
	// Content manipulation methods
	mockElem.Set("cloneNode", func(call goja.FunctionCall) goja.Value {
		return env.runtime.ToValue(env.createMockElement())
	})
	mockElem.Set("remove", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	
	// Framework-specific properties for better compatibility
	mockElem.Set("_reactInternalFiber", goja.Undefined()) // React
	mockElem.Set("__reactInternalInstance", goja.Undefined()) // React
	mockElem.Set("_vnode", goja.Undefined()) // Vue
	mockElem.Set("__vue__", goja.Undefined()) // Vue
	mockElem.Set("_ngcontent", goja.Undefined()) // Angular
	
	return mockElem
}

// createMockParentNode creates a realistic parent node with comprehensive methods
func (env *JSEnvironment) createMockParentNode() *goja.Object {
	mockParent := env.runtime.NewObject()
	
	// Parent node properties
	mockParent.Set("nodeName", "BODY")
	mockParent.Set("nodeType", 1)
	mockParent.Set("tagName", "BODY")
	
	// Child manipulation methods
	mockParent.Set("appendChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockParent.Set("insertBefore", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockParent.Set("removeChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return goja.Undefined()
	})
	mockParent.Set("replaceChild", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 1 {
			return call.Arguments[1]
		}
		return goja.Undefined()
	})
	
	// Query methods
	mockParent.Set("querySelector", func(call goja.FunctionCall) goja.Value {
		return env.runtime.ToValue(env.createMockElement())
	})
	mockParent.Set("querySelectorAll", func(call goja.FunctionCall) goja.Value {
		return env.runtime.NewArray()
	})
	
	// Event handling
	mockParent.Set("addEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	mockParent.Set("removeEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	
	return mockParent
}

// setupDOMStubs creates document object and DOM manipulation methods
func (env *JSEnvironment) setupDOMStubs() {
	env.document = env.runtime.NewObject()
	
	// getElementById stub
	env.document.Set("getElementById", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: getElementById called")
		return env.runtime.ToValue(env.createMockElement())
	})
	
	// createElement stub
	env.document.Set("createElement", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: createElement called")
		return env.runtime.ToValue(env.createMockElement())
	})
	
	// getElementsByTagName stub
	env.document.Set("getElementsByTagName", func(call goja.FunctionCall) goja.Value {
		mockArray := env.runtime.NewArray()
		mockArray.Set("0", env.createMockElement())
		mockArray.Set("length", 1)
		log.Println("Stub: getElementsByTagName called")
		return env.runtime.ToValue(mockArray)
	})
	
	// querySelector stubs
	env.document.Set("querySelector", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: querySelector called")
		return env.runtime.ToValue(env.createMockElement())
	})
	env.document.Set("querySelectorAll", func(call goja.FunctionCall) goja.Value {
		mockArray := env.runtime.NewArray()
		log.Println("Stub: querySelectorAll called")
		return env.runtime.ToValue(mockArray)
	})
	
	// Document event handling
	env.document.Set("addEventListener", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: document.addEventListener called")
		return goja.Undefined()
	})
	
	env.runtime.Set("document", env.document)
}

// setupBrowserStubs creates window, navigator, and location objects
func (env *JSEnvironment) setupBrowserStubs() {
	env.window = env.runtime.NewObject()
	
	// Location object
	location := env.runtime.NewObject()
	location.Set("protocol", "https:")
	location.Set("host", "localhost")
	location.Set("pathname", "/")
	env.window.Set("location", location)
	
	// Window event handling
	env.window.Set("addEventListener", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: window.addEventListener called")
		return goja.Undefined()
	})
	
	// TCF API stub
	env.window.Set("__tcfapiLocator", env.runtime.NewObject())
	
	// Navigator object
	navigator := env.runtime.NewObject()
	navigator.Set("userAgent", "Brauser/1.0")
	env.runtime.Set("navigator", navigator)
	
	// Storage APIs
	env.setupStorageAPIs()
	
	env.runtime.Set("window", env.window)
}

// setupStorageStubs creates storage-related stubs
func (env *JSEnvironment) setupStorageStubs() {
	env.setupStorageAPIs()
}

// setupStorageAPIs creates localStorage and sessionStorage
func (env *JSEnvironment) setupStorageAPIs() {
	// localStorage
	localStorage := env.runtime.NewObject()
	localStorage.Set("getItem", func(call goja.FunctionCall) goja.Value { return goja.Null() })
	localStorage.Set("setItem", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	localStorage.Set("removeItem", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	localStorage.Set("clear", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	env.runtime.Set("localStorage", localStorage)
	env.window.Set("localStorage", localStorage)
	
	// sessionStorage
	sessionStorage := env.runtime.NewObject()
	sessionStorage.Set("getItem", func(call goja.FunctionCall) goja.Value { return goja.Null() })
	sessionStorage.Set("setItem", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	sessionStorage.Set("removeItem", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	sessionStorage.Set("clear", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	env.runtime.Set("sessionStorage", sessionStorage)
	env.window.Set("sessionStorage", sessionStorage)
}

// setupWebAPIStubs creates modern web API stubs
func (env *JSEnvironment) setupWebAPIStubs() {
	// matchMedia API
	env.window.Set("matchMedia", func(call goja.FunctionCall) goja.Value {
		matchMediaResult := env.runtime.NewObject()
		matchMediaResult.Set("matches", false)
		if len(call.Arguments) > 0 {
			matchMediaResult.Set("media", call.Arguments[0].String())
		} else {
			matchMediaResult.Set("media", "")
		}
		matchMediaResult.Set("addListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		matchMediaResult.Set("removeListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		matchMediaResult.Set("addEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		matchMediaResult.Set("removeEventListener", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		return env.runtime.ToValue(matchMediaResult)
	})
	
	// CustomEvent constructor
	env.runtime.Set("CustomEvent", func(call goja.ConstructorCall) *goja.Object {
		event := env.runtime.NewObject()
		if len(call.Arguments) > 0 {
			event.Set("type", call.Arguments[0].String())
		}
		event.Set("detail", goja.Null())
		event.Set("bubbles", false)
		event.Set("cancelable", false)
		return event
	})
	
	// URLSearchParams constructor
	env.runtime.Set("URLSearchParams", func(call goja.ConstructorCall) *goja.Object {
		params := env.runtime.NewObject()
		params.Set("get", func(call goja.FunctionCall) goja.Value { return goja.Null() })
		params.Set("set", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		params.Set("has", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
		params.Set("append", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		params.Set("delete", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
		return params
	})
}

// setupFrameworkStubs creates framework-specific stubs with realistic property access patterns
func (env *JSEnvironment) setupFrameworkStubs() {
	// Enhanced jQuery stub with comprehensive chainable methods and realistic behavior
	jqueryStub := func(call goja.FunctionCall) goja.Value {
		jqObj := env.createJQueryObject()
		return env.runtime.ToValue(jqObj)
	}
	
	// Add jQuery static methods
	jqueryStatic := env.runtime.NewObject()
	jqueryStatic.Set("extend", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			return call.Arguments[0]
		}
		return env.runtime.NewObject()
	})
	jqueryStatic.Set("each", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	jqueryStatic.Set("map", func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() })
	jqueryStatic.Set("grep", func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() })
	jqueryStatic.Set("inArray", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(-1) })
	jqueryStatic.Set("isArray", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	jqueryStatic.Set("isFunction", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	jqueryStatic.Set("isPlainObject", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	jqueryStatic.Set("parseJSON", func(call goja.FunctionCall) goja.Value { return env.runtime.NewObject() })
	jqueryStatic.Set("ajax", func(call goja.FunctionCall) goja.Value {
		// Return a mock jqXHR object
		jqXHR := env.runtime.NewObject()
		jqXHR.Set("done", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(jqXHR) })
		jqXHR.Set("fail", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(jqXHR) })
		jqXHR.Set("always", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(jqXHR) })
		return env.runtime.ToValue(jqXHR)
	})
	jqueryStatic.Set("get", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(env.createJQueryObject()) })
	jqueryStatic.Set("post", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(env.createJQueryObject()) })
	jqueryStatic.Set("fn", env.runtime.NewObject()) // jQuery prototype
	
	// Copy static methods to the main jQuery function
	for _, prop := range []string{"extend", "each", "map", "grep", "inArray", "isArray", "isFunction", "isPlainObject", "parseJSON", "ajax", "get", "post", "fn"} {
		if val := jqueryStatic.Get(prop); val != nil {
			env.runtime.Set(prop, val)
		}
	}
	
	env.runtime.Set("$", jqueryStub)
	env.runtime.Set("jQuery", jqueryStub)
	
	// Add modern framework globals for better compatibility
	env.setupModernFrameworkStubs()
}

// createJQueryObject creates a comprehensive jQuery object with realistic chaining
func (env *JSEnvironment) createJQueryObject() *goja.Object {
	jqObj := env.runtime.NewObject()
	
	// Core properties
	jqObj.Set("length", 0)
	jqObj.Set("selector", "")
	jqObj.Set("context", env.document)
	
	// Traversal methods (all chainable)
	chainableMethods := map[string]func() goja.Value{
		"find": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"filter": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"not": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"eq": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"first": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"last": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"parent": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"parents": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"children": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"siblings": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"next": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"prev": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"closest": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
	}
	
	// Event methods (chainable)
	eventMethods := map[string]func() goja.Value{
		"on": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"off": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"one": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"trigger": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"click": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"dblclick": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"mousedown": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"mouseup": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"mouseover": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"mouseout": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"keydown": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"keyup": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"focus": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"blur": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"change": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"submit": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"resize": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"scroll": func() goja.Value { return env.runtime.ToValue(jqObj) },
	}
	
	// CSS and styling methods (chainable)
	styleMethods := map[string]func() goja.Value{
		"addClass": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"removeClass": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"toggleClass": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"hasClass": func() goja.Value { return env.runtime.ToValue(false) },
		"css": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"show": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"hide": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"toggle": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"fadeIn": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"fadeOut": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"slideUp": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"slideDown": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"animate": func() goja.Value { return env.runtime.ToValue(jqObj) },
	}
	
	// DOM manipulation methods (chainable)
	domMethods := map[string]func() goja.Value{
		"append": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"prepend": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"after": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"before": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"remove": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"empty": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"clone": func() goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"wrap": func() goja.Value { return env.runtime.ToValue(jqObj) },
		"unwrap": func() goja.Value { return env.runtime.ToValue(jqObj) },
	}
	
	// Content methods (getters/setters)
	contentMethods := map[string]func(call goja.FunctionCall) goja.Value{
		"html": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
		"text": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
		"val": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
		"attr": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 1 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
		"prop": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 1 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
		"data": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 1 {
				return env.runtime.ToValue(jqObj) // setter
			}
			return env.runtime.ToValue("") // getter
		},
	}
	
	// Utility methods
	utilityMethods := map[string]func(call goja.FunctionCall) goja.Value{
		"each": func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(jqObj) },
		"map": func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(env.createJQueryObject()) },
		"is": func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) },
		"index": func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(-1) },
		"size": func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(0) },
		"toArray": func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() },
		"get": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				return env.runtime.ToValue(env.createMockElement())
			}
			return env.runtime.NewArray()
		},
	}
	
	// Special methods
	jqObj.Set("ready", func(call goja.FunctionCall) goja.Value {
		// Execute callback immediately for compatibility
		if len(call.Arguments) > 0 {
			if callback, ok := goja.AssertFunction(call.Arguments[0]); ok {
				// Call the callback function with jQuery object as 'this'
				callback(env.runtime.ToValue(jqObj))
			}
		}
		return env.runtime.ToValue(jqObj)
	})
	
	// Add all methods to the jQuery object
	for name, method := range chainableMethods {
		jqObj.Set(name, func(call goja.FunctionCall) goja.Value { return method() })
	}
	for name, method := range eventMethods {
		jqObj.Set(name, func(call goja.FunctionCall) goja.Value { return method() })
	}
	for name, method := range styleMethods {
		jqObj.Set(name, func(call goja.FunctionCall) goja.Value { return method() })
	}
	for name, method := range domMethods {
		jqObj.Set(name, func(call goja.FunctionCall) goja.Value { return method() })
	}
	for name, method := range contentMethods {
		jqObj.Set(name, method)
	}
	for name, method := range utilityMethods {
		jqObj.Set(name, method)
	}
	
	return jqObj
}

// setupModernFrameworkStubs adds stubs for modern JavaScript frameworks
func (env *JSEnvironment) setupModernFrameworkStubs() {
	// React-like globals
	env.runtime.Set("React", env.runtime.NewObject())
	env.runtime.Set("ReactDOM", env.runtime.NewObject())
	
	// Vue-like globals
	vue := env.runtime.NewObject()
	vue.Set("config", env.runtime.NewObject())
	vue.Set("component", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	vue.Set("directive", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	env.runtime.Set("Vue", vue)
	
	// Angular-like globals
	angular := env.runtime.NewObject()
	angular.Set("module", func(call goja.FunctionCall) goja.Value {
		module := env.runtime.NewObject()
		module.Set("controller", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(module) })
		module.Set("service", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(module) })
		module.Set("directive", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(module) })
		return env.runtime.ToValue(module)
	})
	env.runtime.Set("angular", angular)
	
	// Lodash/Underscore-like utilities
	underscore := env.runtime.NewObject()
	underscore.Set("each", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	underscore.Set("map", func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() })
	underscore.Set("filter", func(call goja.FunctionCall) goja.Value { return env.runtime.NewArray() })
	underscore.Set("find", func(call goja.FunctionCall) goja.Value { return goja.Undefined() })
	underscore.Set("extend", func(call goja.FunctionCall) goja.Value { return env.runtime.NewObject() })
	underscore.Set("isArray", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	underscore.Set("isObject", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	underscore.Set("isFunction", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(false) })
	env.runtime.Set("_", underscore)
	
	// Moment.js-like date library
	moment := func(call goja.FunctionCall) goja.Value {
		mom := env.runtime.NewObject()
		mom.Set("format", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue("") })
		mom.Set("add", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(mom) })
		mom.Set("subtract", func(call goja.FunctionCall) goja.Value { return env.runtime.ToValue(mom) })
		return env.runtime.ToValue(mom)
	}
	env.runtime.Set("moment", moment)
}

// setupSiteSpecificStubs creates stubs for specific websites and platforms
func (env *JSEnvironment) setupSiteSpecificStubs() {
	// Common global functions
	env.runtime.Set("loadScript", func(call goja.FunctionCall) goja.Value {
		log.Println("Stub: loadScript called")
		return goja.Undefined()
	})
	env.runtime.Set("IOMm", env.runtime.NewObject())
	
	// CMS and platform globals
	env.runtime.Set("wp", env.runtime.NewObject()) // WordPress
	env.runtime.Set("StackExchange", env.runtime.NewObject()) // Stack Overflow
	
	// Analytics
	dataLayer := env.runtime.NewArray()
	env.runtime.Set("dataLayer", dataLayer)
	env.window.Set("dataLayer", dataLayer)
}

// ExecuteScript runs JavaScript with error handling and timeout
func (env *JSEnvironment) ExecuteScript(script string) {
	// Get timeout from configuration
	timeoutSeconds := time.Duration(env.config.JavaScriptCompatibility.TimeoutSeconds) * time.Second
	
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds)
	defer cancel()
	
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("panic: %v", r)
			}
		}()
		
		_, err := env.runtime.RunString(script)
		done <- err
	}()
	
	select {
	case err := <-done:
		if err != nil {
			// Categorize and handle the error
			jsErr := env.errorHandler.CategorizeError(err, script)
			
			// Attempt graceful degradation for compatibility errors
			if jsErr.Type == CompatibilityError {
				env.handleCompatibilityError(jsErr)
			}
			
			// Log the error with appropriate severity
			env.errorHandler.LogError(jsErr)
		}
	case <-ctx.Done():
		// Handle timeout as a specific error type
		timeoutErr := fmt.Errorf("execution timeout after %v", timeoutSeconds)
		jsErr := env.errorHandler.CategorizeError(timeoutErr, script)
		env.errorHandler.LogError(jsErr)
	}
}

// executeJS extracts and runs JavaScript from the document using Goja with enhanced handling.
func executeJS(doc *goquery.Document) {
	// Check if JavaScript is enabled via configuration
	config, err := loadJSConfig("js_config.json")
	if err != nil {
		log.Printf("Failed to load JS config, using defaults: %v", err)
		config = loadDefaultJSConfig()
	}
	
	if !config.JavaScriptCompatibility.Enabled {
		log.Println("JavaScript execution disabled via configuration")
		return
	}
	
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		script := s.Text()
		if script == "" {
			return
		}
		// Selective execution: Skip potentially nasty JS with 'eval'
		if strings.Contains(script, "eval") {
			log.Println("Skipping potentially nasty JS containing 'eval'")
			return
		}
		// Create new environment for each script
		env := NewJSEnvironment()
		// Set document title from actual DOM
		if env.document != nil {
			env.document.Set("title", doc.Find("title").Text())
			env.document.Set("write", func(call goja.FunctionCall) goja.Value {
				for _, arg := range call.Arguments {
					fmt.Println("[JS Document Write]:", arg.String())
					// TODO: Actually append to DOM and re-render
				}
				return goja.Undefined()
			})
		}
		env.ExecuteScript(script)
	})
	// After JS execution, potentially re-render if DOM changed
	// For now, just print a message
	fmt.Println("JS execution completed. Dynamic content may have been added.")
}