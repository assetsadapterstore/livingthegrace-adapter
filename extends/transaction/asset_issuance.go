package transaction

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	AssetIssuanceType    = 2
	AssetIssuanceSubType = 0
)

type AssetIssuance struct {
	*pb.AssetIssuance
}

func EmptyAssetIssuance() *AssetIssuance {
	return &AssetIssuance{
		AssetIssuance: &pb.AssetIssuance{
			Attachment: &pb.AssetIssuance_Attachment{},
		},
	}
}

func (tx *AssetIssuance) WriteAttachmentBytes(e encoding.Encoder) {
	e.WriteUint8(uint8(len(tx.Attachment.Name)))
	e.WriteBytes(tx.Attachment.Name)
	e.WriteUint16(uint16(len(tx.Attachment.Description)))
	e.WriteBytes(tx.Attachment.Description)
	e.WriteUint64(tx.Attachment.Quantity)
	e.WriteUint8(uint8(tx.Attachment.Decimals))
}

func (tx *AssetIssuance) ReadAttachmentBytes(d encoding.Decoder) {
	nameLen := d.ReadUint8()
	tx.Attachment.Name = d.ReadBytes(int(nameLen))
	descritpionLen := d.ReadUint16()
	tx.Attachment.Description = d.ReadBytes(int(descritpionLen))
	tx.Attachment.Quantity = d.ReadUint64()
	tx.Attachment.Decimals = uint32(d.ReadUint8())
}

func (tx *AssetIssuance) AttachmentSizeInBytes() int {
	return 1 + len(tx.Attachment.Name) + 2 + len(tx.Attachment.Description) + 8 + 1
}

func (tx *AssetIssuance) GetType() uint16 {
	return AssetIssuanceSubType<<8 | AssetIssuanceType
}

func (tx *AssetIssuance) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *AssetIssuance) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
