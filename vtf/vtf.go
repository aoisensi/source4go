package vtf

import (
	"encoding/binary"
	"errors"
)

var (
	headerSignature uint32 = 0x00465456
	order                  = binary.LittleEndian
)

var (
	ErrNotVtfFile = errors.New("the file is not vtf file.")
)

const (
	ImageFormatNone     = -1
	ImageFormatRGBA8888 = iota
	ImageFormatABGR8888
	ImageFormatRGB888
	ImageFormatBGR888
	ImageFormatRGB565
	ImageFormatI8
	ImageFormatIA88
	ImageFormatP8
	ImageFormatA8
	ImageFormatRGB888BlueScreen
	ImageFormatBGR888BlueScreen
	ImageFormatARGB8888
	ImageFormatBGRA8888
	ImageFormatDXT1
	ImageFormatDXT3
	ImageFormatDXT5
	ImageFormatBGRX8888
	ImageFormatBGR565
	ImageFormatBGRX5551
	ImageFormatBGRA4444
	ImageFormatDXT1OneBitAlpha
	ImageFormatBGRA5551
	ImageFormatUV88
	ImageFormatUVWQ8888
	ImageFormatRGBA16161616F
	ImageFormatRGBA16161616
	ImageFormatUVLX8888
)

const (
	TextureFlagsPointSample uint32 = 1 << iota
	TextureFlagsTrilinear
	TextureFlagsClampS
	TextureFlagsClampT
	TextureFlagsAnisotropic
	TextureFlagsHintDXT5
	TextureFlagsPWLCorrected
	TextureFlagsNormal
	TextureFlagsNoMip
	TextureFlagsProcedural

	TextureFlagsOneBitAlpha
	TextureFlagsEightBitAlpha
	TextureFlagsENVMap
	TextureFlagsRenderTarget
	TextureFlagsDepthRenderTarget
	TextureFlagsNoDebugOverride
	TextureFlagsSingleCopy
	TextureFlagsPreSRGB

	TextureFlagsUnused00100000
	TextureFlagsUnused00200000
	TextureFlagsUnused00400000

	TextureFlagsNoDepthBuffer

	TextureFlagsUnused01000000

	TextureFlagsClampU
	TextureFlagsVertexTexture
	TextureFlagsSSBump

	TextureFlagsUnused10000000

	TextureFlagsBorder

	TextureFlagsUnused40000000
	TextureFlagsUnused80000000
)

type vtfHeader struct {
	signature          uint32 //little-endian integer, 0x00465456little-endian integer, 0x00465456
	version            [2]uint32
	headerSize         uint32
	width              uint16
	height             uint16
	flags              uint32
	frames             uint16
	firstFrame         uint16
	padding0           [4]byte
	reflectivity       [3]float32
	padding1           [4]byte
	bumpmapScale       float32
	highResImageFormat uint32
	mipmapCount        byte
	lowResImageFormat  uint32
	lowResImageWidth   byte
	lowResImageHeight  byte
	depth              uint16
}
