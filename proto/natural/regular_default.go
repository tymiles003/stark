// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package natural

import "github.com/xconstruct/stark/proto"

const defaultRegularText = `[
{
	"example": "I started to [play] [music]",
	"msg": {
		"action": "event/new",
		"p": {
			"subject": ".I",
			"verb": "{{.play}}",
			"object": "{{.music}}",
			"status": "started"
		}
	}
},
{
	"example": "I finished to [play] [music]",
	"msg": {
		"action": "event/new",
		"p": {
			"subject": "I",
			"verb": "{{.play}}",
			"object": "{{.music}}",
			"status": "ended"
		}
	}
},
{
	"example": "I [drink] [coffee]",
	"msg": {
		"action": "event/new",
		"p": {
			"subject": "I",
			"verb": "{{.drink}}",
			"object": "{{.coffee}}",
			"status": "singular"
		}
	}
},
{
	"example": "When did I last visit [Big City]",
	"msg": {
		"action": "location/last",
		"p": {
			"address": "{{.BigCity}}"
		}
	}
},
{
	"example": "When did [I] last [drink] [coffee]",
	"msg": {
		"action": "event/last",
		"p": {
			"subject": "{{.I}}",
			"verb": "{{.drink}}",
			"object": "{{.coffee}}"
		}
	}
},
{
	"example": "Push [this long text] to [phone]",
	"msg": {
		"action": "push/text",
		"p": {
			"device": "{{.phone}}"
		}
	}
},
{
	"example": "Remind me in [some duration] to [make coffee]",
	"msg": {
		"action": "schedule/duration",
		"p": {
			"duration": "{{.someduration}}"
		},
		"text": "{{.makecoffee}}"
	}
},
{
	"example": "Remind me to [make coffee] in [duration]",
	"msg": {
		"action": "schedule/duration",
		"p": {
			"duration": "{{.someduration}}"
		},
		"text": "{{.makecoffee}}"
	}
},
{
	"example": "Remind me in [some duration] that [something is happening]",
	"msg": {
		"action": "schedule/duration",
		"p": {
			"duration": "{{.someduration}}"
		},
		"text": "{{.somethingishappening}}"
	}
},
{
	"example": "Remind me in [some duration]",
	"msg": {
		"action": "schedule/duration",
		"p": {
			"duration": "{{.someduration}}"
		}
	}
},
{
	"example": "Create a geofence named [home] at [Baker Street 221b]",
	"msg": {
		"action": "location/fence/create",
		"p": {
			"name": "{{.home}}",
			"address": "{{.BakerStreet221b}}"
		}
	}
},
{
	"example": "Create a geofence at [friends house]",
	"msg": {
		"action": "location/fence/create",
		"p": {
			"address": "{{.friendshouse}}"
		}
	}
},
{
	"example": "What is [the birth day of Tuomas Holopainen]",
	"msg": {
		"action": "knowledge/query",
		"text": "{{.thebirthdayoftuomasholopainen}}"
	}
}
]`

var defaultRegular RegularSchemata

func ParseRegular(text string) (proto.Message, bool) {
	if defaultRegular == nil {
		var err error
		defaultRegular, err = LoadRegularSchemata(defaultRegularText)
		if err != nil {
			panic(err)
		}
	}

	return defaultRegular.Parse(text)
}
