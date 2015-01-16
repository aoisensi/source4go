package vtf

import (
	"encoding/binary"
	"image"
	"io"

	"github.com/k0kubun/pp"
)

type decoder struct {
	r      io.Reader
	header vtfHeader
}

func Decode(r io.Reader) (image.Image, error) {
	d := &decoder{
		r: r,
	}
	if err := d.checkHeader(); err != nil {
		return nil, err
	}
	return nil, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	c, _, err := decodeConfig(r)
	return c, err
}

func decodeConfig(r io.Reader) (image.Config, *decoder, error) {
	d := &decoder{
		r: r,
	}
	cfg := image.Config{}
	if err := d.checkHeader(); err != nil {
		return cfg, d, err
	}
	cfg.Width = int(d.header.Width)
	cfg.Height = int(d.header.Height)
	return cfg, d, nil
}

func (d *decoder) checkHeader() error {
	if err := binary.Read(d.r, order, &d.header); err != nil {
		return err
	}
	if d.header.Signature != headerSignature {
		return ErrNotVtfFile
	}
	pp.Println(d.header)
	return nil
}
