package js

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"brauser/config"
)

// ExecuteJS processes and executes JavaScript from HTML document
func ExecuteJS(doc *goquery.Document, title string) {
	// Load JavaScript configuration
	jsConfig, err := config.LoadJSConfig("js_config.json")
	if err != nil {
		log.Printf("Failed to load JS config, using defaults: %v", err)
		jsConfig = config.LoadDefaultJSConfig()
	}

	if !jsConfig.JavaScriptCompatibility.Enabled {
		log.Println("JavaScript execution is disabled")
		return
	}

	log.Println("Processing JavaScript...")

	// Find and execute all script tags
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()
		if scriptContent == "" {
			return
		}

		// Skip scripts that contain eval (potential security risk)
		if strings.Contains(strings.ToLower(scriptContent), "eval") {
			log.Println("Skipping script containing 'eval' for security reasons")
			return
		}

		log.Printf("Executing script %d...", i+1)

		// Create a new JavaScript environment for each script
		env := NewJSEnvironment(jsConfig)
		env.SetupAllStubs()

		// Set document title in the JavaScript environment
		if title != "" {
			env.vm.Set("document.title", title)
		}

		// Provide a document.write stub
		env.vm.Set("document.write", func(content string) {
			log.Printf("document.write called: %s", content)
		})

		// Execute the script
		if err := env.ExecuteScript(scriptContent); err != nil {
			log.Printf("Script %d execution failed: %v", i+1, err)
		} else {
			log.Printf("Script %d executed successfully", i+1)
		}
	})
}