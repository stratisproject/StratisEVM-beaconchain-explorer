package rpc

import (
	"encoding/hex"
	"errors"
	"eth2-exporter/utils"
	"fmt"
	"strconv"
)

var errNotFound = errors.New("not found 404")

func uint64List(li []uint64Str) []uint64 {
	out := make([]uint64, len(li))
	for i, v := range li {
		out[i] = uint64(v)
	}
	return out
}

func syncCommitteeParticipation(bits []byte) float64 {
	participating := 0
	for i := 0; i < int(utils.Config.Chain.ClConfig.SyncCommitteeSize); i++ {
		if utils.BitAtVector(bits, i) {
			participating++
		}
	}
	return float64(participating) / float64(utils.Config.Chain.ClConfig.SyncCommitteeSize)
}

type bytesHexStr []byte

func (s *bytesHexStr) UnmarshalText(b []byte) error {
	if s == nil {
		return fmt.Errorf("cannot unmarshal bytes into nil")
	}
	if len(b) >= 2 && b[0] == '0' && b[1] == 'x' {
		b = b[2:]
	}
	out := make([]byte, len(b)/2)
	hex.Decode(out, b)
	*s = out
	return nil
}

type uint64Str uint64

func (s *uint64Str) UnmarshalJSON(b []byte) error {
	return Uint64Unmarshal((*uint64)(s), b)
}

// Parse a uint64, with or without quotes, in any base, with common prefixes accepted to change base.
func Uint64Unmarshal(v *uint64, b []byte) error {
	if v == nil {
		return errors.New("nil dest in uint64 decoding")
	}
	if len(b) == 0 {
		return errors.New("empty uint64 input")
	}
	if b[0] == '"' || b[0] == '\'' {
		if len(b) == 1 || b[len(b)-1] != b[0] {
			return errors.New("uneven/missing quotes")
		}
		b = b[1 : len(b)-1]
	}
	n, err := strconv.ParseUint(string(b), 0, 64)
	if err != nil {
		return err
	}
	*v = n
	return nil
}

