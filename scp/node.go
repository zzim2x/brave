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

// 충분한 데이터가 쌓이지 않은 상태에서는 판단을 보류하는 응답을 한다.
// 히스토리가 쌓이지 않았거나, 모르는 quorumSetHash 인 경우 등이 있다.
func (o *LocalNode) IsNodeInQuorum(nodeId PublicKey, qfun func(Statement) *QuorumSet, stats map[PublicKey][]Statement) TriBool {
	res := TriBoolFalse
	backlog := make(map[PublicKey]bool)
	visited := make(map[PublicKey]bool)

	// 현재 node 는 무사 통과
	backlog[o.nodeId] = true

	for len(backlog) > 0 {
		for item := range backlog {
			if item == nodeId {
				return TriBoolTrue
			}
			delete(backlog, item)
			visited[item] = true

			if history, exists := stats[item]; !exists {
				res = TriBoolMaybe
				continue
			} else {
				for _, s := range history {
					if qset := qfun(s); qset == nil {
						res = TriBoolMaybe
						continue
					} else {
						forEachNodes(*qset, func(n PublicKey) {
							if _, exists := visited[n]; !exists {
								backlog[n] = true
							}
						})
					}
				}
			}
		}
	}

	return res
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

// 나중에 코드로 작성하겠지만,
// statements 엔 타 peer 에서 받은 envelopes 뿐이며, 그 또한 quorum 에 속하지 않으면 채워지지 않는다.
func IsVBlockingWithFilter(quorumSet QuorumSet, statements map[PublicKey]Envelope, filter func(statement Statement) bool) bool {
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
