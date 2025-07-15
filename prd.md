brauser - a minimalistic terminal web browser that can handle modern web pages

## GOAL & Requirements

browse modern web pages
able to handle basic javascript
ascii image rendering

## IDEATION

Ideation on tools and modules:

- Existing terminal browsers: W3M, Lynx, Links2 for inspiration.
- Libraries: ncurses or termbox for terminal UI.
- For web handling: Use Go's net/http for HTTP requests, goquery or similar for HTML parsing.
- JavaScript: QuickJS or Otto for basic JS execution.
- Image rendering: Libraries like sixel or chafa for ASCII art images.

## PLAN / STATUS (=CHECKMARKED ITEMS)

- Architecture: Modular Go application with components for fetching, parsing, rendering, and execution.
- Key Components:
  - HTTP Client: Using net/http to fetch pages.
  - HTML Parser: goquery for DOM manipulation.
  - Terminal Renderer: termbox or tview for UI.
  - JS Engine: Integrate QuickJS for basic scripting.
  - Image Converter: chafa for ASCII rendering.
- Development Steps:
  1. [x] Set up Go project with modules.
  2. [x] Implement basic URL fetching and text rendering.
  3. [x] Add HTML parsing and structured display.
  4. [x] Integrate JS execution for dynamic content.
  - [x] Research and select a lightweight JS runtime library for Go (e.g., goja or otto).
  - [x] Install the selected library.
  - [x] Update code to execute JS scripts after parsing HTML.
  - [x] Implement basic DOM manipulation support.
  - [x] Test JS execution with a sample dynamic page.
  - [x] Implement sandboxing using isolated Goja Runtimes per script.
  - [x] Add try-catch wrappers to skip unrenderable JS with error logging.
  - [x] Provide stubs for browser globals like window and localStorage.
  - [x] Add timeouts to interrupt long-running or nasty JS.
  - [x] Develop selective execution to render useful JS while filtering problematic ones.
  - [x] Test enhanced JS handling on dynamic sites.
  - [x] Add comprehensive DOM stubs for common methods (createElement, getElementById, getElementsByTagName, addEventListener, appendChild).
  - [x] Enhance window and navigator objects with realistic properties (location.protocol, userAgent).
  - [x] Add stubs for common globals like loadScript and IOMm to reduce undefined reference errors.
  - [x] Add remaining DOM methods like insertBefore, removeChild, and querySelector for better compatibility.
  - [x] Add sessionStorage, className, classList, and dataLayer stubs for modern web compatibility.
  - [x] Add matchMedia API for responsive design compatibility.
  - [x] Add CustomEvent, URLSearchParams, jQuery ($), and site-specific globals for broader compatibility.
  - [x] Refactored JavaScript stub architecture with organized structure
