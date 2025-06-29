# Planning Poker Documentation Build System
# Requires: pdflatex, bibtex

.PHONY: all clean pdf view help

# Default target
all: pdf

# Build the PDF documentation
pdf: design.pdf

design.pdf: design.tex references.bib
	@echo "Building Planning Poker documentation..."
	pdflatex design.tex
	bibtex design
	pdflatex design.tex
	pdflatex design.tex
	@echo "Documentation built successfully: design.pdf"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f *.aux *.bbl *.blg *.log *.out *.toc *.fdb_latexmk *.fls *.synctex.gz
	@echo "Build artifacts cleaned"

# Clean everything including PDF
distclean: clean
	rm -f design.pdf
	@echo "All generated files removed"

# View the PDF (macOS)
view: design.pdf
	open design.pdf

# Check for LaTeX installation
check:
	@echo "Checking LaTeX installation..."
	@command -v pdflatex >/dev/null 2>&1 || { echo >&2 "pdflatex not found. Please install LaTeX."; exit 1; }
	@command -v bibtex >/dev/null 2>&1 || { echo >&2 "bibtex not found. Please install LaTeX."; exit 1; }
	@echo "LaTeX installation found"

# Build and view in one command
build-view: pdf view

# Quick build without bibliography (for fast iteration)
quick: design.tex
	pdflatex design.tex

# Help target
help:
	@echo "Planning Poker Documentation Build System"
	@echo "========================================"
	@echo ""
	@echo "Available targets:"
	@echo "  pdf        Build the complete PDF documentation (default)"
	@echo "  clean      Remove build artifacts (keep PDF)"
	@echo "  distclean  Remove all generated files including PDF"
	@echo "  view       Open the PDF in default viewer (macOS)"
	@echo "  check      Verify LaTeX installation"
	@echo "  quick      Fast build without bibliography processing"
	@echo "  build-view Build and immediately view the PDF"
	@echo "  help       Show this help message"
	@echo ""
	@echo "Requirements:"
	@echo "  - pdflatex (part of LaTeX distribution)"
	@echo "  - bibtex   (part of LaTeX distribution)"
	@echo ""
	@echo "To install LaTeX on macOS:"
	@echo "  brew install --cask mactex"
	@echo ""
	@echo "To install LaTeX on Ubuntu/Debian:"
	@echo "  sudo apt-get install texlive-full"
