package network

import (
	"encoding/binary"
	"fmt"

	"github.com/oxtoacart/bpool"
)

type PacketId byte

type Packet struct {
	id       PacketId
	buffer   []byte
	position int
	pooled   bool
}

const packetBufferSize int = 256

var packetBufferPool = bpool.NewBytePool(4*1024*1024, packetBufferSize)

func GetPacket(buffer []byte) Packet {
	packet := Packet{
		buffer:   buffer,
		position: 3, // skip length and version
		pooled:   false,
	}
	packet.id = PacketId(packet.ReadU8())
	return packet
}

func NewPacket(id PacketId) Packet {
	buffer := packetBufferPool.Get()
	packet := Packet{
		buffer:   buffer,
		position: 2, // skip length
		id:       id,
		pooled:   true,
	}
	packet.WriteU8(0xf0) // version
	packet.WriteU8(byte(packet.id))
	return packet
}

func (p *Packet) Free() {
	if !p.pooled {
		return
	}
	zeroFill(p.buffer)
	packetBufferPool.Put(p.buffer)
}

func zeroFill(b []byte) {
	for i := 0; i < packetBufferSize; i++ {
		b[i] = 0
	}
}

func (p *Packet) ensureBufferSize(size int) {
	if p.position+size > cap(p.buffer) {
		panic("server error: insufficient buffer capacity")
	}
}

func (p *Packet) Skip(size int) {
	p.ensureBufferSize(size)
	p.position += size
}

func (p *Packet) Data() []byte {
	binary.LittleEndian.PutUint16(p.buffer[0:2], uint16(p.position))
	return p.buffer[0:p.position]
}

func (p *Packet) ReadU8() byte {
	p.ensureBufferSize(1)
	value := p.buffer[p.position]
	p.position += 1
	return value
}

func (p *Packet) ReadU16() uint16 {
	p.ensureBufferSize(2)
	value := binary.LittleEndian.Uint16(p.buffer[p.position : p.position+2])
	p.position += 2
	return value
}

func (p *Packet) ReadU32() uint32 {
	p.ensureBufferSize(4)
	value := binary.LittleEndian.Uint32(p.buffer[p.position : p.position+4])
	p.position += 4
	return value
}

func (p *Packet) ReadU64() uint64 {
	p.ensureBufferSize(8)
	value := binary.LittleEndian.Uint64(p.buffer[p.position : p.position+8])
	p.position += 8
	return value
}

func (p *Packet) ReadStringSlice(size int) string {
	const nullChar = '\x00'
	p.ensureBufferSize(size)
	slice := p.buffer[p.position : p.position+size]
	pos := 0
	for ; pos < len(slice) && slice[pos] != nullChar; pos++ {
	}
	if pos == 0 {
		return ""
	}
	str := string(slice[0:pos])
	p.position += size
	return str
}

func (p *Packet) ReadSlice(size int) []byte {
	p.ensureBufferSize(size)
	slice := p.buffer[p.position : p.position+size]
	return slice
}

func (p *Packet) WriteU8(value byte) {
	p.ensureBufferSize(1)
	p.buffer[p.position] = value
	p.position += 1
}

func (p *Packet) WriteU16(value uint16) {
	p.ensureBufferSize(2)
	binary.LittleEndian.PutUint16(p.buffer[p.position:p.position+2], value)
	p.position += 2
}

func (p *Packet) WriteU32(value uint32) {
	p.ensureBufferSize(4)
	binary.LittleEndian.PutUint32(p.buffer[p.position:p.position+4], value)
	p.position += 4
}

func (p *Packet) WriteU64(value uint64) {
	p.ensureBufferSize(8)
	binary.LittleEndian.PutUint64(p.buffer[p.position:p.position+8], value)
	p.position += 8
}

func (p *Packet) WriteStringSlice(value string, size int) {
	p.ensureBufferSize(size)
	slice := p.buffer[p.position : p.position+size]
	copy(slice, value)
	p.position += size
}

func (p *Packet) WriteSlice(data []byte) {
	size := len(data)
	p.ensureBufferSize(size)
	slice := p.buffer[p.position : p.position+size]
	copy(slice, data)
	p.position += size
}

func (p *Packet) WriteI16(value int16) {
	p.ensureBufferSize(2)
	slice := p.buffer[p.position : p.position+2]
	slice[0] = byte(value)
	slice[1] = byte(value >> 8)
	p.position += 2
}

func (p *Packet) WriteI32(value int32) {
	p.ensureBufferSize(4)
	slice := p.buffer[p.position : p.position+4]
	slice[0] = byte(value)
	slice[1] = byte(value >> 8)
	slice[2] = byte(value >> 16)
	slice[3] = byte(value >> 24)
	p.position += 4
}

func (p *Packet) WriteI64(value int64) {
	p.ensureBufferSize(8)
	slice := p.buffer[p.position : p.position+8]
	slice[0] = byte(value)
	slice[1] = byte(value >> 8)
	slice[2] = byte(value >> 16)
	slice[3] = byte(value >> 24)
	slice[4] = byte(value >> 32)
	slice[5] = byte(value >> 40)
	slice[6] = byte(value >> 48)
	slice[7] = byte(value >> 56)
	p.position += 8
}

func (p Packet) Id() PacketId { return p.id }

func (pid PacketId) HexString() string { return fmt.Sprintf("0x%02X", pid) }
