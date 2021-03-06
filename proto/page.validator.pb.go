// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: page.proto

package cloudwalker

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/golang/protobuf/ptypes/empty"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Page) Validate() error {
	if this.PageName == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("PageName", fmt.Errorf(`Page Name cannot be empty.`))
	}
	for _, item := range this.Row {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Row", err)
			}
		}
	}
	for _, item := range this.Carousel {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Carousel", err)
			}
		}
	}
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	if this.UpdatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
		}
	}
	return nil
}

var _regex_Carousel_Index = regexp.MustCompile(`^[0-100]*$`)

func (this *Carousel) Validate() error {
	if this.Target == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Target", fmt.Errorf(`Deeplink cannot be empty.`))
	}
	if this.Package == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Package", fmt.Errorf(`App Package Name cannot be empty.`))
	}
	if this.ImageUrl == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("ImageUrl", fmt.Errorf(`Image Url cannot be empty.`))
	}
	if this.Title == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Title", fmt.Errorf(`Carousel Title cannot be empty.`))
	}
	if !(this.Index > -1) {
		return github_com_mwitkow_go_proto_validators.FieldError("Index", fmt.Errorf(`Carousel Index must be a digit and >= 0 and < 100.`))
	}
	if !(this.Index < 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("Index", fmt.Errorf(`Carousel Index must be a digit and >= 0 and < 100.`))
	}
	return nil
}

var _regex_Row_RowIndex = regexp.MustCompile(`^[0-100]*$`)
var _regex_Row_RowSort = regexp.MustCompile(`^[-1-2]*$`)

func (this *Row) Validate() error {
	if _, ok := RowLayout_name[int32(this.RowLayout)]; !ok {
		return github_com_mwitkow_go_proto_validators.FieldError("RowLayout", fmt.Errorf(`Row Layout must be set. Landscape = 0, Portrait = 1, Square = 2, Circle = 3.`))
	}
	if this.RowName == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("RowName", fmt.Errorf(`Row Name cannot be empty.`))
	}
	if !(this.RowIndex > -1) {
		return github_com_mwitkow_go_proto_validators.FieldError("RowIndex", fmt.Errorf(`Row Index must be a digit and in the range of 0 - 100 .`))
	}
	if !(this.RowIndex < 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("RowIndex", fmt.Errorf(`Row Index must be a digit and in the range of 0 - 100 .`))
	}
	// Validation of proto3 map<> fields is unsupported.
	// Validation of proto3 map<> fields is unsupported.
	if _, ok := RowType_name[int32(this.RowType)]; !ok {
		return github_com_mwitkow_go_proto_validators.FieldError("RowType", fmt.Errorf(`Row Type must be set. Editorial = 0 , Recommendation_CB  = 1, Dynamic  = 2, Recommendation_CF  = 3, Web = 4`))
	}
	return nil
}
func (this *RowFilterValue) Validate() error {
	return nil
}
func (this *GetPageReq) Validate() error {
	if this.PageId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("PageId", fmt.Errorf(`Page id cannot be empty.`))
	}
	return nil
}
func (this *DeletePageReq) Validate() error {
	if this.PageId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("PageId", fmt.Errorf(`Page id cannot be empty.`))
	}
	return nil
}
func (this *ResultPage) Validate() error {
	for _, item := range this.Carousels {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Carousels", err)
			}
		}
	}
	for _, item := range this.Rows {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Rows", err)
			}
		}
	}
	return nil
}
func (this *ResultRow) Validate() error {
	for _, item := range this.Tiles {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Tiles", err)
			}
		}
	}
	return nil
}
func (this *Content) Validate() error {
	for _, item := range this.Play {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Play", err)
			}
		}
	}
	return nil
}
func (this *Play) Validate() error {
	return nil
}
func (this *DropDownReq) Validate() error {
	if this.Field == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Field", fmt.Errorf(`DB filed must exists to get suggested values of it.`))
	}
	return nil
}
func (this *DropDownResp) Validate() error {
	return nil
}
