// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package _1inchlimitorder

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OrderMixinOrder is an auto generated low-level Go binding around an user-defined struct.
type OrderMixinOrder struct {
	Salt           *big.Int
	MakerAsset     common.Address
	TakerAsset     common.Address
	Maker          common.Address
	Receiver       common.Address
	AllowedSender  common.Address
	MakingAmount   *big.Int
	TakingAmount   *big.Int
	MakerAssetData []byte
	TakerAssetData []byte
	GetMakerAmount []byte
	GetTakerAmount []byte
	Predicate      []byte
	Permit         []byte
	Interaction    []byte
}

// OrderRFQMixinOrderRFQ is an auto generated low-level Go binding around an user-defined struct.
type OrderRFQMixinOrderRFQ struct {
	Info          *big.Int
	MakerAsset    common.Address
	TakerAsset    common.Address
	Maker         common.Address
	AllowedSender common.Address
	MakingAmount  *big.Int
	TakingAmount  *big.Int
}

// OneInchLimitOrderMetaData contains all meta data concerning the OneInchLimitOrder contract.
var OneInchLimitOrderMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNonce\",\"type\":\"uint256\"}],\"name\":\"NonceIncreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingRaw\",\"type\":\"uint256\"}],\"name\":\"OrderCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"name\":\"OrderFilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"}],\"name\":\"OrderFilledRFQ\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIMIT_ORDER_RFQ_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIMIT_ORDER_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"amount\",\"type\":\"uint8\"}],\"name\":\"advanceNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"and\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"arbitraryStaticCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderInfo\",\"type\":\"uint256\"}],\"name\":\"cancelOrderRFQ\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"checkPredicate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"oracle1\",\"type\":\"address\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"oracle2\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"spread\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"doublePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"eq\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"thresholdAmount\",\"type\":\"uint256\"}],\"name\":\"fillOrder\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"info\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOrderRFQMixin.OrderRFQ\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"}],\"name\":\"fillOrderRFQ\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"info\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOrderRFQMixin.OrderRFQ\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"fillOrderRFQTo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"info\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOrderRFQMixin.OrderRFQ\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"}],\"name\":\"fillOrderRFQToWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"thresholdAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"fillOrderTo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"thresholdAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"}],\"name\":\"fillOrderToWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderMakerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"swapTakerAmount\",\"type\":\"uint256\"}],\"name\":\"getMakerAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderMakerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"swapMakerAmount\",\"type\":\"uint256\"}],\"name\":\"getTakerAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"gt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"makerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takingAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getMakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"getTakerAmount\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"predicate\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"permit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"interaction\",\"type\":\"bytes\"}],\"internalType\":\"structOrderMixin.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"hashOrder\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increaseNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"slot\",\"type\":\"uint256\"}],\"name\":\"invalidatorForOrderRFQ\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"lt\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerNonce\",\"type\":\"uint256\"}],\"name\":\"nonceEquals\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"or\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"remaining\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"remainingRaw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"orderHashes\",\"type\":\"bytes32[]\"}],\"name\":\"remainingsRaw\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"simulateCalls\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"inverseAndSpread\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"singlePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"timestampBelow\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OneInchLimitOrderABI is the input ABI used to generate the binding from.
// Deprecated: Use OneInchLimitOrderMetaData.ABI instead.
var OneInchLimitOrderABI = OneInchLimitOrderMetaData.ABI

// OneInchLimitOrder is an auto generated Go binding around an Ethereum contract.
type OneInchLimitOrder struct {
	OneInchLimitOrderCaller     // Read-only binding to the contract
	OneInchLimitOrderTransactor // Write-only binding to the contract
	OneInchLimitOrderFilterer   // Log filterer for contract events
}

// OneInchLimitOrderCaller is an auto generated read-only Go binding around an Ethereum contract.
type OneInchLimitOrderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OneInchLimitOrderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OneInchLimitOrderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OneInchLimitOrderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OneInchLimitOrderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OneInchLimitOrderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OneInchLimitOrderSession struct {
	Contract     *OneInchLimitOrder // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OneInchLimitOrderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OneInchLimitOrderCallerSession struct {
	Contract *OneInchLimitOrderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// OneInchLimitOrderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OneInchLimitOrderTransactorSession struct {
	Contract     *OneInchLimitOrderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// OneInchLimitOrderRaw is an auto generated low-level Go binding around an Ethereum contract.
type OneInchLimitOrderRaw struct {
	Contract *OneInchLimitOrder // Generic contract binding to access the raw methods on
}

// OneInchLimitOrderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OneInchLimitOrderCallerRaw struct {
	Contract *OneInchLimitOrderCaller // Generic read-only contract binding to access the raw methods on
}

// OneInchLimitOrderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OneInchLimitOrderTransactorRaw struct {
	Contract *OneInchLimitOrderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOneInchLimitOrder creates a new instance of OneInchLimitOrder, bound to a specific deployed contract.
func NewOneInchLimitOrder(address common.Address, backend bind.ContractBackend) (*OneInchLimitOrder, error) {
	contract, err := bindOneInchLimitOrder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrder{OneInchLimitOrderCaller: OneInchLimitOrderCaller{contract: contract}, OneInchLimitOrderTransactor: OneInchLimitOrderTransactor{contract: contract}, OneInchLimitOrderFilterer: OneInchLimitOrderFilterer{contract: contract}}, nil
}

// NewOneInchLimitOrderCaller creates a new read-only instance of OneInchLimitOrder, bound to a specific deployed contract.
func NewOneInchLimitOrderCaller(address common.Address, caller bind.ContractCaller) (*OneInchLimitOrderCaller, error) {
	contract, err := bindOneInchLimitOrder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderCaller{contract: contract}, nil
}

