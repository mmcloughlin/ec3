#!/bin/bash -ex

useragent='github.com/mmcloughlin/ec3'
url='https://hyperelliptic.org/EFD'
archive='efd.tar.gz'

# Download in a temporary directory.

cwd=$(pwd -P)
workdir=$(mktemp -d)
cd ${workdir}

# Download.
wget \
    --user-agent="${useragent}" \
    --recursive \
    --no-parent \
    --limit-rate=10k \
    --no-host-directories --cut-dirs=1 \
    ${url}

# Clean up files we don't need.
rm -rf EFD oldefd
find . -name '*.html' | xargs rm
rm -rf g1*/auto-sage

# Leave a note.
cat > README <<EOF
Downloaded from ${url}.

The Explicit-Formulas Database is joint work by Daniel J. Bernstein and Tanja
Lange, building on work by many authors.

http://cr.yp.to/djb.html
http://www.hyperelliptic.org/tanja/
https://hyperelliptic.org/EFD/bib.html
EOF

# Archive.
COPYFILE_DISABLE=1 tar czf ${archive} *
mv ${archive} ${cwd}
