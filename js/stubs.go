package js

import (
	"log"
	"github.com/dop251/goja"
)

// setupConsoleStubs creates console object with log method
func (env *JSEnvironment) setupConsoleStubs() {
	consoleObj := env.vm.NewObject()
	consoleObj.Set("log", func(args ...interface{}) {
		log.Println(args...)
	})
	env.vm.Set("console", consoleObj)
}

// setupDOMStubs creates document object with basic DOM methods
func (env *JSEnvironment) setupDOMStubs() {
	documentObj := env.vm.NewObject()
	
	// Basic document methods
	documentObj.Set("getElementById", func(id string) interface{} {
		return env.createMockElement()
	})
	
	documentObj.Set("createElement", func(tagName string) interface{} {
		return env.createMockElement()
	})
	
	documentObj.Set("getElementsByTagName", func(tagName string) interface{} {
		return []interface{}{env.createMockElement()}
	})
	
	documentObj.Set("querySelector", func(selector string) interface{} {
		return env.createMockElement()
	})
	
	documentObj.Set("querySelectorAll", func(selector string) interface{} {
		return []interface{}{env.createMockElement()}
	})
	
	documentObj.Set("addEventListener", func(event string, handler interface{}) {})
	
	env.vm.Set("document", documentObj)
}

// createMockElement creates a comprehensive mock DOM element
func (env *JSEnvironment) createMockElement() *goja.Object {
	element := env.vm.NewObject()
	
	// Core properties
	element.Set("innerHTML", "")
	element.Set("outerHTML", "")
	element.Set("textContent", "")
	element.Set("id", "")
	element.Set("tagName", "DIV")
	element.Set("nodeType", 1)
	
	// Style object
	styleObj := env.vm.NewObject()
	styleObj.Set("display", "")
	styleObj.Set("visibility", "")
	styleObj.Set("position", "")
	styleObj.Set("top", "")
	styleObj.Set("left", "")
	styleObj.Set("width", "")
	styleObj.Set("height", "")
	styleObj.Set("setProperty", func(prop, value string) {})
	styleObj.Set("getPropertyValue", func(prop string) string { return "" })
	styleObj.Set("removeProperty", func(prop string) {})
	element.Set("style", styleObj)
	
	// ClassList object
	classList := env.vm.NewObject()
	classList.Set("add", func(className string) {})
	classList.Set("remove", func(className string) {})
	classList.Set("toggle", func(className string) bool { return true })
	classList.Set("contains", func(className string) bool { return false })
	classList.Set("replace", func(oldClass, newClass string) {})
	element.Set("classList", classList)
	
	// Attributes
	element.Set("setAttribute", func(name, value string) {})
	element.Set("getAttribute", func(name string) interface{} { return nil })
	element.Set("removeAttribute", func(name string) {})
	element.Set("hasAttribute", func(name string) bool { return false })
	element.Set("getAttributeNames", func() []string { return []string{} })
	
	// Dataset
	dataset := env.vm.NewObject()
	element.Set("dataset", dataset)
	
	// Parent/child relationships
	element.Set("parentNode", env.createMockParentNode())
	element.Set("parentElement", env.createMockParentNode())
	element.Set("appendChild", func(child interface{}) interface{} { return child })
	element.Set("insertBefore", func(newNode, referenceNode interface{}) interface{} { return newNode })
	element.Set("removeChild", func(child interface{}) interface{} { return child })
	element.Set("replaceChild", func(newChild, oldChild interface{}) interface{} { return oldChild })
	
	// Child collections
	element.Set("childNodes", []interface{}{})
	element.Set("children", []interface{}{})
	element.Set("firstChild", nil)
	element.Set("lastChild", nil)
	element.Set("firstElementChild", nil)
	element.Set("lastElementChild", nil)
	element.Set("nextSibling", nil)
	element.Set("previousSibling", nil)
	element.Set("nextElementSibling", nil)
	element.Set("previousElementSibling", nil)
	
	// Query methods
	element.Set("querySelector", func(selector string) interface{} { return env.createMockElement() })
	element.Set("querySelectorAll", func(selector string) []interface{} { return []interface{}{env.createMockElement()} })
	element.Set("getElementsByTagName", func(tagName string) []interface{} { return []interface{}{env.createMockElement()} })
	element.Set("getElementsByClassName", func(className string) []interface{} { return []interface{}{env.createMockElement()} })
	
	// Event handling
	element.Set("addEventListener", func(event string, handler interface{}) {})
	element.Set("removeEventListener", func(event string, handler interface{}) {})
	element.Set("dispatchEvent", func(event interface{}) bool { return true })
	
	// Common event properties
	element.Set("onclick", nil)
	element.Set("onload", nil)
	element.Set("onerror", nil)
	element.Set("onchange", nil)
	element.Set("onsubmit", nil)
	
	// Geometry and positioning
	element.Set("offsetWidth", 0)
	element.Set("offsetHeight", 0)
	element.Set("offsetTop", 0)
	element.Set("offsetLeft", 0)
	element.Set("offsetParent", nil)
	element.Set("clientWidth", 0)
	element.Set("clientHeight", 0)
	element.Set("scrollWidth", 0)
	element.Set("scrollHeight", 0)
	element.Set("scrollTop", 0)
	element.Set("scrollLeft", 0)
	
	element.Set("getBoundingClientRect", func() *goja.Object {
		rect := env.vm.NewObject()
		rect.Set("top", 0)
		rect.Set("left", 0)
		rect.Set("bottom", 0)
		rect.Set("right", 0)
		rect.Set("width", 0)
		rect.Set("height", 0)
		return rect
	})
	
	// Form-specific properties
	element.Set("value", "")
	element.Set("checked", false)
	element.Set("disabled", false)
	element.Set("readonly", false)
	element.Set("selected", false)
	element.Set("type", "")
	element.Set("name", "")
	element.Set("form", nil)
	
	// Focus and interaction
	element.Set("focus", func() {})
	element.Set("blur", func() {})
	element.Set("click", func() {})
	
	// Content manipulation
	element.Set("cloneNode", func(deep bool) interface{} { return env.createMockElement() })
	element.Set("remove", func() {})
	
	// Framework-specific properties
	element.Set("__reactInternalInstance", nil) // React
	element.Set("__vue__", nil)                  // Vue
	element.Set("$$watchers", nil)               // Angular
	
	return element
}

// createMockParentNode creates a mock parent node
func (env *JSEnvironment) createMockParentNode() *goja.Object {
	parent := env.vm.NewObject()
	parent.Set("nodeName", "BODY")
	parent.Set("nodeType", 1)
	parent.Set("tagName", "BODY")
	parent.Set("appendChild", func(child interface{}) interface{} { return child })
	parent.Set("insertBefore", func(newNode, referenceNode interface{}) interface{} { return newNode })
	parent.Set("removeChild", func(child interface{}) interface{} { return child })
	parent.Set("replaceChild", func(newChild, oldChild interface{}) interface{} { return oldChild })
	parent.Set("querySelector", func(selector string) interface{} { return env.createMockElement() })
	parent.Set("querySelectorAll", func(selector string) []interface{} { return []interface{}{env.createMockElement()} })
	parent.Set("getElementsByTagName", func(tagName string) []interface{} { return []interface{}{env.createMockElement()} })
	parent.Set("getElementsByClassName", func(className string) []interface{} { return []interface{}{env.createMockElement()} })
	parent.Set("addEventListener", func(event string, handler interface{}) {})
	parent.Set("removeEventListener", func(event string, handler interface{}) {})
	return parent
}