package transaction

import (
	pb "github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/encoding"
)

const (
	RewardRecipientAssignmentType    = 20
	RewardRecipientAssignmentSubType = 0
)

type RewardRecipientAssignment struct {
	*pb.RewardRecipientAssignment
}

func EmptyRewardRecipientAssignment() *RewardRecipientAssignment {
	return &RewardRecipientAssignment{
		RewardRecipientAssignment: &pb.RewardRecipientAssignment{},
	}
}

func (tx *RewardRecipientAssignment) WriteAttachmentBytes(e encoding.Encoder) {}

func (tx *RewardRecipientAssignment) AttachmentSizeInBytes() int {
	return 0
}

func (tx *RewardRecipientAssignment) ReadAttachmentBytes(d encoding.Decoder) {}

func (tx *RewardRecipientAssignment) GetType() uint16 {
	return RewardRecipientAssignmentSubType<<8 | RewardRecipientAssignmentType
}

func (tx *RewardRecipientAssignment) SetAppendix(a *pb.Appendix) {
	tx.Appendix = a
}

func (tx *RewardRecipientAssignment) SetHeader(h *pb.TransactionHeader) {
	tx.Header = h
}
