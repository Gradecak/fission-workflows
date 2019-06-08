package graph

import (
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/sirupsen/logrus"
	"hash/fnv"
)

func GenProvenance(wfi *types.WorkflowInvocation) *Provenance {
	provenance := NewProvenance()
	wfID := provenance.FromWorkflow(wfi)
	// If the workflow invocation comes with a consent identifier, we use
	// this as the unique identifier for a 'user' who's data has been
	// operated on
	if consentId := wfi.GetSpec().GetConsentId(); consentId != "" {
		provenance.AddEntity(consentId, wfID)
	}
	return provenance
}

func NewIDArray() *IDs {
	return &IDs{Ids: []int64{}}

}

func (m *IDs) AddID(id int64) {
	m.Ids = append(m.Ids, id)
}

func (m *IDs) Merge(ids *IDs) {
	m.Ids = append(m.Ids, ids.Ids...)
}

func NewProvenance() *Provenance {
	return &Provenance{
		Nodes:          make(map[int64]*Node),
		WfTasks:        make(map[int64]*IDs),
		WfPredecessors: make(map[int64]*IDs),
		Executed:       make(map[string]int64),
	}
}

func (p *Provenance) AddTasks(tasks map[string]*types.Task) *IDs {
	ids := NewIDArray()
	for id, task := range tasks {
		// TODO only add tasks that are suceeded
		n := &Node{Type: Node_TASK}
		if pMeta := task.GetSpec().GetProvenanceMeta(); pMeta != nil {
			n.addTaskProvMeta(pMeta)
		}
		tId := createID(id)
		n.FnName = task.GetSpec().GetFunctionRef()
		p.Nodes[tId] = n
		ids.AddID(tId)
	}
	return ids
}

func (p *Provenance) AddEntity(consentID string, wfIDs int64) {
	if _, ok := p.Executed[consentID]; !ok {
		p.Executed[consentID] = wfIDs
	}
}

func (p *Provenance) FromWorkflow(wfi *types.WorkflowInvocation) int64 {
	// Add Tasks from workflow
	logrus.Debugf("DYNAMIC TASKS: %+v", wfi.GetStatus().GetDynamicTasks())
	ids := p.AddTasks(wfi.Tasks())
	wfId := fingerprintWfi(wfi)
	p.WfTasks[wfId] = ids
	return wfId
}

func (p *Provenance) GetWorkflowTaskIds(wfid int64) []int64 {
	return p.GetWfTasks()[wfid].GetIds()
}
func (p *Provenance) GetWorkflowPredecessors(wfid int64) []int64 {
	return p.GetWfPredecessors()[wfid].GetIds()
}

func (n *Node) addTaskProvMeta(pMeta *types.ProvenanceMetadata) {
	switch pMeta.OpType {
	case types.ProvenanceMetadata_TRANSFORM:
		n.Op = Node_TRANSFORM
	case types.ProvenanceMetadata_READ:
		n.Op = Node_READ
	case types.ProvenanceMetadata_WRITE:
		n.Op = Node_WRITE
	case types.ProvenanceMetadata_CONTROL:
		n.Op = Node_CONTROL
	default:
		n.Op = Node_NOP
	}
	n.Meta = pMeta.GetMeta()
}

// Given a workflow invocation Generate a unique "fingerprint" so that we know
// if all of the tasks executed in an invocation are the same
// TODO make it actually work :)
func fingerprintWfi(wfi *types.WorkflowInvocation) int64 {
	id := createID(wfi.Spec.Workflow.Metadata.Id)
	logrus.Debugf("Worfklow id: %v Workflow Fingerprint %v", wfi.Spec.Workflow.Metadata.Id, id)
	return id
}

func createID(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}
