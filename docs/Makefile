all: lint references.md references.bib

%.md: %.yml
	refs gen -type markdown -out $@ $<

%.bib: %.yml
	refs bib -out $@ $<

.PHONY: lint
lint:
	refs lint references.yml
