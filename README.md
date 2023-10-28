# README

## About

`klc` is a small utility which, when executed, alters the keyboard layout based on a
predefined locale rotation. It does so by wrapping the `localectl` & `setxkbmap` binaries
found on various Linux distributions.

## Why I wrote it?

In my day to day life, I frequently need to swap between keyboard layouts (Greek, English, etc.).
However, Xubuntu's handling of keypresses is unwilling to abide by my very simple rules:

* Alt+Ctrl+Shift+Arrow: moves windows to different workspaces
* Alt+Shift: should rotate the current keyboard layout to a different locale.

The reason Xubuntu won't do this very simple thing through it's builtin layout switcher is probably
because instead of handling events on key press, they handle events on key down and this creates 
collisions as multiple shortcuts which use similar keys start to race against each other.

It's either that or whatever.

## Usage

`klc us gr fr`: will rotate they keyboard layout (when executed) respecting the supplied rotation of 
locales.

## Fixing Xubuntu's layout switcher

* Remove any shortcut you may have applied to the `Layout` section of the `Keyboard` settings' 
applet.
* Add a shortcut to `klc` in the `Application Shortcuts` tab of the same `Keyboard` settings' applet
appending your locales and assign a shortcut to it.
* Profit.
