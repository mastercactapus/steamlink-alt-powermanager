package main

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

const warnTimeout = 14 * time.Minute
const sleepTimeout = 15 * time.Minute

const intro = `
<node>
	<interface name="system.powermanager">
		<method name="Sleep" />
		<method name="SetActivity">
			<arg name="state" type="u" direction="in" />
		</method>
		<signal name="UpdateSuspendCountdown">
			<arg name="seconds" type="u" />
		</signal>
		<signal name="CancelSuspendCountdown" />
	</interface>` + introspect.IntrospectDataString + `</node>`

const introRoot = `<node><node name="powermanager" />` + introspect.IntrospectDataString + `</node>`

type PowerManager struct {
	activityCh chan bool
	conn       *dbus.Conn
}

func NewPowerManager(conn *dbus.Conn) *PowerManager {
	pm := &PowerManager{
		activityCh: make(chan bool), //no buffer, method_return will only happen after it's actually processed
		conn:       conn,
	}
	go pm.loop()
	return pm
}

func (pm *PowerManager) loop() {
	t := time.NewTicker(time.Millisecond * 200)
	act := false
	hasWarned := false
	lastAct := time.Now()
	for {
		select {
		case setAct := <-pm.activityCh:
			if hasWarned {
				pm.conn.Emit("/powermanager", "system.powermanager.CancelSuspendCountdown")
			}
			hasWarned = false
			act = setAct
			lastAct = time.Now()
		case <-t.C:
			if act {
				lastAct = time.Now()
				break
			}

			elapsed := time.Since(lastAct)

			if elapsed > warnTimeout {
				hasWarned = true
				seconds := (sleepTimeout - elapsed) / time.Second
				if seconds < 0 {
					seconds = 0
				}
				pm.conn.Emit("/powermanager", "system.powermanager.UpdateSuspendCountdown", uint32(seconds))
			}

			if elapsed > sleepTimeout {
				pm.Sleep()
			}

		}
	}
}

func (pm *PowerManager) Sleep() {
	// TODO: this isn't right, need to turn off controllers
	exec.Command("/bin/poweroff", "-f").Run()
}

func (pm *PowerManager) SetActivity(state uint32) {
	pm.activityCh <- state == 1
}

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalln(err)
	}
	reply, err := conn.RequestName("system.powermanager", dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalln(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalln("system.powermanager already registered")
	}
	log.Infoln(conn.Names())
	pm := NewPowerManager(conn)
	err = conn.Export(&pm, "/powermanager", "system.powermanager")
	if err != nil {
		log.Fatalln("export powermanager:", err)
	}
	err = conn.Export(introspect.Introspectable(introRoot), "/", "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Fatalln("export powermanager (introspectable-root):", err)
	}
	err = conn.Export(introspect.Introspectable(intro), "/powermanager", "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Fatalln("export powermanager (introspectable):", err)
	}

	log.Fatalln(exec.Command("/home/steam/app_run.sh").Run())
}