- [x] Extracted configuration to external JSON file with runtime customization.
  - [x] Implement more realistic property access patterns to handle complex JS frameworks.
  5. [x] Add ASCII image support.
  - [x] Research and select a Go library for converting images to ASCII art (e.g., ascii-image-converter).
  - [x] Install the selected library.
  - [x] Update code to fetch images from HTML, convert to ASCII, and render in terminal.
  - [x] Test ASCII image rendering with a sample page containing images.
  6. [x] **REFACTORING PHASE** - Improve code organization and maintainability
  - [x] Split monolithic main.go into separate modules
    - [x] Create browser/ package for core browser functionality
    - [x] Create js/ package for JavaScript execution and stubs
    - [x] Create renderer/ package for HTML and image rendering
    - [x] Create config/ package for configuration management
  - [x] Implement proper error handling patterns
  - [x] Add comprehensive unit tests
  - [ ] Optimize performance bottlenecks
  - [ ] Improve memory management
  7. [x] **CONTENT DETECTION & DYNAMIC LOADING IMPROVEMENTS**
  - [x] Enhanced Content Detection: Detect when sites show loading pages vs actual content
  - [x] Improved Image Handling: Filter out problematic image formats (SVG, tracking pixels) to prevent rendering errors
  - [x] Wait Mechanisms: Add optional delays for sites that load content dynamically
  - [x] Content Validation: Verify if extracted content represents actual page or loading screen
  - [x] Site-Specific Handling: Add CodePen-specific handling for better content processing
  - [x] ~~Hacker News Handler~~ → Removed: General HTML renderer handles all sites effectively
  - [x] Loading State Recognition: Identify common loading indicators ("Just a moment...", spinners, etc.)
  - [x] Retry Logic: Implement smart retry mechanisms for pages that require time to load
  - [x] Handle cookie/adblock banners and other interstitial pages
  8. [x] **VISUAL ENHANCEMENTS & GZIP SUPPORT**
  - [x] Enhanced HTML renderer with beautiful formatting
  - [x] Hierarchical heading display with emojis
  - [x] Structured content sections (navigation, links, lists)
  - [x] Content summary statistics
  - [x] GZIP decompression support for HTTP responses
  - [x] Improved ASCII art image rendering
  - [x] Visual content boundaries and separators
  - [x] **Output Formatting Improvements**: Implemented buffered output system to compress multiple consecutive empty lines
  - [x] **Regex-based Whitespace Normalization**: Added pattern matching to clean up excessive whitespace in HTML output
  - [x] **Consistent Output Formatting**: Created helper methods (printf, println, flushOutput) for uniform output handling
  - [x] **Improved Terminal Display**: Enhanced readability by removing redundant empty lines while preserving content structure
  - [x] **Link Text Cleanup**: Fixed excessive empty lines in content links by implementing regex-based whitespace normalization for extracted link text
  9. [x] **INTERACTIVE NAVIGATION SYSTEM** - Complete browser-like navigation experience
  - [x] Link Navigation: Implemented clickable links with numbered selection (e.g., [1], [2], [3])
  - [x] Back/Forward History: Added browser-like navigation with history stack (50 page limit)
  - [x] User Input Handling: Created interactive terminal interface for link selection
  - [x] URL Bar: Allow users to enter new URLs without restarting the application
  - [x] Navigation Menu: Comprehensive command interface (back, forward, history, links, url, refresh, quit)
  - [x] Link Organization: Categorized links by type (navigation, content, stories) for better UX
  - [x] History Management: Full browsing history with current page indicators
  - [x] Cached Content Display: Show cached pages for back/forward navigation
  - [x] URL Resolution: Proper handling of relative URLs against base URL
  - [x] **Back Navigation Fix**: Fixed issue where back navigation showed incorrect links by ensuring cached pages re-extract links from their HTML content
  10. Test on MacOS for cross-platform compatibility.
  11. Optimize for performance and minimalism.

## Project Structure

The brauser project is organized into focused packages for maintainability:

```
brauser/
├── main.go              # Entry point and interactive browsing loop
├── browser/             # HTTP client functionality
│   ├── client.go        # Web page fetching with timeout/headers and GZIP support
│   ├── content_detector.go # Content analysis engine for dynamic loading detection
│   └── site_handlers.go # Site-specific handlers for popular websites
├── js/                  # JavaScript execution environment
│   ├── environment.go   # JS VM setup and stub management
│   ├── executor.go      # Script processing and execution
│   ├── stubs.go         # Core DOM/console stubs
│   └── browser_stubs.go # Browser APIs and framework stubs
├── renderer/            # Content rendering
│   ├── html.go          # HTML parsing and display
│   └── image.go         # ASCII art image conversion
├── navigation/          # Interactive navigation system
│   └── navigator.go     # Link extraction, history, user input handling
├── config/              # Configuration management
│   └── config.go        # JS compatibility settings
└── js_config.json       # Runtime configuration file
```

**Key Design Principles:**

- **Single Responsibility**: Each package handles one specific concern
- **Loose Coupling**: Packages communicate through well-defined interfaces
- **Testability**: Modular structure enables comprehensive unit testing
- **Extensibility**: New features can be added without modifying existing packages

## Test Websites

lets test with websites such as (add more throughout the project, test randomly through them)

https://www.orf.at
https://diepresse.com
https://derstandard.at - successfully renders, shows adblock/js wall, good test case for content detection
https://news.ycombinator.com/
https://techcrunch.com
https://codepen.io/
https://cnn.com

## LEARNINGS, CODING GUIDELINES & CODING RULES

### Coding Guidelines

