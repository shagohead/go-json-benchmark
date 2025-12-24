package main

import "time"

type Event struct {
	SDK         SDK
	Platform    string
	ServerName  string
	Environment string
	Release     string
	Level       string
	Contexts    map[string]any
	Extra       map[string]any
	User        map[string]any
	Tags        map[string]string
	EventID     string
	Message     string
	Exception   Exceptions
	Timestamp   time.Time
}

type SDK struct {
	Name    string
	Version string
}

type Exceptions []Exception

type Exception struct {
	Module string
	Type   string
	Value  string
	Frames Frames
}

type Frames []Frame

type Frame struct {
	Filename string
	AbsPath  string
	Module   string
	Function string
	LineNum  int
	CtxLine  string
	PreCtx   []string
	PostCtx  []string
	Vars     map[string]any
	InApp    bool
}
