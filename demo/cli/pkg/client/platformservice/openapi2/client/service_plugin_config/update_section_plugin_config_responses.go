// Code generated by go-swagger; DO NOT EDIT.

package service_plugin_config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"accelbyte.net/rotating-shop-items-cli/pkg/client/platformservice/openapi2/models"
)

// UpdateSectionPluginConfigReader is a Reader for the UpdateSectionPluginConfig structure.
type UpdateSectionPluginConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateSectionPluginConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateSectionPluginConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 422:
		result := NewUpdateSectionPluginConfigUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateSectionPluginConfigOK creates a UpdateSectionPluginConfigOK with default headers values
func NewUpdateSectionPluginConfigOK() *UpdateSectionPluginConfigOK {
	return &UpdateSectionPluginConfigOK{}
}

/*
UpdateSectionPluginConfigOK describes a response with status code 200, with default header values.

successful operation
*/
type UpdateSectionPluginConfigOK struct {
	Payload *models.SectionPluginConfigInfo
}

// IsSuccess returns true when this update section plugin config o k response has a 2xx status code
func (o *UpdateSectionPluginConfigOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update section plugin config o k response has a 3xx status code
func (o *UpdateSectionPluginConfigOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update section plugin config o k response has a 4xx status code
func (o *UpdateSectionPluginConfigOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update section plugin config o k response has a 5xx status code
func (o *UpdateSectionPluginConfigOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update section plugin config o k response a status code equal to that given
func (o *UpdateSectionPluginConfigOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update section plugin config o k response
func (o *UpdateSectionPluginConfigOK) Code() int {
	return 200
}

func (o *UpdateSectionPluginConfigOK) Error() string {
	return fmt.Sprintf("[PUT /admin/namespaces/{namespace}/catalog/plugins/section][%d] updateSectionPluginConfigOK  %+v", 200, o.Payload)
}

func (o *UpdateSectionPluginConfigOK) String() string {
	return fmt.Sprintf("[PUT /admin/namespaces/{namespace}/catalog/plugins/section][%d] updateSectionPluginConfigOK  %+v", 200, o.Payload)
}

func (o *UpdateSectionPluginConfigOK) GetPayload() *models.SectionPluginConfigInfo {
	return o.Payload
}

func (o *UpdateSectionPluginConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SectionPluginConfigInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateSectionPluginConfigUnprocessableEntity creates a UpdateSectionPluginConfigUnprocessableEntity with default headers values
func NewUpdateSectionPluginConfigUnprocessableEntity() *UpdateSectionPluginConfigUnprocessableEntity {
	return &UpdateSectionPluginConfigUnprocessableEntity{}
}

/*
UpdateSectionPluginConfigUnprocessableEntity describes a response with status code 422, with default header values.

<table><tr><td>ErrorCode</td><td>ErrorMessage</td></tr><tr><td>20002</td><td>validation error</td></tr></table>
*/
type UpdateSectionPluginConfigUnprocessableEntity struct {
	Payload *models.ValidationErrorEntity
}

// IsSuccess returns true when this update section plugin config unprocessable entity response has a 2xx status code
func (o *UpdateSectionPluginConfigUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update section plugin config unprocessable entity response has a 3xx status code
func (o *UpdateSectionPluginConfigUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update section plugin config unprocessable entity response has a 4xx status code
func (o *UpdateSectionPluginConfigUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this update section plugin config unprocessable entity response has a 5xx status code
func (o *UpdateSectionPluginConfigUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this update section plugin config unprocessable entity response a status code equal to that given
func (o *UpdateSectionPluginConfigUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the update section plugin config unprocessable entity response
func (o *UpdateSectionPluginConfigUnprocessableEntity) Code() int {
	return 422
}

func (o *UpdateSectionPluginConfigUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PUT /admin/namespaces/{namespace}/catalog/plugins/section][%d] updateSectionPluginConfigUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdateSectionPluginConfigUnprocessableEntity) String() string {
	return fmt.Sprintf("[PUT /admin/namespaces/{namespace}/catalog/plugins/section][%d] updateSectionPluginConfigUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdateSectionPluginConfigUnprocessableEntity) GetPayload() *models.ValidationErrorEntity {
	return o.Payload
}

func (o *UpdateSectionPluginConfigUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationErrorEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
