package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// sdk代码中，有且仅有一个main()方法
func main() {
	// main()方法中，下面的代码为必须代码，不建议修改main()方法当中的代码
	// 其中，TestContract为用户实现合约的具体名称
	err := sandbox.Start(new(FactContract))
	if err != nil {
		log.Fatal(err)
	}
}

//合约必要代码：
// 合约结构体，合约名称需要写入main()方法当中
type FactContract struct {
}

// -------------------合约必须实现下面两个方法：-----------------------
// InitContract() protogo.Response
// UpgradeContract() protogo.Response
// InvokeContract(method string) protogo.Response

// 用于合约的部署
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// 用于合约的升级
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// 用于合约的调用
// @param method: 交易请求调用的方法
// @return: 	合约返回结果，包括Success和Error
func (f *FactContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return f.save()
	case "findByFileHash":
		return f.findByFileHash()
	default:
		return sdk.Error("invalid method")
	}
}
//------------------存证合约实现--------------------------
//见contract_save.go