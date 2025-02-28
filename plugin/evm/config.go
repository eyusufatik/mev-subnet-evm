// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
)

const (
	defaultAcceptorQueueLimit                     = 64 // Provides 2 minutes of buffer (2s block target) for a commit delay
	defaultPruningEnabled                         = true
	defaultCommitInterval                         = 4096
	defaultSnapshotAsync                          = true
	defaultRpcGasCap                              = 50_000_000 // Default to 50M Gas Limit
	defaultRpcTxFeeCap                            = 100        // 100 AVAX
	defaultMetricsExpensiveEnabled                = false
	defaultApiMaxDuration                         = 0 // Default to no maximum API call duration
	defaultWsCpuRefillRate                        = 0 // Default to no maximum WS CPU usage
	defaultWsCpuMaxStored                         = 0 // Default to no maximum WS CPU usage
	defaultMaxBlocksPerRequest                    = 0 // Default to no maximum on the number of blocks per getLogs request
	defaultContinuousProfilerFrequency            = 15 * time.Minute
	defaultContinuousProfilerMaxFiles             = 5
	defaultRegossipFrequency                      = 1 * time.Minute
	defaultRegossipMaxTxs                         = 16
	defaultRegossipTxsPerAddress                  = 1
	defaultPriorityRegossipFrequency              = 1 * time.Second
	defaultPriorityRegossipMaxTxs                 = 32
	defaultPriorityRegossipTxsPerAddress          = 16
	defaultOfflinePruningBloomFilterSize   uint64 = 512 // Default size (MB) for the offline pruner to use
	defaultLogLevel                               = "info"
	defaultLogJSONFormat                          = false
	defaultMaxOutboundActiveRequests              = 8
	defaultPopulateMissingTriesParallelism        = 1024
)

var defaultEnabledAPIs = []string{
	"public-eth",
	"public-eth-filter",
	"net",
	"web3",
	"internal-public-eth",
	"internal-public-blockchain",
	"internal-public-tx-pool",
	"internal-public-transaction-pool",
}

type Duration struct {
	time.Duration
}

// Config ...
type Config struct {
	// Subnet EVM APIs
	SnowmanAPIEnabled bool   `json:"snowman-api-enabled"`
	AdminAPIEnabled   bool   `json:"admin-api-enabled"`
	AdminAPIDir       string `json:"admin-api-dir"`

	// EnabledEthAPIs is a list of Ethereum services that should be enabled
	// If none is specified, then we use the default list [defaultEnabledAPIs]
	EnabledEthAPIs []string `json:"eth-apis"`

	// Continuous Profiler
	ContinuousProfilerDir       string   `json:"continuous-profiler-dir"`       // If set to non-empty string creates a continuous profiler
	ContinuousProfilerFrequency Duration `json:"continuous-profiler-frequency"` // Frequency to run continuous profiler if enabled
	ContinuousProfilerMaxFiles  int      `json:"continuous-profiler-max-files"` // Maximum number of files to maintain

	// Gas/Price Caps
	RPCGasCap   uint64  `json:"rpc-gas-cap"`
	RPCTxFeeCap float64 `json:"rpc-tx-fee-cap"`

	// Eth Settings
	Preimages      bool `json:"preimages-enabled"`
	SnapshotAsync  bool `json:"snapshot-async"`
	SnapshotVerify bool `json:"snapshot-verification-enabled"`

	// Pruning Settings
	Pruning                         bool    `json:"pruning-enabled"`                    // If enabled, trie roots are only persisted every 4096 blocks
	AcceptorQueueLimit              int     `json:"accepted-queue-limit"`               // Maximum blocks to queue before blocking during acceptance
	CommitInterval                  uint64  `json:"commit-interval"`                    // Specifies the commit interval at which to persist EVM and atomic tries.
	AllowMissingTries               bool    `json:"allow-missing-tries"`                // If enabled, warnings preventing an incomplete trie index are suppressed
	PopulateMissingTries            *uint64 `json:"populate-missing-tries,omitempty"`   // Sets the starting point for re-populating missing tries. Disables re-generation if nil.
	PopulateMissingTriesParallelism int     `json:"populate-missing-tries-parallelism"` // Number of concurrent readers to use when re-populating missing tries on startup.

	// Metric Settings
	MetricsExpensiveEnabled bool `json:"metrics-expensive-enabled"` // Debug-level metrics that might impact runtime performance

	// API Settings
	LocalTxsEnabled         bool     `json:"local-txs-enabled"`
	APIMaxDuration          Duration `json:"api-max-duration"`
	WSCPURefillRate         Duration `json:"ws-cpu-refill-rate"`
	WSCPUMaxStored          Duration `json:"ws-cpu-max-stored"`
	MaxBlocksPerRequest     int64    `json:"api-max-blocks-per-request"`
	AllowUnfinalizedQueries bool     `json:"allow-unfinalized-queries"`
	AllowUnprotectedTxs     bool     `json:"allow-unprotected-txs"`

	// Keystore Settings
	KeystoreDirectory             string `json:"keystore-directory"` // both absolute and relative supported
	KeystoreExternalSigner        string `json:"keystore-external-signer"`
	KeystoreInsecureUnlockAllowed bool   `json:"keystore-insecure-unlock-allowed"`

	// Gossip Settings
	RemoteGossipOnlyEnabled       bool             `json:"remote-gossip-only-enabled"`
	RegossipFrequency             Duration         `json:"regossip-frequency"`
	RegossipMaxTxs                int              `json:"regossip-max-txs"`
	RegossipTxsPerAddress         int              `json:"regossip-txs-per-address"`
	PriorityRegossipFrequency     Duration         `json:"priority-regossip-frequency"`
	PriorityRegossipMaxTxs        int              `json:"priority-regossip-max-txs"`
	PriorityRegossipTxsPerAddress int              `json:"priority-regossip-txs-per-address"`
	PriorityRegossipAddresses     []common.Address `json:"priority-regossip-addresses"`

	// Log
	LogLevel      string `json:"log-level"`
	LogJSONFormat bool   `json:"log-json-format"`

	// Address for Tx Fees (must be empty if not supported by blockchain)
	FeeRecipient string `json:"feeRecipient"`

	// Offline Pruning Settings
	OfflinePruning                bool   `json:"offline-pruning-enabled"`
	OfflinePruningBloomFilterSize uint64 `json:"offline-pruning-bloom-filter-size"`
	OfflinePruningDataDirectory   string `json:"offline-pruning-data-directory"`

	// VM2VM network
	MaxOutboundActiveRequests int64 `json:"max-outbound-active-requests"`
}

