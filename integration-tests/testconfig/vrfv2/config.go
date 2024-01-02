package testconfig

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
	ErrDeviationShouldBeLessThanOriginal = "`RandomnessRequestCountPerRequestDeviation` should be less than `RandomnessRequestCountPerRequest`"
)

type Config struct {
	Common            *Common            `toml:"Common"`
	General           *General           `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig `toml:"ExistingEnv"`
	NewEnvConfig      *NewEnvConfig      `toml:"NewEnv"`
	Performance       *PerformanceConfig `toml:"Performance"`
}

func (c *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}
	if c.Common == nil && from.Common != nil {
		c.Common = from.Common
	} else if c.Common != nil && from.Common != nil {
		if err := c.Common.ApplyOverrides(from.Common); err != nil {
			return err
		}
	}
	if c.General == nil && from.General != nil {
		c.General = from.General
	} else if c.General != nil && from.General != nil {
		if err := c.General.ApplyOverrides(from.General); err != nil {
			return err
		}
	}
	if c.ExistingEnvConfig == nil && from.ExistingEnvConfig != nil {
		c.ExistingEnvConfig = from.ExistingEnvConfig
	} else if c.ExistingEnvConfig != nil && from.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.ApplyOverrides(from.ExistingEnvConfig); err != nil {
			return err
		}
	}
	if c.NewEnvConfig == nil && from.NewEnvConfig != nil {
		c.NewEnvConfig = from.NewEnvConfig
	} else if c.NewEnvConfig != nil && from.NewEnvConfig != nil {
		if err := c.NewEnvConfig.ApplyOverrides(from.NewEnvConfig); err != nil {
			return err
		}
	}
	if c.Performance == nil && from.Performance != nil {
		c.Performance = from.Performance
	} else if c.Performance != nil && from.Performance != nil {
		if err := c.Performance.ApplyOverrides(from.Performance); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Common != nil {
		if err := c.Common.Validate(); err != nil {
			return err
		}
	}
	if c.General != nil {
		if err := c.General.Validate(); err != nil {
			return err
		}
	}
	if c.Performance != nil {
		if err := c.Performance.Validate(); err != nil {
			return err
		}
		if *c.Performance.UseExistingEnv {
			if c.ExistingEnvConfig != nil {
				if err := c.ExistingEnvConfig.Validate(); err != nil {
					return err
				}
			}
		} else {
			if c.NewEnvConfig != nil {
				if err := c.NewEnvConfig.Validate(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type Common struct {
	CancelSubsAfterTestRun *bool `toml:"cancel_subs_after_test_run"`
}

func (c *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}
	if from.CancelSubsAfterTestRun != nil {
		c.CancelSubsAfterTestRun = from.CancelSubsAfterTestRun
	}

	return nil
}

func (c *Common) Validate() error {
	return nil
}

type PerformanceConfig struct {
	TestDuration          *models.Duration `toml:"test_duration"`
	RPS                   *int64           `toml:"rps"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`

	// Using existing environment and contracts
	UseExistingEnv     *bool `toml:"use_existing_env"`
	CoordinatorAddress *string
	ConsumerAddress    *string
	LinkAddress        *string
	SubID              *uint64
	KeyHash            *string
}

func (c *PerformanceConfig) ApplyOverrides(from *PerformanceConfig) error {
	if from == nil {
		return nil
	}
	if from.TestDuration != nil {
		c.TestDuration = from.TestDuration
	}
	if from.RPS != nil {
		c.RPS = from.RPS
	}
	if from.RateLimitUnitDuration != nil {
		c.RateLimitUnitDuration = from.RateLimitUnitDuration
	}
	if from.UseExistingEnv != nil {
		c.UseExistingEnv = from.UseExistingEnv
	}

	return nil
}

func (c *PerformanceConfig) Validate() error {
	if c.TestDuration == nil || c.TestDuration.Duration() == 0 {
		return errors.New("test_duration must be set to a positive value")
	}
	if c.RPS == nil || *c.RPS == 0 {
		return errors.New("rps must be set to a positive value")
	}
	if c.RateLimitUnitDuration == nil {
		return errors.New("rate_limit_unit_duration must be set ")
	}
	if c.UseExistingEnv == nil {
		return errors.New("use_existing_env must be set ")
	}

	return nil
}

type ExistingEnvConfig struct {
	CoordinatorAddress            *string  `toml:"coordinator_address"`
	ConsumerAddress               *string  `toml:"consumer_address"`
	LinkAddress                   *string  `toml:"link_address"`
	SubID                         *uint64  `toml:"sub_id"`
	KeyHash                       *string  `toml:"key_hash"`
	CreateFundSubsAndAddConsumers *bool    `toml:"create_fund_subs_and_add_consumers"`
	NodeSendingKeys               []string `toml:"node_sending_keys"`
	*Funding
}

func (c *ExistingEnvConfig) ApplyOverrides(from *ExistingEnvConfig) error {
	if from == nil {
		return nil
	}
	if from.CoordinatorAddress != nil {
		c.CoordinatorAddress = from.CoordinatorAddress
	}
	if from.ConsumerAddress != nil {
		c.ConsumerAddress = from.ConsumerAddress
	}
	if from.LinkAddress != nil {
		c.LinkAddress = from.LinkAddress
	}
	if from.SubID != nil {
		c.SubID = from.SubID
	}
	if from.KeyHash != nil {
		c.KeyHash = from.KeyHash
	}
	if from.CreateFundSubsAndAddConsumers != nil {
		c.CreateFundSubsAndAddConsumers = from.CreateFundSubsAndAddConsumers
	}
	if from.NodeSendingKeys != nil {
		c.NodeSendingKeys = from.NodeSendingKeys
	}
	if from.Funding != nil && c.Funding == nil {
		c.Funding = from.Funding
	} else if from.Funding != nil && c.Funding != nil {
		if err := c.Funding.ApplyOverrides(from.Funding); err != nil {
			return err
		}
	}

	return nil
}

func (c *ExistingEnvConfig) Validate() error {
	if c.CreateFundSubsAndAddConsumers == nil {
		return errors.New("create_fund_subs_and_add_consumers must be set ")
	}
	if c.CoordinatorAddress == nil {
		return errors.New("coordinator_address must be set when using existing environment")
	}
	if !common.IsHexAddress(*c.CoordinatorAddress) {
		return errors.New("coordinator_address must be a valid hex address")
	}
	if c.KeyHash == nil {
		return errors.New("key_hash must be set when using existing environment")
	}
	if *c.KeyHash == "" {
		return errors.New("key_hash must be a non-empty string")
	}
	if c.LinkAddress != nil && !common.IsHexAddress(*c.LinkAddress) {
		return errors.New("link_address must be a valid hex address")
	}

	if *c.CreateFundSubsAndAddConsumers {
		if c.ConsumerAddress == nil {
			return errors.New("consumer_address must be set when using existing environment")
		}
		if !common.IsHexAddress(*c.ConsumerAddress) {
			return errors.New("consumer_address must be a valid hex address")
		}
		if c.SubID == nil {
			return errors.New("sub_id must be set when using existing environment")
		}
		if *c.SubID == 0 {
			return errors.New("sub_id must be a positive value")
		}
		if c.Funding != nil {
			if err := c.Funding.Validate(); err != nil {
				return err
			}
		}
		if c.Funding == nil && c.Funding.SubFunding == nil {
			return errors.New("sub_funds_link must be set when using existing environment")
		}
		if err := c.Funding.SubFunding.Validate(); err != nil {
			return err
		}
	}

	if c.NodeSendingKeys != nil {
		for _, key := range c.NodeSendingKeys {
			if !common.IsHexAddress(key) {
				return errors.New("node_sending_keys must be a valid hex address")
			}
		}
	}

	return nil
}

type NewEnvConfig struct {
	*Funding
}

func (c *NewEnvConfig) ApplyOverrides(from *NewEnvConfig) error {
	if from == nil {
		return nil
	}

	return c.Funding.ApplyOverrides(from.Funding)
}

func (c *NewEnvConfig) Validate() error {
	if c.Funding != nil {
		return c.Funding.Validate()
	}

	return nil
}

type Funding struct {
	*SubFunding
	NodeSendingKeyFunding    *float64 `toml:"node_sending_key_funding"`
	NodeSendingKeyFundingMin *float64 `toml:"node_sending_key_funding_min"`
}

func (c *Funding) ApplyOverrides(from *Funding) error {
	if from == nil {
		return nil
	}
	if from.NodeSendingKeyFunding != nil {
		c.NodeSendingKeyFunding = from.NodeSendingKeyFunding
	}
	if from.NodeSendingKeyFundingMin != nil {
		c.NodeSendingKeyFundingMin = from.NodeSendingKeyFundingMin
	}
	if from.SubFundsLink != nil && c.SubFundsLink == nil {
		c.SubFundsLink = from.SubFundsLink
	} else if from.SubFundsLink != nil && c.SubFundsLink != nil {
		if err := c.SubFunding.ApplyOverrides(from.SubFunding); err != nil {
			return err
		}
	}

	return nil
}

func (c *Funding) Validate() error {
	if c.NodeSendingKeyFunding != nil && *c.NodeSendingKeyFunding <= 0 {
		return errors.New("when set node_sending_key_funding must be a positive value")
	}
	if c.NodeSendingKeyFundingMin != nil && *c.NodeSendingKeyFundingMin <= 0 {
		return errors.New("when set node_sending_key_funding_min must be a positive value")
	}

	return nil
}

type SubFunding struct {
	SubFundsLink *float64 `toml:"sub_funds_link"`
}

func (c *SubFunding) ApplyOverrides(from *SubFunding) error {
	if from == nil {
		return nil
	}
	if from.SubFundsLink != nil {
		c.SubFundsLink = from.SubFundsLink
	}

	return nil
}

func (c *SubFunding) Validate() error {
	if c.SubFundsLink == nil || *c.SubFundsLink == 0 {
		return errors.New("when set sub_funds_link must be a positive value")
	}

	return nil
}

type General struct {
	CLNodeMaxGasPriceGWei          *int64   `toml:"max_gas_price_gwei"`                  // Max gas price in GWei for the chainlink node default:"1000"
	LinkNativeFeedResponse         *int64   `toml:"link_native_feed_response"`           // Response of the LINK/ETH feed default:"1000000000000000000"
	MinimumConfirmations           *uint16  `toml:"minimum_confirmations" `              // Minimum number of confirmations for the VRF Coordinator default:"3"
	SubscriptionFundingAmountLink  *float64 `toml:"subscription_funding_amount_link"`    // Amount of LINK to fund the subscription with default:"5"
	NumberOfWords                  *uint32  `toml:"number_of_words" `                    // Number of words to request default:"3"
	CallbackGasLimit               *uint32  `toml:"callback_gas_limit" `                 // Gas limit for the callback default:"1000000"
	MaxGasLimitCoordinatorConfig   *uint32  `toml:"max_gas_limit_coordinator_config"`    // Max gas limit for the VRF Coordinator config  default:"2500000"
	FallbackWeiPerUnitLink         *int64   `toml:"fallback_wei_per_unit_link"`          // Fallback wei per unit LINK for the VRF Coordinator config  default:"60000000000000000"
	StalenessSeconds               *uint32  `toml:"staleness_seconds" `                  // Staleness in seconds for the VRF Coordinator config default:"86400"
	GasAfterPaymentCalculation     *uint32  `toml:"gas_after_payment_calculation" `      // Gas after payment calculation for the VRF Coordinator config default:"33825"
	FulfillmentFlatFeeLinkPPMTier1 *uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_1"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier2 *uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_2"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier3 *uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_3"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier4 *uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_4"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier5 *uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_5"` //default:"500"
	ReqsForTier2                   *int64   `toml:"reqs_for_tier_2"`                     // default:"0"
	ReqsForTier3                   *int64   `toml:"reqs_for_tier_3"`                     // default:"0"
	ReqsForTier4                   *int64   `toml:"reqs_for_tier_4"`                     // default:"0"
	ReqsForTier5                   *int64   `toml:"reqs_for_tier_5"`                     // default:"0"

	NumberOfSubToCreate *int `toml:"number_of_sub_to_create"` // Number of subscriptions to create default:"1"

	RandomnessRequestCountPerRequest          *uint16 `toml:"randomness_request_count_per_request"`           // How many randomness requests to send per request default:"1"
	RandomnessRequestCountPerRequestDeviation *uint16 `toml:"randomness_request_count_per_request_deviation"` // How many randomness requests to send per request  default:"0"

	RandomWordsFulfilledEventTimeout *models.Duration `toml:"random_words_fulfilled_event_timeout"` // How long to wait for the RandomWordsFulfilled event to be emitted default:"2m"

	// Wrapper Config
	WrapperGasOverhead                      *uint32  `toml:"wrapped_gas_overhead"`                         // default:"50000"
	CoordinatorGasOverhead                  *uint32  `toml:"coordinator_gas_overhead"`                     // default:"52000"
	WrapperPremiumPercentage                *uint8   `toml:"wrapper_premium_percentage"`                   // default:"25"
	WrapperMaxNumberOfWords                 *uint8   `toml:"wrapper_max_number_of_words"`                  // default:"10"
	WrapperConsumerFundingAmountNativeToken *float64 `toml:"wrapper_consumer_funding_amount_native_token"` // default:"1"
	WrapperConsumerFundingAmountLink        *int64   `toml:"wrapper_consumer_funding_amount_link"`         // default:"10"
}

