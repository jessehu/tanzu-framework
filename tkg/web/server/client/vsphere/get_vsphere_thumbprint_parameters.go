// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetVsphereThumbprintParams creates a new GetVsphereThumbprintParams object
// with the default values initialized.
func NewGetVsphereThumbprintParams() *GetVsphereThumbprintParams {
	var ()
	return &GetVsphereThumbprintParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetVsphereThumbprintParamsWithTimeout creates a new GetVsphereThumbprintParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetVsphereThumbprintParamsWithTimeout(timeout time.Duration) *GetVsphereThumbprintParams {
	var ()
	return &GetVsphereThumbprintParams{

		timeout: timeout,
	}
}

// NewGetVsphereThumbprintParamsWithContext creates a new GetVsphereThumbprintParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetVsphereThumbprintParamsWithContext(ctx context.Context) *GetVsphereThumbprintParams {
	var ()
	return &GetVsphereThumbprintParams{

		Context: ctx,
	}
}

// NewGetVsphereThumbprintParamsWithHTTPClient creates a new GetVsphereThumbprintParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetVsphereThumbprintParamsWithHTTPClient(client *http.Client) *GetVsphereThumbprintParams {
	var ()
	return &GetVsphereThumbprintParams{
		HTTPClient: client,
	}
}

/*
GetVsphereThumbprintParams contains all the parameters to send to the API endpoint
for the get vsphere thumbprint operation typically these are written to a http.Request
*/
type GetVsphereThumbprintParams struct {

	/*Host
	  vSphere host

	*/
	Host string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) WithTimeout(timeout time.Duration) *GetVsphereThumbprintParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) WithContext(ctx context.Context) *GetVsphereThumbprintParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) WithHTTPClient(client *http.Client) *GetVsphereThumbprintParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithHost adds the host to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) WithHost(host string) *GetVsphereThumbprintParams {
	o.SetHost(host)
	return o
}

// SetHost adds the host to the get vsphere thumbprint params
func (o *GetVsphereThumbprintParams) SetHost(host string) {
	o.Host = host
}

// WriteToRequest writes these params to a swagger request
func (o *GetVsphereThumbprintParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param host
	qrHost := o.Host
	qHost := qrHost
	if qHost != "" {
		if err := r.SetQueryParam("host", qHost); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
