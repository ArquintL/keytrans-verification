package client

import (
	"errors"

	"github.com/felixlinker/keytrans-verification/pkg/proofs"
)

type TreeHead struct {
	Tree_size uint64
	Signature []byte
}

/*@
pred (t TreeHead) Inv() {
	acc(t.Signature)
}
@*/

type FullTreeHead struct {
	Tree_head TreeHead
	// TODO: AuditorTreeHead auditor_tree_head
}

/*@
pred (f FullTreeHead) Inv() {
	f.Tree_head.Inv()
}

ghost
decreases
requires f.Inv()
pure func (f FullTreeHead) Size() uint64 {
	return unfolding f.Inv() in f.Tree_head.Tree_size
}
@*/

type SearchRequest struct {
	Last  *uint32
	Label []byte
	// TODO: optional<uint32> version
}

/*@
pred (s SearchRequest) Inv() {
	acc(s.Last) && acc(s.Label)
}
@*/

type SearchResponse struct {
	Full_tree_head FullTreeHead
	Version        *uint32 // version; only present for latest-key queries
	Binary_ladder  []proofs.BinaryLadderStep
	Search         proofs.CombinedTreeProof
	Inclusion      proofs.InclusionProof
	Opening        []byte
	Value          proofs.UpdateValue // value associated with queried label
}

//@ requires noPerm < p
//@ preserves st.Inv()
//@ preserves acc(query.Inv(), p) && acc(resp.Inv(), p)
//@ ensures err == nil ==> acc(res) && res.Inv()
func (st *UserState) VerifyLatest(query SearchRequest, resp SearchResponse) (*proofs.UpdateValue, error) {
	//@ unfold acc(resp.Inv(), p)
	if err := st.UpdateView(resp.Full_tree_head, resp.Search); err != nil {
		//@ fold acc(resp.Inv(), p)
		return nil, err
		} else if resp.Version != nil {
		//@ fold acc(resp.Inv(), p)
		return nil, errors.New("no version provided")
	} else if len(resp.Search.Prefix_roots) != 0 {
		//@ fold acc(resp.Inv(), p)
		return nil, errors.New("prefix roots provided")
	}

	ladderIndices := proofs.FullBinaryLadderSteps(*resp.Version)
	if len(resp.Binary_ladder) != len(ladderIndices) {
		return nil, errors.New("length of binary ladder does not match greatest version")
	}

	trees := make([]*proofs.PrefixTree, 0, len(resp.Search.Prefix_proofs))
	for _, prf := range(resp.Search.Prefix_proofs) {
		if tree, err := prf.ToTree(resp.Binary_ladder); err != nil {
			return nil, err
		} else {
			trees = append(trees, tree)
		}
	}

	// TODO: Verify proof of inclusion in all trees

	//@ fold acc(resp.Inv(), p)
	return nil, nil
}
