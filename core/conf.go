// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

var ErrNoFile = errors.New("config: No file specified")

type ErrConfNotFound struct {
	Section string
}

func (e ErrConfNotFound) Error() string {
	return "config: section '" + e.Section + "' not found"
}

type Config struct {
	path     string
	modified bool
	sections map[string]*json.RawMessage
}

func NewConfig(file string) *Config {
	return &Config{
		file,
		false,
		make(map[string]*json.RawMessage),
	}
}

func (cfg *Config) Get(section string, v interface{}) error {
	raw, ok := cfg.sections[section]
	if !ok {
		cfg.Set(section, v)
		return nil
	}
	return json.Unmarshal(*raw, v)
}

func (cfg *Config) Exists(section string) bool {
	_, ok := cfg.sections[section]
	return ok
}

func (cfg *Config) Set(section string, v interface{}) error {
	raw, err := json.MarshalIndent(&v, "", "\t")
	if err != nil {
		return err
	}
	rawjson := json.RawMessage(raw)
	cfg.sections[section] = &rawjson
	cfg.modified = true
	return nil
}

func (cfg *Config) IsModified() bool {
	return cfg.modified
}

func (cfg *Config) Write() error {
	if cfg.path == "" {
		return ErrNoFile
	}
	if !cfg.IsModified() {
		return nil
	}

	if err := os.MkdirAll(path.Dir(cfg.path), 0700); err != nil {
		return err
	}

	f, err := os.Create(cfg.path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoded, err := json.MarshalIndent(&cfg.sections, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(encoded)
	return err
}

func (cfg *Config) Path() string {
	return cfg.path
}

func OpenConfig(file string, create bool) (*Config, error) {
	cfg := NewConfig(file)
	f, err := os.Open(file)
	if err != nil {
		if !create || !os.IsNotExist(err) {
			return cfg, err
		}
		if err := cfg.Write(); err != nil {
			return cfg, err
		}
		return cfg, nil
	}

	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&cfg.sections)
	return cfg, err
}

func GetDefaultDir(name string) string {
	path := os.Getenv("XDG_CONFIG_HOME")
	if path != "" {
		return path + "/" + name
	}

	home := os.Getenv("HOME")
	if home != "" {
		return home + "/.config/" + name
	}

	return "."
}
