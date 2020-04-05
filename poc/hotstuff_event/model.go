package hotstuff_event

import (
	"fmt"
	"strconv"
)

//go:generate msgp

type MsgType int

func (m MsgType) String() string {
	switch m {
	case Proposal:
		return "Proposal"
	case Vote:
		return "Vote"
	case Timeout:
		return "Timeout"
	//case LocalTimeout:
	//	return "LocalTimeout"
	case String:
		return "String"
	default:
		return "NAMsg"
	}
}

const (
	Proposal MsgType = iota
	Vote
	Timeout
	//LocalTimeout
	String
)

//msgp:tuple Block
type Block struct {
	Round    int    // the round that generated this proposal
	Payload  string // proposed transactions
	ParentQC *QC    // qc for parent block
	Id       string // unique digest of round, payload and parent_qc.id
}

func (b Block) String() string {
	return fmt.Sprintf("[Block round=%d payload=%s parentQC=%s id=%s]", b.Round, b.Payload, b.ParentQC, b.Id)
}

// WireMessage is for transportation.
//msgp:tuple WireMessage
type WireMessage struct {
	MsgType      int    // what ContentByte is.
	ContentBytes []byte // this Byte will be recovered to implementation of Content interface
	SenderId     string
	Signature    Signature
}

func (z *WireMessage) String() string {
	return fmt.Sprintf("WM: Type=%d SenderId=%s ContentBytes=%s", z.MsgType, PrettyId(z.SenderId), ToBriefHex(z.ContentBytes, 100))
}

//msgp:tuple ContentString
type ContentString struct {
	Content string
}

func (z *ContentString) String() string {
	return z.Content
}

func (z *ContentString) SignatureTarget() string {
	return z.Content
}

//msgp:tuple ContentProposal
type ContentProposal struct {
	Proposal Block
	TC       *TC
}

func (c ContentProposal) String() string {
	return fmt.Sprintf("[CProposal: p=%s tc=%s]", c.Proposal, c.TC)
}

func (c ContentProposal) SignatureTarget() string {
	return c.Proposal.String()
}

//msgp:tuple ContentTimeout
type ContentTimeout struct {
	Round  int
	HighQC *QC
	TC     *TC
}

func (c ContentTimeout) SignatureTarget() string {
	return strconv.Itoa(c.Round)
}

func (c ContentTimeout) String() string {
	return fmt.Sprintf("[CTimeout: round=%d highQC=%s tc=%s]", c.Round, c.HighQC, c.TC)
}

//msgp:tuple ContentVote
type ContentVote struct {
	VoteInfo         VoteInfo
	LedgerCommitInfo LedgerCommitInfo
	QC               *QC
	TC               *TC
}

func (c ContentVote) String() string {
	return fmt.Sprintf("[CVote: voteinfo=%s ledgercinfo=%s qc=%s tc=%s]", c.VoteInfo, c.LedgerCommitInfo, c.QC, c.TC)
}

func (c ContentVote) SignatureTarget() string {
	return fmt.Sprintf("[CVote: voteinfo=%s ledgercinfo=%s qc=%s tc=%s]", c.VoteInfo, c.LedgerCommitInfo, c.QC, c.TC)
}

//msgp:tuple QC
type QC struct {
	VoteData   VoteInfo
	Signatures []Signature // simulate sig by give an id to the vote, use joint signature instead
}

func (q *QC) String() string {
	if q == nil {
		return "nil"
	}
	return fmt.Sprintf("[QC: voteInfo=%s sigs=%d]", q.VoteData, len(q.Signatures))
}

//msgp:tuple TC
type TC struct {
	Round      int         // round of the block
	Signatures []Signature // simulate sig by give an id to the vote, use joint signature instead
}

func (t *TC) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("[TC: round=%d, sigs=%d]", t.Round, len(t.Signatures))
}

//msgp:tuple Signature
type Signature struct {
	PartnerIndex int
	Signature    string
}

func (s Signature) String() string {
	return fmt.Sprintf("[Sig: partnerId %d Sig %s]", s.PartnerIndex, s.Signature)
}

//msgp:tuple VoteInfo
type VoteInfo struct {
	Id               string // Id of the block
	Round            int    // round of the block
	ParentId         string // Id of the parent
	ParentRound      int    // round of the parent
	GrandParentId    string // Id of the grandParent
	GrandParentRound int    // round of the grandParent
	ExecStateId      string // speculated execution state
}

func (i VoteInfo) String() string {
	return fmt.Sprintf("<%s %d> <%s %d> <%s %d> %s", i.Id, i.Round, i.ParentId, i.ParentRound, i.GrandParentId, i.GrandParentRound, i.ExecStateId)
}

func (i VoteInfo) GetHashContent() string {
	return fmt.Sprintf("%s %d %s %d %s %d %s", i.Id, i.Round, i.ParentId, i.ParentRound, i.GrandParentId, i.GrandParentRound, i.ExecStateId)
}

//msgp:tuple LedgerCommitInfo
type LedgerCommitInfo struct {
	CommitStateId string // nil if no commit happens when this vote is aggregated to QC. Usually the merkle root
	VoteInfoHash  string // hash of VoteMsg.voteInfo
}

func (l LedgerCommitInfo) String() string {
	return fmt.Sprintf("[LedgerCommitInfo: commitStateId %s VoteInfoHash %s]", l.CommitStateId, l.VoteInfoHash)
}

func (l LedgerCommitInfo) GetHashContent() string {
	return fmt.Sprintf("%s %s", l.CommitStateId, l.VoteInfoHash)
}

//msgp:tuple VoteMsg
type VoteMsg struct {
	VoteInfo         VoteInfo
	LedgerCommitInfo LedgerCommitInfo
	Sender           int
	Signature        Signature
}