// EthAPIs returns an array of strings representing the Eth APIs that should be enabled
func (c Config) EthAPIs() []string {
	return c.EnabledEthAPIs
}

func (c Config) EthBackendSettings() eth.Settings {
	return eth.Settings{MaxBlocksPerRequest: c.MaxBlocksPerRequest}
}

func (c *Config) SetDefaults() {
	c.EnabledEthAPIs = defaultEnabledAPIs
	c.RPCGasCap = defaultRpcGasCap
	c.RPCTxFeeCap = defaultRpcTxFeeCap
	c.MetricsExpensiveEnabled = defaultMetricsExpensiveEnabled
	c.APIMaxDuration.Duration = defaultApiMaxDuration
	c.WSCPURefillRate.Duration = defaultWsCpuRefillRate
	c.WSCPUMaxStored.Duration = defaultWsCpuMaxStored
	c.MaxBlocksPerRequest = defaultMaxBlocksPerRequest
	c.ContinuousProfilerFrequency.Duration = defaultContinuousProfilerFrequency
	c.ContinuousProfilerMaxFiles = defaultContinuousProfilerMaxFiles
	c.Pruning = defaultPruningEnabled
	c.AcceptorQueueLimit = defaultAcceptorQueueLimit
	c.CommitInterval = defaultCommitInterval
	c.SnapshotAsync = defaultSnapshotAsync
	c.RegossipFrequency.Duration = defaultRegossipFrequency
	c.RegossipMaxTxs = defaultRegossipMaxTxs
	c.RegossipTxsPerAddress = defaultRegossipTxsPerAddress
	c.PriorityRegossipFrequency.Duration = defaultPriorityRegossipFrequency
	c.PriorityRegossipMaxTxs = defaultPriorityRegossipMaxTxs
	c.PriorityRegossipTxsPerAddress = defaultPriorityRegossipTxsPerAddress
	c.OfflinePruningBloomFilterSize = defaultOfflinePruningBloomFilterSize
	c.LogLevel = defaultLogLevel
	c.LogJSONFormat = defaultLogJSONFormat
	c.MaxOutboundActiveRequests = defaultMaxOutboundActiveRequests
	c.PopulateMissingTriesParallelism = defaultPopulateMissingTriesParallelism
}

func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.Duration, err = cast.ToDurationE(v)
	return err
}

// String implements the stringer interface.
func (d Duration) String() string {
	return d.Duration.String()
}

// String implements the stringer interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// Validate returns an error if this is an invalid config.
func (c *Config) Validate() error {
	if c.PopulateMissingTries != nil && (c.OfflinePruning || c.Pruning) {
		return fmt.Errorf("cannot enable populate missing tries while offline pruning (enabled: %t)/pruning (enabled: %t) are enabled", c.OfflinePruning, c.Pruning)
	}
	if c.PopulateMissingTries != nil && c.PopulateMissingTriesParallelism < 1 {
		return fmt.Errorf("cannot enable populate missing tries without at least one reader (parallelism: %d)", c.PopulateMissingTriesParallelism)
	}

	if !c.Pruning && c.OfflinePruning {
		return fmt.Errorf("cannot run offline pruning while pruning is disabled")
	}

	// If pruning is enabled, the commit interval must be non-zero so the node commits state tries every CommitInterval blocks.
	if c.Pruning && c.CommitInterval == 0 {
		return fmt.Errorf("cannot use commit interval of 0 with pruning enabled")
	}
	return nil
}
