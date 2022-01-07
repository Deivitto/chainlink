package fluxmonitorv2

import (
	"math/big"

	"github.com/pkg/errors"
	eth "github.com/smartcontractkit/chainlink/core/chains/evm/eth"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --name ContractSubmitter --output ./mocks/ --case=underscore

// FluxAggregatorABI initializes the Flux Aggregator ABI
var FluxAggregatorABI = eth.MustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

// ContractSubmitter defines an interface to submit an eth tx.
type ContractSubmitter interface {
	Submit(roundID *big.Int, submission *big.Int, qopts ...pg.QOpt) error
}

// FluxAggregatorContractSubmitter submits the polled answer in an eth tx.
type FluxAggregatorContractSubmitter struct {
	flux_aggregator_wrapper.FluxAggregatorInterface
	orm      ORM
	keyStore KeyStoreInterface
	gasLimit uint64
}

// NewFluxAggregatorContractSubmitter constructs a new NewFluxAggregatorContractSubmitter
func NewFluxAggregatorContractSubmitter(
	contract flux_aggregator_wrapper.FluxAggregatorInterface,
	orm ORM,
	keyStore KeyStoreInterface,
	gasLimit uint64,
) *FluxAggregatorContractSubmitter {
	return &FluxAggregatorContractSubmitter{
		FluxAggregatorInterface: contract,
		orm:                     orm,
		keyStore:                keyStore,
		gasLimit:                gasLimit,
	}
}

// Submit submits the answer by writing a EthTx for the bulletprooftxmanager to
// pick up
func (c *FluxAggregatorContractSubmitter) Submit(roundID *big.Int, submission *big.Int, qopts ...pg.QOpt) error {
	fromAddress, err := c.keyStore.GetRoundRobinAddress()
	if err != nil {
		return err
	}

	payload, err := FluxAggregatorABI.Pack("submit", roundID, submission)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(
		c.orm.CreateEthTransaction(fromAddress, c.Address(), payload, c.gasLimit, qopts...),
		"failed to send Eth transaction",
	)
}
