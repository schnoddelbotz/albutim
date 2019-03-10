#!/bin/sh -e

make clean albutim && ./albutim build --root testalbum

echo building release
make clean ziprelease
scp albutim_* jan@zotac:htdocs/albutim
ssh jan@zotac 'htdocs/albutim/albutim_linux-amd64 --root htdocs/albutim/testalbum build'

echo pushing pics to webserver
rsync -av ~/Pictures/Tim\ Emanuel\ Hacker/ tim@zotac:public_html

echo Finally running albutim...
ssh tim@zotac '/home/jan/htdocs/albutim/albutim_linux-amd64 --root public_html build'

# cd /opt/traefik && docker-compose up -d --force-recreate www
make clean
