package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

//结构体
type PointsTransferChaincode struct {
}

//文件结构体
type fileChaincode struct {
	ObjectType     string
	FileID         string
	FileRecordHash string
	FileName       string
	FromUser       string
	ToUser         string
	TxID           string
	IsDeleteRecode string
	// FileTimecode   string
	FileTimecode   int64
}

// // TransactionDetail获取了交易的具体信息
// type TransactionDetail struct {
// 	TransactionId string
// 	CreateTime    string
// 	Args          []string
// }

// Init => 链码初始化接口
func (t *PointsTransferChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke => 链码调用接口
func (t *PointsTransferChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createFileChaincode" {
		return t.createFileChaincode(stub, args)
	} else if function == "readFileChaincode" {
		return t.readFileChaincode(stub, args)
	} else if function == "deleteFileChaincode" {
		return t.deleteFileChaincode(stub, args)
	} else if function == "transferFileChaincode" {
		return t.transferFileChaincode(stub, args)
	} else if function == "queryFileChaincodeByOwner" {
		return t.queryFileChaincodeByOwner(stub, args)
	} else if function == "getHistoryFromFileChaincode" {
		return t.getHistoryFromFileChaincode(stub, args)
	} else if function == "queryFileChaincodeByTxID" {
		return t.queryFileChaincodeByTxID(stub, args)
	} else if function == "queryFileChaincodeByFromUser" {
		return t.queryFileChaincodeByFromUser(stub, args)
	} else if function == "queryFileChaincodeByTime" {
		return t.queryFileChaincodeByTime(stub, args)
	}

	return shim.Error("无效交易方法，仅支持：transfer|query")
}

//创建一个文件信息
func (t *PointsTransferChaincode) createFileChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	txid := stub.GetTxID()
	fileID := args[0]
	// if err != nil {
	// 	return shim.Error("fileID 必须是数字字符串")
	// }
	//判断该文件是否存在
	fileRecodeBytes, err := stub.GetState(fileID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fileRecodeBytes != nil {
		return shim.Error("FileRecode已存在！")
	}
	fileRecodeHash := args[1]
	fileName := args[2]
	fromUser := args[3]
	toUser := args[4]
	isDeleteRecode := "false"
	objectType := "fileChaincode"
	// fileTimecode := time.Now().Format("2006-01-02 15:04:05")
	fileTimecode := time.Now().Unix()
	fileChaincode := &fileChaincode{objectType, fileID, fileRecodeHash, fileName, fromUser, toUser, txid, isDeleteRecode, fileTimecode}
	fileChaincodeJSONAsBytesgo, err := json.Marshal(fileChaincode)
	//在账本中添加键值对
	err = stub.PutState(fileID, fileChaincodeJSONAsBytesgo)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//读取一个文件信息
func (t *PointsTransferChaincode) readFileChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fileID := args[0]
	fileChaincodeAsBytes, err := stub.GetState(fileID)
	if err != nil {
		return shim.Error(err.Error())
	} else if fileChaincodeAsBytes == nil {
		return shim.Error("fileChaincode信息不存在")
	}

	return shim.Success(fileChaincodeAsBytes)
}

