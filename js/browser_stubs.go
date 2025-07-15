package js

import "github.com/dop251/goja"

// setupBrowserStubs creates window, navigator, and location objects
func (env *JSEnvironment) setupBrowserStubs() {
	// Window object
	windowObj := env.vm.NewObject()
	windowObj.Set("innerWidth", 1024)
	windowObj.Set("innerHeight", 768)
	windowObj.Set("outerWidth", 1024)
	windowObj.Set("outerHeight", 768)
	windowObj.Set("addEventListener", func(event string, handler interface{}) {})
	windowObj.Set("removeEventListener", func(event string, handler interface{}) {})
	windowObj.Set("setTimeout", func(callback interface{}, delay int) int { return 1 })
	windowObj.Set("clearTimeout", func(id int) {})
	windowObj.Set("setInterval", func(callback interface{}, delay int) int { return 1 })
	windowObj.Set("clearInterval", func(id int) {})
	env.vm.Set("window", windowObj)
	
	// Navigator object
	navigatorObj := env.vm.NewObject()
	navigatorObj.Set("userAgent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	navigatorObj.Set("platform", "MacIntel")
	navigatorObj.Set("language", "en-US")
	navigatorObj.Set("cookieEnabled", true)
	navigatorObj.Set("onLine", true)
	env.vm.Set("navigator", navigatorObj)
	
	// Location object
	locationObj := env.vm.NewObject()
	locationObj.Set("href", "http://localhost")
	locationObj.Set("protocol", "http:")
	locationObj.Set("host", "localhost")
	locationObj.Set("hostname", "localhost")
	locationObj.Set("port", "")
	locationObj.Set("pathname", "/")
	locationObj.Set("search", "")
	locationObj.Set("hash", "")
	locationObj.Set("reload", func() {})
	env.vm.Set("location", locationObj)
	
	// TCF API stub for consent management
	env.vm.Set("__tcfapiLocator", env.vm.NewObject())
}

// setupStorageStubs creates localStorage and sessionStorage
func (env *JSEnvironment) setupStorageStubs() {
	env.setupStorageAPIs()
}

// setupStorageAPIs implements localStorage and sessionStorage
func (env *JSEnvironment) setupStorageAPIs() {
	createStorage := func() *goja.Object {
		storage := env.vm.NewObject()
		storage.Set("getItem", func(key string) interface{} { return nil })
		storage.Set("setItem", func(key, value string) {})
		storage.Set("removeItem", func(key string) {})
		storage.Set("clear", func() {})
		storage.Set("length", 0)
		return storage
	}
	
	env.vm.Set("localStorage", createStorage())
	env.vm.Set("sessionStorage", createStorage())
}

// setupWebAPIStubs creates modern web API stubs
func (env *JSEnvironment) setupWebAPIStubs() {
	// matchMedia
	env.vm.Set("matchMedia", func(query string) *goja.Object {
		mql := env.vm.NewObject()
		mql.Set("matches", false)
		mql.Set("media", query)
		mql.Set("addListener", func(handler interface{}) {})
		mql.Set("removeListener", func(handler interface{}) {})
		mql.Set("addEventListener", func(event string, handler interface{}) {})
		mql.Set("removeEventListener", func(event string, handler interface{}) {})
		return mql
	})
	
	// CustomEvent constructor
	env.vm.Set("CustomEvent", func(call goja.ConstructorCall) *goja.Object {
		event := env.vm.NewObject()
		if len(call.Arguments) > 0 {
			event.Set("type", call.Arguments[0].String())
		}
		if len(call.Arguments) > 1 {
			if options := call.Arguments[1].ToObject(env.vm); options != nil {
				if detail := options.Get("detail"); detail != nil {
					event.Set("detail", detail)
				}
				if bubbles := options.Get("bubbles"); bubbles != nil {
					event.Set("bubbles", bubbles.ToBoolean())
				}
				if cancelable := options.Get("cancelable"); cancelable != nil {
					event.Set("cancelable", cancelable.ToBoolean())
				}
			}
		}
		return event
	})
	
	// URLSearchParams constructor
	env.vm.Set("URLSearchParams", func(call goja.ConstructorCall) *goja.Object {
		urlParams := env.vm.NewObject()
		urlParams.Set("get", func(name string) interface{} { return nil })
		urlParams.Set("set", func(name, value string) {})
		urlParams.Set("has", func(name string) bool { return false })
		urlParams.Set("append", func(name, value string) {})
		urlParams.Set("delete", func(name string) {})
		return urlParams
	})
}

// setupFrameworkStubs creates framework-specific stubs
func (env *JSEnvironment) setupFrameworkStubs() {
	if env.config.JavaScriptCompatibility.Categories.Frameworks.JQuery.Enabled {
		env.setupJQueryStubs()
	}
	env.setupModernFrameworkStubs()
}

// setupJQueryStubs creates comprehensive jQuery stubs
func (env *JSEnvironment) setupJQueryStubs() {
	// Create main jQuery function
	jQuery := func(call goja.FunctionCall) goja.Value {
		return env.vm.ToValue(env.createJQueryObject())
	}
	
	// Add static methods to jQuery
	jQueryObj := env.vm.NewObject()
	jQueryObj.Set("extend", func(args ...interface{}) interface{} { return args[0] })
	jQueryObj.Set("each", func(obj interface{}, callback interface{}) {})
	jQueryObj.Set("map", func(obj interface{}, callback interface{}) []interface{} { return []interface{}{} })
	jQueryObj.Set("grep", func(array []interface{}, callback interface{}) []interface{} { return array })
	jQueryObj.Set("inArray", func(value interface{}, array []interface{}) int { return -1 })
	jQueryObj.Set("isArray", func(obj interface{}) bool { return false })
	jQueryObj.Set("isFunction", func(obj interface{}) bool { return false })
	jQueryObj.Set("isPlainObject", func(obj interface{}) bool { return false })
	jQueryObj.Set("parseJSON", func(json string) interface{} { return nil })
	
	// AJAX methods
	jQueryObj.Set("ajax", func(options interface{}) *goja.Object {
		jqXHR := env.vm.NewObject()
		jqXHR.Set("done", func(callback interface{}) *goja.Object { return jqXHR })
		jqXHR.Set("fail", func(callback interface{}) *goja.Object { return jqXHR })
		jqXHR.Set("always", func(callback interface{}) *goja.Object { return jqXHR })
		return jqXHR
	})
	jQueryObj.Set("get", func(url string, callback interface{}) *goja.Object {
		return jQueryObj.Get("ajax").(*goja.Object)
	})
	jQueryObj.Set("post", func(url string, data, callback interface{}) *goja.Object {
		return jQueryObj.Get("ajax").(*goja.Object)
	})
	
	// fn property for plugins
	jQueryObj.Set("fn", env.vm.NewObject())
	
	env.vm.Set("$", jQuery)
	env.vm.Set("jQuery", jQuery)
}

// createJQueryObject creates a jQuery object with chainable methods
func (env *JSEnvironment) createJQueryObject() *goja.Object {
	jq := env.vm.NewObject()
	
	// Core properties
	jq.Set("length", 1)
	jq.Set("selector", "")
	
	// Chainable methods that return jQuery object
	chainableMethods := []string{
		"addClass", "removeClass", "toggleClass", "hasClass",
		"attr", "removeAttr", "prop", "removeProp",
		"css", "show", "hide", "toggle",
		"on", "off", "trigger", "click", "focus", "blur",
		"append", "prepend", "after", "before", "remove", "empty",
		"parent", "parents", "children", "siblings", "next", "prev",
		"find", "filter", "not", "eq", "first", "last",
		"animate", "fadeIn", "fadeOut", "slideUp", "slideDown",
	}
	
	for _, method := range chainableMethods {
		jq.Set(method, func(args ...interface{}) *goja.Object { return jq })
	}
	
	// Content methods (getters/setters)
	jq.Set("html", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			return jq // setter
		}
		return "" // getter
	})
	jq.Set("text", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			return jq // setter
		}
		return "" // getter
	})
	jq.Set("val", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			return jq // setter
		}
		return "" // getter
	})
	
	// Utility methods
	jq.Set("each", func(callback interface{}) *goja.Object { return jq })
	jq.Set("map", func(callback interface{}) *goja.Object { return jq })
	jq.Set("is", func(selector interface{}) bool { return false })
	jq.Set("index", func(element interface{}) int { return 0 })
	jq.Set("size", func() int { return 1 })
	jq.Set("toArray", func() []interface{} { return []interface{}{} })
	jq.Set("get", func(index interface{}) interface{} { return env.createMockElement() })
	
	// Ready function
	jq.Set("ready", func(callback interface{}) *goja.Object {
		// Immediately execute the callback
		if fn, ok := goja.AssertFunction(env.vm.ToValue(callback)); ok {
			fn(goja.Undefined())
		}
		return jq
	})
	
	return jq
}

