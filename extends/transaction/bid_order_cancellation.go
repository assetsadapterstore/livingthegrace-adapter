package transaction

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	BidOrderCancellationType    = 2
	BidOrderCancellationSubType = 5
)

type BidOrderCancellation struct {
	*pb.BidOrderCancellation
}

func EmptyBidOrderCancellation() *BidOrderCancellation {
	return &BidOrderCancellation{
		BidOrderCancellation: &pb.BidOrderCancellation{
			Attachment: &pb.BidOrderCancellation_Attachment{},
		},
	}
}

func (tx *BidOrderCancellation) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Order)
}

func (tx *BidOrderCancellation) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Order = d.ReadUint64()
}

func (tx *BidOrderCancellation) AttachmentSizeInBytes() int {
	return 8
}

func (tx *BidOrderCancellation) GetType() uint16 {
	return BidOrderCancellationSubType<<8 | BidOrderCancellationType
}

func (tx *BidOrderCancellation) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *BidOrderCancellation) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
