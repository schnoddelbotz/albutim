# albutim
Yet another, simple web photo album generator and server

## features

album generator ...

- generates content suitable for static web serving
- preserves your directory structure for album
- builds thumbs/previews without requirement for external tools
- reads EXIF data
- supports built-in and custom templates

server ...

- can serve your album instead of apache/nginx/whatever
- supports basic auth and https
- builds tumbs on demand, caches on disk
- integrates nicely with docker

web ui ...

- full keyboard navigation
- fullscreen mode

## todo

- [ ] add zap all text fullscreenPro option
- [ ] add search/filter
- [ ] /index.html and /albumdata.js in <- there .data
- [ ] build cmd does not overwrite existing assets!
- [ ] originals/ links broken in static album builds
- [ ] tests...?!!!
- [ ] gui help / kbd shortcut info
- [ ] Dockerfile, allows thumbs/preview/assets mounts!
- [ ] cmdline/config options:
  - [ ] serve --no-thumbs, --no-previews, --no-cache
  - [ ] foo
- [ ] mobile view/test?
- [ ] make templates exportable (doc-only? build does...)
- [ ] enable custom templates
- [ ] autoplay feature (dropdown: whole album, current folder, folder+subfolders; delay in s)
- [x] windows: fwd vs back-slashes in paths -> thumbs ok, but previews not... test ec2/public! https://stackoverflow.com/questions/9371031/how-do-i-create-crossplatform-file-paths-in-go
- [x] windows: use https://github.com/inconshreveable/mousetrap ? ~~or is~~ cobra is already ~~?~~.
- [ ] electron: https://hackernoon.com/how-to-add-a-gui-to-your-golang-app-in-5-easy-steps-c25c99d4d8e0
- [ ] rename to zuaberalbum.com