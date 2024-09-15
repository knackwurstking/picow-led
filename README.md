# PicoW LED

## TODOs

- [x] Test build, current platform (macos)
- [ ] Test ios build
- [x] Test android build
- [x] Add and use custom fonts "Recursive"
- [x] Set theme colors from my [ui lib](https://github.com/knackwurstking/ui)
  - [x] Background
  - [x] Foreground
- [x] Change app icon resolution

## Build for android

**Current** Platform:

```bash
cd cmd/picow-led-gui
fyne package -icon ./icon.png
```

Platform **MacOS**:

```bash
cd cmd/picow-led-gui
export ANDROID_NDK_HOME=$HOME/Library/Android/sdk/ndk/26.2.11394342
fyne package -os android -icon ./icon.png -appID picowled.knackwurstking.com -release -name "PicoW LED"
```
