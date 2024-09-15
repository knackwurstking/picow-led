# PicoW LED

## TODOs

- [x] Test macos build
- [ ] Test android build
- [ ] Test ios build

## Build for android

Platform **MacOS**:

```bash
cd cmd/picow-led-gui
export ANDROID_NDK_HOME=$HOME/Library/Android/sdk/ndk/26.2.11394342
fyne package -os android -icon ./icon.png -appID picowled.knackwurstking.com -release -name "PicoW LED"
```
