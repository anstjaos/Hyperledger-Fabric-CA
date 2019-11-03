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
	TicketCode   string `json:"ticket_code"`
	AttendeeId   string `json:"attendee_id"`
	EventName    string `json:"event_name"`
	Venue        string `json:"venue"`
	EventDate    string `json:"event_date"` // 2019-10-22
	EventTime    string `json:"event_time"` // 19:00
	TicketIssuer string `json:"ticket_issuer"`
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
	} else if function == "createNewTicket" { // args = ["ticket_code", "attendee_id", "event_id", "venue", "event_date", "event_time", "ticket_issuer"], 티켓을 새로 생성, 결과타입 boolean
		return s.createNewTicket(APIstub, args)
	} else if function == "queryUserTickets" { // args = ["attendee_id"], 특정 회원이 가지고 있는 모든 티켓을 조회, 결과 타입 StringArray
		return s.queryUserTickets(APIstub, args)
	} else if function == "deleteTicket" { // args = ["ticket_code"], 특정 티켓을 삭제, 결과 타입 boolean
		return s.deleteTicket(APIstub, args)
	} else if function == "queryOneTicket" { // args = ["ticket_code"] 티켓 하나 정보 조회
		return s.queryOneTicket(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	tickets := []Ticket{
		Ticket{TicketCode: "EVN001", AttendeeId: "owen1994", EventName: "IU Concert", Venue: "coex_conference_room", EventDate: "2019-10-22", EventTime: "19:00", TicketIssuer: "interpark"},
		Ticket{TicketCode: "CON222", AttendeeId: "chris88", EventName: "Mammamia", Venue: "sejong_art_hall", EventDate: "2019-10-24", EventTime: "13:00", TicketIssuer: "auction"},
	}

	i := 0
	for i < len(tickets) {
		ticketAsBytes, _ := json.Marshal(tickets[i])
		APIstub.PutState(tickets[i].TicketCode, ticketAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createNewTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response { // 티켓 하나 생성

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	startKey := "000000"
	endKey := "ZZZZZZ"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if queryResponse.Key == args[0] {
			return shim.Success([]byte("false"))
		}
	}
	var ticket = Ticket{TicketCode: args[0], AttendeeId: args[1], EventName: args[2], Venue: args[3], EventDate: args[4], EventTime: args[5], TicketIssuer: args[6]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)
	return shim.Success([]byte("true"))
}

func (s *SmartContract) queryUserTickets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response { // 특정 사용자가 가진 티켓 조회

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
		ticketCode := fmt.Sprintf("%v", raw["ticket_code"])
		eventName := fmt.Sprintf("%v", raw["event_name"])
		eventDate := fmt.Sprintf("%v", raw["event_date"])
		if id == args[0] {
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(", ")
			}
			buffer.WriteString("{ \"TicketCode\" : \"")
			buffer.WriteString(ticketCode)
			buffer.WriteString("\", \"EventName\" : \"")
			buffer.WriteString(eventName)
			buffer.WriteString("\", \"EventDate\" : \"")
			buffer.WriteString(eventDate)
			buffer.WriteString("\" }")
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryOneTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response { // 티켓 하나 정보 조회

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

func (s *SmartContract) deleteTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response { // 티켓 삭제

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

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if queryResponse.Key == args[0] {
			APIstub.DelState(args[0])
			return shim.Success([]byte("true"))
		}
	}
	return shim.Success([]byte("false"))
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
