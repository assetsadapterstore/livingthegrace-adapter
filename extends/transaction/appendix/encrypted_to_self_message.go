package appendix

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

type EncryptedToSelfMessage struct {
	*pb.Message
}

func (a *EncryptedToSelfMessage) WriteBytes(e encoding.Encoder) {
	e.WriteBytesWithInt32Len(tx.Attachment.IsText, []byte(tx.Attachment.Data))
	e.WriteBytesWithInt32Len(tx.Attachment.IsText, []byte(tx.Attachment.Nonce))
}

func (a *EncryptedToSelfMessage) SizeInBytes() int {
	return 4 + len(tx.Attachment.Data) + len(tx.Attachment.Nonce)
}
