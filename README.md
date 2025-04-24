# Templ Echo PWA Starter

Project Structure

- `./components/` — templ components
- `./public/` — Web director containing all the content which will be copied
  to the "./dist" directory
- `./web/` — Packages for generating & serving files (uses the
  "./internal/web" package)
- `./main.go` — Main Go application (commands: server, generate)
- `Makefile` — The makefile makes your project :)
- `pwa-assets.config.js` — PWA assets configuration for the
  `make generate-pwa-assets`, generate the icons

> Yes in `./web/`, I'm using golangs "html/template" even for javascript files :)

Install Templ:

```build
go install -u github.com/a-h/templ/cmd/templ@latest
```

Clone & Initialize:

```bash
git clone https://github.com/knackwurstking/picow-led
cd picow-led
make init
```

Build executable to `./bin/`

```bash
make build
```

Generate HTML to `./dist/index.html`:

```bash
make generate-dist
```

Or use as starter with `degit` command (This will clone the repo without
the git folder)

```bash
mkdir <project-name>
cd <project-name>
#git init
#git remote add origin <url>
npx degit knackwurstking/picow-led .
```

## TODO:

- [x] Update `./public/manifest.json`
- [x] Update `./public/service-worker.js`
    - [ ] ~Use vite to bundle CSS/TypeScript and JavaScript~
        - [ ] ~Use vite to bundle manifest.json and service-worker.js~
        - [ ] ~Update the makefie, don't forget to set the env vars needed~
            - [ ] ~server path prefix~
            - [ ] ~version~
    - [x] I could also just use typescript with the tsc command
    - [x] In web/js and web/pwa is way to much copy/paste stuff, get rid of that
- [x] Rename this project to something like this: `picow-led` (?)
- [ ] ~Should i really include this go.sum file here?~
- [x] Try to get rid of some errors when offline (service-worker, manifest.json)
- [x] Finish this project with version v1.0.0 for now
- [ ] Create a script for rename "picow-led", ...
