package models

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//)
//
//var _ fmt.Stringer = (*InputResources)(nil)
//
//type InputResources struct {
//	Id                   string    `json:"id"`
//	Name                 string    `json:"name"`
//	Url                  string    `json:"url"`
//	CreatedDateTimestamp int64     `json:"created_date_timestamp"`
//	Type                 string    `json:"type"`
//	Source               string    `json:"source"`
//	NumberOfLayers       int       `json:"number_of_layers"`
//	Architecture         string    `json:"architecture"`
//	LastPush             time.Time `json:"last_push"`
//	Size                 int       `json:"size"`
//}
//
//func (ir *InputResources) String() string {
//	name := "InputResources"
//	res, err := json.Marshal(ir)
//	if err != nil {
//		newErr := fmt.Errorf("%s json.Marshal error: %w", name, err)
//		panic(newErr.Error())
//	}
//
//	return string(res)
//}
