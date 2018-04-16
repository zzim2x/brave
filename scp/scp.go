package scp

import (
	"sync"
	"math"
)

type EnvelopeState int
type TriBool int

const (
	EnvelopeStateInvalid EnvelopeState = 1
	EnvelopeStateValid   EnvelopeState = 2
	TriBoolTrue          TriBool       = 1
	TriBoolFalse         TriBool       = 2
	TriBoolMaybe         TriBool       = 3
)

type Driver interface {
	VerifyEnvelope(envelope Envelope) bool
}

type SCP struct {
	sync.RWMutex

	driver       Driver
	localNode    *LocalNode
	knownSlotIds []uint64
	knownSlots   map[uint64]*slot
}

func NewSCP(nodeId PublicKey, isValidator bool, quorumSet QuorumSet) *SCP {
	scp := &SCP{
		knownSlotIds: make([]uint64, 0),
		knownSlots:   make(map[uint64]*slot),
	}
	scp.localNode = newLocalNode(nodeId, isValidator, quorumSet, scp)
	return scp
}

func (o *SCP) Nominate(slotId uint64, value *Value, previousValue *Value) bool {
	return o.GetSlot(slotId, true).nominate(value, previousValue, false)
}

func (o *SCP) StopNominate(slotId uint64) {
	s := o.GetSlot(slotId, false)
	if s != nil {
		s.stopNomination()
	}
}

func (o *SCP) ReceiveSCPEnvelope(envelope Envelope) EnvelopeState {
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
