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

- [ ] exclude dotfiles in scanner
- [ ] cmdline/config options:
  - [ ] serve --no-thumbs, --no-previews, --no-cache
  - [ ] foo
- [ ] make templates exportable (doc-only? build does...)
- [ ] enable custom templates
- [ ] autoplay feature (dropdown: whole album, current folder, folder+subfolders; delay in s)
