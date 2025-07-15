# Brauser 🌐

```
 ____  ____    __    __  __  ___  ____  ____ 
(  _ \(  _ \  /__\  (  )(  )/ __)( ___)(  _ \
 ) _ < )   / /(__)\  )(__)( \__ \ )__)  )   /
(____/(_)\_)(__)(__)(______)(___/(____)(_)\_)
```

> **The next-generation terminal web browser built for developers, power users, and AI agents**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Terminal](https://img.shields.io/badge/Platform-Terminal-green.svg)](#)
[![AI Ready](https://img.shields.io/badge/AI-Ready-purple.svg)](#)

## 🚀 Why Brauser?

Brauser isn't just another terminal browser—it's a **modern, JavaScript-capable web client** designed for the era of AI agents and terminal-first workflows. While traditional terminal browsers struggle with modern web content, Brauser bridges the gap between terminal efficiency and web compatibility.

### ✨ Key Features

- 🧠 **AI Agent Ready**: Perfect for LLM-driven web automation and content extraction
- ⚡ **Modern JS Support**: Handles dynamic content with sandboxed JavaScript execution
- 🎨 **Smart Rendering**: ASCII art images, structured content display, and intelligent formatting
- 🧭 **Interactive Navigation**: Browser-like history, numbered link selection, and intuitive commands
- 🔍 **Content Intelligence**: Advanced content detection, loading state recognition, and retry mechanisms
- 🏗️ **Modular Architecture**: Clean, testable Go codebase with separated concerns
- 🌐 **Real-World Ready**: Handles GZIP compression, relative URLs, and complex modern websites

## 🎯 Perfect For

- **🤖 AI Agents & Automation**: Programmatic web browsing for LLMs and autonomous systems
- **👨‍💻 Terminal Power Users**: Efficient web browsing without leaving your terminal workflow
- **🔧 DevOps & SysAdmins**: Quick web content inspection and monitoring
- **📊 Data Scientists**: Web scraping and content analysis in terminal environments
- **🚀 CI/CD Pipelines**: Automated web testing and content validation

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/brauser.git
cd brauser

# Build and run
go build .
./brauser https://example.com

# Or run directly
go run . https://news.ycombinator.com
```

## 🎮 Quick Start

```bash
# Start browsing
./brauser https://github.com

# Interactive commands:
# [1-50]     - Follow numbered links
# b/back     - Navigate back
# f/forward  - Navigate forward  
# h/history  - View browsing history
# u/url      - Enter new URL
# r/refresh  - Reload current page
# q/quit     - Exit
```

## 🏗️ Architecture

Brauser features a **clean, modular architecture** designed for maintainability and extensibility:

```
brauser/
├── browser/     # HTTP client & content detection
├── js/          # JavaScript execution environment
├── renderer/    # HTML & ASCII image rendering
├── navigation/  # Interactive navigation system
└── config/      # Configuration management
```

### 🧩 Core Components

- **🌐 Smart HTTP Client**: GZIP support, timeout handling, and robust error recovery
- **🔧 JS Execution Engine**: Sandboxed JavaScript with comprehensive DOM stubs
- **🎨 Advanced Renderer**: Structured HTML display with ASCII art image conversion
- **🧭 Navigation System**: Browser-like history, link extraction, and user interaction
- **🔍 Content Detector**: Intelligent loading state recognition and retry logic

## 🤖 AI Agent Integration

Brauser is specifically designed for **AI agent workflows**:

```go
// Example: Programmatic browsing for AI agents
navigator := navigation.NewNavigator()
client := browser.NewClient()

// Navigate and extract structured content
content, err := client.FetchPage("https://example.com")
links := navigator.ExtractLinks(content)

// Perfect for LLM content processing
for _, link := range links {
    fmt.Printf("Link %d: %s -> %s\n", link.Number, link.Text, link.URL)
}
```

## 🌟 What Makes It Special

### 🚀 Modern Web Compatibility
Unlike traditional terminal browsers, Brauser handles:
- Dynamic JavaScript content
- Modern web frameworks
- AJAX-loaded content
- Complex CSS layouts

### 🧠 Intelligent Content Processing
- **Loading State Detection**: Recognizes when pages are still loading
- **Content Validation**: Distinguishes between actual content and loading screens
- **Smart Retry Logic**: Automatically retries for dynamically loaded content
- **Link Categorization**: Organizes links by type (navigation, content, stories)

### 🎨 Terminal-Optimized Display
- **ASCII Art Images**: Converts images to beautiful terminal art
- **Structured Output**: Clean, hierarchical content display
- **Compressed Formatting**: Intelligent whitespace management
- **Visual Indicators**: Emojis and separators for better readability

## 🔧 Configuration

Customize JavaScript compatibility and rendering behavior:

```json
{
  "timeout_ms": 5000,
  "enable_dom_stubs": true,
  "enable_console_stubs": true,
  "enable_browser_stubs": true,
  "max_execution_time_ms": 3000
}
```

## 🧪 Testing

Brauser has been tested on diverse websites:
- 📰 News sites (ORF, Der Standard, TechCrunch)
- 💬 Social platforms (Hacker News)
- 🛠️ Developer tools (CodePen)
- 🏢 Corporate sites (CNN)

## 🤝 Contributing

We welcome contributions! Brauser follows clean architecture principles:

1. **Modular Design**: Each package has a single responsibility
2. **Comprehensive Testing**: Critical functionality is well-tested
3. **Error Handling**: Graceful degradation for robust operation
4. **Documentation**: Clear code comments and architectural decisions

## 📈 Roadmap

- [ ] **Plugin System**: Extensible architecture for custom handlers
- [ ] **API Mode**: RESTful API for programmatic access
- [ ] **Enhanced AI Integration**: Built-in LLM content summarization
- [ ] **Performance Optimization**: Caching and parallel processing
- [ ] **Mobile Responsive**: Better handling of mobile-first websites

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Built with ❤️ for the terminal-first future**

*Brauser: Where terminal efficiency meets modern web compatibility*
