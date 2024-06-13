package model

type EnvConfigModel struct{
	Server struct {
		Mode       string
		Port       string
	}
	Workspace struct {
		Cache   string
		Key     string
	}
	Performance struct {
		MaxCpuCore int
		MaxMemory  int
		TaskLimit  int
	}
}