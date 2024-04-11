package models

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//)
//
//var _ fmt.Stringer = (*InputScan)(nil)
//
//type InputScan struct {
//	ScanId            int       `json:"scan_id"`
//	ResourceId        string    `json:"resource_id"`
//	ResourceType      string    `json:"resource_type"`
//	HighestSeverity   string    `json:"highest_severity"`
//	TotalFindings     int       `json:"total_findings"`
//	ScanDateTimestamp time.Time `json:"scan_date_timestamp"`
//}
//
//func (is *InputScan) String() string {
//	name := "InputScan"
//	res, err := json.Marshal(is)
//	if err != nil {
//		newErr := fmt.Errorf("%s json.Marshal error: %w", name, err)
//		panic(newErr.Error())
//	}
//
//	return string(res)
//}