- **Modular Design**: Keep components loosely coupled and highly cohesive
- **Error Handling**: Always handle errors gracefully and provide meaningful feedback
- **Testing**: Write tests for critical functionality, especially content detection logic
- **Performance**: Implement timeouts and resource limits to prevent hanging
- **User Experience**: Provide clear feedback during long operations
- **Code Organization**: Use packages to separate concerns (browser, renderer, etc.)
- **Import Management**: Use absolute imports in Go modules, avoid relative imports
- **Content Validation**: Always validate extracted content before processing
- **Retry Logic**: Implement exponential backoff for retry mechanisms
- **Resource Management**: Properly close HTTP connections and clean up resources
- **Generalization First**: Prefer building robust general-purpose components over site-specific solutions. Only add specialized handlers when the general approach fails.
- **Graceful Degradation**: When optional features (like image rendering) fail, continue with core functionality and provide informative error messages rather than breaking the entire experience.

### Key Learnings

1. **JavaScript Compatibility**: Many websites use modern JS features that need to be stubbed out for compatibility
2. **Error Handling**: Graceful degradation is crucial - websites should still be usable even when JS fails
3. **Image Processing**: ASCII art conversion adds visual appeal while maintaining terminal compatibility
4. **Site-Specific Handling**: Different websites have unique requirements (e.g., HackerNews story parsing)
5. **Generalization Over Specialization**: A well-designed general HTML renderer can handle most websites effectively, reducing the need for site-specific handlers. Only create specialized handlers when truly necessary.
6. **Non-Semantic HTML Handling**: Many websites use table-based layouts and CSS classes instead of semantic HTML elements. The renderer must be flexible enough to extract content from common patterns like `.athing`, `.titleline`, etc.
7. **Performance**: Timeout mechanisms prevent hanging on problematic scripts
8. **Code Organization**: Refactoring monolithic code into packages improves maintainability and testability
9. **Go Module System**: Proper import paths are crucial - relative imports don't work in module mode
10. **Struct Definition**: Complex nested structs require careful attention to closing braces and field access patterns
11. **Content Detection**: Modern websites often show loading screens, cookie banners, or adblock messages before actual content
12. **Dynamic Loading**: Many sites load content asynchronously, requiring retry mechanisms with appropriate wait times
13. **Site-Specific Patterns**: Different sites have unique loading patterns and content structures that benefit from specialized handling
14. **User Feedback**: Clear visual indicators help users understand what's happening during page loading and content analysis.
15. **HTTP Compression**: Modern websites use GZIP compression extensively - proper decompression is essential for content parsing.
16. **Visual Terminal Design**: Well-structured terminal output with emojis, separators, and hierarchical formatting significantly improves readability.
17. **Content Extraction Strategy**: Different content types (headings, paragraphs, navigation, lists) require different extraction and display strategies.
18. **ASCII Art Limitations**: SVG and complex image formats often fail ASCII conversion - graceful error handling is important.
19. **Interactive Navigation Design**: Numbered link selection provides intuitive navigation in terminal environments - users can easily select links without complex keyboard navigation.
20. **History Management**: Browser-like back/forward functionality significantly improves user experience - caching content enables instant navigation through history.
21. **Link Categorization**: Organizing links by type (navigation, content, stories) helps users understand page structure and find relevant links faster.
22. **User Input Patterns**: Simple command patterns (single letters, numbers) work best in terminal interfaces - complex commands should be avoided for better usability.
23. **URL Resolution**: Proper relative URL resolution is crucial for navigation - many websites use relative links that need to be resolved against the base URL.
24. **Interactive Loop Design**: Separating page loading from user interaction allows for responsive navigation without reloading pages unnecessarily.
25. **Content Rendering Verification**: Always test content rendering with actual websites to ensure proper display
26. **Debug Output Analysis**: Comprehensive logging helps identify content detection and rendering issues
27. **User Perception vs Reality**: Sometimes content is rendered correctly but user interface issues can create perception of missing content
28. **Output Formatting & User Experience**: Buffer output before displaying to enable post-processing, compress multiple consecutive empty lines for cleaner display, and use regex patterns to normalize whitespace in terminal output
29. **Site Handler HTML Integration**: Site handlers must return properly formatted HTML (not plain text) to integrate correctly with the HTML renderer - the renderer expects HTML structure with headings, paragraphs, and links for proper content extraction and display
30. **Image Format Support**: ASCII image conversion has limitations with certain formats (SVG, some GIFs) - consider filtering out problematic image types or providing fallback text descriptions to improve user experience
31. **Link Text Normalization**: HTML link text often contains excessive whitespace and newlines from the DOM structure - use regex patterns to normalize whitespace (`\s+` → single space) and trim edges for clean terminal display

### Configuration Management

