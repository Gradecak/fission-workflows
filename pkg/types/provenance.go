package types

const (
	DEFAULT_EDGE_TYPE = "OP"
)

func (spec *WorkflowSpec) ToProvenance() *Node {
	node := &Node{Type: Node_DATAFLOW}
	node.Tag = spec.GetName()
	// Add predecssor tag
	if prov := spec.GetProvenanceMeta(); prov != nil {
		if predec := prov.GetPredecessor(); predec != "" {
			node.Predecessor = predec
		}
	}
	// convert tasks to Child nodes
	for _, task := range spec.GetTasks() {
		edge := &Edge{EdgeType: DEFAULT_EDGE_TYPE, Dst: task.ToProvenance()}

		// convert additional metadata annotated by user into the provenance graph
		if prov := spec.GetProvenanceMeta(); prov != nil {
			if opType := prov.GetOpType(); opType != "" {
				edge.EdgeType = opType
			}
		}

		node.Edges = append(node.Edges, edge)
	}

	return node
}

func (spec *TaskSpec) ToProvenance() *Node {
	node := &Node{Type: Node_TASK, Tag: spec.GetFunctionRef()}

	if prov := spec.GetProvenanceMeta(); prov != nil {
		if meta := prov.GetMeta(); meta != "" {
			node.Tag = meta
		}
	}

	return node
}
