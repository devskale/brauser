package js

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dop251/goja"
	"brauser/config"
)

// JSEnvironment manages the JavaScript runtime and stubs
type JSEnvironment struct {
	vm     *goja.Runtime
	config *config.JSConfig
}

// NewJSEnvironment creates a new JavaScript environment with the given configuration
func NewJSEnvironment(jsConfig *config.JSConfig) *JSEnvironment {
	return &JSEnvironment{
		vm:     goja.New(),
		config: jsConfig,
	}
}

// SetupAllStubs sets up all JavaScript stubs based on the configuration
func (env *JSEnvironment) SetupAllStubs() {
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

// ExecuteScript executes JavaScript code with timeout and error handling
func (env *JSEnvironment) ExecuteScript(script string) error {
	done := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("JavaScript panic: %v", r)
			}
		}()

		_, err := env.vm.RunString(script)
		done <- err
	}()

	timeout := time.Duration(env.config.JavaScriptCompatibility.TimeoutSeconds) * time.Second
	select {
	case err := <-done:
		if err != nil {
			errorType := env.categorizeError(err)
			if errorType == "compatibility" {
				env.handleCompatibilityError(err, script)
				return nil // Continue execution after graceful degradation
			}
			log.Printf("JavaScript error (%s): %v", errorType, err)
			return err
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("JavaScript execution timeout after %v", timeout)
	}
}

// categorizeError categorizes JavaScript errors for better handling
func (env *JSEnvironment) categorizeError(err error) string {
	errorStr := strings.ToLower(err.Error())

	if strings.Contains(errorStr, "is not defined") ||
		strings.Contains(errorStr, "is not a function") ||
		strings.Contains(errorStr, "is not a constructor") ||
		strings.Contains(errorStr, "cannot read property") ||
		strings.Contains(errorStr, "cannot read properties") {
		return "compatibility"
	}

	if strings.Contains(errorStr, "syntax") {
		return "syntax"
	}

	if strings.Contains(errorStr, "timeout") {
		return "timeout"
	}

	return "runtime"
}

// handleCompatibilityError attempts graceful degradation for compatibility errors
func (env *JSEnvironment) handleCompatibilityError(err error, script string) {
	errorStr := err.Error()
	log.Printf("Attempting graceful degradation for: %v", err)

	// Extract the undefined API name from the error
	var apiName string
	if strings.Contains(errorStr, "is not defined") {
		parts := strings.Split(errorStr, " ")
		for i, part := range parts {
			if part == "is" && i > 0 {
				apiName = strings.Trim(parts[i-1], "'\"")
				break
			}
		}
	}

	// Add common missing APIs dynamically
	commonAPIs := map[string]string{
		"dispatchEvent": "function dispatchEvent() { return true; }",
		"requestAnimationFrame": "function requestAnimationFrame(callback) { setTimeout(callback, 16); }",
		"fetch": "function fetch() { return Promise.resolve({json: function() { return Promise.resolve({}); }}); }",
		"IntersectionObserver": "function IntersectionObserver() { this.observe = function() {}; this.disconnect = function() {}; }",
		"MutationObserver": "function MutationObserver() { this.observe = function() {}; this.disconnect = function() {}; }",
	}

	if apiName != "" {
		if stubCode, exists := commonAPIs[apiName]; exists {
			log.Printf("Adding stub for missing API: %s", apiName)
			env.vm.RunString(stubCode)
		} else {
			// Generic stub
			genericStub := fmt.Sprintf("var %s = function() { return {}; };", apiName)
			log.Printf("Adding generic stub: %s", genericStub)
			env.vm.RunString(genericStub)
		}
	}
}