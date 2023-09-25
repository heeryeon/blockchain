package chaincode




import ( 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"time"
	"strconv"
	"encoding/json"
	"fmt"
	"strings"
)




const (
    YYYYMMDD = "2006-01-02"
)





type SmartContract struct{
	contractapi.Contract
}

type Person struct{					
	PersonId	string			
	Name		string		
	Birth		string		
	PhoneNumber string				
	FluInfo		[]FluData		
}

type FluData struct{
	FluName string
	FluDate string
}




func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface) error {	
	
	// personcount 값 읽어오기
	personcountJSON, err := ctx.GetStub().GetState("personcount")
	if err !=nil {
		 return err
	}
	if personcountJSON != nil{
		return fmt.Errorf("personcount is already set")
	}
	
	//  personcount 0으로 초기화
	err = ctx.GetStub().PutState("personcount", []byte(strconv.Itoa(0)))
	if err !=nil {
		return err
   	}
	return nil
}






func (s *SmartContract)RegisterPerson(ctx contractapi.TransactionContextInterface, name string, birth string, phonenumber string) error  {

	// New 사람 data 등록
	// PersonId : 현재시간 + 휴대폰 번호 뒷 네자리
	// Name : 이 요청을 보낸 사람의 ClientID
	// Birth : 생일
	// PhoneNumber : 휴대폰 번호


	nowtime := time.Now()
	unixtime := nowtime.Unix()


	
	
	splitedPN := strings.Split(phonenumber,"-")
	person := Person{
		PersonId: strconv.Itoa(int(unixtime)) + splitedPN[2],
		Name: name,
		Birth: birth,
		PhoneNumber: phonenumber,
	}	
	
	
	personJSON, err := json.Marshal(person)
	if err !=nil {
		return err
   	}

	// ctx를 활용하여 Client ID를 알아내기
	clientID, err := ctx.GetClientIdentity().GetID()
	if err !=nil {
		return fmt.Errorf("failed to get clientID : %v", err)
   	}

	
	err = ctx.GetStub().PutState(clientID, personJSON )
	if err !=nil {
		return err
   	}

	// personcount 값 읽어오기
	personcountJSON, err := ctx.GetStub().GetState("personcount")
	if err !=nil {
		 return err
	}
	
	personIdINT,_ := strconv.Atoi(string(personcountJSON))
	personIdINT += 1 
	err = ctx.GetStub().PutState("personcount", []byte(strconv.Itoa(personIdINT)))
	if err !=nil {
		return err
   	}

	return nil
}




func (s *SmartContract) GetFluDateByPersonId(ctx contractapi.TransactionContextInterface) ([]FluData, error) {

	// ctx를 활용하여 Client ID를 알아내기
	clientID, err := ctx.GetClientIdentity().GetID()
	if err !=nil {
		return nil, fmt.Errorf("failed to get clientID : %v", err)
   	}

	// queryString := fmt.Sprintf( `{"selector":{"name":"%s"}}`,clientID)

	// queryResult := ctx.GetStub().GetQueryResult(queryString)

	countBytes, err := ctx.GetStub().GetState(clientID)
	if err !=nil {
		return nil, err
   	}

	var person Person

	err = json.Unmarshal(countBytes, &person)
	if err !=nil {
		return nil,err
	}

	return person.FluInfo, nil
}


func (s *SmartContract)UpdateFluDate(ctx contractapi.TransactionContextInterface, clientID string, fluName string) error {

	countBytes, err := ctx.GetStub().GetState(clientID)
	if err !=nil {
		return err
   	}

	var person Person

	err = json.Unmarshal(countBytes, &person)
	if err !=nil {
		return err
	}

	
	nowtime := time.Now()
	nowtime_yyyymmdd := nowtime.UTC().Format(YYYYMMDD)
	

	// person.FluDate=string(nowtime_yyyymmdd)

	var fludata FluData
	fludata.FluName = fluName
	fludata.FluDate = string(nowtime_yyyymmdd)

	person.FluInfo = append(person.FluInfo, fludata)

	personJSON, err := json.Marshal(person)
	if err !=nil {
		return err
   	}
	
	err = ctx.GetStub().PutState(clientID, personJSON )
	if err !=nil {
		return err
   	}
	return nil
}