# Minimal makefile for Sphinx documentation
#

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?=
VENVDIR       = _venv
BINDIR        = $(VENVDIR)/bin
SPHINXBUILD   = $(BINDIR)/sphinx-build
SOURCEDIR     = .
BUILDDIR      = _build

# Put it first so that "make" without argument is like "make help".
help: $(SPHINXBUILD)
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

$(SPHINXBUILD): $(VENVDIR)
	$(VENVDIR)/bin/pip install -r requirements.txt

$(VENVDIR):
	python3 -m venv $(VENVDIR)

.PHONY: help Makefile

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
