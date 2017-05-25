/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
	"strconv"
)

type ManageAllocations struct {
}

type Berth struct{							// Attributes of a Berth 				
	VesselID string `json:"vesselID"`
	VesselName string `json:"vesselName"`					
	VesselType string `json:"vesselType"`
	VesselClass string `json:"vesselClass"`
	ShippingLine string `json:"shippingLine"`
	AgentRefNumber string `json:"agentRefNumber"`
	ArrivalPort string `json:"arrivalPort"`
	InboundVoyageNo string `json:"inboundVoyageNo"`
	OutboundVoyageNo string `json:"outboundVoyageNo"`
	ArriveFrom string `json:"arriveFrom"`
	Terminal string `json:"terminal"`
	Remarks string `json:"remarks"`
	BerthBookingStatus string `json:"berthBookingStatus"`
	RotationNumber string `json:"rotationNumber"`
	TOID string `json:"toID"`
	ApproverID string `json:"approverID"`
	MMSInumber string `json:"mmsiNumber"`
	PortOfRegisteration string `json:"portOfRegisteration"`
	OwnerName string `json:"ownerName"`
	OwnerPhoneNumber string `json:"ownerPhoneNumber"`
}

type Vessel struct{							// Attributes of a Vessel 				
	VesselID string `json:"vesselID"`
	VesselName string `json:"vesselName"`					
	VesselType string `json:"vesselType"`
	SIN string `json:"sin"`
	MMSInumber string `json:"mmsiNumber"`
	PortOfRegisteration string `json:"portOfRegisteration"`
	OwnerName string `json:"ownerName"`
	OwnerPhoneNumber string `json:"ownerPhoneNumber"`
	OwnerAddressLine1 string `json:"ownerAddressLine1"`
	OwnerAddressLine2 string `json:"ownerAddressLine2"`
	OwnerAddressLine3 string `json:"ownerAddressLine3"`
	OwnerCity string `json:"ownerCity"`
	OwnerState string `json:"ownerState"`
	OwnerPostCode string `json:"ownerPostCode"`
	OwnerCountry string `json:"ownerCountry"`
	VesselClass string `json:"vesselClass"`
	BerthBookingStatus string `json:"berthBookingStatus"`
}

