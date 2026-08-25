// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package script

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

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
	_ = abi.ConvertType
	_ = time.Tick
	_ = context.Background
)

// MyCounterMetaData contains all meta data concerning the MyCounter contract.
var MyCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Increased\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506101aa8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c80638ada066e14610038578063e8927fbc14610056575b5f5ffd5b610040610060565b60405161004d91906100e7565b60405180910390f35b61005e610068565b005b5f5f54905090565b5f5f5f81546100769061012d565b91905081905590503373ffffffffffffffffffffffffffffffffffffffff167f071c8af8707bfeb7b8186295479bbffcaff15c8cca8e9727046e8b0215c01fcb826040516100c491906100e7565b60405180910390a250565b5f819050919050565b6100e1816100cf565b82525050565b5f6020820190506100fa5f8301846100d8565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610137826100cf565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361016957610168610100565b5b60018201905091905056fea2646970667358221220f3db6cfdd9e858b87523837e9212acbcc1759bcec0d1e92b216525b09d7fa6ed64736f6c63430008240033",
}

// MyCounterABI is the input ABI used to generate the binding from.
// Deprecated: Use MyCounterMetaData.ABI instead.
var MyCounterABI = MyCounterMetaData.ABI

// MyCounterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MyCounterMetaData.Bin instead.
var MyCounterBin = MyCounterMetaData.Bin

// DeployMyCounter deploys a new Ethereum contract, binding an instance of MyCounter to it.
func DeployMyCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MyCounter, error) {
	parsed, err := MyCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MyCounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MyCounter{MyCounterCaller: MyCounterCaller{contract: contract}, MyCounterTransactor: MyCounterTransactor{contract: contract}, MyCounterFilterer: MyCounterFilterer{contract: contract}}, nil
}

// MyCounter is an auto generated Go binding around an Ethereum contract.
type MyCounter struct {
	MyCounterCaller     // Read-only binding to the contract
	MyCounterTransactor // Write-only binding to the contract
	MyCounterFilterer   // Log filterer for contract events
}

// MyCounterCaller is an auto generated read-only Go binding around an Ethereum contract.
type MyCounterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyCounterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MyCounterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyCounterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MyCounterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyCounterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MyCounterSession struct {
	Contract     *MyCounter        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MyCounterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MyCounterCallerSession struct {
	Contract *MyCounterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MyCounterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MyCounterTransactorSession struct {
	Contract     *MyCounterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MyCounterRaw is an auto generated low-level Go binding around an Ethereum contract.
type MyCounterRaw struct {
	Contract *MyCounter // Generic contract binding to access the raw methods on
}

// MyCounterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MyCounterCallerRaw struct {
	Contract *MyCounterCaller // Generic read-only contract binding to access the raw methods on
}

// MyCounterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MyCounterTransactorRaw struct {
	Contract *MyCounterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMyCounter creates a new instance of MyCounter, bound to a specific deployed contract.
func NewMyCounter(address common.Address, backend bind.ContractBackend) (*MyCounter, error) {
	contract, err := bindMyCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MyCounter{MyCounterCaller: MyCounterCaller{contract: contract}, MyCounterTransactor: MyCounterTransactor{contract: contract}, MyCounterFilterer: MyCounterFilterer{contract: contract}}, nil
}

// NewMyCounterCaller creates a new read-only instance of MyCounter, bound to a specific deployed contract.
func NewMyCounterCaller(address common.Address, caller bind.ContractCaller) (*MyCounterCaller, error) {
	contract, err := bindMyCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MyCounterCaller{contract: contract}, nil
}

// NewMyCounterTransactor creates a new write-only instance of MyCounter, bound to a specific deployed contract.
func NewMyCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*MyCounterTransactor, error) {
	contract, err := bindMyCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MyCounterTransactor{contract: contract}, nil
}

// NewMyCounterFilterer creates a new log filterer instance of MyCounter, bound to a specific deployed contract.
func NewMyCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*MyCounterFilterer, error) {
	contract, err := bindMyCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MyCounterFilterer{contract: contract}, nil
}

