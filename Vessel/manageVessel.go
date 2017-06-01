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

// ManageVessel example simple Chaincode implementation
type ManageVessel struct {
}

var VesselIndexStr = "_Vesselindex"				//name for the key/value that will store a list of all known Vessel

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
// Main - start the chaincode for Vessel management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageVessel))
	if err != nil {
		fmt.Printf("Error starting Vessel management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageVessel) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the chaincode
	msg = args[0]
	fmt.Println("ManageVessel chaincode is deployed successfully.");
	
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(VesselIndexStr, jsonAsBytes)
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
	func (t *ManageVessel) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("run is running " + function)
		return t.Invoke(stub, function, args)
	}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
	func (t *ManageVessel) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_vessel" {											//create a new Vessel
		return t.create_vessel(stub, args)
	}else if function == "delete_vessel" {									// delete a Vessel
		return t.delete_vessel(stub, args)
	}else if function == "update_vessel" {									//update a Vessel
		return t.update_vessel(stub, args)
	}else if function == "update_vessel_allocationStatus" {									//update a Vessel
		return t.update_vessel_allocationStatus(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation")
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManageVessel) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getVessel_byID" {													//Read a Vessel by transId
		return t.getVessel_byID(stub, args)
	} else if function == "getVessel_byOwner" {													//Read all Vessels
		return t.getVessel_byOwner(stub, args)
	} else if function == "get_AllVessel" {													//Read all Vessels
		return t.get_AllVessel(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query")
}
// ============================================================================================================================
// getVessel_byID - get Vessel details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageVessel) getVessel_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var vesselID, jsonResp string
	var err error
	fmt.Println("start getVessel_byID")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the vessel to query")
	}
	// set vesselID
	vesselID = args[0]
	valAsbytes, err := stub.GetState(vesselID)									//get the vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("valAsbytes : ")
	//fmt.Println(valAsbytes)
	fmt.Println("end getVessel_byID")
	return valAsbytes, nil													//send it onward
}

// ============================================================================================================================
// getVessel_byOwner - get Vessel details for a specific ID from chaincode state
// ============================================================================================================================

