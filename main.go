package main

import (
	"bytes"
	"encoding/json"
	log2 "github.com/cybriq/log"
	"github.com/cybriq/opts/version"
	"github.com/cybriq/qu"
	"os"
)

func main() {
	<-Main()
}

func Main() (quit qu.C) {
	quit = qu.T()
	go func() {
		log2.SetLogLevel("trace")
		log.T.Ln(os.Args)
		log.T.Ln(version.Get())
		var cx *state.State
		var e error
		if cx, e = state.GetNew(podcfgs.GetDefaultConfig(), podhelp.HelpFunction, quit); E.Chk(e) {
			fail()
		}

		// fail()
		// if e = debugConfig(cx.Config); E.Chk(e) {
		// }

		log.D.Ln("running command:", cx.Config.RunningCommand.Name)
		if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
			fail()
		}
		quit.Q()
	}()
	return quit
}

func fail() {
	os.Exit(1)
}

func debugConfig(c *config.Config) (e error) {
	c.ShowAll = true
	defer func() { c.ShowAll = false }()
	var j []byte
	if j, e = c.MarshalJSON(); E.Chk(e) {
		return
	}
	var b []byte
	jj := bytes.NewBuffer(b)
	if e = json.Indent(jj, j, "", "\t"); log.E.Chk(e) {
		return
	}
	log.T.Ln("\n" + jj.String())
	return
}