- **External Configuration**: Move complex settings to external JSON/YAML files for runtime customization
- **Granular Control**: Enable/disable specific feature categories independently (console, DOM, browser, storage, webapi, frameworks, site_specific)
- **Fallback Strategy**: Always provide sensible defaults when configuration loading fails
- **Runtime Flexibility**: Allow users to customize compatibility levels without code changes
- **Validation**: Load and validate configuration at startup with clear error messages
- **Category-based Organization**: Group related stubs into logical categories for easier management

## Configuration Management

**Status**: ✅ COMPLETED

Implemented external configuration system for JavaScript compatibility with runtime customization.

### Features

- **External JSON Configuration**: `js_config.json` for all JavaScript stub definitions
- **Granular Control**: Individual enable/disable toggles for each category (Console, DOM, Browser, Storage, Web API, Frameworks, Site-specific)
- **Fallback Strategy**: Automatic fallback to sensible defaults if configuration file is missing or invalid
- **Runtime Flexibility**: Configurable timeouts and execution limits
- **Validation**: Robust error handling for malformed configuration files
- **Category-based Organization**: Group related stubs into logical categories for easier management

## Enhanced Property Access Patterns

### Implemented Features

- **Comprehensive DOM Element Stubs**: Enhanced `createMockElement` with realistic property access patterns including:

  - Dynamic getters/setters for element properties (innerHTML, textContent, className, etc.)
  - Complete style object with common CSS properties and methods
  - Enhanced classList with full API (add, remove, toggle, contains, replace)
  - Comprehensive attribute handling (setAttribute, getAttribute, hasAttribute, etc.)
  - Dataset support for data-\* attributes
  - Child node management with realistic collections and navigation
  - Query methods (querySelector, querySelectorAll, getElementsByTagName, etc.)
  - Event handling properties and methods
  - Geometry and positioning properties (offsetWidth, clientHeight, getBoundingClientRect, etc.)
  - Form-specific properties and focus methods
  - Framework-specific properties for React, Vue, and Angular compatibility

- **Enhanced jQuery Stubs**: Comprehensive jQuery object implementation with:

  - Chainable methods for traversal, manipulation, and effects
  - Static jQuery methods (extend, each, map, grep, etc.)
  - Proper callback handling for ready() method
  - AJAX method stubs
  - Event handling methods

- **Modern Framework Support**: Added stubs for popular JavaScript frameworks:

  - React and ReactDOM globals
  - Vue.js with component and directive methods
  - Angular with module system
  - Lodash/Underscore utilities
  - Moment.js date library

- **Realistic Parent Node**: Enhanced parent node implementation with comprehensive child manipulation methods

### Development Guidelines

- Use comprehensive property access patterns to handle complex JavaScript frameworks
- Implement chainable methods for jQuery-like libraries
- Provide realistic geometry and positioning properties for layout calculations
- Support framework-specific properties for better compatibility
- Maintain consistent API patterns across all stub implementations

## Interactive Navigation System

**Status**: ✅ COMPLETED

Implemented a comprehensive interactive navigation system that transforms brauser from a simple page viewer into a fully functional terminal web browser.

### Features

- **Numbered Link Selection**: All clickable links are automatically numbered [1-50] for easy selection
- **Link Categorization**: Links are organized by type:
  - 🧭 Navigation: Menu and navigation links
  - 📰 Stories: News articles and story links (e.g., Hacker News)
  - 📄 Content Links: General page content links
- **Browser History**: Full back/forward navigation with 50-page history limit
- **Interactive Commands**:
  - Numbers [1-50]: Navigate to numbered link
  - 'b' or 'back': Go to previous page
  - 'f' or 'forward': Go to next page
  - 'h' or 'history': View browsing history
  - 'l' or 'links': Show links again
  - 'u' or 'url': Enter new URL
  - 'r' or 'refresh': Reload current page
  - 'q' or 'quit': Exit browser
- **URL Bar Functionality**: Enter new URLs without restarting the application
- **Cached Navigation**: Instant back/forward using cached content
- **URL Resolution**: Proper handling of relative URLs against base URL

### Implementation Details

- **Navigator Package**: Dedicated `navigation/navigator.go` handles all interactive functionality
- **Link Extraction**: Comprehensive link detection from HTML documents with goquery
- **History Management**: Efficient history stack with current page tracking
- **User Input Processing**: Robust command parsing with error handling
- **Content Caching**: Store page content for instant back/forward navigation
- **Interactive Loop**: Separate page loading from user interaction for responsive UX