// bindMyCounter binds a generic wrapper to an already deployed contract.
func bindMyCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MyCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MyCounter *MyCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MyCounter.Contract.MyCounterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MyCounter *MyCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MyCounter.Contract.MyCounterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MyCounter *MyCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MyCounter.Contract.MyCounterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MyCounter *MyCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MyCounter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MyCounter *MyCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MyCounter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MyCounter *MyCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MyCounter.Contract.contract.Transact(opts, method, params...)
}

// GetCounter is a free data retrieval call binding the contract method 0x8ada066e.
//
// Solidity: function getCounter() view returns(uint256)
func (_MyCounter *MyCounterCaller) GetCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MyCounter.contract.Call(opts, &out, "getCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCounter is a free data retrieval call binding the contract method 0x8ada066e.
//
// Solidity: function getCounter() view returns(uint256)
func (_MyCounter *MyCounterSession) GetCounter() (*big.Int, error) {
	return _MyCounter.Contract.GetCounter(&_MyCounter.CallOpts)
}

// GetCounter is a free data retrieval call binding the contract method 0x8ada066e.
//
// Solidity: function getCounter() view returns(uint256)
func (_MyCounter *MyCounterCallerSession) GetCounter() (*big.Int, error) {
	return _MyCounter.Contract.GetCounter(&_MyCounter.CallOpts)
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_MyCounter *MyCounterTransactor) Increase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MyCounter.contract.Transact(opts, "increase")
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_MyCounter *MyCounterSession) Increase() (*types.Transaction, error) {
	return _MyCounter.Contract.Increase(&_MyCounter.TransactOpts)
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_MyCounter *MyCounterTransactorSession) Increase() (*types.Transaction, error) {
	return _MyCounter.Contract.Increase(&_MyCounter.TransactOpts)
}

// MyCounterIncreasedIterator is returned from FilterIncreased and is used to iterate over the raw logs and unpacked data for Increased events raised by the MyCounter contract.
type MyCounterIncreasedIterator struct {
	Event *MyCounterIncreased // Event containing the contract specifics and raw log

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
func (it *MyCounterIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyCounterIncreased)
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
		it.Event = new(MyCounterIncreased)
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
func (it *MyCounterIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyCounterIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyCounterIncreased represents a Increased event raised by the MyCounter contract.
type MyCounterIncreased struct {
	Caller common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIncreased is a free log retrieval operation binding the contract event 0x071c8af8707bfeb7b8186295479bbffcaff15c8cca8e9727046e8b0215c01fcb.
//
// Solidity: event Increased(address indexed caller, uint256 value)
func (_MyCounter *MyCounterFilterer) FilterIncreased(opts *bind.FilterOpts, caller []common.Address) (*MyCounterIncreasedIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _MyCounter.contract.FilterLogs(opts, "Increased", callerRule)
	if err != nil {
		return nil, err
	}
	return &MyCounterIncreasedIterator{contract: _MyCounter.contract, event: "Increased", logs: logs, sub: sub}, nil
}

// WatchIncreased is a free log subscription operation binding the contract event 0x071c8af8707bfeb7b8186295479bbffcaff15c8cca8e9727046e8b0215c01fcb.
//
// Solidity: event Increased(address indexed caller, uint256 value)
func (_MyCounter *MyCounterFilterer) WatchIncreased(opts *bind.WatchOpts, sink chan<- *MyCounterIncreased, caller []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _MyCounter.contract.WatchLogs(opts, "Increased", callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyCounterIncreased)
				if err := _MyCounter.contract.UnpackLog(event, "Increased", log); err != nil {
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

// ParseIncreased is a log parse operation binding the contract event 0x071c8af8707bfeb7b8186295479bbffcaff15c8cca8e9727046e8b0215c01fcb.
//
// Solidity: event Increased(address indexed caller, uint256 value)
func (_MyCounter *MyCounterFilterer) ParseIncreased(log types.Log) (*MyCounterIncreased, error) {
	event := new(MyCounterIncreased)
	if err := _MyCounter.contract.UnpackLog(event, "Increased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