// setupModernFrameworkStubs creates stubs for modern frameworks
func (env *JSEnvironment) setupModernFrameworkStubs() {
	// React stubs
	reactObj := env.vm.NewObject()
	reactObj.Set("createElement", func(args ...interface{}) interface{} { return env.vm.NewObject() })
	reactObj.Set("Component", func() {})
	env.vm.Set("React", reactObj)
	
	reactDOMObj := env.vm.NewObject()
	reactDOMObj.Set("render", func(element, container interface{}) {})
	env.vm.Set("ReactDOM", reactDOMObj)
	
	// Vue stubs
	vueObj := env.vm.NewObject()
	vueConfig := env.vm.NewObject()
	vueConfig.Set("silent", true)
	vueObj.Set("config", vueConfig)
	vueObj.Set("component", func(name string, options interface{}) {})
	vueObj.Set("directive", func(name string, options interface{}) {})
	env.vm.Set("Vue", vueObj)
	
	// Angular stubs
	angularObj := env.vm.NewObject()
	angularObj.Set("module", func(name string, deps []string) *goja.Object {
		module := env.vm.NewObject()
		module.Set("controller", func(name string, fn interface{}) *goja.Object { return module })
		module.Set("service", func(name string, fn interface{}) *goja.Object { return module })
		module.Set("directive", func(name string, fn interface{}) *goja.Object { return module })
		return module
	})
	env.vm.Set("angular", angularObj)
	
	// Lodash/Underscore stubs
	lodashObj := env.vm.NewObject()
	lodashObj.Set("each", func(collection, iteratee interface{}) {})
	lodashObj.Set("map", func(collection, iteratee interface{}) []interface{} { return []interface{}{} })
	lodashObj.Set("filter", func(collection, predicate interface{}) []interface{} { return []interface{}{} })
	lodashObj.Set("find", func(collection, predicate interface{}) interface{} { return nil })
	lodashObj.Set("extend", func(destination interface{}, sources ...interface{}) interface{} { return destination })
	lodashObj.Set("isArray", func(value interface{}) bool { return false })
	lodashObj.Set("isObject", func(value interface{}) bool { return false })
	lodashObj.Set("isFunction", func(value interface{}) bool { return false })
	env.vm.Set("_", lodashObj)
	
	// Moment.js stubs
	momentObj := env.vm.NewObject()
	momentObj.Set("format", func(format string) string { return "" })
	momentObj.Set("add", func(amount int, unit string) *goja.Object { return momentObj })
	momentObj.Set("subtract", func(amount int, unit string) *goja.Object { return momentObj })
	env.vm.Set("moment", func(args ...interface{}) *goja.Object { return momentObj })
}

// setupSiteSpecificStubs creates site-specific global stubs
func (env *JSEnvironment) setupSiteSpecificStubs() {
	// Common global functions
	env.vm.Set("loadScript", func(url string, callback interface{}) {})
	env.vm.Set("IOMm", func(args ...interface{}) {})
	
	// CMS/Platform globals
	wpObj := env.vm.NewObject()
	wpObj.Set("ajax", env.vm.NewObject())
	env.vm.Set("wp", wpObj)
	
	stackExchangeObj := env.vm.NewObject()
	stackExchangeObj.Set("ready", func(callback interface{}) {})
	env.vm.Set("StackExchange", stackExchangeObj)
	
	// Analytics
	env.vm.Set("dataLayer", []interface{}{})
}