### Usage Examples

```bash
# Start interactive browsing
./brauser https://news.ycombinator.com

# Navigate using numbered links
brauser> 8  # Opens link number 8

# Use navigation commands
brauser> b  # Go back
brauser> f  # Go forward
brauser> h  # Show history
brauser> u  # Enter new URL
brauser> q  # Quit
```

### Testing Results

- **Hacker News**: Successfully extracts 50+ story links with proper numbering
- **ORF.at**: Categorizes navigation and content links effectively
- **GitHub**: Handles complex page structures with multiple link types
- **General Sites**: Robust URL resolution for relative links

### User Experience Improvements

- **Intuitive Navigation**: Number-based link selection is faster than scrolling
- **Context Awareness**: Link categorization helps users understand page structure
- **Efficient Browsing**: Back/forward with cached content provides instant navigation
- **Error Handling**: Clear error messages for invalid commands or link numbers
- **Visual Feedback**: Emoji indicators and clear status messages guide user interaction

## Code Quality & Maintainability Improvements

**Status**: 📋 RECOMMENDATIONS

Based on the current codebase analysis, here are key suggestions to enhance code quality and maintainability:

### Architecture Improvements

1. **Interface Segregation**: Consider splitting large interfaces into smaller, more focused ones
   - `SiteHandler` interface could be split into `ContentProcessor` and `RetryHandler`
   - `ContentDetector` could implement separate interfaces for different detection types

2. **Dependency Injection**: Implement proper DI container for better testability
   - Create interfaces for all external dependencies (HTTP client, file system, etc.)
   - Use constructor injection instead of direct instantiation

3. **Configuration Validation**: Add comprehensive config validation
   - Validate JSON schema on startup
   - Provide clear error messages for invalid configurations
   - Add default fallbacks for missing optional settings

### Error Handling Enhancements

4. **Structured Error Types**: Create custom error types for different failure scenarios
   ```go
   type ContentError struct {
       Type    string // "network", "parsing", "timeout"
       URL     string
       Message string
       Cause   error
   }
   ```

5. **Error Context**: Add more context to errors using `fmt.Errorf` with `%w` verb
   - Include URL, operation type, and relevant parameters in error messages
   - Implement error wrapping for better debugging

6. **Graceful Degradation**: Implement fallback mechanisms
   - If JavaScript fails, continue with HTML-only rendering
   - If image rendering fails, show placeholder text
   - If site handler fails, use generic processing

### Testing Strategy

7. **Unit Test Coverage**: Add comprehensive unit tests
   - Test all public methods with various input scenarios
   - Mock external dependencies (HTTP requests, file system)
   - Test error conditions and edge cases

8. **Integration Tests**: Add end-to-end testing
   - Test with real websites (using cached responses)
   - Verify navigation flow and user interactions
   - Test JavaScript execution with various frameworks

9. **Benchmark Tests**: Add performance benchmarks
   - Measure content processing time for different page sizes
   - Track memory usage during navigation sessions
   - Monitor JavaScript execution performance

### Code Organization

10. **Package Structure**: Refine package boundaries
    - Move shared types to a `types` package
    - Create a `utils` package for common utilities
    - Separate HTTP client logic into dedicated package

11. **Constants Management**: Centralize magic numbers and strings
    ```go
    const (
        MaxHistorySize = 50
        MaxLinksDisplay = 50
        DefaultTimeout = 30 * time.Second
        MinContentLength = 500
    )
    ```

12. **Documentation**: Enhance code documentation
    - Add package-level documentation with usage examples
    - Document all public APIs with clear descriptions
    - Include performance characteristics and limitations

### Performance Optimizations

13. **Memory Management**: Implement better memory handling
    - Use object pools for frequently allocated objects
    - Implement LRU cache for page content with size limits
    - Add memory usage monitoring and cleanup

14. **Concurrent Processing**: Add safe concurrency where beneficial
    - Parallel image processing for multiple images
    - Concurrent JavaScript execution for independent scripts
    - Background content prefetching for linked pages

15. **Caching Strategy**: Implement intelligent caching
    - Cache parsed HTML documents to avoid re-parsing
    - Cache JavaScript execution results for repeated scripts
    - Implement cache invalidation based on content changes

### Security Enhancements

16. **Input Validation**: Strengthen input validation
    - Validate all URLs before processing
    - Sanitize user input for navigation commands
    - Implement rate limiting for HTTP requests

