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
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

var EVENT_COUNTER = "event_counter"

// ManageBerth example simple Chaincode implementation
type ManageBerth struct {
}

var BerthIndexStr = "_Berthindex"				//name for the key/value that will store a list of all known Berth

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


// ============================================================================================================================
// Main - start the chaincode for Berth management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageBerth))
	if err != nil {
		fmt.Printf("Error starting Berth management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageBerth) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the chaincode
	msg = args[0]
	fmt.Println("ManageBerth chaincode is deployed successfully.");
	
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(BerthIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(EVENT_COUNTER, []byte("1"))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
	func (t *ManageBerth) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("run is running " + function)
		return t.Invoke(stub, function, args)
	}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
	func (t *ManageBerth) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_berth" {											//create a new Berth
		return t.create_berth(stub, args)
	}else if function == "delete_berth" {									// delete a Berth
		return t.delete_berth(stub, args)
	}else if function == "update_berth" {									//update a Berth
		return t.update_berth(stub, args)
	}else if function == "update_berth_allocationStatus" {									//update a Berth
		return t.update_berth_allocationStatus(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation")
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManageBerth) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getBerth_byVesselID" {													//Read a Berth by transId
		return t.getBerth_byVesselID(stub, args)
	} else if function == "getBerth_byTO" {													//Read all Berths
		return t.getBerth_byTO(stub, args)
	} else if function == "getBerth_byOwner" {													//Read all Berths
		return t.getBerth_byOwner(stub, args)
	} else if function == "getBerth_bySA" {													//Read all Berths
		return t.getBerth_bySA(stub, args)
	} else if function == "getBerth_byPA" {													//Read all Berths
		return t.getBerth_byPA(stub, args)
	} else if function == "get_AllBerth" {													//Read all Berths
		return t.get_AllBerth(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query")
}
// ============================================================================================================================
// getBerth_byVesselID - get Berth details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageBerth) getBerth_byVesselID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var vesselID, jsonResp string
	var err error
	fmt.Println("start getBerth_byVesselID")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the vessel to query")
	}
	// set berthID
	vesselID = args[0]
	valAsbytes, err := stub.GetState(vesselID)									//get the vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("valAsbytes : ")
	//fmt.Println(valAsbytes)
	fmt.Println("end getBerth_byVesselID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// getBerth_byTO - get Berth details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageBerth) getBerth_byTO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, toID, errResp string
	var berthIndex []string
	var valIndex Berth
	fmt.Println("start getBerth_byTO")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set buyer's name
	toID = args[0]
	//fmt.Println("buyerName" + buyerName)
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	count := 0
	jsonResp = "{"
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getBerth_byTO")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.TOID == toID{
			if count > 0 {
				jsonResp = jsonResp + ","
			}
			fmt.Println("TO found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			count++
		}
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getBerth_byTO")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
// getBerth_byID - get Berth details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageBerth) getBerth_byOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, ownerName, errResp string
	var berthIndex []string
	var valIndex Berth
	fmt.Println("start getBerth_byOwner")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set buyer's name
	ownerName = args[0]
	//fmt.Println("buyerName" + buyerName)
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getBerth_byOwner")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.OwnerName == ownerName{
			fmt.Println("ownerName found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(berthIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getBerth_byOwner")
	return []byte(jsonResp), nil											//send it onward	
}
// ============================================================================================================================
// getBerth_byID - get Berth details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageBerth) getBerth_bySA(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, agentRefNumber, errResp string
	var berthIndex []string
	var valIndex Berth
	fmt.Println("start getBerth_bySA")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set buyer's name
	agentRefNumber = args[0]
	//fmt.Println("buyerName" + buyerName)
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getBerth_bySA")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.AgentRefNumber == agentRefNumber{
			fmt.Println("agentRefNumber found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(berthIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getBerth_bySA")
	return []byte(jsonResp), nil
}
// ============================================================================================================================
// getBerth_byID - get Berth details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageBerth) getBerth_byPA(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, approverID, errResp string
	var berthIndex []string
	var valIndex Berth
	fmt.Println("start getBerth_byPA")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set buyer's name
	approverID = args[0]
	//fmt.Println("buyerName" + buyerName)
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getBerth_byPA")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.ApproverID == approverID{
			fmt.Println("approverID found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(berthIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} 
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getBerth_byPA")
	return []byte(jsonResp), nil
}
// ============================================================================================================================
//  get_AllBerth- get details of all Berth from chaincode state
// ============================================================================================================================
func (t *ManageBerth) get_AllBerth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var berthIndex []string
	fmt.Println("start get_AllBerth")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index")
	}
	//fmt.Print("berthAsBytes : ")
	//fmt.Println(berthAsBytes)
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	jsonResp = "{"
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Berth")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(berthIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllBerth")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// Delete - remove a Berth from chain
// ============================================================================================================================
func (t *ManageBerth) delete_berth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// set berthID
	vesselID := args[0]
	err := stub.DelState(vesselID)													//remove the Berth from chaincode
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the Berth index
	berthAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index")
	}
	//fmt.Println("poAsBytes in delete po")
	//fmt.Println(poAsBytes);
	var berthIndex []string
	json.Unmarshal(berthAsBytes, &berthIndex)								//un stringify it aka JSON.parse()
	//fmt.Println("poIndex in delete po")
	//fmt.Println(poIndex);
	//remove marble from index
	for i,val := range berthIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + vesselID)
		if val == vesselID{															//find the correct Berth
			fmt.Println("found Berth with matching berthID")
			berthIndex = append(berthIndex[:i], berthIndex[i+1:]...)			//remove it
			for x:= range berthIndex{											//debug prints...
				fmt.Println(string(x) + " - " + berthIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(berthIndex)									//save new index
	err = stub.PutState(BerthIndexStr, jsonAsBytes)
	return nil, nil
}
// ============================================================================================================================
// Write - update Berth into chaincode state
// ============================================================================================================================
func (t *ManageBerth) update_berth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_berth")
	if len(args) != 19 {
		return nil, errors.New("Incorrect number of arguments. Expecting 15.")
	}
	// set vesselID
	vesselID := args[0]
	berthAsBytes, err := stub.GetState(vesselID)									//get the Berth for the specified vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("berthAsBytes in update berth")
	//fmt.Println(berthAsBytes);
	res := Berth{}
	json.Unmarshal(berthAsBytes, &res)
	if res.VesselID == vesselID{
		fmt.Println("Berth found with vesselID : " + vesselID)
		//fmt.Println(res);
		res.VesselName = args[1]
		res.VesselType = args[2]
		res.VesselClass = args[3]
		res.ShippingLine = args[4]
		res.AgentRefNumber = args[5]
		res.ArrivalPort = args[6]
		res.InboundVoyageNo = args[7]
		res.OutboundVoyageNo = args[8]
		res.ArriveFrom = args[9]
		res.Terminal = args[10]
		res.Remarks = args[11]
		res.BerthBookingStatus = "New"
		res.RotationNumber = args[12]
		res.TOID = args[13]
		res.ApproverID = args[14]
		res.MMSInumber = args[15]
		res.PortOfRegisteration = args[16]
		res.OwnerName = args[17]
		res.OwnerPhoneNumber = args[18]
	}
	
	//build the Berth json string manually
	berthDetails := 	`{`+
		`"vesselID": "` + res.VesselID + `" , `+
		`"vesselName": "` + res.VesselName + `" , `+
		`"vesselType": "` + res.VesselType + `" , `+
		`"vesselClass": "` + res.VesselClass + `" , `+
		`"shippingLine": "` + res.ShippingLine + `" , `+ 
		`"agentRefNumber": "` + res.AgentRefNumber + `" , `+ 
		`"arrivalPort": "` + res.ArrivalPort + `" , `+ 
		`"inboundVoyageNo": "` + res.InboundVoyageNo + `" , `+ 
		`"outboundVoyageNo": "` + res.OutboundVoyageNo + `" , `+ 
		`"arriveFrom": "` + res.ArriveFrom + `" , `+ 
		`"terminal": "` + res.Terminal + `" , `+ 
		`"remarks": "` +  res.Remarks + `" , `+ 
		`"berthBookingStatus": "` + res.BerthBookingStatus + `" , `+ 
		`"rotationNumber": "` + res.RotationNumber + `" , `+ 
		`"toID": "` +  res.TOID + `" , `+ 
		`"approverID": "` + res.ApproverID + `" , `+ 
		`"mmsiNumber": "` + res.MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + res.PortOfRegisteration + `" , `+ 
		`"ownerName": "` + res.OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + res.OwnerPhoneNumber + `" `+ 
		`}`
	err = stub.PutState(vesselID, []byte(berthDetails))									//store Berth with id as key
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// create Berth - create a new Berth, store into chaincode state
// ============================================================================================================================
func (t *ManageBerth) create_berth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 19 {
		return nil, errors.New("Incorrect number of arguments. Expecting 15")
	}
	fmt.Println("start create_berth")

	VesselID := args[0]
	VesselName := args[1]
	VesselType := args[2]
	VesselClass := args[3]
	ShippingLine := args[4]
	AgentRefNumber := args[5]
	ArrivalPort := args[6]
	InboundVoyageNo := args[7]
	OutboundVoyageNo := args[8]
	ArriveFrom := args[9]
	Terminal := args[10]
	Remarks := args[11]
	BerthBookingStatus := "New"
	RotationNumber := args[12]
	TOID := args[13]
	ApproverID := args[14]
	MMSInumber := args[15]
	PortOfRegisteration := args[16]
	OwnerName := args[17]
	OwnerPhoneNumber := args[18]
	
	berthAsBytes, err := stub.GetState(VesselID)
	if err != nil {
		return nil, errors.New("Failed to get Berth VesselID")
	}
	//fmt.Print("berthAsBytes: ")
	//fmt.Println(berthAsBytes)
	res := Berth{}
	json.Unmarshal(berthAsBytes, &res)
	//fmt.Print("res: ")
	//fmt.Println(res)
	if res.VesselID == VesselID{
		//fmt.Println("This Berth arleady exists: " + BerthID)
		//fmt.Println(res);
		return nil, errors.New("This Berth arleady exists")				//all stop a Berth by this name exists
	}
	
	//build the Berth json string manually
	berthDetails := 	`{`+
		`"vesselID": "` + VesselID + `" , `+
		`"vesselName": "` + VesselName + `" , `+
		`"vesselType": "` + VesselType + `" , `+
		`"vesselClass": "` + VesselClass + `" , `+
		`"shippingLine": "` + ShippingLine + `" , `+ 
		`"agentRefNumber": "` + AgentRefNumber + `" , `+ 
		`"arrivalPort": "` + ArrivalPort + `" , `+ 
		`"inboundVoyageNo": "` + InboundVoyageNo + `" , `+ 
		`"outboundVoyageNo": "` + OutboundVoyageNo + `" , `+ 
		`"arriveFrom": "` + ArriveFrom + `" , `+ 
		`"terminal": "` + Terminal + `" , `+ 
		`"remarks": "` +  Remarks + `" , `+ 
		`"berthBookingStatus": "` + BerthBookingStatus + `" , `+ 
		`"rotationNumber": "` + RotationNumber + `" , `+ 
		`"toID": "` +  TOID + `" , `+ 
		`"approverID": "` + ApproverID + `" , `+ 
		`"mmsiNumber": "` + MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + PortOfRegisteration + `" , `+ 
		`"ownerName": "` + OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + OwnerPhoneNumber + `" `+ 
		`}`

		//fmt.Println("berthDetails: " + berthDetails)
		fmt.Print("Berth details in bytes array: ")
		fmt.Println([]byte(berthDetails))
	err = stub.PutState(VesselID, []byte(berthDetails))									//store Berth with BerthID as key
	if err != nil {
		return nil, err
	}
	//get the Berth index
	berthIndexAsBytes, err := stub.GetState(BerthIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Berth index")
	}
	var berthIndex []string
	//fmt.Print("berthIndexAsBytes: ")
	//fmt.Println(berthIndexAsBytes)
	
	json.Unmarshal(berthIndexAsBytes, &berthIndex)							//un stringify it aka JSON.parse()
	//fmt.Print("poIndex after unmarshal..before append: ")
	//fmt.Println(poIndex)
	//append
	berthIndex = append(berthIndex, VesselID)									//add Berth transID to index list
	//fmt.Println("! Berth index after appending transId: ", poIndex)
	jsonAsBytes, _ := json.Marshal(berthIndex)
	//fmt.Print("jsonAsBytes: ")
	//fmt.Println(jsonAsBytes)
	err = stub.PutState(BerthIndexStr, jsonAsBytes)						//store name of Berth
	if err != nil {
		return nil, err
	}

	fmt.Println("end create_berth")
	return nil, nil
}

// ============================================================================================================================
// Write - update Berth into chaincode state
// ============================================================================================================================
func (t *ManageBerth) update_berth_allocationStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_berth_allocationStatus")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}
	// set vesselID
	vesselID := args[0]
	berthAsBytes, err := stub.GetState(vesselID)									//get the Berth for the specified vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("berthAsBytes in update berth")
	//fmt.Println(berthAsBytes);
	res := Berth{}
	json.Unmarshal(berthAsBytes, &res)
	if res.VesselID == vesselID{
		fmt.Println("Berth found with vesselID : " + vesselID)
		res.BerthBookingStatus = args[1]
		res.ApproverID = args[2]
	}
	
	//build the Berth json string manually
	berthDetails := 	`{`+
		`"vesselID": "` + res.VesselID + `" , `+
		`"vesselName": "` + res.VesselName + `" , `+
		`"vesselType": "` + res.VesselType + `" , `+
		`"vesselClass": "` + res.VesselClass + `" , `+
		`"shippingLine": "` + res.ShippingLine + `" , `+ 
		`"agentRefNumber": "` + res.AgentRefNumber + `" , `+ 
		`"arrivalPort": "` + res.ArrivalPort + `" , `+ 
		`"inboundVoyageNo": "` + res.InboundVoyageNo + `" , `+ 
		`"outboundVoyageNo": "` + res.OutboundVoyageNo + `" , `+ 
		`"arriveFrom": "` + res.ArriveFrom + `" , `+ 
		`"terminal": "` + res.Terminal + `" , `+ 
		`"remarks": "` +  res.Remarks + `" , `+ 
		`"berthBookingStatus": "` + res.BerthBookingStatus + `" , `+ 
		`"rotationNumber": "` + res.RotationNumber + `" , `+ 
		`"toID": "` +  res.TOID + `" , `+ 
		`"approverID": "` + res.ApproverID + `" , `+ 
		`"mmsiNumber": "` + res.MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + res.PortOfRegisteration + `" , `+ 
		`"ownerName": "` + res.OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + res.OwnerPhoneNumber + `" `+ 
		`}`
	err = stub.PutState(vesselID, []byte(berthDetails))									//store Berth with id as key
	if err != nil {
		return nil, err
	}
	return nil, nil
}