// NewOneInchLimitOrderTransactor creates a new write-only instance of OneInchLimitOrder, bound to a specific deployed contract.
func NewOneInchLimitOrderTransactor(address common.Address, transactor bind.ContractTransactor) (*OneInchLimitOrderTransactor, error) {
	contract, err := bindOneInchLimitOrder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderTransactor{contract: contract}, nil
}

// NewOneInchLimitOrderFilterer creates a new log filterer instance of OneInchLimitOrder, bound to a specific deployed contract.
func NewOneInchLimitOrderFilterer(address common.Address, filterer bind.ContractFilterer) (*OneInchLimitOrderFilterer, error) {
	contract, err := bindOneInchLimitOrder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderFilterer{contract: contract}, nil
}

// bindOneInchLimitOrder binds a generic wrapper to an already deployed contract.
func bindOneInchLimitOrder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OneInchLimitOrderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OneInchLimitOrder *OneInchLimitOrderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OneInchLimitOrder.Contract.OneInchLimitOrderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OneInchLimitOrder *OneInchLimitOrderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.OneInchLimitOrderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OneInchLimitOrder *OneInchLimitOrderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.OneInchLimitOrderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OneInchLimitOrder *OneInchLimitOrderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OneInchLimitOrder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OneInchLimitOrder *OneInchLimitOrderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OneInchLimitOrder *OneInchLimitOrderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "DOMAIN_SEPARATOR")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.DOMAINSEPARATOR(&_OneInchLimitOrder.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.DOMAINSEPARATOR(&_OneInchLimitOrder.CallOpts)
}

// LIMITORDERRFQTYPEHASH is a free data retrieval call binding the contract method 0x06bf53d0.
//
// Solidity: function LIMIT_ORDER_RFQ_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) LIMITORDERRFQTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "LIMIT_ORDER_RFQ_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// LIMITORDERRFQTYPEHASH is a free data retrieval call binding the contract method 0x06bf53d0.
//
// Solidity: function LIMIT_ORDER_RFQ_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderSession) LIMITORDERRFQTYPEHASH() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.LIMITORDERRFQTYPEHASH(&_OneInchLimitOrder.CallOpts)
}

// LIMITORDERRFQTYPEHASH is a free data retrieval call binding the contract method 0x06bf53d0.
//
// Solidity: function LIMIT_ORDER_RFQ_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) LIMITORDERRFQTYPEHASH() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.LIMITORDERRFQTYPEHASH(&_OneInchLimitOrder.CallOpts)
}

// LIMITORDERTYPEHASH is a free data retrieval call binding the contract method 0x54dd5f74.
//
// Solidity: function LIMIT_ORDER_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) LIMITORDERTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "LIMIT_ORDER_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// LIMITORDERTYPEHASH is a free data retrieval call binding the contract method 0x54dd5f74.
//
// Solidity: function LIMIT_ORDER_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderSession) LIMITORDERTYPEHASH() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.LIMITORDERTYPEHASH(&_OneInchLimitOrder.CallOpts)
}

// LIMITORDERTYPEHASH is a free data retrieval call binding the contract method 0x54dd5f74.
//
// Solidity: function LIMIT_ORDER_TYPEHASH() view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) LIMITORDERTYPEHASH() ([32]byte, error) {
	return _OneInchLimitOrder.Contract.LIMITORDERTYPEHASH(&_OneInchLimitOrder.CallOpts)
}

// And is a free data retrieval call binding the contract method 0x961d5b1e.
//
// Solidity: function and(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) And(opts *bind.CallOpts, targets []common.Address, data [][]byte) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "and", targets, data)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// And is a free data retrieval call binding the contract method 0x961d5b1e.
//
// Solidity: function and(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) And(targets []common.Address, data [][]byte) (bool, error) {
	return _OneInchLimitOrder.Contract.And(&_OneInchLimitOrder.CallOpts, targets, data)
}

// And is a free data retrieval call binding the contract method 0x961d5b1e.
//
// Solidity: function and(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) And(targets []common.Address, data [][]byte) (bool, error) {
	return _OneInchLimitOrder.Contract.And(&_OneInchLimitOrder.CallOpts, targets, data)
}

// ArbitraryStaticCall is a free data retrieval call binding the contract method 0xbf15fcd8.
//
// Solidity: function arbitraryStaticCall(address target, bytes data) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) ArbitraryStaticCall(opts *bind.CallOpts, target common.Address, data []byte) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "arbitraryStaticCall", target, data)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// ArbitraryStaticCall is a free data retrieval call binding the contract method 0xbf15fcd8.
//
// Solidity: function arbitraryStaticCall(address target, bytes data) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) ArbitraryStaticCall(target common.Address, data []byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.ArbitraryStaticCall(&_OneInchLimitOrder.CallOpts, target, data)
}

// ArbitraryStaticCall is a free data retrieval call binding the contract method 0xbf15fcd8.
//
// Solidity: function arbitraryStaticCall(address target, bytes data) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) ArbitraryStaticCall(target common.Address, data []byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.ArbitraryStaticCall(&_OneInchLimitOrder.CallOpts, target, data)
}

// CheckPredicate is a free data retrieval call binding the contract method 0xa65a0e71.
//
// Solidity: function checkPredicate((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) CheckPredicate(opts *bind.CallOpts, order OrderMixinOrder) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "checkPredicate", order)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// CheckPredicate is a free data retrieval call binding the contract method 0xa65a0e71.
//
// Solidity: function checkPredicate((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) CheckPredicate(order OrderMixinOrder) (bool, error) {
	return _OneInchLimitOrder.Contract.CheckPredicate(&_OneInchLimitOrder.CallOpts, order)
}