17. **Resource Limits**: Add resource consumption limits
    - Limit maximum page size to prevent memory exhaustion
    - Set timeouts for all network operations
    - Implement maximum execution time for JavaScript

### Monitoring & Observability

18. **Structured Logging**: Implement structured logging
    ```go
    logger.Info("Page loaded",
        "url", pageURL,
        "size", contentLength,
        "duration", loadTime,
        "links_found", linkCount)
    ```

19. **Metrics Collection**: Add basic metrics
    - Track page load times and success rates
    - Monitor JavaScript execution success/failure rates
    - Collect user interaction patterns

20. **Health Checks**: Implement health monitoring
    - Verify external dependencies are accessible
    - Monitor memory and CPU usage
    - Check for resource leaks during long sessions

### User Experience Improvements

21. **Progressive Loading**: Show content as it becomes available
    - Display basic HTML structure immediately
    - Load and render images asynchronously
    - Show loading indicators for slow operations

22. **Accessibility**: Enhance terminal accessibility
    - Support screen readers with descriptive text
    - Provide keyboard shortcuts for common actions
    - Implement high contrast mode for better visibility

23. **Configuration UI**: Add interactive configuration
    - Allow runtime configuration changes
    - Provide configuration validation and preview
    - Save user preferences persistently

### Implementation Priority

