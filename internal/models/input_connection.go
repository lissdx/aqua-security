package models

//
//import (
//	"encoding/json"
//	"fmt"
//)
//
//var _ fmt.Stringer = (*InputConnection)(nil)
//
//type InputConnection struct {
//	RepositoryId string `json:"repository_id"`
//	ImageId      string `json:"image_id"`
//}
//
//func (ii *InputConnection) String() string {
//	name := "InputConnection"
//	res, err := json.Marshal(ii)
//	if err != nil {
//		newErr := fmt.Errorf("%s json.Marshal error: %w", name, err)
//		panic(newErr.Error())
//	}
//
//	return string(res)
//}