// ============================================================================================================================
// Main - start the chaincode for Allocation management
// ============================================================================================================================
func main() {
	err := shim.Start(new(ManageAllocations))
	if err != nil {
		fmt.Printf("Error starting Allocation management chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageAllocations) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	// Initialize the chaincode
	msg = args[0]
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg)) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState("_init", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	tosend := "{ \"message\" : \"ManageAllocations chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry Dealint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageAllocations) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry Dealint for Invocations
// ============================================================================================================================
func (t *ManageAllocations) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { // Initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "berth_allocation" { // Create a new Allocation
		return t.berth_allocation(stub, args)
	} else if function == "approve_allocation" { // Secondary Fire when Longbox account is updated
		return t.approve_allocation(stub, args)
	} else if function == "reject_allocation" { // Secondary Fire when Longbox account is updated
		return t.reject_allocation(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Query - Our entry Dealint for Queries
// ============================================================================================================================

func (t *ManageAllocations) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	/*if function == "nil" {
		return t.nil(stub, args)
	}*/
	fmt.Println("Allocation does not support query functions.")
	errMsg := "{ \"message\" : \"Allocation does not support query functions.\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// ============================================================================================================================
// Start Allocation - create a new Allocation, store into chaincode state
// ============================================================================================================================
func (t *ManageAllocations) berth_allocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	fmt.Println("start start_allocation")

	// Alloting Params
	VesselChaincode := args[0]
	BerthChainCode := args[1]
	VesselID := args[2]

	// Json to create report
	reportInJson := `{`

	//-----------------------------------------------------------------------------

	// Fetch Vessel details from Blockchain
	f := "getVessel_byID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	vesselAsBytes, err := stub.QueryChaincode(VesselChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	VesselData := Vessel{}
	json.Unmarshal(vesselAsBytes, &VesselData)
	fmt.Println(VesselData)
	if VesselData.VesselID == VesselID {
		fmt.Println("Vessel found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}


	// Fetch Berth details from Blockchain
	f := "getBerth_byVesselID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	berthAsBytes, err := stub.QueryChaincode(BerthChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	BerthData := Berth{}
	json.Unmarshal(berthAsBytes, &BerthData)
	fmt.Println(BerthData)
	if BerthData.VesselID == VesselID {
		fmt.Println("Berth found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Update allocation status to "Allocation in progress"
	function = "update_vessel_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "P")
	result, err := stub.InvokeChaincode(VesselChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Vessel' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	// Update allocation status to "Allocation in progress"
	function = "update_berth_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "P")
	result, err := stub.InvokeChaincode(BerthChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Berth' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	fmt.Println("end start_allocation")
	return nil, nil
}

// ============================================================================================================================
// Start Allocation - create a new Allocation, store into chaincode state
// ============================================================================================================================
func (t *ManageAllocations) approve_allocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	fmt.Println("start approve_allocation")

	// Alloting Params
	VesselChaincode := args[0]
	BerthChainCode := args[1]
	VesselID := args[2]

	// Json to create report
	reportInJson := `{`

	//-----------------------------------------------------------------------------

	// Fetch Vessel details from Blockchain
	f := "getVessel_byID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	vesselAsBytes, err := stub.QueryChaincode(VesselChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	VesselData := Vessel{}
	json.Unmarshal(vesselAsBytes, &VesselData)
	fmt.Println(VesselData)
	if VesselData.VesselID == VesselID {
		fmt.Println("Vessel found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}


	// Fetch Berth details from Blockchain
	f := "getBerth_byVesselID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	berthAsBytes, err := stub.QueryChaincode(BerthChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	BerthData := Berth{}
	json.Unmarshal(berthAsBytes, &BerthData)
	fmt.Println(BerthData)
	if BerthData.VesselID == VesselID {
		fmt.Println("Berth found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Update allocation status to "Allocation in progress"
	function = "update_vessel_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "A")
	result, err := stub.InvokeChaincode(VesselChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Vessel' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	// Update allocation status to "Allocation in progress"
	function = "update_berth_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "A")
	result, err := stub.InvokeChaincode(BerthChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Berth' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	fmt.Println("end approve_allocation")
	return nil, nil
}

// ============================================================================================================================
// Start Allocation - create a new Allocation, store into chaincode state
// ============================================================================================================================
func (t *ManageAllocations) reject_allocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	fmt.Println("start approve_allocation")

	// Alloting Params
	VesselChaincode := args[0]
	BerthChainCode := args[1]
	VesselID := args[2]

	// Json to create report
	reportInJson := `{`

	//-----------------------------------------------------------------------------

	// Fetch Vessel details from Blockchain
	f := "getVessel_byID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	vesselAsBytes, err := stub.QueryChaincode(VesselChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	VesselData := Vessel{}
	json.Unmarshal(vesselAsBytes, &VesselData)
	fmt.Println(VesselData)
	if VesselData.VesselID == VesselID {
		fmt.Println("Vessel found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}


	// Fetch Berth details from Blockchain
	f := "getBerth_byVesselID"
	queryArgs := util.ToChaincodeArgs(f, VesselID)
	berthAsBytes, err := stub.QueryChaincode(BerthChaincode, queryArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	BerthData := Berth{}
	json.Unmarshal(berthAsBytes, &BerthData)
	fmt.Println(BerthData)
	if BerthData.VesselID == VesselID {
		fmt.Println("Berth found with VesselID : " + VesselID)
	} else {
		errMsg := "{ \"message\" : \"" + VesselID + " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Update allocation status to "Allocation in progress"
	function = "update_vessel_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "R")
	result, err := stub.InvokeChaincode(VesselChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Vessel' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	// Update allocation status to "Allocation in progress"
	function = "update_berth_allocationStatus"
	invokeArgs := util.ToChaincodeArgs(function, VesselID, "R")
	result, err := stub.InvokeChaincode(BerthChaincode, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update Transaction status from 'Berth' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	fmt.Print("Transaction hash returned: ")
	fmt.Println(result)
	fmt.Println("Successfully updated allocation status to 'In progress'")

	fmt.Println("end approve_allocation")
	return nil, nil
}