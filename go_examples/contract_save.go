package main

import (
	"encoding/json"
	"fmt" // print
	"log"
	"math/big"
	"strconv" //字符类型转换
	"crypto/rand"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type FactContract struct {
}

// 存证对象
type Fact struct {
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
	Time     int32  `json:"time"`
}

// 新建存证对象
func NewFact(fileHash string, fileName string, time int32) *Fact {
	fact := &Fact{
		FileHash: fileHash,
		FileName: fileName,
		Time:     time,
	}
	return fact
}

// 必须实现的方法1：InitContract() protogo.Response
// 用于合约的部署
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// 必须实现的方法2： UpgradeContract() protogo.Response
// 用于合约的升级
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// 必须实现的方法3： InvokeContract(method string) protogo.Response
// 用于合约的调用
// @param method: 交易请求调用的方法
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return f.save()
	case "findByFileHash":
		return f.findByFileHash()
	case "random":
		return f.random()
	default:
		return sdk.Error("invalid method")
	}
}
// 存证合约实现
//## 使用cmc调用合约
//./cmc client contract user invoke \
//--contract-name=fact \
//--method=save \
//--sdk-conf-path=./testdata/sdk_config.yml \
//--params="{\"file_name\":\"name007\",\"file_hash\":\"ab3456df5799b87c77e7f88\",\"time\":\"6543234\"}" \
//--sync-result=true
//使用sdk调用合约
// 调用合约
// 调用或者查询合约时，method参数请设置为 invoke_contract，此方法会调用合约的InvokeContract方法，再通过参数获得具体方法
//func testUserContractInvoke(client *sdk.ChainClient, method string, withSyncResult bool) (string, error) {
//	curTime := strconv.FormatInt(time.Now().Unix(), 10)
//	fileHash := uuid.GetUUID()
//	kvs := []*common.KeyValuePair{
//		{
//			Key: "method",
//			Value: []byte("save"),
//		},
//		{
//			Key:   "time",
//			Value: []byte(curTime),
//		},
//		{
//			Key:   "file_hash",
//			Value: []byte(fileHash),
//		},
//		{
//			Key:   "file_name",
//			Value: []byte(fmt.Sprintf("file_%s", curTime)),
//		},
//	}
//	err := invokeUserContract(client, factContractName, method, "", kvs, withSyncResult)
//	if err != nil {
//		return "", err
//	}
//	return fileHash, nil
//}


func (f *FactContract) random() protogo.Response {
	// 随机数逻辑
	number, _ := rand.Int(rand.Reader, big.NewInt(10000000000))
	fmt.Println(number)
	number_str := fmt.Sprint(number)
	//number_str := strconv.Itoa(number)
	// 返回结果
	return sdk.Success([]byte(number_str))

}
func (f *FactContract) save() protogo.Response {
	// GetArgs get arg from transaction parameters
	// @return: 参数map
	params := sdk.Instance.GetArgs()
	// 获取参数
	fileHash := string(params["file_hash"])
	fileName := string(params["file_name"])
	timeStr := string(params["time"])
	//string 转换为 int(Ascii to int)
	//由于 string 可能无法转换为int，所以这个函数有两个返回值： 第一个返回值是转换成 int 的值 第二个返回值判断是否转换成功
	time, err := strconv.Atoi(timeStr)
	//如果转换失败
	if err != nil {
		msg := "time is [" + timeStr + "] not int"
		sdk.Instance.Errorf(msg)
		return sdk.Error(msg)
	}
	// 构建结构体 新建存证对象
	fact := NewFact(fileHash, fileName, int32(time))

	// 序列化 json.Marshal将数据结构转换为json字符串
	factBytes, err := json.Marshal(fact)
	if err != nil {
		return sdk.Error(fmt.Sprintf("marshal fact failed, err: %s", err))
	}
	factBytesJson := strconv.Quote(string(factBytes))
	sdk.Instance.Infof("[save] factBytes=" + factBytesJson)
	// 发送事件
	sdk.Instance.EmitEvent("topic_vx", []string{fact.FileHash, fact.FileName})

	// 存储数据
	err = sdk.Instance.PutStateByte("fact_bytes", fact.FileHash, factBytes)
	if err != nil {
		return sdk.Error("fail to save fact bytes")
	}
	// 记录日志
	sdk.Instance.Infof("[save] fileHash=" + fact.FileHash)
	sdk.Instance.Infof("[save] fileName=" + fact.FileName)
	// 返回结果
	return sdk.Success([]byte(fact.FileName + fact.FileHash))
}

func (f *FactContract) findByFileHash() protogo.Response {
	// 获取参数
	fileHash := string(sdk.Instance.GetArgs()["file_hash"])

	// 查询结果
	result, err := sdk.Instance.GetStateByte("fact_bytes", fileHash)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}
	sdk.Instance.Infof("[GetStateByte] result=" + string(result))
	// 反序列化
	var fact Fact
	if err = json.Unmarshal(result, &fact); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal fact failed, err: %s", err))
	}
	// 记录日志
	sdk.Instance.Infof("[find_by_file_hash] fileHash=" + fact.FileHash)
	sdk.Instance.Infof("[find_by_file_hash] fileName=" + fact.FileName)
	// 返回结果
	return sdk.Success(result)
}

func main() {
	err := sandbox.Start(new(FactContract))
	if err != nil {
		log.Fatal(err)
	}
}

