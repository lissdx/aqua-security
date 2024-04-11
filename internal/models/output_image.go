package models

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//)
//
//var _ fmt.Stringer = (*OutputImage)(nil)
//
//type OutputImage struct {
//	Id                    string    `json:"id"`
//	Name                  string    `json:"name"`
//	Url                   string    `json:"url"`
//	CreatedDateTimestamp  time.Time `json:"created_date_timestamp"`
//	NumberOfLayers        int       `json:"number_of_layers"`
//	Architecture          string    `json:"architecture"`
//	UpdatedDateTimestamp  time.Time `json:"updated_date_timestamp"`
//	ScanId                int       `json:"scan_id"`
//	HighestSeverity       string    `json:"highest_severity"`
//	TotalFindings         int       `json:"total_findings"`
//	ScanDateTimestamp     time.Time `json:"scan_date_timestamp"`
//	Source                string    `json:"source"`
//	ConnectedRepositoryId string    `json:"connected_repository_id"`
//}
//
//func (oi *OutputImage) String() string {
//	name := "OutputImage"
//	res, err := json.Marshal(oi)
//	if err != nil {
//		newErr := fmt.Errorf("%s json.Marshal error: %w", name, err)
//		panic(newErr.Error())
//	}
//
//	return string(res)
//}