// CheckPredicate is a free data retrieval call binding the contract method 0xa65a0e71.
//
// Solidity: function checkPredicate((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) CheckPredicate(order OrderMixinOrder) (bool, error) {
	return _OneInchLimitOrder.Contract.CheckPredicate(&_OneInchLimitOrder.CallOpts, order)
}

// DoublePrice is a free data retrieval call binding the contract method 0x36006bf3.
//
// Solidity: function doublePrice(address oracle1, address oracle2, uint256 spread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) DoublePrice(opts *bind.CallOpts, oracle1 common.Address, oracle2 common.Address, spread *big.Int, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "doublePrice", oracle1, oracle2, spread, amount)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// DoublePrice is a free data retrieval call binding the contract method 0x36006bf3.
//
// Solidity: function doublePrice(address oracle1, address oracle2, uint256 spread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) DoublePrice(oracle1 common.Address, oracle2 common.Address, spread *big.Int, amount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.DoublePrice(&_OneInchLimitOrder.CallOpts, oracle1, oracle2, spread, amount)
}

// DoublePrice is a free data retrieval call binding the contract method 0x36006bf3.
//
// Solidity: function doublePrice(address oracle1, address oracle2, uint256 spread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) DoublePrice(oracle1 common.Address, oracle2 common.Address, spread *big.Int, amount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.DoublePrice(&_OneInchLimitOrder.CallOpts, oracle1, oracle2, spread, amount)
}

// Eq is a free data retrieval call binding the contract method 0x32565d61.
//
// Solidity: function eq(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Eq(opts *bind.CallOpts, value *big.Int, target common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "eq", value, target, data)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// Eq is a free data retrieval call binding the contract method 0x32565d61.
//
// Solidity: function eq(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Eq(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Eq(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// Eq is a free data retrieval call binding the contract method 0x32565d61.
//
// Solidity: function eq(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Eq(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Eq(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// GetMakerAmount is a free data retrieval call binding the contract method 0xf4a215c3.
//
// Solidity: function getMakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapTakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) GetMakerAmount(opts *bind.CallOpts, orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapTakerAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "getMakerAmount", orderMakerAmount, orderTakerAmount, swapTakerAmount)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetMakerAmount is a free data retrieval call binding the contract method 0xf4a215c3.
//
// Solidity: function getMakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapTakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) GetMakerAmount(orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapTakerAmount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.GetMakerAmount(&_OneInchLimitOrder.CallOpts, orderMakerAmount, orderTakerAmount, swapTakerAmount)
}

// GetMakerAmount is a free data retrieval call binding the contract method 0xf4a215c3.
//
// Solidity: function getMakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapTakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) GetMakerAmount(orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapTakerAmount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.GetMakerAmount(&_OneInchLimitOrder.CallOpts, orderMakerAmount, orderTakerAmount, swapTakerAmount)
}

// GetTakerAmount is a free data retrieval call binding the contract method 0x296637bf.
//
// Solidity: function getTakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapMakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) GetTakerAmount(opts *bind.CallOpts, orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapMakerAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "getTakerAmount", orderMakerAmount, orderTakerAmount, swapMakerAmount)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetTakerAmount is a free data retrieval call binding the contract method 0x296637bf.
//
// Solidity: function getTakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapMakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) GetTakerAmount(orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapMakerAmount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.GetTakerAmount(&_OneInchLimitOrder.CallOpts, orderMakerAmount, orderTakerAmount, swapMakerAmount)
}

// GetTakerAmount is a free data retrieval call binding the contract method 0x296637bf.
//
// Solidity: function getTakerAmount(uint256 orderMakerAmount, uint256 orderTakerAmount, uint256 swapMakerAmount) pure returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) GetTakerAmount(orderMakerAmount *big.Int, orderTakerAmount *big.Int, swapMakerAmount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.GetTakerAmount(&_OneInchLimitOrder.CallOpts, orderMakerAmount, orderTakerAmount, swapMakerAmount)
}

// Gt is a free data retrieval call binding the contract method 0x057702e9.
//
// Solidity: function gt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Gt(opts *bind.CallOpts, value *big.Int, target common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "gt", value, target, data)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// Gt is a free data retrieval call binding the contract method 0x057702e9.
//
// Solidity: function gt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Gt(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Gt(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// Gt is a free data retrieval call binding the contract method 0x057702e9.
//
// Solidity: function gt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Gt(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Gt(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// HashOrder is a free data retrieval call binding the contract method 0xfa1cb9f2.
//
// Solidity: function hashOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) HashOrder(opts *bind.CallOpts, order OrderMixinOrder) ([32]byte, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "hashOrder", order)
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// HashOrder is a free data retrieval call binding the contract method 0xfa1cb9f2.
//
// Solidity: function hashOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderSession) HashOrder(order OrderMixinOrder) ([32]byte, error) {
	return _OneInchLimitOrder.Contract.HashOrder(&_OneInchLimitOrder.CallOpts, order)
}

// HashOrder is a free data retrieval call binding the contract method 0xfa1cb9f2.
//
// Solidity: function hashOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) view returns(bytes32)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) HashOrder(order OrderMixinOrder) ([32]byte, error) {
	return _OneInchLimitOrder.Contract.HashOrder(&_OneInchLimitOrder.CallOpts, order)
}

