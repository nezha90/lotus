package main

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var MethodMap = make(map[abi.MethodNum]struct{})

func init() {
	MethodMap[builtin.MethodsMiner.SubmitWindowedPoSt] = struct{}{}

	MethodMap[builtin.MethodsMiner.PreCommitSector] = struct{}{}
	MethodMap[builtin.MethodsMiner.PreCommitSectorBatch] = struct{}{}
	MethodMap[builtin.MethodsMiner.PreCommitSectorBatch2] = struct{}{}

	MethodMap[builtin.MethodsMiner.ProveCommitSector] = struct{}{}
	MethodMap[builtin.MethodsMiner.ProveCommitAggregate] = struct{}{}
	MethodMap[builtin.MethodsMiner.ProveCommitSectors3] = struct{}{}

	MethodMap[builtin.MethodsMiner.ExtendSectorExpiration] = struct{}{}
	MethodMap[builtin.MethodsMiner.ExtendSectorExpiration2] = struct{}{}

	MethodMap[builtin.MethodsMarket.PublishStorageDeals] = struct{}{}
	MethodMap[builtin.MethodsMarket.WithdrawBalance] = struct{}{}
	MethodMap[builtin.MethodsMarket.AddBalance] = struct{}{}

	MethodMap[builtin.MethodsVerifiedRegistry.ExtendClaimTerms] = struct{}{}

	MethodMap[builtin.MethodsMiner.ChangePeerID] = struct{}{}
	MethodMap[builtin.MethodsMiner.ChangeMultiaddrs] = struct{}{}
	MethodMap[builtin.MethodsMultisig.Approve] = struct{}{}
	MethodMap[builtin.MethodsMultisig.Propose] = struct{}{}

	MethodMap[builtin.MethodsMiner.DeclareFaultsRecovered] = struct{}{}

	MethodMap[builtin.MethodsMiner.TerminateSectors] = struct{}{}
}

func checkMethod(method abi.MethodNum) bool {
	if _, ok := MethodMap[method]; !ok {
		return false
	} else {
		return true
	}
}
