package dem

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	demSignon byte = iota
	demPacket
	demSyncTick
	demConsoleCmd
	demUserCmd
	demDataTables
	demStop
	demCustomData
	demStringTables
	demLastCommand = demStringTables
)

var (
	signature = "HL2DEMO"
	ErrParse  = errors.New("this is not dem file")
	ErrBroken = errors.New("this file is broken")
	order     = binary.BigEndian
)

type demo struct {
	header *Header
	frames []Frame
}

type SourceDemo interface {
	Header() *Header
}

type Header struct {
	DemoProtocol    int32
	NetworkProtocol int32
	ServerName      string
	ClientName      string
	MapName         string
	GameDirectory   string
	PlaybackTime    float32
	Ticks           int32
	Frames          int32
	SignOnLength    int32
}

type Frame interface {
}

type FrameSignon struct {
	Frame
	data []byte
}

type FramePacket struct {
	Frame
	data []byte
}

type FrameSyncTick struct {
	Frame
}

type FrameConsoleCmd struct {
	Frame
	data []byte
}

type FrameUserCmd struct {
	Frame
	data []byte
}

type FrameDataTables struct {
	Frame
	data []byte
}

type FrameStop struct {
	Frame
	data []byte
}

type FrameCustomData struct {
	Frame
	data []byte
}

type FrameStringTables struct {
	Frame
	data []byte
}

type decoder struct {
	r *bufio.Reader
}

func NewDemo(r io.Reader) (SourceDemo, error) {
	d := new(decoder)
	d.r = bufio.NewReader(r)
	demo := new(demo)
	var err error
	demo.header, err = d.readHeader()
	demo.frames = []Frame{}
	for {
		f, err := d.readFrame(demo.header)
		demo.frames = append(demo.frames, f)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return demo, err
}

func (d *decoder) readHeader() (*Header, error) {
	if sig, err := d.readString(8); err != nil {
		return nil, err
	} else if sig != signature {
		return nil, fmt.Errorf("%s, %d", sig, len(sig))
	}
	h := new(Header)
	binary.Read(d.r, order, &h.DemoProtocol)
	binary.Read(d.r, order, &h.NetworkProtocol)
	var err error
	h.ServerName, err = d.readString(260)
	h.ClientName, err = d.readString(260)
	h.MapName, err = d.readString(260)
	h.GameDirectory, err = d.readString(260)
	if err != nil {
		return nil, err
	}
	binary.Read(d.r, order, &h.PlaybackTime)
	binary.Read(d.r, order, &h.Ticks)
	binary.Read(d.r, order, &h.Frames)
	binary.Read(d.r, order, &h.SignOnLength)
	return h, nil
}

func (d *decoder) readString(size int) (string, error) {
	slice := make([]byte, size)
	_, err := d.r.Read(slice)
	if err != nil {
		return "", err
	}
	for i, s := range slice {
		if s == 0 {
			return string(slice[:i]), nil
		}
	}
	return "", ErrBroken
}

func (d *decoder) readFrame(head *Header) (Frame, error) {
	t, _ := d.r.ReadByte()
	var tick int32
	err := binary.Read(d.r, order, &tick)
	if err != nil {
		return nil, err
	}
	switch t {
	case demSignon:
		return d.readFrameSignon(head.SignOnLength)
	case demPacket:
		return d.readFramePacket()
	case demSyncTick:
		return d.readFrameSyncTick()
	case demConsoleCmd:
		return d.readFrameConsoleCmd()
	case demUserCmd:
		return d.readFrameUserCmd()
	case demDataTables:
		return d.readFrameDataTables()
	case demStringTables:
		return d.readFrameStringTables()
	}
	return nil, nil
}

func (d *decoder) readFrameSignon(size int32) (*FrameSignon, error) {
	if _, err := d.r.Read(make([]byte, size)); err != nil {
		return nil, err
	}
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FrameSignon{data: data}, nil
}

func (d *decoder) readFramePacket() (*FramePacket, error) {
	if _, err := d.r.Read(make([]byte, 0x54)); err != nil {
		return nil, err
	}
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FramePacket{data: data}, nil
}

func (d *decoder) readFrameSyncTick() (*FrameSyncTick, error) {
	return &FrameSyncTick{}, nil
}

func (d *decoder) readFrameConsoleCmd() (*FrameConsoleCmd, error) {
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FrameConsoleCmd{data: data}, nil
}

func (d *decoder) readFrameUserCmd() (*FrameUserCmd, error) {
	if _, err := d.r.Read(make([]byte, 0x4)); err != nil {
		return nil, err
	}
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FrameUserCmd{data: data}, nil
}

func (d *decoder) readFrameDataTables() (*FrameDataTables, error) {
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FrameDataTables{data: data}, nil
}

func (d *decoder) readFrameStringTables() (*FrameStringTables, error) {
	data, err := d.readData()
	if err != nil {
		return nil, err
	}
	return &FrameStringTables{data: data}, nil
}

func (d *decoder) readData() ([]byte, error) {
	var size int32
	if err := binary.Read(d.r, binary.BigEndian, &size); err != nil {
		return nil, err
	}
	if size < 0 {
		return []byte{}, nil
	}
	data := make([]byte, int(size))
	_, err := d.r.Read(data)
	return data, err
}

func (d *demo) Header() *Header {
	return d.header
}
