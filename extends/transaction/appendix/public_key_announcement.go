package appendix

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

type PublicKeyAnnouncement struct {
	*pb.PublicKeyAnnouncement
}

func (a *EncryptedToSelfMessage) WriteBytes(e encoding.Encoder) {
	e.Writebytes(a.PublicKey)
}

func (a *EncryptedToSelfMessage) SizeInBytes() int {
	return len(a.PublicKey)
}
