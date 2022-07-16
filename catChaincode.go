package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CatChaincode struct {
	contractapi.Contract
}

type Cat struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Owner string `json:"owner"`
}

func (cc *CatChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cats := []Cat{
		{ID: "1", Name: "米米", Color: "black", Owner: "郑雅菱"},
		{ID: "2", Name: "小黄", Color: "green", Owner: "水成渊"},
		{ID: "3", Name: "花花", Color: "red", Owner: "零零自"},
		{ID: "4", Name: "艾灸", Color: "blue", Owner: "雅菱二"},
	}

	for _, cat := range cats {
		catJson, err := json.Marshal(cat)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(cat.ID, catJson)
		if err != nil {
			return fmt.Errorf("Failed to create cat %v.", cat.ID)
		}
	}
	return nil
}

func (cc *CatChaincode) Exist(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	catJson, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read %s cat from world state with error %v.", ID, err)
	}
	return catJson != nil, nil
}

func (cc *CatChaincode) CreateCat(ctx contractapi.TransactionContextInterface, ID, name, color, owner string) error {
	exist, err := cc.Exist(ctx, ID)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("cat %v already exists", ID)
	}

	cat := Cat{
		ID:    ID,
		Name:  name,
		Color: color,
		Owner: owner,
	}
	catJson, _ := json.Marshal(cat)
	return ctx.GetStub().PutState(ID, catJson)
}

func (cc *CatChaincode) ReadCat(ctx contractapi.TransactionContextInterface, ID string) (*Cat, error) {
	catJson, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state with error %v.", err)
	}
	if catJson == nil {
		return nil, fmt.Errorf("the asset %v cat does not exist.", ID)
	}

	var cat Cat
	err = json.Unmarshal(catJson, &cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (cc *CatChaincode) UpdateCat(ctx contractapi.TransactionContextInterface, id, name, color, owner string) error {
	exist, err := cc.Exist(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("the cat %v does not exist.", id)
	}

	cat := Cat{
		ID:    id,
		Name:  name,
		Color: color,
		Owner: owner,
	}

	catJson, err := json.Marshal(cat)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, catJson)
}

func (cc *CatChaincode) DeleteCat(ctx contractapi.TransactionContextInterface, id string) error {
	exist, err := cc.Exist(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("the cat %v does not exist.", id)
	}
	return ctx.GetStub().DelState(id)
}

func (cc *CatChaincode) TransferCat(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	cat, err := cc.ReadCat(ctx, id)
	if err != nil {
		return err
	}
	cat.Owner = newOwner
	catJson, _ := json.Marshal(cat)
	return ctx.GetStub().PutState(cat.ID, catJson)
}

func (cc *CatChaincode) GetAllCats(ctx contractapi.TransactionContextInterface) ([]*Cat, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cats []*Cat
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var cat Cat
		err = json.Unmarshal(queryResponse.Value, &cat)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &cat)
	}
	return cats, nil
}

func main() {
	mc, err := contractapi.NewChaincode(&CatChaincode{})
	if err != nil {
		log.Panicf("error creating new chaincode, %v.", err)
	}
	if err = mc.Start(); err != nil {
		log.Panicf("error starting my chaincode, %v.", err)
	}
}
