# PicoW LED

- [TODOs](#todos)
- [Build](#build)
  - [**Current** Platform](#current-platform)
  - [Platform **MacOS**](#platform-android)
  - [Platform **IOS**](#platform-ios)

<a id="todos"></a>
## TODOs

- [x] Test build, current platform (macos)
- [x] Test ios build
- [x] Test android build
- [x] Add and use custom fonts "Recursive"
- [x] Set theme colors from my [ui lib](https://github.com/knackwurstking/ui)
  - [x] Background
  - [x] Foreground
- [x] Change app icon resolution

<a id="build"></a>
## Build

<a id="current-platform"></a>

**Current** Platform (tested on "MacOS"):

```bash
cd cmd/picow-led-gui
fyne package -icon ./icon.png
```

<a id="platform-android"></a>

Platform **Android**:

```bash
cd cmd/picow-led-gui
export ANDROID_NDK_HOME=$HOME/Library/Android/sdk/ndk/26.2.11394342
fyne package -os android -icon ./icon.png -appID picowled.knackwurstking.com -release -name "PicoW LED"
```

<a id="platform-ios"></a>

Platform **IOS**:

```bash
cd cmd/picow-led-gui
fyne package -os ios -icon ./icon.png -appID picowled.knackwurstking.com -release -name "PicoW LED"

mkdir Payload
mv *.app Payload/
# Compress Payload directory
zip -r PicoW-LED.ipa Payload/
```

- Upload ipa file to the iphone (ex.: via telegram)
- Connect iphone, via cable, to the macbooc, start the "AltServer"
- Install "AltStore" to iphone if needed
- Start "AltStore" app on iphone and import ipa file