// InvalidatorForOrderRFQ is a free data retrieval call binding the contract method 0x56f16124.
//
// Solidity: function invalidatorForOrderRFQ(address maker, uint256 slot) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) InvalidatorForOrderRFQ(opts *bind.CallOpts, maker common.Address, slot *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "invalidatorForOrderRFQ", maker, slot)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// InvalidatorForOrderRFQ is a free data retrieval call binding the contract method 0x56f16124.
//
// Solidity: function invalidatorForOrderRFQ(address maker, uint256 slot) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) InvalidatorForOrderRFQ(maker common.Address, slot *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.InvalidatorForOrderRFQ(&_OneInchLimitOrder.CallOpts, maker, slot)
}

// InvalidatorForOrderRFQ is a free data retrieval call binding the contract method 0x56f16124.
//
// Solidity: function invalidatorForOrderRFQ(address maker, uint256 slot) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) InvalidatorForOrderRFQ(maker common.Address, slot *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.InvalidatorForOrderRFQ(&_OneInchLimitOrder.CallOpts, maker, slot)
}

// Lt is a free data retrieval call binding the contract method 0x871919d5.
//
// Solidity: function lt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Lt(opts *bind.CallOpts, value *big.Int, target common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "lt", value, target, data)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// Lt is a free data retrieval call binding the contract method 0x871919d5.
//
// Solidity: function lt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Lt(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Lt(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// Lt is a free data retrieval call binding the contract method 0x871919d5.
//
// Solidity: function lt(uint256 value, address target, bytes data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Lt(value *big.Int, target common.Address, data []byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Lt(&_OneInchLimitOrder.CallOpts, value, target, data)
}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Nonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "nonce", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Nonce(arg0 common.Address) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.Nonce(&_OneInchLimitOrder.CallOpts, arg0)
}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Nonce(arg0 common.Address) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.Nonce(&_OneInchLimitOrder.CallOpts, arg0)
}

// NonceEquals is a free data retrieval call binding the contract method 0xcf6fc6e3.
//
// Solidity: function nonceEquals(address makerAddress, uint256 makerNonce) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) NonceEquals(opts *bind.CallOpts, makerAddress common.Address, makerNonce *big.Int) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "nonceEquals", makerAddress, makerNonce)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// NonceEquals is a free data retrieval call binding the contract method 0xcf6fc6e3.
//
// Solidity: function nonceEquals(address makerAddress, uint256 makerNonce) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) NonceEquals(makerAddress common.Address, makerNonce *big.Int) (bool, error) {
	return _OneInchLimitOrder.Contract.NonceEquals(&_OneInchLimitOrder.CallOpts, makerAddress, makerNonce)
}

// NonceEquals is a free data retrieval call binding the contract method 0xcf6fc6e3.
//
// Solidity: function nonceEquals(address makerAddress, uint256 makerNonce) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) NonceEquals(makerAddress common.Address, makerNonce *big.Int) (bool, error) {
	return _OneInchLimitOrder.Contract.NonceEquals(&_OneInchLimitOrder.CallOpts, makerAddress, makerNonce)
}

// Or is a free data retrieval call binding the contract method 0xe6133301.
//
// Solidity: function or(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Or(opts *bind.CallOpts, targets []common.Address, data [][]byte) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "or", targets, data)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// Or is a free data retrieval call binding the contract method 0xe6133301.
//
// Solidity: function or(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Or(targets []common.Address, data [][]byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Or(&_OneInchLimitOrder.CallOpts, targets, data)
}

// Or is a free data retrieval call binding the contract method 0xe6133301.
//
// Solidity: function or(address[] targets, bytes[] data) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Or(targets []common.Address, data [][]byte) (bool, error) {
	return _OneInchLimitOrder.Contract.Or(&_OneInchLimitOrder.CallOpts, targets, data)
}

// Remaining is a free data retrieval call binding the contract method 0xbc1ed74c.
//
// Solidity: function remaining(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) Remaining(opts *bind.CallOpts, orderHash [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "remaining", orderHash)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Remaining is a free data retrieval call binding the contract method 0xbc1ed74c.
//
// Solidity: function remaining(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) Remaining(orderHash [32]byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.Remaining(&_OneInchLimitOrder.CallOpts, orderHash)
}

// Remaining is a free data retrieval call binding the contract method 0xbc1ed74c.
//
// Solidity: function remaining(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) Remaining(orderHash [32]byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.Remaining(&_OneInchLimitOrder.CallOpts, orderHash)
}

// RemainingRaw is a free data retrieval call binding the contract method 0x7e54f092.
//
// Solidity: function remainingRaw(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) RemainingRaw(opts *bind.CallOpts, orderHash [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "remainingRaw", orderHash)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RemainingRaw is a free data retrieval call binding the contract method 0x7e54f092.
//
// Solidity: function remainingRaw(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) RemainingRaw(orderHash [32]byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.RemainingRaw(&_OneInchLimitOrder.CallOpts, orderHash)
}

// RemainingRaw is a free data retrieval call binding the contract method 0x7e54f092.
//
// Solidity: function remainingRaw(bytes32 orderHash) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) RemainingRaw(orderHash [32]byte) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.RemainingRaw(&_OneInchLimitOrder.CallOpts, orderHash)
}

// RemainingsRaw is a free data retrieval call binding the contract method 0x942461bb.
//
// Solidity: function remainingsRaw(bytes32[] orderHashes) view returns(uint256[])
func (_OneInchLimitOrder *OneInchLimitOrderCaller) RemainingsRaw(opts *bind.CallOpts, orderHashes [][32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "remainingsRaw", orderHashes)
	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err
}

