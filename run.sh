#!/bin/bash -e

echo "build templates"

/www/projects/weasel2/bin/build-html  -src="/www/projects/weasel2/templates/src" -dst="/www/projects/weasel2/templates/html"
rm -f /www/projects/weasel2/templates/html/layout.html

echo "generate bindata"

mkdir --parents ./src/app/bindata/assets

go-bindata -nomemcopy -prefix "templates/html" -o ./src/app/bindata/templates/a.go ./templates/html/...
go-bindata-assetfs -o ./src/app/bindata/assets/a.go ./assets/...

grep 'package main' -P -R -I -l  ./src/app/bindata/templates/* | xargs sed -i 's/package main/package templates/g'
grep 'package main' -P -R -I -l  ./src/app/bindata/assets/* | xargs sed -i 's/package main/package assets/g'
grep 'func assetFS' -P -R -I -l  ./src/app/bindata/assets/* | xargs sed -i 's/func assetFS/func AssetFS/g'

echo "build app"

#env GOPATH=$GOPATH:/srv/src/weasel2/ go build -race -v -o bin/eve-industry

gb build

echo "run"

env CONFIG="/www/projects/weasel2/conf.d" GODEBUG=gctrace=1 bin/server -host 127.0.0.1 -port 8082 -withbinstatic