**High Priority** (Immediate Impact):
- Error handling enhancements (#4, #5, #6)
- Unit test coverage (#7)
- Constants management (#11)
- Input validation (#16)

**Medium Priority** (Quality Improvements):
- Interface segregation (#1)
- Documentation (#12)
- Structured logging (#18)
- Memory management (#13)

**Low Priority** (Future Enhancements):
- Dependency injection (#2)
- Concurrent processing (#14)
- Metrics collection (#19)
- Progressive loading (#21)

## Enhanced Error Handling

**Status**: ✅ COMPLETED

Implemented comprehensive JavaScript error handling with categorization, suppression, and graceful degradation.

### Features

- **Error Categorization**: Automatic classification of JavaScript errors into types:
  - Syntax Errors: Malformed JavaScript code
  - Runtime Errors: Execution-time failures
  - Compatibility Errors: Missing APIs or browser features
  - Timeout Errors: Script execution timeouts
  - Harmless Errors: Analytics, ads, and tracking scripts
- **Error Suppression**: Pattern-based suppression of known harmless errors to reduce log noise
- **Graceful Degradation**: Automatic addition of missing API stubs when compatibility errors are detected
- **Intelligent Logging**: Severity-based logging (DEBUG, INFO, WARN, ERROR) based on error type
- **Dynamic API Injection**: Runtime addition of common missing APIs (fetch, requestAnimationFrame, IntersectionObserver, etc.)
- **Regex Pattern Matching**: Sophisticated error pattern recognition for analytics, social media, and ad-related scripts

### Error Handling Best Practices

- Categorize errors by type (Syntax, Runtime, Compatibility, Timeout, Harmless) for appropriate handling
- Implement error suppression patterns to reduce noise from known harmless errors
- Use graceful degradation to automatically add missing API stubs when compatibility errors occur
- Log errors with appropriate severity levels (DEBUG, INFO, WARN, ERROR) based on error type
- Separate error handling logic into dedicated components for maintainability
- Pattern matching enables intelligent error classification and automated responses

### JavaScript Compatibility

- Modern websites require extensive DOM API stubs for proper functionality
- Event handling (addEventListener) is critical for interactive content
- Storage APIs (localStorage, sessionStorage) are commonly used
- Framework detection and stubs (jQuery, site-specific globals) improve compatibility
- Error wrapping and timeouts prevent infinite loops and crashes
- Cross-site testing reveals different compatibility requirements per site

### Architecture & Code Organization

- **Structured stub architecture**: Organize JavaScript stubs into logical categories (DOM, Browser, WebAPI, Framework, Site-specific)
- **Separation of concerns**: Use dedicated types and methods for different functionality areas
- **Reusable components**: Create factory methods for common objects (mock elements, storage APIs)
- **Centralized configuration**: Group related stubs together for easier maintenance
- **Error handling isolation**: Implement timeout and panic recovery at the environment level
- **Clean interfaces**: Provide simple, focused public methods for complex internal logic

### Error Handling

- Graceful degradation when images fail to load
- JavaScript errors should not crash the browser
- Timeout mechanisms prevent hanging on long-running scripts
- Panic recovery should be implemented at the runtime level

### Performance

- ASCII image conversion should be cached for repeated requests
- JavaScript execution should have reasonable timeouts
- HTTP requests should have proper timeout configurations
- Reuse JavaScript environments when possible to reduce initialization overhead

### Development Guidelines

- Always include function-level comments to improve code readability and maintainability.
- Modularize code by separating concerns into distinct functions for easier debugging and extension.
- Implement robust error handling using logging to track and diagnose issues effectively.
- Write unit tests for key functions to verify behavior and prevent regressions.
- Regularly run `go mod tidy` to manage and clean up module dependencies.
- use timeouts, but not longer than 10s
- Add a User-Agent header to HTTP requests to improve compatibility with websites that restrict access based on client identification.
- Integrate a lightweight JavaScript runtime like Goja to execute embedded scripts and handle dynamic content.
- Bind Go objects to JavaScript for basic DOM interactions, enabling scripts to manipulate document elements.
- The current JS integration executes simple scripts but encounters errors on sites requiring browser globals like 'window'; extend bindings for better compatibility.
- Test on multiple JavaScript-heavy sites to identify and address compatibility issues early.
- Integrated ascii-image-converter library to convert fetched images to colored ASCII art for terminal display.
- Handle relative image URLs by resolving against the base page URL, use temporary files for conversion, and clean up afterward.
- Ensure ASCII rendering dimensions fit terminal constraints for optimal viewing; tested on pages with images to confirm functionality.
- Enhanced JS handling with isolated Runtimes for sandboxing, try-catch for skipping errors, stubs for globals like window/localStorage to render useful JS.
- Implemented timeouts using Interrupt to tackle long-running/nasty JS, and selective skipping based on patterns like 'eval' for safety.
- Added comprehensive DOM method stubs (createElement, getElementById, getElementsByTagName, addEventListener, appendChild, setAttribute) to reduce JS errors.
- Enhanced browser object stubs with realistic properties: window.location with protocol/host, navigator.userAgent, and common globals like loadScript.
- Testing on diepresse.com shows significant reduction in JS errors; insertBefore now working, remaining issues mainly property access patterns and syntax errors.
- Successfully implemented comprehensive DOM method stubs: createElement, getElementById, getElementsByTagName, addEventListener, appendChild, setAttribute, insertBefore, removeChild, querySelector, querySelectorAll.
- Enhanced parent node objects with proper method implementations to handle complex DOM manipulation patterns used by modern websites.
- Added modern web API stubs: sessionStorage, className, classList (add/remove), dataLayer for Google Analytics, and wp global for WordPress sites.
- Added responsive design API: matchMedia with proper MediaQueryList object structure.
- Added web standards APIs: CustomEvent constructor, URLSearchParams for URL manipulation.
- Added jQuery compatibility: comprehensive $ and jQuery stubs with chainable methods.
- Added site-specific globals: StackExchange for Stack Overflow compatibility.
- Cross-site testing results:
  - GitHub.com: Only syntax errors (excellent compatibility) ✅
  - Medium.com: Minimal errors, mostly syntax issues ✅
  - React.dev: Clean execution after matchMedia addition ✅
  - StackOverflow.com: Major APIs working, only method-specific issues remain ✅
- Testing shows excellent DOM stub coverage across diverse modern websites; remaining errors are primarily syntax errors and complex framework patterns, which is expected for a terminal browser.

## Content Detection & Dynamic Loading

**Status**: ✅ COMPLETED

Implemented comprehensive content detection and dynamic loading capabilities to handle modern websites that show loading screens, cookie banners, or load content asynchronously.

### Features

- **Content Analysis Engine**: Analyzes page content to determine loading state, detect banners, and validate content completeness
- **Smart Retry Logic**: Automatically retries page loading with appropriate wait times when content appears incomplete
- **Site-Specific Handlers**: Specialized processing for different websites (CodePen, DerStandard, SPAs)
- **Loading State Detection**: Recognizes common loading indicators and interstitial pages
- **Banner Detection**: Identifies cookie consent and adblock detection banners
- **User Feedback**: Provides detailed analysis results to help users understand content state
- **Configurable Timeouts**: Customizable retry limits and wait times to prevent hanging

### Implementation Details

- **ContentDetector**: Core analysis engine that examines HTML content for loading indicators, banners, and content quality
- **SiteHandlerManager**: Manages site-specific processing rules and handlers
- **Enhanced Client**: Browser client with integrated retry logic and content validation
- **Analysis Reporting**: User-friendly display of content analysis results with emoji indicators

### Site-Specific Handlers

- **CodePen Handler**: Extracts pen titles, descriptions, and code snippets with 3-second retry logic
- **DerStandard Handler**: Detects adblock banners and extracts article content
- **Generic SPA Handler**: Handles single-page applications with dynamic content loading

### Content Detection Patterns

- **Loading Indicators**: "just a moment", "loading", "please wait", "checking your browser"
- **Cookie Banners**: "accept cookies", "cookie policy", "we use cookies", "privacy policy"
- **AdBlock Banners**: "disable adblock", "ad blocker detected", "whitelist this site"
- **Security Checks**: "cloudflare", "ddos protection", "security check"

### Usage

```bash
# Default mode with content detection and retry logic
./brauser https://example.com

# Disable content detection for faster loading
./brauser https://example.com --no-retry
```

### Testing Results

- **DerStandard.at**: Successfully detects content without false positives for adblock banners
- **CodePen.io**: Properly handles dynamic loading with 3-second retries
- **Various Sites**: Accurate detection of loading states and content completeness

### Development Guidelines

- Implement pattern-based detection for common loading indicators and banners
- Use site-specific handlers for websites with unique content structures
- Provide clear user feedback about content analysis and retry attempts
- Balance retry logic to avoid excessive delays while ensuring content completeness
- Design modular handlers that can be easily extended for new sites

## Visual Enhancements & GZIP Support Implementation

**Status**: ✅ COMPLETED

Implemented comprehensive visual enhancements and GZIP support to improve content display and handle modern web compression.

### Features

- **Enhanced HTML Renderer**: Beautiful terminal formatting with emojis and visual separators
- **Hierarchical Content Display**: Different styling for H1, H2, and other headings
- **Structured Content Sections**: Organized display of navigation, links, lists, and main content
- **Content Summary Statistics**: Shows counts of headings, paragraphs, and links
- **GZIP Decompression**: Proper handling of compressed HTTP responses
- **Improved ASCII Art**: Better image rendering with error handling
- **Visual Boundaries**: Clear content separation with decorative borders

### Implementation Details

- **Enhanced HTMLRenderer**: Updated `renderer/html.go` with comprehensive content extraction
- **GZIP Support**: Added compression handling in `browser/client.go`
- **Content Categorization**: Separate handling for titles, headings, paragraphs, navigation, links, lists, and images
- **Visual Formatting**: Emoji-based indicators and structured layout
- **Error Handling**: Graceful fallbacks for failed image conversions

### Content Display Features

1. **Title Display**: Prominent title with decorative underline
2. **Hierarchical Headings**: H1 with double-line underline, H2 with single-line, others with bullet points
3. **Navigation Section**: Dedicated section for menu and navigation links
4. **Main Content**: Extraction and display of primary content areas
5. **Links Section**: Organized display of important page links
6. **List Items**: Structured display of bulleted and numbered lists
7. **Image Gallery**: ASCII art conversion with metadata display
8. **Content Summary**: Statistical overview of page elements

### Visual Output Example

```
============================================================
           BRAUSER - TERMINAL WEB CONTENT
============================================================

📄 TITLE: Example Domain
------------------------

🔸 Example Domain
==============

This domain is for use in illustrative examples...

🧭 NAVIGATION:
  • Home (/)
  • About (/about)

🔗 LINKS:
  → More information... (https://www.iana.org/domains/example)

🖼️  IMAGES:
  Image 1: logo.png (alt: Company Logo)
    ASCII Art:
    ████████
    ██    ██
    ████████

============================================================
📊 CONTENT SUMMARY: 1 headings, 2 paragraphs, 1 links
============================================================
```

### Testing Results

- **Example.com**: Clean display of simple content structure
- **Hacker News**: Proper extraction of story titles and links
- **GitHub**: Complex page with headings, navigation, and content sections
- **StackOverflow**: Handling of dynamic content and cookie banners

### Development Guidelines

- Always handle GZIP compression in HTTP responses
- Use emoji indicators for different content types
- Limit content display to avoid terminal spam (e.g., max 15 links, 5 images)
- Provide graceful fallbacks for failed operations
- Structure output with clear visual boundaries
- Include summary statistics for user awareness
