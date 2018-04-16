package scp

type LocalNode struct {
	scp           *SCP
	nodeId        PublicKey
	isValidator   bool
	quorumSet     QuorumSet
	quorumSetHash Hash
}

func newLocalNode(nodeId PublicKey, isValidator bool, quorumSet QuorumSet, scp *SCP) *LocalNode {
	return &LocalNode{
		scp:           scp,
		nodeId:        nodeId,
		isValidator:   isValidator,
		quorumSet:     quorumSet,
		quorumSetHash: quorumSet.Hash(),
	}
}
