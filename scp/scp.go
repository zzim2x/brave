package scp

import (
	"sync"
	"math"
	"crypto/sha256"
	"hash"
	"github.com/davecgh/go-xdr/xdr"
)

type EnvelopeState int
type ValidationLevel int

// true 또는 false 는 확실하게 계산을 통해 확인 되는 경우 응답하는 결과이며
// maybe 는 현재까지의 정보를 기준으로는 판단할 수 없는 경우 응답하는 값이다.
// 예) IsNodeInQuorum : slot 여러개를 보면서 true|false 로 확신될 때까지 확인하고 안되면 maybe 넘기고 충분히 데이터가 쌓인 후 다음 연산에 결과를 답변하는 구조
type TriBool int

// 아래 상수 값은 변경될 수 있다. 코드 포팅하면서 바뀔 수 있음. stellar-core 서버와 통신할 때 문제가 확인될 때 변경된다. 아직까진 무슨값이든 문제될 것 없다.
const (
	EnvelopeStateInvalid               EnvelopeState   = 1
	EnvelopeStateValid                 EnvelopeState   = 2
	ValidationLevelInvalidValue        ValidationLevel = 1
	ValidationLevelFullyValidatedValue ValidationLevel = 2
	ValidationLevelMaybeValidValue     ValidationLevel = 3
	TriBoolTrue                        TriBool         = 1
	TriBoolFalse                       TriBool         = 2
	TriBoolMaybe                       TriBool         = 3
)

// SCP 본연의 기능만을 수행하기위해서 필요한 driver
// scp <-> application 를 이어주는 interface
type Driver interface {
	SignEnvelope(envelope *Envelope)

	VerifyEnvelope(envelope Envelope) bool

	ValidateValue(slotId uint64, value Value, nomination bool) ValidationLevel

	GetQuorumSet(hash Hash) *QuorumSet

	EmitEnvelope(envelope Envelope)

	NominatingValue(slotIndex uint64, value Value)

	CombineCandidates(slotIndex uint64) Value

	ComputeHashNode(slotIndex uint64, prev Value, isPriority bool, roundNumber int32, nodeId PublicKey) uint64

	ComputeHashValue(slotIndex uint64, prev Value, roundNumber int32, value Value) uint64
}

func hashHelper(slotIndex uint64, prev Value, extra func(hash.Hash)) uint64 {
	h := sha256.New()
	if b, err := xdr.Marshal(slotIndex); err == nil {
		h.Write(b)
	}
	h.Write(prev)
	extra(h)
	t := h.Sum(nil)

	res := uint64(0)
	for i := 0; i < 8; i++ {
		res = (res << 8) | uint64(t[i])
	}
	return res
}

type SCP struct {
	sync.RWMutex

	driver       Driver
	localNode    *LocalNode
	knownSlotIds []uint64
	knownSlots   map[uint64]*slot
}

func NewSCP(driver Driver, nodeId PublicKey, isValidator bool, quorumSet QuorumSet) *SCP {
	scp := &SCP{
		driver:       driver,
		knownSlotIds: make([]uint64, 0),
		knownSlots:   make(map[uint64]*slot),
	}
	scp.localNode = newLocalNode(nodeId, isValidator, quorumSet, scp)
	return scp
}

func (o *SCP) Nominate(slotIndex uint64, value Value, previousValue Value) bool {
	return o.GetSlot(slotIndex, true).nominate(value, previousValue, false)
}

func (o *SCP) StopNominate(slotId uint64) {
	s := o.GetSlot(slotId, false)
	if s != nil {
		s.stopNomination()
	}
}

func (o *SCP) ReceiveEnvelope(envelope Envelope) EnvelopeState {
	if !o.driver.VerifyEnvelope(envelope) {
		return EnvelopeStateInvalid
	}
	return o.GetSlot(envelope.Statement.SlotIndex, false).processEnvelope(envelope, false)
}

func (o *SCP) PurgeSlots(maxSlotId uint64) {
	i := 0

	knownSlotIds := make([]uint64, 0)
	for _, slotId := range o.knownSlotIds {
		if slotId < maxSlotId {
			delete(o.knownSlots, slotId)
		} else {
			knownSlotIds = append(knownSlotIds, slotId)
		}
		i += 1
	}
	o.knownSlotIds = knownSlotIds
}

func (o *SCP) validateValue(slotIndex uint64, value Value, nomination bool) ValidationLevel {
	return ValidationLevelMaybeValidValue
}

func (o *SCP) GetLocalID() PublicKey {
	return o.localNode.nodeId
}

func (o *SCP) GetLowSlotIndex() uint64 {
	if len(o.knownSlotIds) == 0 {
		return 0
	}
	slotId := uint64(math.MaxInt64)
	for _, c := range o.knownSlotIds {
		if c < slotId {
			slotId = c
		}
	}
	return uint64(slotId)
}

// isNodeInQuorum 함수는 매우 중요하다. 무수히 많은 peer 에서 envelopes 가 넘쳐난다.
// 불특정 다수가 보내는 statement 중 무시해도 되는 정보인지 아는 것은 중요하다.
func (o *SCP) IsNodeInQuorum(nodeId PublicKey) TriBool {
	res := TriBoolMaybe
	for _, s := range o.knownSlots {
		res = s.IsNodeInQuorum(nodeId)
		if res != TriBoolMaybe {
			break
		}
	}
	return res
}

func (o *SCP) GetHighSlotIndex() uint64 {
	if len(o.knownSlotIds) == 0 {
		return 0
	}
	slotId := uint64(0)
	for _, c := range o.knownSlotIds {
		if c > slotId {
			slotId = c
		}
	}
	return uint64(slotId)
}

func (o *SCP) GetSlot(slotId uint64, create bool) *slot {
	o.RLock()
	defer o.RUnlock()

	if current, exists := o.knownSlots[slotId]; exists {
		return current
	} else if !create {
		return nil
	}

	o.knownSlots[slotId] = newSlot(o, slotId)
	o.knownSlotIds = append(o.knownSlotIds, slotId)

	return o.knownSlots[slotId]
}
