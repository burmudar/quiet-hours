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
	*Preamble
	Whoami string
}

type QuietReponse struct {
	*Preamble
	IsQuietTime bool
	WakeUpHour  uint16
}
