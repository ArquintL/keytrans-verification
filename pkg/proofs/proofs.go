package proofs

import (
	"crypto/sha256"
)

type NodeValue = [sha256.Size]byte

type UpdateValue struct {
	Value []byte
}

type CommitmentValue struct {
	Opening []byte
	Label   []byte      // pseudonym; max length 2^8-1 bytes
	Update  UpdateValue // value associated with label, e.g., public key
}

type BinaryLadderStep struct {
	Proof      [sha256.Size]byte
	Commitment [sha256.Size]byte
}

type InclusionProof struct {
}

const (
	Reserved           = 0
	Inclusion          = 1
	NonInclusionLeaf   = 2
	NonInclusionParent = 3
)

type PrefixLeaf struct {
	Vrf_output [sha256.Size]byte
	Commitment [sha256.Size]byte
}

type PrefixSearchResults struct {
	Result_type int
	Leaf        *PrefixLeaf // only present when result_type == NonInclusionLeaf
	Depth       uint8
}

type PrefixProof struct {
	Results  []PrefixSearchResults
	Elements []NodeValue
}

type CombinedTreeProof struct {
	Timestamps    []uint64
	Prefix_proofs []PrefixProof
	Prefix_roots  []NodeValue
}
