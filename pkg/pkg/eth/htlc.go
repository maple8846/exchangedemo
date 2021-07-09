// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// HTLCABI is the input ABI used to generate the binding from.
const HTLCABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"timelock\",\"type\":\"uint256\"}],\"name\":\"LogHTLCNew\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"LogHTLCWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"LogHTLCRefund\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_hashlock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"newContract\",\"outputs\":[{\"name\":\"contractId\",\"type\":\"bytes32\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"},{\"name\":\"_preimage\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"}],\"name\":\"getContract\",\"outputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\"},{\"name\":\"refunded\",\"type\":\"bool\"},{\"name\":\"preimage\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// HTLC is an auto generated Go binding around an Ethereum contract.
type HTLC struct {
	HTLCCaller     // Read-only binding to the contract
	HTLCTransactor // Write-only binding to the contract
	HTLCFilterer   // Log filterer for contract events
}

// HTLCCaller is an auto generated read-only Go binding around an Ethereum contract.
type HTLCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HTLCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HTLCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HTLCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HTLCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HTLCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HTLCSession struct {
	Contract     *HTLC             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HTLCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HTLCCallerSession struct {
	Contract *HTLCCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HTLCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HTLCTransactorSession struct {
	Contract     *HTLCTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HTLCRaw is an auto generated low-level Go binding around an Ethereum contract.
type HTLCRaw struct {
	Contract *HTLC // Generic contract binding to access the raw methods on
}

// HTLCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HTLCCallerRaw struct {
	Contract *HTLCCaller // Generic read-only contract binding to access the raw methods on
}

// HTLCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HTLCTransactorRaw struct {
	Contract *HTLCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHTLC creates a new instance of HTLC, bound to a specific deployed contract.
func NewHTLC(address common.Address, backend bind.ContractBackend) (*HTLC, error) {
	contract, err := bindHTLC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HTLC{HTLCCaller: HTLCCaller{contract: contract}, HTLCTransactor: HTLCTransactor{contract: contract}, HTLCFilterer: HTLCFilterer{contract: contract}}, nil
}

// NewHTLCCaller creates a new read-only instance of HTLC, bound to a specific deployed contract.
func NewHTLCCaller(address common.Address, caller bind.ContractCaller) (*HTLCCaller, error) {
	contract, err := bindHTLC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HTLCCaller{contract: contract}, nil
}

// NewHTLCTransactor creates a new write-only instance of HTLC, bound to a specific deployed contract.
func NewHTLCTransactor(address common.Address, transactor bind.ContractTransactor) (*HTLCTransactor, error) {
	contract, err := bindHTLC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HTLCTransactor{contract: contract}, nil
}

// NewHTLCFilterer creates a new log filterer instance of HTLC, bound to a specific deployed contract.
func NewHTLCFilterer(address common.Address, filterer bind.ContractFilterer) (*HTLCFilterer, error) {
	contract, err := bindHTLC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HTLCFilterer{contract: contract}, nil
}

// bindHTLC binds a generic wrapper to an already deployed contract.
func bindHTLC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HTLCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HTLC *HTLCRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HTLC.Contract.HTLCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HTLC *HTLCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HTLC.Contract.HTLCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HTLC *HTLCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HTLC.Contract.HTLCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HTLC *HTLCCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HTLC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HTLC *HTLCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HTLC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HTLC *HTLCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HTLC.Contract.contract.Transact(opts, method, params...)
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(_contractId bytes32) constant returns(sender address, receiver address, amount uint256, hashlock bytes32, timelock uint256, withdrawn bool, refunded bool, preimage bytes32)
func (_HTLC *HTLCCaller) GetContract(opts *bind.CallOpts, _contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	ret := new(struct {
		Sender    common.Address
		Receiver  common.Address
		Amount    *big.Int
		Hashlock  [32]byte
		Timelock  *big.Int
		Withdrawn bool
		Refunded  bool
		Preimage  [32]byte
	})
	out := ret
	err := _HTLC.contract.Call(opts, out, "getContract", _contractId)
	return *ret, err
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(_contractId bytes32) constant returns(sender address, receiver address, amount uint256, hashlock bytes32, timelock uint256, withdrawn bool, refunded bool, preimage bytes32)
func (_HTLC *HTLCSession) GetContract(_contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	return _HTLC.Contract.GetContract(&_HTLC.CallOpts, _contractId)
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(_contractId bytes32) constant returns(sender address, receiver address, amount uint256, hashlock bytes32, timelock uint256, withdrawn bool, refunded bool, preimage bytes32)
func (_HTLC *HTLCCallerSession) GetContract(_contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	return _HTLC.Contract.GetContract(&_HTLC.CallOpts, _contractId)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(_receiver address, _hashlock bytes32, _timelock uint256) returns(contractId bytes32)
func (_HTLC *HTLCTransactor) NewContract(opts *bind.TransactOpts, _receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _HTLC.contract.Transact(opts, "newContract", _receiver, _hashlock, _timelock)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(_receiver address, _hashlock bytes32, _timelock uint256) returns(contractId bytes32)
func (_HTLC *HTLCSession) NewContract(_receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _HTLC.Contract.NewContract(&_HTLC.TransactOpts, _receiver, _hashlock, _timelock)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(_receiver address, _hashlock bytes32, _timelock uint256) returns(contractId bytes32)
func (_HTLC *HTLCTransactorSession) NewContract(_receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _HTLC.Contract.NewContract(&_HTLC.TransactOpts, _receiver, _hashlock, _timelock)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_contractId bytes32) returns(bool)
func (_HTLC *HTLCTransactor) Refund(opts *bind.TransactOpts, _contractId [32]byte) (*types.Transaction, error) {
	return _HTLC.contract.Transact(opts, "refund", _contractId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_contractId bytes32) returns(bool)
func (_HTLC *HTLCSession) Refund(_contractId [32]byte) (*types.Transaction, error) {
	return _HTLC.Contract.Refund(&_HTLC.TransactOpts, _contractId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_contractId bytes32) returns(bool)
func (_HTLC *HTLCTransactorSession) Refund(_contractId [32]byte) (*types.Transaction, error) {
	return _HTLC.Contract.Refund(&_HTLC.TransactOpts, _contractId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(_contractId bytes32, _preimage bytes32) returns(bool)
func (_HTLC *HTLCTransactor) Withdraw(opts *bind.TransactOpts, _contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _HTLC.contract.Transact(opts, "withdraw", _contractId, _preimage)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(_contractId bytes32, _preimage bytes32) returns(bool)
func (_HTLC *HTLCSession) Withdraw(_contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _HTLC.Contract.Withdraw(&_HTLC.TransactOpts, _contractId, _preimage)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(_contractId bytes32, _preimage bytes32) returns(bool)
func (_HTLC *HTLCTransactorSession) Withdraw(_contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _HTLC.Contract.Withdraw(&_HTLC.TransactOpts, _contractId, _preimage)
}

// HTLCLogHTLCNewIterator is returned from FilterLogHTLCNew and is used to iterate over the raw logs and unpacked data for LogHTLCNew events raised by the HTLC contract.
type HTLCLogHTLCNewIterator struct {
	Event *HTLCLogHTLCNew // Event containing the contract specifics and raw log

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
func (it *HTLCLogHTLCNewIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HTLCLogHTLCNew)
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
		it.Event = new(HTLCLogHTLCNew)
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
func (it *HTLCLogHTLCNewIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HTLCLogHTLCNewIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HTLCLogHTLCNew represents a LogHTLCNew event raised by the HTLC contract.
type HTLCLogHTLCNew struct {
	ContractId [32]byte
	Sender     common.Address
	Receiver   common.Address
	Amount     *big.Int
	Hashlock   [32]byte
	Timelock   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCNew is a free log retrieval operation binding the contract event 0x329a8316ed9c3b2299597538371c2944c5026574e803b1ec31d6113e1cd67bde.
//
// Solidity: e LogHTLCNew(contractId indexed bytes32, sender indexed address, receiver indexed address, amount uint256, hashlock bytes32, timelock uint256)
func (_HTLC *HTLCFilterer) FilterLogHTLCNew(opts *bind.FilterOpts, contractId [][32]byte, sender []common.Address, receiver []common.Address) (*HTLCLogHTLCNewIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _HTLC.contract.FilterLogs(opts, "LogHTLCNew", contractIdRule, senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &HTLCLogHTLCNewIterator{contract: _HTLC.contract, event: "LogHTLCNew", logs: logs, sub: sub}, nil
}

// WatchLogHTLCNew is a free log subscription operation binding the contract event 0x329a8316ed9c3b2299597538371c2944c5026574e803b1ec31d6113e1cd67bde.
//
// Solidity: e LogHTLCNew(contractId indexed bytes32, sender indexed address, receiver indexed address, amount uint256, hashlock bytes32, timelock uint256)
func (_HTLC *HTLCFilterer) WatchLogHTLCNew(opts *bind.WatchOpts, sink chan<- *HTLCLogHTLCNew, contractId [][32]byte, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _HTLC.contract.WatchLogs(opts, "LogHTLCNew", contractIdRule, senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HTLCLogHTLCNew)
				if err := _HTLC.contract.UnpackLog(event, "LogHTLCNew", log); err != nil {
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

// HTLCLogHTLCRefundIterator is returned from FilterLogHTLCRefund and is used to iterate over the raw logs and unpacked data for LogHTLCRefund events raised by the HTLC contract.
type HTLCLogHTLCRefundIterator struct {
	Event *HTLCLogHTLCRefund // Event containing the contract specifics and raw log

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
func (it *HTLCLogHTLCRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HTLCLogHTLCRefund)
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
		it.Event = new(HTLCLogHTLCRefund)
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
func (it *HTLCLogHTLCRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HTLCLogHTLCRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HTLCLogHTLCRefund represents a LogHTLCRefund event raised by the HTLC contract.
type HTLCLogHTLCRefund struct {
	ContractId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCRefund is a free log retrieval operation binding the contract event 0x989b3a845197c9aec15f8982bbb30b5da714050e662a7a287bb1a94c81e2e70e.
//
// Solidity: e LogHTLCRefund(contractId indexed bytes32)
func (_HTLC *HTLCFilterer) FilterLogHTLCRefund(opts *bind.FilterOpts, contractId [][32]byte) (*HTLCLogHTLCRefundIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _HTLC.contract.FilterLogs(opts, "LogHTLCRefund", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &HTLCLogHTLCRefundIterator{contract: _HTLC.contract, event: "LogHTLCRefund", logs: logs, sub: sub}, nil
}

// WatchLogHTLCRefund is a free log subscription operation binding the contract event 0x989b3a845197c9aec15f8982bbb30b5da714050e662a7a287bb1a94c81e2e70e.
//
// Solidity: e LogHTLCRefund(contractId indexed bytes32)
func (_HTLC *HTLCFilterer) WatchLogHTLCRefund(opts *bind.WatchOpts, sink chan<- *HTLCLogHTLCRefund, contractId [][32]byte) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _HTLC.contract.WatchLogs(opts, "LogHTLCRefund", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HTLCLogHTLCRefund)
				if err := _HTLC.contract.UnpackLog(event, "LogHTLCRefund", log); err != nil {
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

// HTLCLogHTLCWithdrawIterator is returned from FilterLogHTLCWithdraw and is used to iterate over the raw logs and unpacked data for LogHTLCWithdraw events raised by the HTLC contract.
type HTLCLogHTLCWithdrawIterator struct {
	Event *HTLCLogHTLCWithdraw // Event containing the contract specifics and raw log

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
func (it *HTLCLogHTLCWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HTLCLogHTLCWithdraw)
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
		it.Event = new(HTLCLogHTLCWithdraw)
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
func (it *HTLCLogHTLCWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HTLCLogHTLCWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HTLCLogHTLCWithdraw represents a LogHTLCWithdraw event raised by the HTLC contract.
type HTLCLogHTLCWithdraw struct {
	ContractId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCWithdraw is a free log retrieval operation binding the contract event 0xd6fd4c8e45bf0c70693141c7ce46451b6a6a28ac8386fca2ba914044e0e23916.
//
// Solidity: e LogHTLCWithdraw(contractId indexed bytes32)
func (_HTLC *HTLCFilterer) FilterLogHTLCWithdraw(opts *bind.FilterOpts, contractId [][32]byte) (*HTLCLogHTLCWithdrawIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _HTLC.contract.FilterLogs(opts, "LogHTLCWithdraw", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &HTLCLogHTLCWithdrawIterator{contract: _HTLC.contract, event: "LogHTLCWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogHTLCWithdraw is a free log subscription operation binding the contract event 0xd6fd4c8e45bf0c70693141c7ce46451b6a6a28ac8386fca2ba914044e0e23916.
//
// Solidity: e LogHTLCWithdraw(contractId indexed bytes32)
func (_HTLC *HTLCFilterer) WatchLogHTLCWithdraw(opts *bind.WatchOpts, sink chan<- *HTLCLogHTLCWithdraw, contractId [][32]byte) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _HTLC.contract.WatchLogs(opts, "LogHTLCWithdraw", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HTLCLogHTLCWithdraw)
				if err := _HTLC.contract.UnpackLog(event, "LogHTLCWithdraw", log); err != nil {
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