// RemainingsRaw is a free data retrieval call binding the contract method 0x942461bb.
//
// Solidity: function remainingsRaw(bytes32[] orderHashes) view returns(uint256[])
func (_OneInchLimitOrder *OneInchLimitOrderSession) RemainingsRaw(orderHashes [][32]byte) ([]*big.Int, error) {
	return _OneInchLimitOrder.Contract.RemainingsRaw(&_OneInchLimitOrder.CallOpts, orderHashes)
}

// RemainingsRaw is a free data retrieval call binding the contract method 0x942461bb.
//
// Solidity: function remainingsRaw(bytes32[] orderHashes) view returns(uint256[])
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) RemainingsRaw(orderHashes [][32]byte) ([]*big.Int, error) {
	return _OneInchLimitOrder.Contract.RemainingsRaw(&_OneInchLimitOrder.CallOpts, orderHashes)
}

// SinglePrice is a free data retrieval call binding the contract method 0xc05435f1.
//
// Solidity: function singlePrice(address oracle, uint256 inverseAndSpread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) SinglePrice(opts *bind.CallOpts, oracle common.Address, inverseAndSpread *big.Int, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "singlePrice", oracle, inverseAndSpread, amount)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// SinglePrice is a free data retrieval call binding the contract method 0xc05435f1.
//
// Solidity: function singlePrice(address oracle, uint256 inverseAndSpread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) SinglePrice(oracle common.Address, inverseAndSpread *big.Int, amount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.SinglePrice(&_OneInchLimitOrder.CallOpts, oracle, inverseAndSpread, amount)
}

// SinglePrice is a free data retrieval call binding the contract method 0xc05435f1.
//
// Solidity: function singlePrice(address oracle, uint256 inverseAndSpread, uint256 amount) view returns(uint256)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) SinglePrice(oracle common.Address, inverseAndSpread *big.Int, amount *big.Int) (*big.Int, error) {
	return _OneInchLimitOrder.Contract.SinglePrice(&_OneInchLimitOrder.CallOpts, oracle, inverseAndSpread, amount)
}

// TimestampBelow is a free data retrieval call binding the contract method 0x63592c2b.
//
// Solidity: function timestampBelow(uint256 time) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCaller) TimestampBelow(opts *bind.CallOpts, time *big.Int) (bool, error) {
	var out []interface{}
	err := _OneInchLimitOrder.contract.Call(opts, &out, "timestampBelow", time)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// TimestampBelow is a free data retrieval call binding the contract method 0x63592c2b.
//
// Solidity: function timestampBelow(uint256 time) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderSession) TimestampBelow(time *big.Int) (bool, error) {
	return _OneInchLimitOrder.Contract.TimestampBelow(&_OneInchLimitOrder.CallOpts, time)
}

// TimestampBelow is a free data retrieval call binding the contract method 0x63592c2b.
//
// Solidity: function timestampBelow(uint256 time) view returns(bool)
func (_OneInchLimitOrder *OneInchLimitOrderCallerSession) TimestampBelow(time *big.Int) (bool, error) {
	return _OneInchLimitOrder.Contract.TimestampBelow(&_OneInchLimitOrder.CallOpts, time)
}

// AdvanceNonce is a paid mutator transaction binding the contract method 0x72c244a8.
//
// Solidity: function advanceNonce(uint8 amount) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) AdvanceNonce(opts *bind.TransactOpts, amount uint8) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "advanceNonce", amount)
}

// AdvanceNonce is a paid mutator transaction binding the contract method 0x72c244a8.
//
// Solidity: function advanceNonce(uint8 amount) returns()
func (_OneInchLimitOrder *OneInchLimitOrderSession) AdvanceNonce(amount uint8) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.AdvanceNonce(&_OneInchLimitOrder.TransactOpts, amount)
}

// AdvanceNonce is a paid mutator transaction binding the contract method 0x72c244a8.
//
// Solidity: function advanceNonce(uint8 amount) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) AdvanceNonce(amount uint8) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.AdvanceNonce(&_OneInchLimitOrder.TransactOpts, amount)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xb244b450.
//
// Solidity: function cancelOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) CancelOrder(opts *bind.TransactOpts, order OrderMixinOrder) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xb244b450.
//
// Solidity: function cancelOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) returns()
func (_OneInchLimitOrder *OneInchLimitOrderSession) CancelOrder(order OrderMixinOrder) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.CancelOrder(&_OneInchLimitOrder.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xb244b450.
//
// Solidity: function cancelOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) CancelOrder(order OrderMixinOrder) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.CancelOrder(&_OneInchLimitOrder.TransactOpts, order)
}

// CancelOrderRFQ is a paid mutator transaction binding the contract method 0x825caba1.
//
// Solidity: function cancelOrderRFQ(uint256 orderInfo) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) CancelOrderRFQ(opts *bind.TransactOpts, orderInfo *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "cancelOrderRFQ", orderInfo)
}

// CancelOrderRFQ is a paid mutator transaction binding the contract method 0x825caba1.
//
// Solidity: function cancelOrderRFQ(uint256 orderInfo) returns()
func (_OneInchLimitOrder *OneInchLimitOrderSession) CancelOrderRFQ(orderInfo *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.CancelOrderRFQ(&_OneInchLimitOrder.TransactOpts, orderInfo)
}

// CancelOrderRFQ is a paid mutator transaction binding the contract method 0x825caba1.
//
// Solidity: function cancelOrderRFQ(uint256 orderInfo) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) CancelOrderRFQ(orderInfo *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.CancelOrderRFQ(&_OneInchLimitOrder.TransactOpts, orderInfo)
}

