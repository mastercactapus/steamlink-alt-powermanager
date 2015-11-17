package main

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

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
		conn: conn,
	}
	return pm
}

func (pm *PowerManager) Sleep() {
	// TODO: this isn't right, need to turn off controllers
	exec.Command("/bin/poweroff", "-f").Run()
}

func (pm *PowerManager) SetActivity(state uint32) {
	// do nothing
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
