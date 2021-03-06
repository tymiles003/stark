// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"flag"
	"os"
	"os/signal"

	"github.com/xconstruct/stark/proto"
)

var verbose = flag.Bool("v", false, "verbose debug output")
var vverbose = flag.Bool("vv", false, "very verbose debug output: db, individual messages")
var configPath = flag.String("config", "", "path to config file")

type App struct {
	AppName    string
	ModuleName string
	Config     *Config
	Log        *Logger
}

func NewApp(appName, moduleName string) *App {
	if appName == "" {
		appName = "stark"
	}
	if !flag.Parsed() {
		flag.Parse()
	}

	app := &App{
		AppName:    appName,
		ModuleName: moduleName,
		Log:        DefaultLog,
	}
	if *vverbose {
		app.Log.SetLevel(LogLevelTrace)
	} else if *verbose {
		app.Log.SetLevel(LogLevelDebug)
	} else {
		app.Log.SetLevel(LogLevelInfo)
	}

	return app
}

func (app *App) Init() {
	if err := app.initConfig(); err != nil {
		app.Log.Fatalln(err)
	}

	app.WriteConfig()
}

func (app *App) initConfig() error {
	path := *configPath
	if path == "" {
		path = GetDefaultDir(app.AppName) + "/" + app.ModuleName + ".json"
	}
	cfg, err := OpenConfig(path, true)
	if err != nil {
		return err
	}
	app.Log.Debugf("[core] reading config from '%s'", cfg.Path())
	app.Config = cfg
	return nil
}

func (app *App) WriteConfig() {
	if app.Config.IsModified() {
		app.Log.Infof("[core] writing config to '%s'", app.Config.Path())
		app.Must(app.Config.Write())
	}
}

func (app *App) Close() {
	app.WriteConfig()
}

func (app *App) Must(err error) {
	if err != nil {
		app.Log.Fatalln(err)
	}
}

func (app *App) SetupInjector(inj *Injector, name string) {
	inj.Instance(app.Config)
	inj.Instance(app.Log)
	inj.Factory(func() proto.Logger {
		return app.Log
	})
}

func (app *App) Inject(name string, container interface{}) error {
	inj := NewInjector()
	app.SetupInjector(inj, name)
	return inj.Inject(container)
}

func (app *App) Dial() proto.Conn {
	cfg := proto.NetConfig{
		Address: "tcp://localhost:" + proto.DefaultPort,
	}
	app.Config.Get("dial", &cfg)
	app.WriteConfig()
	conn, err := proto.Dial(&cfg)
	if err != nil {
		app.Log.Fatal(err)
	}
	return conn
}

func WaitUntilInterrupt() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
