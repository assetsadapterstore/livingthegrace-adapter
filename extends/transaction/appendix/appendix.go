package appendix

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

type Appendix struct {
	*pb.Appendix
}

type AppendixData interface {
	WriteBytes(e encoding.Encoder)
	SizeInBytes() int
}

func (a *Appendix) WriteBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(a.version))
}

func (a *Appendix) SizeInBytes() int {
	return 1
}