//删除一个文件信息
func (t *PointsTransferChaincode) deleteFileChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	txid := stub.GetTxID()
	//判断文件记录信息是否存在
	fileID := args[0]
	fileChaincodeAsBytes, err := stub.GetState(fileID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fileChaincodeAsBytes != nil {
		//删除
		err = stub.DelState(fileID)
		if err != nil {
			return shim.Error(string(err.Error()))
		}

		//创建删除文件记录信息
		fileRecodeHash := args[1]
		fileName := args[2]
		fromUser := args[3]
		toUser := args[4]
		isDeleteRecode := "true"
		// fileTimecode,_ := shim.GetTxTimestamp()
		// fileTimecode := time.Now().Format("2006-01-02 15:04:05")
		fileTimecode := time.Now().Unix()
		objectType := "fileChaincode"
		fileChaincode := &fileChaincode{objectType, fileID, fileRecodeHash, fileName, fromUser, toUser, txid, isDeleteRecode, fileTimecode}
		fileChaincodeJSONAsBytesgo, err := json.Marshal(fileChaincode)
		//在账本中添加键值对
		err = stub.PutState(txid, fileChaincodeJSONAsBytesgo)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

//转移文件所有者
func (t *PointsTransferChaincode) transferFileChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fileID := args[0]
	fileRecodeHash := args[1]
	fromUser := args[2]
	toUser := args[3]
	txID := stub.GetTxID()
	// fileTimecode,_ := shim.GetTxTimestamp()
	// fileTimecode := time.Now().Format("2006-01-02 15:04:05")
	fileTimecode := time.Now().Unix()
	//判断file是否存在
	fileChaincodeAsBytes, err := stub.GetState(fileID)
	if err != nil {
		return shim.Error(err.Error())
	} else if fileChaincodeAsBytes == nil {
		return shim.Error("fileChaincode信息不存在")
	}
	fileChaincodeInfo := fileChaincode{}
	err = json.Unmarshal(fileChaincodeAsBytes, &fileChaincodeInfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	fileChaincodeInfo.FileRecordHash = fileRecodeHash
	fileChaincodeInfo.FromUser = fromUser
	fileChaincodeInfo.ToUser = toUser
	fileChaincodeInfo.TxID = txID
	fileChaincodeInfo.FileTimecode = fileTimecode
	fileChaincodeAsBytes, err = json.Marshal(fileChaincodeInfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(fileID, fileChaincodeAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//查询指定拥有者拥有的所有文件
func (t *PointsTransferChaincode) queryFileChaincodeByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	toUser := args[0]
	queryStr := fmt.Sprintf("{\"selector\": {\"ToUser\":\"%s\"}}", toUser)
	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(",\"record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//查询指定FromUser的所有文件
func (t *PointsTransferChaincode) queryFileChaincodeByFromUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fromUser := args[0]
	queryStr := fmt.Sprintf("{\"selector\": {\"FromUser\":\"%s\"}}", fromUser)
	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(",\"record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//查询指定TxID文件信息
func (t *PointsTransferChaincode) queryFileChaincodeByTxID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	txID := args[0]
	queryStr := fmt.Sprintf("{\"selector\": {\"TxID\":\"%s\"}}", txID)
	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(",\"record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//查询指定时间区间、指定fromUser、指定toUser、指定isDelete的信息
func (t *PointsTransferChaincode) queryFileChaincodeByTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	timeStart := args[0]
	timeEnd := args[1]
	fromUser := args[2]
	toUser := args[3]
	isDelete := args[4]
	// selectStr := "{\"selector\": {\"FileTimecode\":{\"$gte\":" + timeStart + ",\"$lte\":" + timeEnd + ",\"FromUser\":" + fromUser + ",\"ToUser\":" + toUser + ",\"IsDeleteRecode\":" + isDelete + "}}}"

	selectStr := "{\"selector\": {\"FileTimecode\":{\"$gte\":" + timeStart + ",\"$lte\":" + timeEnd + "}"

	if fromUser != "" {
		selectStr = selectStr + ",\"FromUser\":\"" + fromUser + "\""
	}
	if toUser != "" {
		selectStr = selectStr + ",\"ToUser\":\"" + toUser + "\""
	}
	if isDelete != "" {
		selectStr = selectStr + ",\"IsDeleteRecode\":\"" + isDelete + "\""
	}
	selectStr = selectStr + "}}"
	// queryStr := fmt.Sprintf(selectStr)
	resultIterator, err := stub.GetQueryResult(selectStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(",\"record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")
	// buffer.WriteString(selectStr)

	return shim.Success(buffer.Bytes())
}

//查询一个文件的所有历史记录
func (t *PointsTransferChaincode) getHistoryFromFileChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fileID := args[0]
	resultIterator, err := stub.GetHistoryForKey(fileID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString(queryResponse.TxId)

		buffer.WriteString(",\"Timestamp\":")
		buffer.WriteString(time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String())

		// buffer.WriteString(",\"key\":")
		// buffer.WriteString(string(queryResponse.Key))

		buffer.WriteString(",\"Value\":")
		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString(",\"IsDelete\":")
		buffer.WriteString(strconv.FormatBool(queryResponse.IsDelete))

		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(PointsTransferChaincode))
	if err != nil {
		fmt.Printf("链码启动失败: %s\n", err)
	}
}
