package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// # UnPack#############################################################################################
type UnPack struct {
	m_bin    []byte
	location int
	m_endian binary.ByteOrder
}

func NewUnPack(BigOrLittle string) *UnPack {
	var endian binary.ByteOrder
	if BigOrLittle == "little" {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}
	return &UnPack{
		m_bin:    []byte{},
		location: 0,
		m_endian: endian,
	}
}

func (up *UnPack) Empty() {
	up.m_bin = []byte{}
	up.location = 0
}

func (up *UnPack) SetData(data []byte) {
	up.m_bin = data
}

func (up *UnPack) GetAll() []byte {
	return up.m_bin[up.location:]
}

func (up *UnPack) Length() int {
	return len(up.m_bin) - up.location
}

func (up *UnPack) GetHex(length int) string {
	return hex.EncodeToString(up.GetBin(length))
}

func (up *UnPack) GetBin(length int) []byte {
	if length <= 0 {
		return []byte{}
	}
	retbin := up.m_bin[up.location : up.location+length]
	up.location += length
	return retbin
}

func (up *UnPack) GetByte() []byte {
	retbin := up.m_bin[up.location : up.location+1]
	up.location++
	return retbin
}

func (up *UnPack) GetInt() uint32 {
	retbin := up.m_bin[up.location : up.location+4]
	up.location += 4
	var ret uint32
	buf := bytes.NewReader(retbin)
	binary.Read(buf, up.m_endian, &ret)
	return ret
}

func (up *UnPack) GetShort() uint16 {
	retbin := up.m_bin[up.location : up.location+2]
	up.location += 2
	var ret uint16
	buf := bytes.NewReader(retbin)
	binary.Read(buf, up.m_endian, &ret)
	return ret
}

func (up *UnPack) GetLong() uint64 {
	retbin := up.m_bin[up.location : up.location+8]
	up.location += 8
	var ret uint64
	buf := bytes.NewReader(retbin)
	binary.Read(buf, up.m_endian, &ret)
	return ret
}

func (up *UnPack) GetShortStr() string {
	tempdata := up.GetToken()
	return string(tempdata)
}

func (up *UnPack) GetToken() []byte {
	length := int(up.GetShort())
	return up.GetBin(length)
}

func (up *UnPack) GetLToken() []byte {
	length := int(up.GetInt())
	return up.GetBin(length - 4)
}

func (up *UnPack) GetIntStr() string {
	tempdata := up.GetLToken()
	return string(tempdata)
}

// # Pack#############################################################################################

type Pack struct {
	m_bin    []byte
	m_endian binary.ByteOrder
}

func NewPack(BigOrLittle string) *Pack {
	p := &Pack{
		m_bin: []byte{},
	}
	if BigOrLittle == "little" {
		p.m_endian = binary.LittleEndian
	} else {
		p.m_endian = binary.BigEndian
	}
	return p
}

func (p *Pack) Empty() {
	p.m_bin = []byte{}
}

func (p *Pack) GetAll() []byte {
	return p.m_bin
}

func (p *Pack) length() int {
	return len(p.m_bin)
}

func (p *Pack) SetData(data []byte) {
	p.m_bin = data
}

func (p *Pack) SetBin(data []byte) {
	p.m_bin = append(p.m_bin, data...)
}

// 长度为 1 的字节
func (p *Pack) SetByte(data byte) {
	buf := make([]byte, 1)
	buf[0] = data
	p.m_bin = append(p.m_bin, buf...)
}

func (p *Pack) SetInt(data int32) {
	buf := make([]byte, 4)
	p.m_endian.PutUint32(buf, uint32(data))
	p.m_bin = append(p.m_bin, buf...)
}

func (p *Pack) SetShort(data int16) {
	buf := make([]byte, 2)
	p.m_endian.PutUint16(buf, uint16(data))
	p.m_bin = append(p.m_bin, buf...)
}

func (p *Pack) SetLong(data int64) {
	buf := make([]byte, 8)
	p.m_endian.PutUint64(buf, uint64(data))
	p.m_bin = append(p.m_bin, buf...)
}

func (p *Pack) SetHex(data string) error {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	p.m_bin = append(p.m_bin, buf...)
	return nil
}

func (p *Pack) SetShortStr(data string) {
	tmpdata := []byte(data)
	p.SetToken(tmpdata)
}

func (p *Pack) SetToken(data []byte) {
	lengthBuf := make([]byte, 2)
	p.m_endian.PutUint16(lengthBuf, uint16(len(data)))
	p.m_bin = append(p.m_bin, lengthBuf...)
	p.m_bin = append(p.m_bin, data...)
}

func (p *Pack) SetLToken(data []byte) {
	lengthBuf := make([]byte, 4)
	p.m_endian.PutUint32(lengthBuf, uint32(len(data)+4))
	p.m_bin = append(p.m_bin, lengthBuf...)
	p.m_bin = append(p.m_bin, data...)
}

func (p *Pack) SetIntStr(data string) {
	tmpdata := []byte(data)
	p.SetLToken(tmpdata)
}
