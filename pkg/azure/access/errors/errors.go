// Copyright 2023 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/gardener/machine-controller-manager-provider-azure/pkg/azure/utils"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

const (
	// ZonalAllocationFailedAzErrorCode is an Azure error code indicating that there is insufficient capacity in the target zone.
	ZonalAllocationFailedAzErrorCode = "ZonalAllocationFailed"
	// CorrelationRequestIDAzHeaderKey is the Azure API response header key whose value is a request correlation ID.
	CorrelationRequestIDAzHeaderKey = "x-ms-correlation-request-id"
	// RequestIDAzHeaderKey is the Azure API response header key whose value is the request ID.
	RequestIDAzHeaderKey = "x-ms-request-id"
	// ErrorCodeAzHeaderKey is the Azure API response header key whose value is the underline error code as set by the server.
	ErrorCodeAzHeaderKey = "x-ms-error-code"
	// ClientRequestIDAzHeaderKey is the Azure API response header key whose value is the client request ID.
	ClientRequestIDAzHeaderKey = "x-ms-client-request-id"
)

var (
	// Raised https://github.com/Azure/azure-sdk-for-go/issues/21094 to prevent hard coding these here and instead
	// use well-maintained constants defined in the Azure SDK.
	lookupResponseHeaderKeys = sets.New(
		CorrelationRequestIDAzHeaderKey,
		RequestIDAzHeaderKey,
		ErrorCodeAzHeaderKey,
		ClientRequestIDAzHeaderKey,
	)
)

// IsNotFoundAzAPIError checks if error is an AZ API error and if it is a 404 response code.
func IsNotFoundAzAPIError(err error) bool {
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		return respErr.StatusCode == http.StatusNotFound
	}
	return false
}

// LogAzAPIError collects additional information from AZ response and logs it as part of the error log message.
func LogAzAPIError(err error, format string, v ...any) {
	if err == nil {
		return
	}
	respHeaders := traceResponseHeaders(err)
	errMsg := fmt.Sprintf(format, v...)
	if len(respHeaders) == 0 {
		klog.Errorf("%s: %+v\n", errMsg, err)
		return
	}
	klog.Errorf("%s : Azure API Response-Headers: %+v Err: %+v\n", errMsg, respHeaders, err)
}

func traceResponseHeaders(err error) map[string]string {
	var respErr *azcore.ResponseError
	headers := make(map[string]string)
	if errors.As(err, &respErr) {
		respHeader := respErr.RawResponse.Header
		for headerKey := range lookupResponseHeaderKeys {
			headerValue := respHeader.Get(headerKey)
			if !utils.IsEmptyString(headerValue) {
				headers[headerKey] = headerValue
			}
		}
	}
	return headers
}

// GetMatchingErrorCode gets a matching codes.Code for the given azure error code.
func GetMatchingErrorCode(err error) codes.Code {
	azHeaders := traceResponseHeaders(err)
	azErrorCode, ok := azHeaders[ErrorCodeAzHeaderKey]
	if ok {
		switch azErrorCode {
		case ZonalAllocationFailedAzErrorCode:
			return codes.ResourceExhausted
		default:
			return codes.Internal
		}
	}
	return codes.Internal
}
