package main

type Packet[T any] struct {
	Preamble Preamble
	Frame    T
}

type Preamble struct {
	Version    uint
	PacketType PacketType
}

type QuietQuery struct {
	Preamble *Preamble
	Whoami   string
}

type QuietReponse struct {
	Preamble    *Preamble
	IsQuietTime bool
	WakeUpHour  uint
	Whoru       string
}
