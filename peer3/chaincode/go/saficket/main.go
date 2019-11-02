package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Ticket struct {
	TicketNo     string `json:"ticket_no"`
	AttendeeId   string `json:"attendee_id"`
	EventId      string `json:"event_id"`
	Venue        string `json:venue`
	EventDate    string `json:event_date` // 2019-10-22
	EventTime    string `json:event_time` // 19:00
	TicketIssuer string `json:ticket_issuer`
}

/*
 * The Init method is called when the Smart Contract "fabticket" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabticket"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createNewTicket" { // args = ["ticket_no", "attendee_id", "event_id", "venue", "event_date", "event_time", "ticket_issuer"]
		return s.createNewTicket(APIstub, args)
	} else if function == "queryUserTickets" { // args = ["attendee_id"]
		return s.queryUserTickets(APIstub, args)
	} else if function == "deleteTicket" { // args = ["ticket_no"]
		return s.deleteTicket(APIstub, args)
	} else if function == "queryOneTicket" { // args = ["ticket_no"]
		return s.queryOneTicket(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tickets := []Ticket{
		Ticket{TicketNo: "EVN001", AttendeeId: "owen1994", EventId: "fk948fsld2", Venue: "coex_conference_room", EventDate: "2019-10-22", EventTime: "19:00", TicketIssuer: "interpark"},
		Ticket{TicketNo: "CON222", AttendeeId: "chris88", EventId: "fk94kh3rsw", Venue: "sejong_art_hall", EventDate: "2019-10-24", EventTime: "13:00", TicketIssuer: "auction"},
	}

	i := 0
	for i < len(tickets) {
		ticketAsBytes, _ := json.Marshal(tickets[i])
		APIstub.PutState(tickets[i].TicketNo, ticketAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createNewTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	var ticket = Ticket{TicketNo: args[0], AttendeeId: args[1], EventId: args[2], Venue: args[3], EventDate: args[4], EventTime: args[5], TicketIssuer: args[6]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)
	return shim.Success(ticketAsBytes)
}

func (s *SmartContract) queryUserTickets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	startKey := "000000"
	endKey := "ZZZZZZ"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var raw map[string]interface{}
		json.Unmarshal(queryResponse.Value, &raw)
		id := fmt.Sprintf("%v", raw["attendee_id"])
		ticketNo := fmt.Sprintf("%v", raw["ticket_no"])
		if id == args[0] {
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(", ")
			}
			buffer.WriteString(ticketNo)
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryOneTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := "000000"
	endKey := "ZZZZZZ"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if queryResponse.Key == args[0] {
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			break
		}
	}

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) deleteTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	APIstub.DelState(args[0])
	return shim.Success(nil)
}
// The main function is only relevant in unit test mode. Only included here for completeness.

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
