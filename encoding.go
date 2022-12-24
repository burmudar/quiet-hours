package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

const SupportedVersion = 1

type byteWriter struct {
	*bytes.Buffer
	enc [binary.MaxVarintLen64]byte
}

func (w *byteWriter) vint(n int) {
	s := binary.PutUvarint(w.enc[:], uint64(n))
	w.Write(w.enc[:s])
}

func (w *byteWriter) str(s string) {
	w.vint(len(s))
	w.WriteString(s)
}

func (w *byteWriter) strSize(s string) int {
	return binary.PutUvarint(w.enc[:], uint64(len(s))) + len(s)
}

func (b *Preamble) Marshal() ([]byte, error) {
	w := byteWriter{
		Buffer: bytes.NewBuffer(nil),
	}

	b.marshalWith(&w)
	return w.Bytes(), nil
}

func (b *Preamble) marshalWith(w *byteWriter) int {
	size := 1 // Version
	size += 1 // Type

	w.Grow(size)
	w.WriteByte(byte(b.Version))
	w.WriteByte(byte(b.PacketType))

	return w.Len()
}

func (b *Preamble) unmarshalWith(r *byteReader) (int, error) {
	version, err := r.uint()
	if err != nil {
		return r.read, fmt.Errorf("failed to read version: %s", err)
	}
	if version != SupportedVersion {
		return r.read, fmt.Errorf("unsupported version: %d", version)
	}

	typ, err := r.uint16()
	if err != nil {
		return r.read, fmt.Errorf("failed to read packet type: %s", err)
	}

	(*b) = Preamble{Version: uint(version), PacketType: PacketType(typ)}

	return r.read, nil
}

func (b *Preamble) Unmarshal(d []byte) (int, error) {
	r := NewByteReader(d)
	return b.unmarshalWith(r)
}

func (q *QuietQuery) Unmarshal(d []byte) (int, error) {
	r := NewByteReader(d)

	if q.Preamble == nil {
		q.Preamble = &Preamble{}
		if n, err := q.Preamble.unmarshalWith(r); err != nil {
			return n, err
		}
	}
	log.Printf("packet: %v", d)

	whoami, err := r.string()
	if err != nil {
		return r.read, fmt.Errorf("malformed packet failed to read whoami: %s", err)
	}

	q.Whoami = whoami
	return r.read, nil
}

func (q *QuietQuery) Type() PacketType {
	return QuietQueryType
}

func (q *QuietQuery) Marshal() ([]byte, error) {
	w := byteWriter{
		Buffer: bytes.NewBuffer(nil),
	}

	if q.Preamble == nil {
		q.Preamble = &Preamble{
			Version:    1,
			PacketType: QuietQueryType,
		}
	}

	n := q.Preamble.marshalWith(&w)

	size := n
	size += w.strSize(q.Whoami)
	w.Grow(size)

	w.str(q.Whoami)

	return w.Bytes(), nil
}

type byteReader struct {
	read int
	data []byte
}

func NewByteReader(d []byte) *byteReader {
	return &byteReader{data: d}
}

func (b *byteReader) uint16() (uint16, error) {
	v, e := b.uint()
	return uint16(v), e
}

func (b *byteReader) uint() (uint64, error) {
	v, n := binary.Uvarint(b.data)
	if n < 0 {
		return 0, fmt.Errorf("malformed packet: failed to read size")
	}
	b.read += n

	b.data = b.data[n:]
	return v, nil
}

func (b *byteReader) string() (string, error) {
	size, err := b.uint()
	if err != nil {
		return "", err
	}
	data := b.data[:size]
	b.read += int(size)
	b.data = b.data[:size]
	return string(data), nil
}

func (q *QuietReponse) Marshal() ([]byte, error) {
	w := byteWriter{
		Buffer: bytes.NewBuffer(nil),
	}

	if q.Preamble == nil {
		q.Preamble = &Preamble{
			Version:    1,
			PacketType: QuietResponseType,
		}
	}

	n := q.Preamble.marshalWith(&w)

	size := n
	size += 1 // IsQuietTime
	size += 1 // uint8
	size += w.strSize(q.Whoru)
	w.Grow(size)

	boolean := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	w.vint(boolean(q.IsQuietTime))
	w.vint(int(q.WakeUpHour))
	w.str(q.Whoru)

	return w.Bytes(), nil
}
