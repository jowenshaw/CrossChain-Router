package tron

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/anyswap/CrossChain-Router/v3/common"
	"github.com/anyswap/CrossChain-Router/v3/log"
	"github.com/anyswap/CrossChain-Router/v3/params"
	"github.com/anyswap/CrossChain-Router/v3/router"
	"github.com/anyswap/CrossChain-Router/v3/rpc/client"
	"github.com/anyswap/CrossChain-Router/v3/tokens"
)

var (
	// ensure Bridge impl tokens.CrossChainBridge
	_ tokens.IBridge = &Bridge{}
)

var TronMainnetChainID = uint64(112233)
var TronShastaChainID = uint64(2494104990)

// Bridge eth bridge
type Bridge struct {
	CustomConfig
	*tokens.CrossChainBridgeBase
	SignerChainID *big.Int
	TronChainID   *big.Int
}

// NewCrossChainBridge new bridge
func NewCrossChainBridge() *Bridge {
	return &Bridge{
		CustomConfig:         NewCustomConfig(),
		CrossChainBridgeBase: tokens.NewCrossChainBridgeBase(),
	}
}

// SupportsChainID supports chainID
func SupportsChainID(chainID *big.Int) bool {
	switch chainID.Uint64() {
	case TronMainnetChainID, TronShastaChainID:
		return true
	default:
		return false
	}
}

// CustomConfig custom config
type CustomConfig struct {
	// some chain's rpc is slow and need config a longer rpc timeout
	RPCClientTimeout int
	// eg. RSK chain do not check mixed case or not same as eth
	DontCheckAddressMixedCase bool
}

// NewCustomConfig new custom config
func NewCustomConfig() CustomConfig {
	return CustomConfig{
		RPCClientTimeout:          client.GetDefaultTimeout(false),
		DontCheckAddressMixedCase: false,
	}
}

// InitAfterConfig init variables (ie. extra members) after loading config
func (b *Bridge) InitAfterConfig() {
	chainID, err := common.GetBigIntFromStr(b.ChainConfig.ChainID)
	if err != nil {
		log.Fatal("wrong chainID", "chainID", b.ChainConfig.ChainID, "blockChain", b.ChainConfig.BlockChain)
	}
	switch chainID.Uint64() {
	case TronMainnetChainID, TronShastaChainID:
		b.TronChainID = chainID
	default:
		log.Fatal("wrong chainID")
	}
	b.InitExtraCustoms()
}

// InitRouterInfo init router info
func (b *Bridge) InitRouterInfo(routerContract, routerVersion string) (err error) {
	if routerContract == "" {
		return nil
	}
	var routerWNative string
	if tokens.IsERC20Router() {
		routerWNative, err = b.GetWNativeAddress(routerContract)
		if err != nil {
			log.Warn("get router wNative address failed", "routerContract", routerContract, "err", err)
		}
	}
	routerMPC, err := b.GetMPCAddress(routerContract)
	if err != nil {
		log.Warn("get router mpc address failed", "routerContract", routerContract, "err", err)
		return err
	}
	if !b.IsValidAddress(routerMPC) {
		log.Warn("get router mpc address return an invalid address", "routerContract", routerContract, "routerMPC", routerMPC)
		return fmt.Errorf("invalid router mpc address")
	}
	log.Info("get router mpc address success", "routerContract", routerContract, "routerMPC", routerMPC)
	routerMPCPubkey, err := router.GetMPCPubkey(routerMPC)
	if err != nil {
		log.Warn("get mpc public key failed", "mpc", routerMPC, "err", err)
		return err
	}
	if err = VerifyMPCPubKey(routerMPC, routerMPCPubkey); err != nil {
		log.Warn("verify mpc public key failed", "mpc", routerMPC, "mpcPubkey", routerMPCPubkey, "err", err)
		return err
	}
	chainID := b.ChainConfig.ChainID
	router.SetRouterInfo(
		routerContract,
		chainID,
		&router.SwapRouterInfo{
			RouterMPC:     routerMPC,
			RouterWNative: routerWNative,
		},
	)
	router.SetMPCPublicKey(routerMPC, routerMPCPubkey)

	log.Info(fmt.Sprintf("[%5v] init router info success", chainID),
		"routerContract", routerContract, "routerMPC", routerMPC,
		"routerWNative", routerWNative)

	return nil
}

// InitExtraCustoms init extra customs
func (b *Bridge) InitExtraCustoms() {
	clientTimeout := params.GetRPCClientTimeout(b.ChainConfig.ChainID)
	if clientTimeout != 0 {
		b.RPCClientTimeout = clientTimeout
	} else {
		timeoutStr := params.GetCustom(b.ChainConfig.ChainID, "sendtxTimeout")
		if timeoutStr != "" {
			timeout, err := common.GetUint64FromStr(timeoutStr)
			if err != nil {
				log.Fatal("get sendtxTimeout failed", "err", err)
			}
			if timeout != 0 {
				b.RPCClientTimeout = int(timeout)
			}
		}
	}
	flag := params.GetCustom(b.ChainConfig.ChainID, "dontCheckAddressMixedCase")
	b.DontCheckAddressMixedCase = strings.EqualFold(flag, "true")
}
