package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("WmOrderForcastChaincode")

type WmOrderForcastChaincode struct {
}

type Forecast struct {
	OldNbr                      string    `json:"old_nbr"`
	Upc                         string    `json:"upc"`
	OrderDeptNbr                string    `json:"order_dept_nbr"`
	PrimaryDesc                 string    `json:"primary_desc"`
	VnpkQty                     int       `json:"vnpk_qty"`
	WhpkQty                     int       `json:"whpk_qty"`
	VendorName                  string    `json:"vendor_name"`
	StoreNbr                    string    `json:"store_nbr"`
	SourceDcNbr                 string    `json:"source_dc_nbr"`
	OrderEach                   int       `json:"order_each"`
	OrderWhpk                   int       `json:"order_whpk"`
	DateThisQtyPlannedToArrive  time.Time `json:"date_this_qty_planned_to_arrive"`
	DateThisOrderShouldBePlaced time.Time `json:"date_this_order_should_be_placed"`
	Status                      string    `json:"status"`
	Comments                    string    `json:"comments"`
}

func (t *WmOrderForcastChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")
	return shim.Success(nil)
}
func (t *WmOrderForcastChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "addForecast" {
		return t.addForecast(stub, args)
	} else if function == "ackForecast" {
		return t.ackForecast(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}
	return pb.Response{Status: 403, Message: "unknown function name"}
}

func (t *WmOrderForcastChaincode) addForecast(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//key is sscc

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error("cannot get creator")
	}

	user, org := getCreator(creatorBytes)
	logger.Debug("User:" + user)
	if org == "" {
		logger.Debug("Org is null")
		return shim.Error("cannot get Org")
	} else if org == "walmart" {

		if len(args) != 13 {
			return pb.Response{Status: 403, Message: "incorrect number of arguments"}
		}

		vnpkQty, _ := strconv.Atoi(args[4])
		whpkQty, _ := strconv.Atoi(args[5])
		orderEach, _ := strconv.Atoi(args[9])
		orderWhpk, _ := strconv.Atoi(args[10])

		location, err := time.LoadLocation("America/Chicago")
		if err != nil {
			fmt.Println(err)
		}
		// Date Format: 'yyyy-mm-dd'
		dateFormat := "2006-01-02"
		dateThisQtyPlannedToArrive, _ := time.ParseInLocation(dateFormat, args[11], location)
		dateThisOrderShouldBePlaced, _ := time.ParseInLocation(dateFormat, args[12], location)

		forecastObj := &Forecast{
			OldNbr:                      args[0],
			Upc:                         args[1],
			OrderDeptNbr:                args[2],
			PrimaryDesc:                 args[3],
			VnpkQty:                     vnpkQty,
			WhpkQty:                     whpkQty,
			VendorName:                  args[6],
			StoreNbr:                    args[7],
			SourceDcNbr:                 args[8],
			OrderEach:                   orderEach,
			OrderWhpk:                   orderWhpk,
			DateThisQtyPlannedToArrive:  dateThisQtyPlannedToArrive,
			DateThisOrderShouldBePlaced: dateThisOrderShouldBePlaced}

		jsonForecastObj, err := json.Marshal(forecastObj)
		if err != nil {
			return shim.Error("Cannot create Json Object")
		}
		logger.Debug("Json Obj: " + string(jsonForecastObj))

		// key is the combination of UPC and StoreNbr ... Need to confirm with Business
		key := args[1] + "-" + args[7]

		err = stub.PutState(key, jsonForecastObj)
		if err != nil {
			return shim.Error("cannot put state")
		}

		logger.Debug("Forecast Created")

	}
	return shim.Success(nil)
}

func (t *WmOrderForcastChaincode) ackForecast(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error("cannot get creator")
	}

	user, org := getCreator(creatorBytes)
	logger.Debug("User:" + user)
	if org == "" {
		logger.Debug("Org is null")
		return shim.Error("cannot get Org")
	} else if org == "unilever" {
		if len(args) < 3 {
			return pb.Response{Status: 403, Message: "incorrect number of arguments"}
		}

		key := args[0] + "-" + args[1]

		forecastBytes, err := stub.GetState(key)
		if err != nil {
			return shim.Error("cannot get state")
		} else if forecastBytes == nil {
			return shim.Error("Cannot get shippment object")
		}

		var forecastObj Forecast
		errUnmarshal := json.Unmarshal([]byte(forecastBytes), &forecastObj)
		if errUnmarshal != nil {
			return shim.Error("Cannot unmarshal Insurance Object")
		}

		logger.Debug("Shipment Object: " + string(forecastBytes))

		forecastObj.Status = args[2]

		forecastObj.Comments = args[3]

		jsonForecastObj, err := json.Marshal(forecastObj)
		if err != nil {
			return shim.Error("Cannot create Json Object")
		}
		logger.Debug("Json Obj: " + string(jsonForecastObj))

		err = stub.PutState(key, jsonForecastObj)
		if err != nil {
			return shim.Error("cannot put state")
		}

		logger.Debug("Forecast Updated with Status")

	}
	return shim.Success(nil)
}

func (t *WmOrderForcastChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if args[0] == "health" {
		logger.Info("Health status Ok")
		return shim.Success(nil)
	}

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error("cannot get creator")
	}

	user, org := getCreator(creatorBytes)
	logger.Debug("User:" + user)
	if org == "" {
		logger.Debug("Org is null")
		return shim.Error("cannot get Org")
	} else {
		if len(args) != 2 {
			return pb.Response{Status: 403, Message: "incorrect number of arguments"}
		}

		key := args[0] + "-" + args[1]

		forcastBytes, err := stub.GetState(key)
		if err != nil {
			return shim.Error("cannot get state")
		} else if forcastBytes == nil {
			return shim.Error("Cannot get shippment object")
		}

		logger.Debug("Forecast Object: " + string(forcastBytes))

		return shim.Success(forcastBytes)
	}

}

var getCreator = func(certificate []byte) (string, string) {
	data := certificate[strings.Index(string(certificate), "-----") : strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization := cert.Issuer.Organization[0]
	commonName := cert.Subject.CommonName
	logger.Debug("commonName: " + commonName + ", organization: " + organization)

	organizationShort := strings.Split(organization, ".")[0]

	return commonName, organizationShort
}

func main() {
	err := shim.Start(new(WmOrderForcastChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
