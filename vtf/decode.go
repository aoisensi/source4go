package vtf

import (
	"encoding/binary"
	"image"
	"io"
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

func (d *decoder) checkHeader() error {
	if err := binary.Read(d.r, order, &d.header); err != nil {
		return err
	}
	if d.header.signature != headerSignature {
		return ErrNotVtfFile
	}
	return nil
}
