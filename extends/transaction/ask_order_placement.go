package transaction

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	AskOrderPlacementType    = 2
	AskOrderPlacementSubType = 2
)

type AskOrderPlacement struct {
	*pb.AskOrderPlacement
}

func EmptyAskOrderPlacement() *AskOrderPlacement {
	return &AskOrderPlacement{
		AskOrderPlacement: &pb.AskOrderPlacement{
			Attachment: &pb.AskOrderPlacement_Attachment{},
		},
	}
}

func (tx *AskOrderPlacement) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Asset)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint64(tx.Attachment.Price)
}

func (tx *AskOrderPlacement) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Asset = d.ReadUint64()
	tx.Attachment.Quantity = d.ReadUint64()
	tx.Attachment.Price = d.ReadUint64()
}

func (tx *AskOrderPlacement) AttachmentSizeInBytes() int {
	return 8 + 8 + 8
}

func (tx *AskOrderPlacement) GetType() uint16 {
	return AskOrderPlacementSubType<<8 | AskOrderPlacementType
}

func (tx *AskOrderPlacement) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *AskOrderPlacement) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