func (c *General) ApplyOverrides(from *General) error {
	if from == nil {
		return nil
	}
	if from.CLNodeMaxGasPriceGWei != nil {
		c.CLNodeMaxGasPriceGWei = from.CLNodeMaxGasPriceGWei
	}
	if from.LinkNativeFeedResponse != nil {
		c.LinkNativeFeedResponse = from.LinkNativeFeedResponse
	}
	if from.MinimumConfirmations != nil {
		c.MinimumConfirmations = from.MinimumConfirmations
	}
	if from.SubscriptionFundingAmountLink != nil {
		c.SubscriptionFundingAmountLink = from.SubscriptionFundingAmountLink
	}
	if from.NumberOfWords != nil {
		c.NumberOfWords = from.NumberOfWords
	}
	if from.CallbackGasLimit != nil {
		c.CallbackGasLimit = from.CallbackGasLimit
	}
	if from.MaxGasLimitCoordinatorConfig != nil {
		c.MaxGasLimitCoordinatorConfig = from.MaxGasLimitCoordinatorConfig
	}
	if from.FallbackWeiPerUnitLink != nil {
		c.FallbackWeiPerUnitLink = from.FallbackWeiPerUnitLink
	}
	if from.StalenessSeconds != nil {
		c.StalenessSeconds = from.StalenessSeconds
	}
	if from.GasAfterPaymentCalculation != nil {
		c.GasAfterPaymentCalculation = from.GasAfterPaymentCalculation
	}
	if from.FulfillmentFlatFeeLinkPPMTier1 != nil {
		c.FulfillmentFlatFeeLinkPPMTier1 = from.FulfillmentFlatFeeLinkPPMTier1
	}
	if from.FulfillmentFlatFeeLinkPPMTier2 != nil {
		c.FulfillmentFlatFeeLinkPPMTier2 = from.FulfillmentFlatFeeLinkPPMTier2
	}
	if from.FulfillmentFlatFeeLinkPPMTier3 != nil {
		c.FulfillmentFlatFeeLinkPPMTier3 = from.FulfillmentFlatFeeLinkPPMTier3
	}
	if from.FulfillmentFlatFeeLinkPPMTier4 != nil {
		c.FulfillmentFlatFeeLinkPPMTier4 = from.FulfillmentFlatFeeLinkPPMTier4
	}
	if from.FulfillmentFlatFeeLinkPPMTier5 != nil {
		c.FulfillmentFlatFeeLinkPPMTier5 = from.FulfillmentFlatFeeLinkPPMTier5
	}
	if from.ReqsForTier2 != nil {
		c.ReqsForTier2 = from.ReqsForTier2
	}
	if from.ReqsForTier3 != nil {
		c.ReqsForTier3 = from.ReqsForTier3
	}
	if from.ReqsForTier4 != nil {
		c.ReqsForTier4 = from.ReqsForTier4
	}
	if from.ReqsForTier5 != nil {
		c.ReqsForTier5 = from.ReqsForTier5
	}
	if from.NumberOfSubToCreate != nil {
		c.NumberOfSubToCreate = from.NumberOfSubToCreate
	}
	if from.RandomnessRequestCountPerRequest != nil {
		c.RandomnessRequestCountPerRequest = from.RandomnessRequestCountPerRequest
	}
	if from.RandomnessRequestCountPerRequestDeviation != nil {
		c.RandomnessRequestCountPerRequestDeviation = from.RandomnessRequestCountPerRequestDeviation
	}
	if from.RandomWordsFulfilledEventTimeout != nil {
		c.RandomWordsFulfilledEventTimeout = from.RandomWordsFulfilledEventTimeout
	}
	if from.WrapperGasOverhead != nil {
		c.WrapperGasOverhead = from.WrapperGasOverhead
	}
	if from.CoordinatorGasOverhead != nil {
		c.CoordinatorGasOverhead = from.CoordinatorGasOverhead
	}
	if from.WrapperPremiumPercentage != nil {
		c.WrapperPremiumPercentage = from.WrapperPremiumPercentage
	}
	if from.WrapperMaxNumberOfWords != nil {
		c.WrapperMaxNumberOfWords = from.WrapperMaxNumberOfWords
	}
	if from.WrapperConsumerFundingAmountNativeToken != nil {
		c.WrapperConsumerFundingAmountNativeToken = from.WrapperConsumerFundingAmountNativeToken
	}
	if from.WrapperConsumerFundingAmountLink != nil {
		c.WrapperConsumerFundingAmountLink = from.WrapperConsumerFundingAmountLink
	}

	return nil
}