// FillOrder is a paid mutator transaction binding the contract method 0x655d13cd.
//
// Solidity: function fillOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrder(opts *bind.TransactOpts, order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrder", order, signature, makingAmount, takingAmount, thresholdAmount)
}

// FillOrder is a paid mutator transaction binding the contract method 0x655d13cd.
//
// Solidity: function fillOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrder(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrder(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount)
}

// FillOrder is a paid mutator transaction binding the contract method 0x655d13cd.
//
// Solidity: function fillOrder((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrder(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrder(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount)
}

// FillOrderRFQ is a paid mutator transaction binding the contract method 0xd0a3b665.
//
// Solidity: function fillOrderRFQ((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrderRFQ(opts *bind.TransactOpts, order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrderRFQ", order, signature, makingAmount, takingAmount)
}

// FillOrderRFQ is a paid mutator transaction binding the contract method 0xd0a3b665.
//
// Solidity: function fillOrderRFQ((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrderRFQ(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQ(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount)
}

// FillOrderRFQ is a paid mutator transaction binding the contract method 0xd0a3b665.
//
// Solidity: function fillOrderRFQ((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrderRFQ(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQ(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount)
}

// FillOrderRFQTo is a paid mutator transaction binding the contract method 0xbaba5855.
//
// Solidity: function fillOrderRFQTo((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrderRFQTo(opts *bind.TransactOpts, order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrderRFQTo", order, signature, makingAmount, takingAmount, target)
}

// FillOrderRFQTo is a paid mutator transaction binding the contract method 0xbaba5855.
//
// Solidity: function fillOrderRFQTo((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrderRFQTo(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQTo(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, target)
}

// FillOrderRFQTo is a paid mutator transaction binding the contract method 0xbaba5855.
//
// Solidity: function fillOrderRFQTo((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrderRFQTo(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQTo(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, target)
}

// FillOrderRFQToWithPermit is a paid mutator transaction binding the contract method 0x4cc4a27b.
//
// Solidity: function fillOrderRFQToWithPermit((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrderRFQToWithPermit(opts *bind.TransactOpts, order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrderRFQToWithPermit", order, signature, makingAmount, takingAmount, target, permit)
}

// FillOrderRFQToWithPermit is a paid mutator transaction binding the contract method 0x4cc4a27b.
//
// Solidity: function fillOrderRFQToWithPermit((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrderRFQToWithPermit(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQToWithPermit(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, target, permit)
}

// FillOrderRFQToWithPermit is a paid mutator transaction binding the contract method 0x4cc4a27b.
//
// Solidity: function fillOrderRFQToWithPermit((uint256,address,address,address,address,uint256,uint256) order, bytes signature, uint256 makingAmount, uint256 takingAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrderRFQToWithPermit(order OrderRFQMixinOrderRFQ, signature []byte, makingAmount *big.Int, takingAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderRFQToWithPermit(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, target, permit)
}

// FillOrderTo is a paid mutator transaction binding the contract method 0xb2610fe3.
//
// Solidity: function fillOrderTo((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrderTo(opts *bind.TransactOpts, order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrderTo", order, signature, makingAmount, takingAmount, thresholdAmount, target)
}

// FillOrderTo is a paid mutator transaction binding the contract method 0xb2610fe3.
//
// Solidity: function fillOrderTo((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrderTo(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderTo(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount, target)
}

// FillOrderTo is a paid mutator transaction binding the contract method 0xb2610fe3.
//
// Solidity: function fillOrderTo((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrderTo(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderTo(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount, target)
}

// FillOrderToWithPermit is a paid mutator transaction binding the contract method 0x6073cc20.
//
// Solidity: function fillOrderToWithPermit((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) FillOrderToWithPermit(opts *bind.TransactOpts, order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "fillOrderToWithPermit", order, signature, makingAmount, takingAmount, thresholdAmount, target, permit)
}

// FillOrderToWithPermit is a paid mutator transaction binding the contract method 0x6073cc20.
//
// Solidity: function fillOrderToWithPermit((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderSession) FillOrderToWithPermit(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderToWithPermit(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount, target, permit)
}

// FillOrderToWithPermit is a paid mutator transaction binding the contract method 0x6073cc20.
//
// Solidity: function fillOrderToWithPermit((uint256,address,address,address,address,address,uint256,uint256,bytes,bytes,bytes,bytes,bytes,bytes,bytes) order, bytes signature, uint256 makingAmount, uint256 takingAmount, uint256 thresholdAmount, address target, bytes permit) returns(uint256, uint256)
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) FillOrderToWithPermit(order OrderMixinOrder, signature []byte, makingAmount *big.Int, takingAmount *big.Int, thresholdAmount *big.Int, target common.Address, permit []byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.FillOrderToWithPermit(&_OneInchLimitOrder.TransactOpts, order, signature, makingAmount, takingAmount, thresholdAmount, target, permit)
}

// IncreaseNonce is a paid mutator transaction binding the contract method 0xc53a0292.
//
// Solidity: function increaseNonce() returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) IncreaseNonce(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "increaseNonce")
}

// IncreaseNonce is a paid mutator transaction binding the contract method 0xc53a0292.
//
// Solidity: function increaseNonce() returns()
func (_OneInchLimitOrder *OneInchLimitOrderSession) IncreaseNonce() (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.IncreaseNonce(&_OneInchLimitOrder.TransactOpts)
}

// IncreaseNonce is a paid mutator transaction binding the contract method 0xc53a0292.
//
// Solidity: function increaseNonce() returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) IncreaseNonce() (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.IncreaseNonce(&_OneInchLimitOrder.TransactOpts)
}