type StandardBeaconHeaderResponse struct {
	Data struct {
		Root   string `json:"root"`
		Header struct {
			Message struct {
				Slot          uint64Str `json:"slot"`
				ProposerIndex uint64Str `json:"proposer_index"`
				ParentRoot    string    `json:"parent_root"`
				StateRoot     string    `json:"state_root"`
				BodyRoot      string    `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
	Finalized bool `json:"finalized"`
}

type StandardBeaconHeadersResponse struct {
	Data []struct {
		Root   string `json:"root"`
		Header struct {
			Message struct {
				Slot          uint64Str `json:"slot"`
				ProposerIndex uint64Str `json:"proposer_index"`
				ParentRoot    string    `json:"parent_root"`
				StateRoot     string    `json:"state_root"`
				BodyRoot      string    `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
	Finalized bool `json:"finalized"`
}

type StandardFinalityCheckpointsResponse struct {
	Data struct {
		PreviousJustified struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"finalized"`
	} `json:"data"`
}

type StreamedBlockEventData struct {
	Slot                uint64Str `json:"slot"`
	Block               string    `json:"block"`
	ExecutionOptimistic bool      `json:"execution_optimistic"`
}

type StandardProposerDuty struct {
	Pubkey         string    `json:"pubkey"`
	ValidatorIndex uint64Str `json:"validator_index"`
	Slot           uint64Str `json:"slot"`
}

type StandardProposerDutiesResponse struct {
	DependentRoot string                 `json:"dependent_root"`
	Data          []StandardProposerDuty `json:"data"`
}

type StandardCommitteeEntry struct {
	Index      uint64Str `json:"index"`
	Slot       uint64Str `json:"slot"`
	Validators []string  `json:"validators"`
}

type StandardCommitteesResponse struct {
	Data []StandardCommitteeEntry `json:"data"`
}

type StandardSyncCommittee struct {
	Validators          []string   `json:"validators"`
	ValidatorAggregates [][]string `json:"validator_aggregates"`
}

type StandardSyncCommitteesResponse struct {
	Data StandardSyncCommittee `json:"data"`
}

type ProposerSlashing struct {
	SignedHeader1 struct {
		Message struct {
			Slot          uint64Str `json:"slot"`
			ProposerIndex uint64Str `json:"proposer_index"`
			ParentRoot    string    `json:"parent_root"`
			StateRoot     string    `json:"state_root"`
			BodyRoot      string    `json:"body_root"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"signed_header_1"`
	SignedHeader2 struct {
		Message struct {
			Slot          uint64Str `json:"slot"`
			ProposerIndex uint64Str `json:"proposer_index"`
			ParentRoot    string    `json:"parent_root"`
			StateRoot     string    `json:"state_root"`
			BodyRoot      string    `json:"body_root"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"signed_header_2"`
}

type AttesterSlashing struct {
	Attestation1 struct {
		AttestingIndices []uint64Str `json:"attesting_indices"`
		Signature        string      `json:"signature"`
		Data             struct {
			Slot            uint64Str `json:"slot"`
			Index           uint64Str `json:"index"`
			BeaconBlockRoot string    `json:"beacon_block_root"`
			Source          struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"source"`
			Target struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"target"`
		} `json:"data"`
	} `json:"attestation_1"`
	Attestation2 struct {
		AttestingIndices []uint64Str `json:"attesting_indices"`
		Signature        string      `json:"signature"`
		Data             struct {
			Slot            uint64Str `json:"slot"`
			Index           uint64Str `json:"index"`
			BeaconBlockRoot string    `json:"beacon_block_root"`
			Source          struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"source"`
			Target struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"target"`
		} `json:"data"`
	} `json:"attestation_2"`
}

type Attestation struct {
	AggregationBits string `json:"aggregation_bits"`
	Signature       string `json:"signature"`
	Data            struct {
		Slot            uint64Str `json:"slot"`
		Index           uint64Str `json:"index"`
		BeaconBlockRoot string    `json:"beacon_block_root"`
		Source          struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"source"`
		Target struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"target"`
	} `json:"data"`
}

type Deposit struct {
	Proof []string `json:"proof"`
	Data  struct {
		Pubkey                string    `json:"pubkey"`
		WithdrawalCredentials string    `json:"withdrawal_credentials"`
		Amount                uint64Str `json:"amount"`
		Signature             string    `json:"signature"`
	} `json:"data"`
}

type VoluntaryExit struct {
	Message struct {
		Epoch          uint64Str `json:"epoch"`
		ValidatorIndex uint64Str `json:"validator_index"`
	} `json:"message"`
	Signature string `json:"signature"`
}

type Eth1Data struct {
	DepositRoot  string    `json:"deposit_root"`
	DepositCount uint64Str `json:"deposit_count"`
	BlockHash    string    `json:"block_hash"`
}

type SyncAggregate struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}

// https://ethereum.github.io/beacon-APIs/#/Beacon/getBlockV2
// https://github.com/ethereum/consensus-specs/blob/v1.1.9/specs/bellatrix/beacon-chain.md#executionpayload
type ExecutionPayload struct {
	ParentHash    bytesHexStr   `json:"parent_hash"`
	FeeRecipient  bytesHexStr   `json:"fee_recipient"`
	StateRoot     bytesHexStr   `json:"state_root"`
	ReceiptsRoot  bytesHexStr   `json:"receipts_root"`
	LogsBloom     bytesHexStr   `json:"logs_bloom"`
	PrevRandao    bytesHexStr   `json:"prev_randao"`
	BlockNumber   uint64Str     `json:"block_number"`
	GasLimit      uint64Str     `json:"gas_limit"`
	GasUsed       uint64Str     `json:"gas_used"`
	Timestamp     uint64Str     `json:"timestamp"`
	ExtraData     bytesHexStr   `json:"extra_data"`
	BaseFeePerGas uint64Str     `json:"base_fee_per_gas"`
	BlockHash     bytesHexStr   `json:"block_hash"`
	Transactions  []bytesHexStr `json:"transactions"`
	// present only after capella
	Withdrawals []WithdrawalPayload `json:"withdrawals"`
	// present only after deneb
	BlobGasUsed   uint64Str `json:"blob_gas_used"`
	ExcessBlobGas uint64Str `json:"excess_blob_gas"`
}

type WithdrawalPayload struct {
	Index          uint64Str   `json:"index"`
	ValidatorIndex uint64Str   `json:"validator_index"`
	Address        bytesHexStr `json:"address"`
	Amount         uint64Str   `json:"amount"`
}

type SignedBLSToExecutionChange struct {
	Message struct {
		ValidatorIndex     uint64Str   `json:"validator_index"`
		FromBlsPubkey      bytesHexStr `json:"from_bls_pubkey"`
		ToExecutionAddress bytesHexStr `json:"to_execution_address"`
	} `json:"message"`
	Signature bytesHexStr `json:"signature"`
}

type AnySignedBlock struct {
	Message struct {
		Slot          uint64Str `json:"slot"`
		ProposerIndex uint64Str `json:"proposer_index"`
		ParentRoot    string    `json:"parent_root"`
		StateRoot     string    `json:"state_root"`
		Body          struct {
			RandaoReveal      string             `json:"randao_reveal"`
			Eth1Data          Eth1Data           `json:"eth1_data"`
			Graffiti          string             `json:"graffiti"`
			ProposerSlashings []ProposerSlashing `json:"proposer_slashings"`
			AttesterSlashings []AttesterSlashing `json:"attester_slashings"`
			Attestations      []Attestation      `json:"attestations"`
			Deposits          []Deposit          `json:"deposits"`
			VoluntaryExits    []VoluntaryExit    `json:"voluntary_exits"`

			// not present in phase0 blocks
			SyncAggregate *SyncAggregate `json:"sync_aggregate,omitempty"`

			// not present in phase0/altair blocks
			ExecutionPayload *ExecutionPayload `json:"execution_payload"`

			// present only after capella
			SignedBLSToExecutionChange []*SignedBLSToExecutionChange `json:"bls_to_execution_changes"`

			// present only after deneb
			BlobKZGCommitments []bytesHexStr `json:"blob_kzg_commitments"`
		} `json:"body"`
	} `json:"message"`
	Signature bytesHexStr `json:"signature"`
}

type StandardV2BlockResponse struct {
	Version             string         `json:"version"`
	ExecutionOptimistic bool           `json:"execution_optimistic"`
	Finalized           bool           `json:"finalized"`
	Data                AnySignedBlock `json:"data"`
}

type StandardV1BlockRootResponse struct {
	Data struct {
		Root string `json:"root"`
	} `json:"data"`
}

type StandardValidatorEntry struct {
	Index     uint64Str `json:"index"`
	Balance   uint64Str `json:"balance"`
	Status    string    `json:"status"`
	Validator struct {
		Pubkey                     string    `json:"pubkey"`
		WithdrawalCredentials      string    `json:"withdrawal_credentials"`
		EffectiveBalance           uint64Str `json:"effective_balance"`
		Slashed                    bool      `json:"slashed"`
		ActivationEligibilityEpoch uint64Str `json:"activation_eligibility_epoch"`
		ActivationEpoch            uint64Str `json:"activation_epoch"`
		ExitEpoch                  uint64Str `json:"exit_epoch"`
		WithdrawableEpoch          uint64Str `json:"withdrawable_epoch"`
	} `json:"validator"`
}

type StandardValidatorsResponse struct {
	Data []StandardValidatorEntry `json:"data"`
}

type StandardSyncingResponse struct {
	Data struct {
		IsSyncing    bool      `json:"is_syncing"`
		HeadSlot     uint64Str `json:"head_slot"`
		SyncDistance uint64Str `json:"sync_distance"`
	} `json:"data"`
}

type StandardValidatorBalancesResponse struct {
	Data []struct {
		Index   uint64Str `json:"index"`
		Balance uint64Str `json:"balance"`
	} `json:"data"`
}

type StandardBlobSidecarsResponse struct {
	Data []struct {
		BlockRoot       bytesHexStr `json:"block_root"`
		Index           uint64Str   `json:"index"`
		Slot            uint64Str   `json:"slot"`
		BlockParentRoot bytesHexStr `json:"block_parent_root"`
		ProposerIndex   uint64Str   `json:"proposer_index"`
		KzgCommitment   bytesHexStr `json:"kzg_commitment"`
		KzgProof        bytesHexStr `json:"kzg_proof"`
		// Blob            string `json:"blob"`
	}
}
