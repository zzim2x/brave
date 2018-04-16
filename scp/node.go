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

func (o *LocalNode) getNodeId() PublicKey {
	return o.nodeId
}

func ForEachNodes(quorumSet QuorumSet, proc func(PublicKey)) {
	nodes := make(map[PublicKey]bool)
	forEachNodes(quorumSet, func(key PublicKey) {
		nodes[key] = true
	})

	for node := range nodes {
		proc(node)
	}
}

func forEachNodes(quorumSet QuorumSet, proc func(PublicKey)) {
	for _, validator := range quorumSet.Validators {
		proc(validator)
	}

	for _, innerSet := range quorumSet.InnerSets {
		forEachNodes(innerSet, proc)
	}
}

func GetNodeWeight(nodeId PublicKey, quorumSet QuorumSet) uint64 {
	return 0
}

func IsVBlocking(quorumSet QuorumSet, nodes []PublicKey) bool {
	return isVBlockingInternal(quorumSet, nodes)
}

func IsVBlockingWithFilter(quorumSet QuorumSet, statements map[PublicKey]Envelope,
		filter func(statement Statement) bool) bool {
	nodes := make([]PublicKey, 0)
	for k := range statements {
		if filter(statements[k].Statement) {
			nodes = append(nodes, k)
		}
	}
	return IsVBlocking(quorumSet, nodes)
}

func isVBlockingInternal(quorumSet QuorumSet, nodes []PublicKey) bool {
	return true
}

func IsQuorumSlice() {
}

func isQuorumSliceInternal(quorumSet QuorumSet, nodes []PublicKey) bool {
	return true
}

func buildSingletonQSet(nodeId PublicKey) QuorumSet {
	return QuorumSet{Threshold: 1, Validators: []PublicKey{nodeId}}
}
