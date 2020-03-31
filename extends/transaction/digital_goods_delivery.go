package transaction

import (
	"encoding/hex"
	"math"

	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	DigitalGoodsDeliveryType    = 3
	DigitalGoodsDeliverySubType = 5
)

type DigitalGoodsDelivery struct {
	*pb.DigitalGoodsDelivery
}

func EmptyDigitalGoodsDelivery() *DigitalGoodsDelivery {
	return &DigitalGoodsDelivery{
		DigitalGoodsDelivery: &pb.DigitalGoodsDelivery{
			Attachment: &pb.DigitalGoodsDelivery_Attachment{},
		},
	}
}

func (tx *DigitalGoodsDelivery) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint64(tx.Attachment.Purchase)
	// java wallet <3
	l := len(tx.Attachment.Data) / 2
	if tx.Attachment.IsText {
		e.WriteInt32(int32(l) | math.MinInt32)
	} else {
		e.WriteInt32(int32(l))
	}
	data := make([]byte, l)
	if _, err := hex.Decode(data, tx.Attachment.Data); err != nil {
		return
	}
	e.WriteBytes(data)
	e.WriteBytes(tx.Attachment.Nonce)
	e.WriteUint64(tx.Attachment.Discount)
}

func (tx *DigitalGoodsDelivery) ReadAttachmentBytes(d encoding.Decoder) {
	tx.Attachment.Purchase = d.ReadUint64()
	dataLen := d.ReadInt32()
	if dataLen < 0 {
		tx.Attachment.IsText = true
		dataLen &= math.MaxInt32

		tx.Attachment.Data = make([]byte, dataLen*2)
		hex.Encode(tx.Attachment.Data, d.ReadBytes(int(dataLen)))
	} else {
		tx.Attachment.Data = make([]byte, dataLen*2)
		hex.Encode(tx.Attachment.Data, d.ReadBytes(int(dataLen)))
	}
	tx.Attachment.Nonce = d.ReadBytes(32)
	tx.Attachment.Discount = d.ReadUint64()
}

func (tx *DigitalGoodsDelivery) AttachmentSizeInBytes() int {
	return 8 + 4 + len(tx.Attachment.Data)/2 + len(tx.Attachment.Nonce) + 8
}

func (tx *DigitalGoodsDelivery) GetType() uint16 {
	return DigitalGoodsDeliverySubType<<8 | DigitalGoodsDeliveryType
}

func (tx *DigitalGoodsDelivery) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *DigitalGoodsDelivery) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
