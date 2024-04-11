package models

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//)
//
//var _ fmt.Stringer = (*OutputRepository)(nil)
//
//type OutputRepository struct {
//	Id                   string    `json:"id"`
//	Name                 string    `json:"name"`
//	Url                  string    `json:"url"`
//	CreatedDateTimestamp time.Time `json:"created_date_timestamp"`
//	Source               string    `json:"source"`
//	LastPush             time.Time `json:"last_push"`
//	UpdatedDateTimestamp time.Time `json:"updated_date_timestamp"`
//	ScanId               int       `json:"scan_id"`
//	HighestSeverity      string    `json:"highest_severity"`
//	TotalFindings        int       `json:"total_findings"`
//	ScanDateTimestamp    time.Time `json:"scan_date_timestamp"`
//	Size                 int       `json:"size"`
//	ConnectedImageId     string    `json:"connected_image_id"`
//}
//
//func (or *OutputRepository) String() string {
//	name := "OutputRepository"
//	res, err := json.Marshal(or)
//	if err != nil {
//		newErr := fmt.Errorf("%s json.Marshal error: %w", name, err)
//		panic(newErr.Error())
//	}
//
//	return string(res)
//}
