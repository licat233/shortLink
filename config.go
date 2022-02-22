package main

type Config struct {
	Admin  *Admin   `yaml:"Admin"`
	Addr   string   `yaml:"Addr"`
	Prizes []*Prize `yaml:"Prizes"`
}
