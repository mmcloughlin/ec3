all: references.md

%.md: %.yml
	refs gen -type markdown -out $@ $<
