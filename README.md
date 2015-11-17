An alternative powermanager implementation while I wait for a new build for the steamlink. It is in no way associated with valve, steamlink, or anything. Just a weekend project, use at your own risk ;)

## Building

Make sure you have go 1.5+ available, clone the repo and run `./build.sh`

The build script just sets `GOARCH=arm` and `GOARM=7` and runs `go build -o powermanager`

## Installing

You need ssh access to the steamlink (not covered here). Backup the old powermanager, copy the new on in it's place, and reboot.

```bash

./build.sh
scp powermanager root@steamlink:
ssh root@steamlink

cp bin/powermanager powermanager.bak
cp powermanager bin/powermanager
reboot -f
```


## Notes

A simple quick-temporary fix; disables suspend timer

Registers `system.powermanager` alias with a single objectPath `/powermanager` and the following methods:

`system.powermanager.Sleep` in turn just runs `poweroff -f` currently
`system.powermanager.SetActivity(state uint32)` does nothing

## Debugging

To debug DBus messages, you can use `dbus-monitor --system --monitor`

I also have another repo so you can explore interfaces: github.com/mastercactapus/dbus-inspector

However, the stock powermanager does not properly implement Introspectable so you have to watch traffic

## Known Issues

- It does not turn of the steam controller automatically.

