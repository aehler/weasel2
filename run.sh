#!/bin/bash -e

echo "build templates"

/www/build-html  -i="/srv/src/monitor/templates/src" -o="/srv/src/monitor/templates/html" -c="" -s=false
rm -f /srv/src/monitor/templates/html/layout.html

echo "generate bindata"

go-bindata -nomemcopy -prefix "templates/html" -o ./src/app/bindata/templates/a.go ./templates/html/...
go-bindata-assetfs ./assets/...

mkdir --parents ./src/app/bindata/assets

mv -f bindata_assetfs.go ./src/app/bindata/assets/a.go

grep 'package main' -P -R -I -l  ./src/app/bindata/templates/* | xargs sed -i 's/package main/package templates/g'
grep 'package main' -P -R -I -l  ./src/app/bindata/assets/* | xargs sed -i 's/package main/package assets/g'
grep 'func assetFS' -P -R -I -l  ./src/app/bindata/assets/* | xargs sed -i 's/func assetFS/func AssetFS/g'

echo "build app"

env GOPATH=$GOPATH:/srv/src/monitor/ go build -race -v
mv -f monitor bin/

echo "run"

sudo env CONFIG="/srv/src/monitor/conf.d" GODEBUG=gctrace=1 bin/monitor -port 808 -withbinstatic
