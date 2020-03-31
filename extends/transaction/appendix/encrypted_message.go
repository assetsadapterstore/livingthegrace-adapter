package appendix

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

type EncrypedMessage struct {
	*pb.Message
}

func (a *EncryptedMessgae) WriteBytes(e encoding.Encoder) {
	e.WriteBytesWithInt32Len(a.IsText, []byte(a.Data))
	e.WriteBytesWithInt32Len(a.IsText, []byte(a.Nonce))
}

func (a *EncryptedMessage) SizeInBytes() int {
	return 4 + len(a.Data) + len(a.Nonce)
}
