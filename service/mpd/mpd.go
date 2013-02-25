package mpd

import (
	"os/exec"
	"strings"

	"github.com/xconstruct/stark"
	"github.com/xconstruct/stark/service"
)

type MPD struct {
	*service.Service
}

func New() *MPD {
	serv := service.New(service.Info{
		Name: "mpd",
		Actions: []string{
			"music.play",
			"music.pause",
			"music.stop",
			"music.prev",
			"music.next",
		},
	})
	m := &MPD{serv}
	serv.Handler = m
	return m
}

func (m *MPD) Handle(msg *stark.Message) (*stark.Message, error) {
	action := strings.Split(msg.Action, ".")
	if action[0] != "music" {
		return nil, nil
	}

	exec.Command("mpc", action[1]).Start()
	reply := stark.NewReply(msg)
	reply.Action = "notify.success"
	reply.Message = "done"
	return reply, nil
}
