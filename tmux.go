package main

type TmuxSession struct {
	Name          string
	DefaultWindow string
	Windows       []TmuxWindow
}

type TmuxWindow struct {
	Title   string
	Path    string
	Command string
}
