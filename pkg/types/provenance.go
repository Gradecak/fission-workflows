package types

func (spec *WorkflowSpec) ToProvenance() *Node {
	node := &Node{Type: Node_DATAFLOW} // TODO
	node.Tag = spec.GetName()
	for _, task := range spec.GetTasks() {
		node.Edges = append(node.Edges, task.ToProvenance())
	}
	return node
}

func (spec *TaskSpec) ToProvenance() *Node {
	node := &Node{Type: Node_TASK}
	node.Tag = spec.GetFunctionRef() // temporary replace with a proper tag later on
	return node
}
