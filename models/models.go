package models

const (
	QuietQueryType    PacketType = 1
	QuietResponseType PacketType = 2
)

type PacketType uint16

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

type QuietResponse struct {
	Preamble    *Preamble
	IsQuietTime bool
	WakeUpHour  uint
	Whoru       string
}
