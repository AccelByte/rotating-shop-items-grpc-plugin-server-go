// Code generated by go-swagger; DO NOT EDIT.

package service_plugin_config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DeleteSectionPluginConfigReader is a Reader for the DeleteSectionPluginConfig structure.
type DeleteSectionPluginConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteSectionPluginConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteSectionPluginConfigNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDeleteSectionPluginConfigNoContent creates a DeleteSectionPluginConfigNoContent with default headers values
func NewDeleteSectionPluginConfigNoContent() *DeleteSectionPluginConfigNoContent {
	return &DeleteSectionPluginConfigNoContent{}
}

/*
DeleteSectionPluginConfigNoContent describes a response with status code 204, with default header values.

Delete successfully
*/
type DeleteSectionPluginConfigNoContent struct {
}

// IsSuccess returns true when this delete section plugin config no content response has a 2xx status code
func (o *DeleteSectionPluginConfigNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete section plugin config no content response has a 3xx status code
func (o *DeleteSectionPluginConfigNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete section plugin config no content response has a 4xx status code
func (o *DeleteSectionPluginConfigNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete section plugin config no content response has a 5xx status code
func (o *DeleteSectionPluginConfigNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete section plugin config no content response a status code equal to that given
func (o *DeleteSectionPluginConfigNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete section plugin config no content response
func (o *DeleteSectionPluginConfigNoContent) Code() int {
	return 204
}

func (o *DeleteSectionPluginConfigNoContent) Error() string {
	return fmt.Sprintf("[DELETE /admin/namespaces/{namespace}/catalog/plugins/section][%d] deleteSectionPluginConfigNoContent ", 204)
}

func (o *DeleteSectionPluginConfigNoContent) String() string {
	return fmt.Sprintf("[DELETE /admin/namespaces/{namespace}/catalog/plugins/section][%d] deleteSectionPluginConfigNoContent ", 204)
}

func (o *DeleteSectionPluginConfigNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
