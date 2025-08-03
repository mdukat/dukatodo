# DukaTODO

Android focused TODO app, made in Go and Fyne with [Unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) in mind: Do one thing and do it well.

## Install

Grab an APK from [latest release](https://github.com/mdukat/dukatodo/releases/latest), and install it like any other Android app!

## Build

Don't trust releases? Good! Make sure you have working Fyne environment (check [Fyne - Getting started](https://docs.fyne.io/started/) and [Fyne - Mobile Packaging](https://docs.fyne.io/started/mobile)). Download Android NDK [here](https://developer.android.com/ndk/downloads#lts-downloads).

Download this repository and build APK with:

```
$ fyne package -os android/arm64 -app-id pl.mdukat.dukatodo -release -name "DukaTODO" -app-version x.x
```

