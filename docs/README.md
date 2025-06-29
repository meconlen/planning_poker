# Planning Poker Documentation

This directory contains the LaTeX-based technical documentation for the Planning Poker system.

## 📚 Documents

- **`design.tex`** - Main system design document covering architecture, protocols, and implementation
- **`references.bib`** - Bibliography with references to standards and academic sources
- **`Makefile`** - Build system for generating PDF documentation
- **`README.md`** - This documentation guide

## 🔧 Prerequisites

To build the documentation, you need a LaTeX distribution installed:

### macOS
```bash
# Install MacTeX (full LaTeX distribution)
brew install --cask mactex

# Or install BasicTeX (minimal) and required packages
brew install --cask basictex
sudo tlmgr update --self
sudo tlmgr install latexmk collection-fontsrecommended
```

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install texlive-full
```

### Windows
Download and install MiKTeX from https://miktex.org/

## 🏗️ Building Documentation

### Quick Start
```bash
# Build the PDF
make pdf

# View the result (macOS)
make view

# Build and view in one command
make build-view
```

### Available Commands
```bash
make pdf        # Build complete PDF with bibliography
make quick      # Fast build without bibliography (for iterations)
make clean      # Remove build artifacts (keep PDF)
make distclean  # Remove all generated files
make view       # Open PDF in default viewer
make check      # Verify LaTeX installation
make help       # Show all available commands
```

## 📄 Document Structure

The main document (`design.tex`) includes:

1. **Introduction** - Purpose, scope, and system overview
2. **System Architecture** - High-level design and technology stack
3. **Session Management** - Lifecycle and user roles
4. **WebSocket Protocol** - Complete message specification
5. **Security Considerations** - Current and planned security measures
6. **Deployment** - Container and CI/CD documentation
7. **Testing Strategy** - Test coverage and methodologies
8. **Performance** - Scalability and resource considerations
9. **Future Enhancements** - Planned features and architecture evolution

## 📖 Bibliography

The `references.bib` file contains citations for:

- **RFC Standards**: WebSocket (6455), JSON (7159), JWT (7519), HTTP (2616), URI (3986), UUID (4122)
- **Agile Methodologies**: Scrum Guide, Planning Poker research, Agile estimation
- **Technical Documentation**: Go language, Gorilla WebSocket, Docker, GitHub Actions
- **Academic Sources**: Software estimation research, microservices architecture
- **Security Standards**: OWASP Top Ten, security best practices

## 🎨 Customization

### Document Styling
The LaTeX document uses:
- Article class with professional formatting
- Syntax highlighting for code listings
- JSON-specific formatting for message examples
- Hyperlinks for cross-references and URLs
- Professional tables and figures

### Adding Content
To add new sections:
1. Edit `design.tex` and add your content
2. Add any new references to `references.bib`
3. Rebuild with `make pdf`

### Code Examples
Use the predefined listing environments:
```latex
% Go code
\begin{lstlisting}[language=go, caption=Example Go Code]
func main() {
    fmt.Println("Hello, World!")
}
\end{lstlisting}

% JSON messages
\begin{lstlisting}[language=json, caption=Example Message]
{
    "type": "example",
    "data": {
        "value": "sample"
    }
}
\end{lstlisting}
```

## 🔄 Development Workflow

For documentation updates:

1. **Edit** the LaTeX source files
2. **Test** with quick builds: `make quick`
3. **Final build** with bibliography: `make pdf`
4. **Review** the generated PDF
5. **Commit** changes to version control

### Continuous Integration

The documentation can be integrated into CI/CD pipelines:

```bash
# In GitHub Actions or similar
- name: Build Documentation
  run: |
    cd docs
    make check
    make pdf
```

## 📋 Troubleshooting

### Common Issues

**Missing packages**: If you get package errors, install the missing LaTeX packages:
```bash
# For MacTeX/BasicTeX
sudo tlmgr install <package-name>

# For Ubuntu/Debian
sudo apt-get install texlive-<collection-name>
```

**Build failures**: Clean and rebuild:
```bash
make distclean
make pdf
```

**Font issues**: Ensure you have the required font packages:
```bash
sudo tlmgr install collection-fontsrecommended
```

### Getting Help

- LaTeX documentation: https://www.latex-project.org/help/
- Package documentation: https://ctan.org/
- BibTeX help: https://www.bibtex.org/

## 📝 Contributing

When contributing to documentation:

1. Follow the existing document structure
2. Add proper citations for any claims or references
3. Use consistent formatting and styling
4. Test builds before committing
5. Update this README if adding new files or processes

The documentation follows academic standards and should maintain professional quality suitable for technical specifications and system documentation.
