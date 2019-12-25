/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type WeShare struct{}

// User example simple Chaincode implementation
type User struct {
	UserId string
	Amount float64
}

const shareAmount string = "100"
const listenAmonunt string = "10"
const rewardPoolId string = "reward_pool_id"

func (t *WeShare) Init(stub shim.ChaincodeStubInterface) pb.Response {
	t.initUser(stub, []string{rewardPoolId})
	return shim.Success(nil)
}

func (t *WeShare) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "initUser" {
		return t.initUser(stub, args)
	} else if function == "completeShare" {
		return t.completeShare(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "shopping" {
		return t.shopping(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *WeShare) initUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]
	oldUserBytes, err := stub.GetState(userId)

	if err != nil {
		return shim.Error(err.Error())
	}

	if oldUserBytes != nil {
		return shim.Success(nil)
	}

	user := User{userId, 0}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Write the state to the ledger
	err = stub.PutState(userId, userBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *WeShare) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]
	// Get the state from the ledger
	userBytes, err := stub.GetState(userId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if userBytes == nil {
		t.initUser(stub, []string{userId})
		jsonResp := "{\"UserId\":\"" + userId + "\",\"Amount\":\"" + strconv.FormatFloat(0, 'E', -1, 64) + "\"}"
		return shim.Success([]byte(jsonResp))
	}

	var userInfo User
	err = json.Unmarshal(userBytes, &userInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "{\"Name\":\"" + userInfo.UserId + "\",\"Amount\":\"" + strconv.FormatFloat(userInfo.Amount, 'E', -1, 64) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(userBytes)
}

func (t *WeShare) completeShare(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]
	var params = []string{userId}
	t.share(stub, params)

	listeners := args[1:]
	t.listen(stub, listeners)
	jsonResp := "{\"shareAmount\":" + shareAmount + ",\"listenAmount\":" + listenAmonunt + "}"
	return shim.Success([]byte(jsonResp))
}

func (t *WeShare) share(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]
	var params = []string{userId, shareAmount}
	t.invoke(stub, params)
	return shim.Success(nil)
}

func (t *WeShare) listen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	for i := 0; i < len(args); i++ {
		userId := args[i]
		var params = []string{userId, listenAmonunt}
		t.invoke(stub, params)
	}
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *WeShare) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userId string // Entities
	var userInfo User // Asset holdings
	var userAmount float64
	var err error

	userId = args[0]

	oldUserInfoBytes, err := stub.GetState(userId)
	err = json.Unmarshal(oldUserInfoBytes, &userInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	userAmount = userInfo.Amount

	// Perform the execution
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	userInfo.Amount = userAmount + amount

	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Write the state to the ledger
	err = stub.PutState(userInfo.UserId, userInfoBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *WeShare) shopping(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	transferUserId := args[0]
	amount := args[1]
	receiverId := rewardPoolId

	params := []string{transferUserId, receiverId, amount}
	t.transfer(stub, params)
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *WeShare) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var transferUserId, receiverId string   // Entities
	var transferUserInfo, receiverInfo User // Asset holdings
	var transferUserAmount, receiverAmount float64
	var err error

	transferUserId = args[0]
	receiverId = args[1]

	oldTransferUserInfoBytes, err := stub.GetState(transferUserId)
	err = json.Unmarshal(oldTransferUserInfoBytes, &transferUserInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	oldReceiverInfoBytes, err := stub.GetState(receiverId)
	err = json.Unmarshal(oldReceiverInfoBytes, &receiverInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	transferUserAmount = transferUserInfo.Amount
	receiverAmount = receiverInfo.Amount

	// Perform the execution
	amount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	transferUserInfo.Amount = transferUserAmount - amount
	receiverInfo.Amount = receiverAmount + amount

	transferUserInfoBytes, err := json.Marshal(transferUserInfo)
	// Write the state to the ledger
	err = stub.PutState(transferUserInfo.UserId, transferUserInfoBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiverInfoBytes, err := json.Marshal(receiverInfo)
	// Write the state to the ledger
	err = stub.PutState(receiverInfo.UserId, receiverInfoBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(WeShare))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
