package transaction

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	EscrowSignType    = 21
	EscrowSignSubType = 1
)

type EscrowSign struct {
	*pb.EscrowSign
}

func EmptyEscrowSign() *EscrowSign {
	return &EscrowSign{
		EscrowSign: &pb.EscrowSign{
			Attachment: &pb.EscrowSign_Attachment{},
		},
	}
}

func (tx *EscrowSign) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Id)
	e.WriteUint8(uint8(tx.Attachment.Decision))
}

func (tx *EscrowSign) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Id = d.ReadUint64()
	tx.Attachment.Decision = pb.DeadlineAction(d.ReadUint8())
}

func (tx *EscrowSign) AttachmentSizeInBytes() int {
	return 8 + 1
}

func (tx *EscrowSign) GetType() uint16 {
	return EscrowSignSubType<<8 | EscrowSignType
}

func (tx *EscrowSign) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *EscrowSign) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