func (t *ManageVessel) getVessel_byOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, ownerName, errResp string
	var vesselIndex []string
	var valIndex Vessel
	fmt.Println("start getVessel_byOwner")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting owner name")
	}
	// set buyer's name
	ownerPhoneNumber = args[0]
	//fmt.Println("buyerName" + buyerName)
	vesselAsBytes, err := stub.GetState(VesselIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Vessel index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(vesselAsBytes, &vesselIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range vesselIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getVessel_byOwner")
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
		if valIndex.OwnerPhoneNumber == ownerPhoneNumber{
			fmt.Println("Owner found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(vesselIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getVessel_byOwner")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  get_AllVessel- get details of all Vessel from chaincode state
// ============================================================================================================================
func (t *ManageVessel) get_AllVessel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var vesselIndex []string
	fmt.Println("start get_AllVessel")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	vesselAsBytes, err := stub.GetState(VesselIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Vessel index")
	}
	//fmt.Print("vesselAsBytes : ")
	//fmt.Println(vesselAsBytes)
	json.Unmarshal(vesselAsBytes, &vesselIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	jsonResp = "{"
	for i,val := range vesselIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Vessel")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(vesselIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllVessel")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// Delete - remove a Vessel from chain
// ============================================================================================================================
func (t *ManageVessel) delete_vessel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// set vesselID
	vesselID := args[0]

	//get the Vessel index
	vesselAsBytes, err := stub.GetState(VesselIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Vessel index")
	}
	//fmt.Println("poAsBytes in delete po")
	//fmt.Println(poAsBytes);
	var vesselIndex []string
	json.Unmarshal(vesselAsBytes, &vesselIndex)								//un stringify it aka JSON.parse()
	//fmt.Println("poIndex in delete po")
	//fmt.Println(poIndex);
	//remove marble from index
	for i,val := range vesselIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + vesselID)
		if val == vesselID{															//find the correct Vessel
			fmt.Println("found Vessel with matching vesselID")
			vesselIndex = append(vesselIndex[:i], vesselIndex[i+1:]...)			//remove it
			for x:= range vesselIndex{											//debug prints...
				fmt.Println(string(x) + " - " + vesselIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(vesselIndex)									//save new index
	err = stub.PutState(VesselIndexStr, jsonAsBytes)
	return nil, nil
}
// ============================================================================================================================
// Write - update Vessel into chaincode state
// ============================================================================================================================
func (t *ManageVessel) update_vessel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_vessel")
	if len(args) != 16 {
		return nil, errors.New("Incorrect number of arguments. Expecting 15.")
	}
	// set vesselID
	vesselID := args[0]
	vesselAsBytes, err := stub.GetState(vesselID)									//get the Vessel for the specified vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("vesselAsBytes in update vessel")
	//fmt.Println(vesselAsBytes);
	res := Vessel{}
	json.Unmarshal(vesselAsBytes, &res)
	if res.VesselID == vesselID{
		fmt.Println("Vessel found with vesselID : " + vesselID)
		//fmt.Println(res);
		res.VesselName = args[1]
		res.VesselType = args[2]
		res.SIN = args[3]
		res.MMSInumber = args[4]
		res.PortOfRegisteration = args[5]
		res.OwnerName = args[6]
		res.OwnerPhoneNumber = args[7]
		res.OwnerAddressLine1 = args[8]
		res.OwnerAddressLine2 = args[9]
		res.OwnerAddressLine3 = args[10]
		res.OwnerCity = args[11]
		res.OwnerState = args[12]
		res.OwnerPostCode = args[13]
		res.OwnerCountry = args[14]
		res.VesselClass = args[15]
		res.BerthBookingStatus = "N"
	}
	
	//build the Vessel json string manually
	vesselDetails := 	`{`+
		`"vesselID": "` + res.VesselID + `" , `+
		`"vesselName": "` + res.VesselName + `" , `+
		`"vesselType": "` + res.VesselType + `" , `+
		`"sin": "` + res.SIN + `" , `+
		`"mmsiNumber": "` + res.MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + res.PortOfRegisteration + `" , `+ 
		`"ownerName": "` + res.OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + res.OwnerPhoneNumber + `" , `+ 
		`"ownerAddressLine1": "` + res.OwnerAddressLine1 + `" , `+ 
		`"ownerAddressLine2": "` + res.OwnerAddressLine2 + `" , `+ 
		`"ownerAddressLine3": "` + res.OwnerAddressLine3 + `" , `+ 
		`"ownerCity": "` +  res.OwnerCity + `", `+ 
		`"ownerState": "` + res.OwnerState + `" , `+ 
		`"ownerPostCode": "` + res.OwnerPostCode + `" , `+ 
		`"ownerCountry": "` +  res.OwnerCountry + `" , `+ 
		`"vesselClass": "` +  res.VesselClass + `" , `+
		`"berthBookingStatus": "` +  res.BerthBookingStatus + `" `+ 
		`}`
	err = stub.PutState(vesselID, []byte(vesselDetails))									//store Vessel with id as key
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// create Vessel - create a new Vessel, store into chaincode state
// ============================================================================================================================
func (t *ManageVessel) create_vessel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 16 {
		return nil, errors.New("Incorrect number of arguments. Expecting 15")
	}
	fmt.Println("start create_vessel")

	VesselID := args[0]
	VesselName := args[1]
	VesselType := args[2]
	SIN := args[3]
	MMSInumber := args[4]
	PortOfRegisteration := args[5]
	OwnerName := args[6]
	OwnerPhoneNumber := args[7]
	OwnerAddressLine1 := args[8]
	OwnerAddressLine2 := args[9]
	OwnerAddressLine3 := args[10]
	OwnerCity := args[11]
	OwnerState := args[12]
	OwnerPostCode := args[13]
	OwnerCountry := args[14]
	VesselClass := args[15]
	BerthBookingStatus := "N"
	
	vesselAsBytes, err := stub.GetState(VesselID)
	if err != nil {
		return nil, errors.New("Failed to get Vessel VesselID")
	}
	//fmt.Print("vesselAsBytes: ")
	//fmt.Println(vesselAsBytes)
	res := Vessel{}
	json.Unmarshal(vesselAsBytes, &res)
	//fmt.Print("res: ")
	//fmt.Println(res)
	if res.VesselID == VesselID{
		//fmt.Println("This Vessel arleady exists: " + VesselID)
		//fmt.Println(res);
		return nil, errors.New("This Vessel arleady exists")				//all stop a Vessel by this name exists
	}
	
	//build the Vessel json string manually
	vesselDetails := 	`{`+
		`"vesselID": "` + VesselID + `" , `+
		`"vesselName": "` + VesselName + `" , `+
		`"vesselType": "` + VesselType + `" , `+
		`"sin": "` + SIN + `" , `+
		`"mmsiNumber": "` + MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + PortOfRegisteration + `" , `+ 
		`"ownerName": "` + OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + OwnerPhoneNumber + `" , `+ 
		`"ownerAddressLine1": "` + OwnerAddressLine1 + `" , `+ 
		`"ownerAddressLine2": "` + OwnerAddressLine2 + `" , `+ 
		`"ownerAddressLine3": "` + OwnerAddressLine3 + `" , `+ 
		`"ownerCity": "` +  OwnerCity + `", `+ 
		`"ownerState": "` + OwnerState + `" , `+ 
		`"ownerPostCode": "` + OwnerPostCode + `" , `+ 
		`"ownerCountry": "` +  OwnerCountry + `" , `+ 
		`"vesselClass": "` + VesselClass + `" , `+
		`"berthBookingStatus": "` + BerthBookingStatus + `" `+
		`}`

		//fmt.Println("vesselDetails: " + vesselDetails)
		fmt.Print("Vessel details in bytes array: ")
		fmt.Println([]byte(vesselDetails))
	err = stub.PutState(VesselID, []byte(vesselDetails))									//store Vessel with VesselID as key
	if err != nil {
		return nil, err
	}
	//get the Vessel index
	vesselIndexAsBytes, err := stub.GetState(VesselIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Vessel index")
	}
	var vesselIndex []string
	//fmt.Print("vesselIndexAsBytes: ")
	//fmt.Println(vesselIndexAsBytes)
	
	json.Unmarshal(vesselIndexAsBytes, &vesselIndex)							//un stringify it aka JSON.parse()
	//fmt.Print("poIndex after unmarshal..before append: ")
	//fmt.Println(poIndex)
	//append
	vesselIndex = append(vesselIndex, VesselID)									//add Vessel transID to index list
	//fmt.Println("! Vessel index after appending transId: ", poIndex)
	jsonAsBytes, _ := json.Marshal(vesselIndex)
	//fmt.Print("jsonAsBytes: ")
	//fmt.Println(jsonAsBytes)
	err = stub.PutState(VesselIndexStr, jsonAsBytes)						//store name of Vessel
	if err != nil {
		return nil, err
	}

	fmt.Println("end create_vessel")
	return nil, nil
}

// ============================================================================================================================
// Write - update Vessel into chaincode state
// ============================================================================================================================
func (t *ManageVessel) update_vessel_allocationStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_vessel_allocationStatus")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}
	// set vesselID
	vesselID := args[0]
	vesselAsBytes, err := stub.GetState(vesselID)									//get the Vessel for the specified vesselID from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + vesselID + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("vesselAsBytes in update vessel")
	//fmt.Println(vesselAsBytes);
	res := Vessel{}
	json.Unmarshal(vesselAsBytes, &res)
	if res.VesselID == vesselID{
		fmt.Println("Vessel found with vesselID : " + vesselID)
		//fmt.Println(res);
		res.BerthBookingStatus = args[1]
	}
	
	//build the Vessel json string manually
	vesselDetails := 	`{`+
		`"vesselID": "` + res.VesselID + `" , `+
		`"vesselName": "` + res.VesselName + `" , `+
		`"vesselType": "` + res.VesselType + `" , `+
		`"sin": "` + res.SIN + `" , `+
		`"mmsiNumber": "` + res.MMSInumber + `" , `+ 
		`"portOfRegisteration": "` + res.PortOfRegisteration + `" , `+ 
		`"ownerName": "` + res.OwnerName + `" , `+ 
		`"ownerPhoneNumber": "` + res.OwnerPhoneNumber + `" , `+ 
		`"ownerAddressLine1": "` + res.OwnerAddressLine1 + `" , `+ 
		`"ownerAddressLine2": "` + res.OwnerAddressLine2 + `" , `+ 
		`"ownerAddressLine3": "` + res.OwnerAddressLine3 + `" , `+ 
		`"ownerCity": "` +  res.OwnerCity + `", `+ 
		`"ownerState": "` + res.OwnerState + `" , `+ 
		`"ownerPostCode": "` + res.OwnerPostCode + `" , `+ 
		`"ownerCountry": "` +  res.OwnerCountry + `" , `+ 
		`"vesselClass": "` +  res.VesselClass + `" , `+
		`"berthBookingStatus": "` +  res.BerthBookingStatus + `" `+ 
		`}`
	err = stub.PutState(vesselID, []byte(vesselDetails))									//store Vessel with id as key
	if err != nil {
		return nil, err
	}
	return nil, nil
}