// SimulateCalls is a paid mutator transaction binding the contract method 0x7f29a59d.
//
// Solidity: function simulateCalls(address[] targets, bytes[] data) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactor) SimulateCalls(opts *bind.TransactOpts, targets []common.Address, data [][]byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.contract.Transact(opts, "simulateCalls", targets, data)
}

// SimulateCalls is a paid mutator transaction binding the contract method 0x7f29a59d.
//
// Solidity: function simulateCalls(address[] targets, bytes[] data) returns()
func (_OneInchLimitOrder *OneInchLimitOrderSession) SimulateCalls(targets []common.Address, data [][]byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.SimulateCalls(&_OneInchLimitOrder.TransactOpts, targets, data)
}

// SimulateCalls is a paid mutator transaction binding the contract method 0x7f29a59d.
//
// Solidity: function simulateCalls(address[] targets, bytes[] data) returns()
func (_OneInchLimitOrder *OneInchLimitOrderTransactorSession) SimulateCalls(targets []common.Address, data [][]byte) (*types.Transaction, error) {
	return _OneInchLimitOrder.Contract.SimulateCalls(&_OneInchLimitOrder.TransactOpts, targets, data)
}

// OneInchLimitOrderNonceIncreasedIterator is returned from FilterNonceIncreased and is used to iterate over the raw logs and unpacked data for NonceIncreased events raised by the OneInchLimitOrder contract.
type OneInchLimitOrderNonceIncreasedIterator struct {
	Event *OneInchLimitOrderNonceIncreased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OneInchLimitOrderNonceIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OneInchLimitOrderNonceIncreased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OneInchLimitOrderNonceIncreased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OneInchLimitOrderNonceIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OneInchLimitOrderNonceIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OneInchLimitOrderNonceIncreased represents a NonceIncreased event raised by the OneInchLimitOrder contract.
type OneInchLimitOrderNonceIncreased struct {
	Maker    common.Address
	NewNonce *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNonceIncreased is a free log retrieval operation binding the contract event 0xfc69110dd11eb791755e4abd6b7d281bae236de95736d38a23782814be5e10db.
//
// Solidity: event NonceIncreased(address indexed maker, uint256 newNonce)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) FilterNonceIncreased(opts *bind.FilterOpts, maker []common.Address) (*OneInchLimitOrderNonceIncreasedIterator, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.FilterLogs(opts, "NonceIncreased", makerRule)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderNonceIncreasedIterator{contract: _OneInchLimitOrder.contract, event: "NonceIncreased", logs: logs, sub: sub}, nil
}

// WatchNonceIncreased is a free log subscription operation binding the contract event 0xfc69110dd11eb791755e4abd6b7d281bae236de95736d38a23782814be5e10db.
//
// Solidity: event NonceIncreased(address indexed maker, uint256 newNonce)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) WatchNonceIncreased(opts *bind.WatchOpts, sink chan<- *OneInchLimitOrderNonceIncreased, maker []common.Address) (event.Subscription, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.WatchLogs(opts, "NonceIncreased", makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OneInchLimitOrderNonceIncreased)
				if err := _OneInchLimitOrder.contract.UnpackLog(event, "NonceIncreased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNonceIncreased is a log parse operation binding the contract event 0xfc69110dd11eb791755e4abd6b7d281bae236de95736d38a23782814be5e10db.
//
// Solidity: event NonceIncreased(address indexed maker, uint256 newNonce)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) ParseNonceIncreased(log types.Log) (*OneInchLimitOrderNonceIncreased, error) {
	event := new(OneInchLimitOrderNonceIncreased)
	if err := _OneInchLimitOrder.contract.UnpackLog(event, "NonceIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OneInchLimitOrderOrderCanceledIterator is returned from FilterOrderCanceled and is used to iterate over the raw logs and unpacked data for OrderCanceled events raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderCanceledIterator struct {
	Event *OneInchLimitOrderOrderCanceled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OneInchLimitOrderOrderCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OneInchLimitOrderOrderCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OneInchLimitOrderOrderCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OneInchLimitOrderOrderCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OneInchLimitOrderOrderCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OneInchLimitOrderOrderCanceled represents a OrderCanceled event raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderCanceled struct {
	Maker        common.Address
	OrderHash    [32]byte
	RemainingRaw *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOrderCanceled is a free log retrieval operation binding the contract event 0xcbfa7d191838ece7ba4783ca3a30afd316619b7f368094b57ee7ffde9a923db1.
//
// Solidity: event OrderCanceled(address indexed maker, bytes32 orderHash, uint256 remainingRaw)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) FilterOrderCanceled(opts *bind.FilterOpts, maker []common.Address) (*OneInchLimitOrderOrderCanceledIterator, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.FilterLogs(opts, "OrderCanceled", makerRule)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderOrderCanceledIterator{contract: _OneInchLimitOrder.contract, event: "OrderCanceled", logs: logs, sub: sub}, nil
}

// WatchOrderCanceled is a free log subscription operation binding the contract event 0xcbfa7d191838ece7ba4783ca3a30afd316619b7f368094b57ee7ffde9a923db1.
//
// Solidity: event OrderCanceled(address indexed maker, bytes32 orderHash, uint256 remainingRaw)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) WatchOrderCanceled(opts *bind.WatchOpts, sink chan<- *OneInchLimitOrderOrderCanceled, maker []common.Address) (event.Subscription, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.WatchLogs(opts, "OrderCanceled", makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OneInchLimitOrderOrderCanceled)
				if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrderCanceled is a log parse operation binding the contract event 0xcbfa7d191838ece7ba4783ca3a30afd316619b7f368094b57ee7ffde9a923db1.
//
// Solidity: event OrderCanceled(address indexed maker, bytes32 orderHash, uint256 remainingRaw)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) ParseOrderCanceled(log types.Log) (*OneInchLimitOrderOrderCanceled, error) {
	event := new(OneInchLimitOrderOrderCanceled)
	if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OneInchLimitOrderOrderFilledIterator is returned from FilterOrderFilled and is used to iterate over the raw logs and unpacked data for OrderFilled events raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderFilledIterator struct {
	Event *OneInchLimitOrderOrderFilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OneInchLimitOrderOrderFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OneInchLimitOrderOrderFilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OneInchLimitOrderOrderFilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OneInchLimitOrderOrderFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OneInchLimitOrderOrderFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OneInchLimitOrderOrderFilled represents a OrderFilled event raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderFilled struct {
	Maker     common.Address
	OrderHash [32]byte
	Remaining *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOrderFilled is a free log retrieval operation binding the contract event 0xb9ed0243fdf00f0545c63a0af8850c090d86bb46682baec4bf3c496814fe4f02.
//
// Solidity: event OrderFilled(address indexed maker, bytes32 orderHash, uint256 remaining)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) FilterOrderFilled(opts *bind.FilterOpts, maker []common.Address) (*OneInchLimitOrderOrderFilledIterator, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.FilterLogs(opts, "OrderFilled", makerRule)
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderOrderFilledIterator{contract: _OneInchLimitOrder.contract, event: "OrderFilled", logs: logs, sub: sub}, nil
}

// WatchOrderFilled is a free log subscription operation binding the contract event 0xb9ed0243fdf00f0545c63a0af8850c090d86bb46682baec4bf3c496814fe4f02.
//
// Solidity: event OrderFilled(address indexed maker, bytes32 orderHash, uint256 remaining)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) WatchOrderFilled(opts *bind.WatchOpts, sink chan<- *OneInchLimitOrderOrderFilled, maker []common.Address) (event.Subscription, error) {
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _OneInchLimitOrder.contract.WatchLogs(opts, "OrderFilled", makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OneInchLimitOrderOrderFilled)
				if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderFilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrderFilled is a log parse operation binding the contract event 0xb9ed0243fdf00f0545c63a0af8850c090d86bb46682baec4bf3c496814fe4f02.
//
// Solidity: event OrderFilled(address indexed maker, bytes32 orderHash, uint256 remaining)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) ParseOrderFilled(log types.Log) (*OneInchLimitOrderOrderFilled, error) {
	event := new(OneInchLimitOrderOrderFilled)
	if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OneInchLimitOrderOrderFilledRFQIterator is returned from FilterOrderFilledRFQ and is used to iterate over the raw logs and unpacked data for OrderFilledRFQ events raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderFilledRFQIterator struct {
	Event *OneInchLimitOrderOrderFilledRFQ // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OneInchLimitOrderOrderFilledRFQIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OneInchLimitOrderOrderFilledRFQ)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OneInchLimitOrderOrderFilledRFQ)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OneInchLimitOrderOrderFilledRFQIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OneInchLimitOrderOrderFilledRFQIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OneInchLimitOrderOrderFilledRFQ represents a OrderFilledRFQ event raised by the OneInchLimitOrder contract.
type OneInchLimitOrderOrderFilledRFQ struct {
	OrderHash    [32]byte
	MakingAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOrderFilledRFQ is a free log retrieval operation binding the contract event 0xc3b639f02b125bfa160e50739b8c44eb2d1b6908e2b6d5925c6d770f2ca78127.
//
// Solidity: event OrderFilledRFQ(bytes32 orderHash, uint256 makingAmount)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) FilterOrderFilledRFQ(opts *bind.FilterOpts) (*OneInchLimitOrderOrderFilledRFQIterator, error) {
	logs, sub, err := _OneInchLimitOrder.contract.FilterLogs(opts, "OrderFilledRFQ")
	if err != nil {
		return nil, err
	}
	return &OneInchLimitOrderOrderFilledRFQIterator{contract: _OneInchLimitOrder.contract, event: "OrderFilledRFQ", logs: logs, sub: sub}, nil
}

// WatchOrderFilledRFQ is a free log subscription operation binding the contract event 0xc3b639f02b125bfa160e50739b8c44eb2d1b6908e2b6d5925c6d770f2ca78127.
//
// Solidity: event OrderFilledRFQ(bytes32 orderHash, uint256 makingAmount)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) WatchOrderFilledRFQ(opts *bind.WatchOpts, sink chan<- *OneInchLimitOrderOrderFilledRFQ) (event.Subscription, error) {
	logs, sub, err := _OneInchLimitOrder.contract.WatchLogs(opts, "OrderFilledRFQ")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OneInchLimitOrderOrderFilledRFQ)
				if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderFilledRFQ", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrderFilledRFQ is a log parse operation binding the contract event 0xc3b639f02b125bfa160e50739b8c44eb2d1b6908e2b6d5925c6d770f2ca78127.
//
// Solidity: event OrderFilledRFQ(bytes32 orderHash, uint256 makingAmount)
func (_OneInchLimitOrder *OneInchLimitOrderFilterer) ParseOrderFilledRFQ(log types.Log) (*OneInchLimitOrderOrderFilledRFQ, error) {
	event := new(OneInchLimitOrderOrderFilledRFQ)
	if err := _OneInchLimitOrder.contract.UnpackLog(event, "OrderFilledRFQ", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
