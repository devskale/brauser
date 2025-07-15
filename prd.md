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
  - [ ] Add comprehensive unit tests
  - [ ] Optimize performance bottlenecks
  - [ ] Improve memory management
  7. [x] **CONTENT DETECTION & DYNAMIC LOADING IMPROVEMENTS**
  - [x] Enhanced Content Detection: Detect when sites show loading pages vs actual content
  - [x] Wait Mechanisms: Add optional delays for sites that load content dynamically
  - [x] Content Validation: Verify if extracted content represents actual page or loading screen
  - [x] Site-Specific Handling: Add CodePen-specific handling for better content processing
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
  9. Implement navigation (links, back/forward) and user input.
  10. Test on MacOS for cross-platform compatibility.
  11. Optimize for performance and minimalism.

## Project Structure

The brauser project is organized into focused packages for maintainability:

```
brauser/
├── main.go              # Entry point and CLI handling
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

### Key Learnings

1. **JavaScript Compatibility**: Many websites use modern JS features that need to be stubbed out for compatibility
2. **Error Handling**: Graceful degradation is crucial - websites should still be usable even when JS fails
3. **Image Processing**: ASCII art conversion adds visual appeal while maintaining terminal compatibility
4. **Site-Specific Handling**: Different websites have unique requirements (e.g., HackerNews story parsing)
5. **Performance**: Timeout mechanisms prevent hanging on problematic scripts
6. **Code Organization**: Refactoring monolithic code into packages improves maintainability and testability
7. **Go Module System**: Proper import paths are crucial - relative imports don't work in module mode
8. **Struct Definition**: Complex nested structs require careful attention to closing braces and field access patterns
9. **Content Detection**: Modern websites often show loading screens, cookie banners, or adblock messages before actual content
10. **Dynamic Loading**: Many sites load content asynchronously, requiring retry mechanisms with appropriate wait times
11. **Site-Specific Patterns**: Different sites have unique loading patterns and content structures that benefit from specialized handling
12. **User Feedback**: Clear visual indicators help users understand what's happening during page loading and content analysis.

13. **HTTP Compression**: Modern websites use GZIP compression extensively - proper decompression is essential for content parsing.

14. **Visual Terminal Design**: Well-structured terminal output with emojis, separators, and hierarchical formatting significantly improves readability.

15. **Content Extraction Strategy**: Different content types (headings, paragraphs, navigation, lists) require different extraction and display strategies.

16. **ASCII Art Limitations**: SVG and complex image formats often fail ASCII conversion - graceful error handling is important.

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