func (c *General) Validate() error {
	if c.CLNodeMaxGasPriceGWei == nil || *c.CLNodeMaxGasPriceGWei == 0 {
		return errors.New("max_gas_price_gwei must be set to a positive value")
	}
	if c.LinkNativeFeedResponse == nil || *c.LinkNativeFeedResponse == 0 {
		return errors.New("link_native_feed_response must be set to a positive value")
	}
	if c.MinimumConfirmations == nil {
		return errors.New("minimum_confirmations must be set to a non-negative value")
	}
	if c.SubscriptionFundingAmountLink == nil || *c.SubscriptionFundingAmountLink == 0 {
		return errors.New("subscription_funding_amount_link must be set to a positive value")
	}
	if c.NumberOfWords == nil || *c.NumberOfWords == 0 {
		return errors.New("number_of_words must be set to a positive value")
	}
	if c.CallbackGasLimit == nil || *c.CallbackGasLimit == 0 {
		return errors.New("callback_gas_limit must be set to a positive value")
	}
	if c.MaxGasLimitCoordinatorConfig == nil || *c.MaxGasLimitCoordinatorConfig == 0 {
		return errors.New("max_gas_limit_coordinator_config must be set to a positive value")
	}
	if c.FallbackWeiPerUnitLink == nil || *c.FallbackWeiPerUnitLink == 0 {
		return errors.New("fallback_wei_per_unit_link must be set to a positive value")
	}
	if c.StalenessSeconds == nil || *c.StalenessSeconds == 0 {
		return errors.New("staleness_seconds must be set to a positive value")
	}
	if c.GasAfterPaymentCalculation == nil || *c.GasAfterPaymentCalculation == 0 {
		return errors.New("gas_after_payment_calculation must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier1 == nil || *c.FulfillmentFlatFeeLinkPPMTier1 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_1 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier2 == nil || *c.FulfillmentFlatFeeLinkPPMTier2 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_2 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier3 == nil || *c.FulfillmentFlatFeeLinkPPMTier3 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_3 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier4 == nil || *c.FulfillmentFlatFeeLinkPPMTier4 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_4 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier5 == nil || *c.FulfillmentFlatFeeLinkPPMTier5 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_5 must be set to a positive value")
	}
	if c.ReqsForTier2 == nil || *c.ReqsForTier2 < 0 {
		return errors.New("reqs_for_tier_2 must be set to a non-negative value")
	}
	if c.ReqsForTier3 == nil || *c.ReqsForTier3 < 0 {
		return errors.New("reqs_for_tier_3 must be set to a non-negative value")
	}
	if c.ReqsForTier4 == nil || *c.ReqsForTier4 < 0 {
		return errors.New("reqs_for_tier_4 must be set to a non-negative value")
	}
	if c.ReqsForTier5 == nil || *c.ReqsForTier5 < 0 {
		return errors.New("reqs_for_tier_5 must be set to a non-negative value")
	}
	if c.NumberOfSubToCreate == nil || *c.NumberOfSubToCreate == 0 {
		return errors.New("number_of_sub_to_create must be set to a positive value")
	}
	if c.RandomnessRequestCountPerRequest == nil || *c.RandomnessRequestCountPerRequest == 0 {
		return errors.New("randomness_request_count_per_request must be set to a positive value")
	}
	if c.RandomnessRequestCountPerRequestDeviation == nil {
		return errors.New("randomness_request_count_per_request_deviation must be set to a non-negative value")
	}
	if c.RandomWordsFulfilledEventTimeout == nil || c.RandomWordsFulfilledEventTimeout.Duration() == 0 {
		return errors.New("random_words_fulfilled_event_timeout must be set to a positive value")
	}
	if c.WrapperGasOverhead == nil {
		return errors.New("wrapped_gas_overhead must be set to a non-negative value")
	}
	if c.CoordinatorGasOverhead == nil || *c.CoordinatorGasOverhead == 0 {
		return errors.New("coordinator_gas_overhead must be set to a non-negative value")
	}
	if c.WrapperPremiumPercentage == nil || *c.WrapperPremiumPercentage == 0 {
		return errors.New("wrapper_premium_percentage must be set to a positive value")
	}
	if c.WrapperMaxNumberOfWords == nil || *c.WrapperMaxNumberOfWords == 0 {
		return errors.New("wrapper_max_number_of_words must be set to a positive value")
	}
	if c.WrapperConsumerFundingAmountNativeToken == nil || *c.WrapperConsumerFundingAmountNativeToken < 0 {
		return errors.New("wrapper_consumer_funding_amount_native_token must be set to a non-negative value")
	}
	if c.WrapperConsumerFundingAmountLink == nil || *c.WrapperConsumerFundingAmountLink < 0 {
		return errors.New("wrapper_consumer_funding_amount_link must be set to a non-negative value")
	}
	if *c.RandomnessRequestCountPerRequest <= *c.RandomnessRequestCountPerRequestDeviation {
		return errors.New(ErrDeviationShouldBeLessThanOriginal)
	}

	return nil
}
