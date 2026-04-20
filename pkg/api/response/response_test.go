// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetHTTPStatusMapsBusinessCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{name: "success", code: CodeSuccess, want: http.StatusOK},
		{name: "validation", code: CodeValidationFailed, want: http.StatusUnprocessableEntity},
		{name: "permission denied", code: CodePermissionDenied, want: http.StatusForbidden},
		{name: "database error", code: CodeDatabaseError, want: http.StatusInternalServerError},
		{name: "unknown", code: 99999, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, GetHTTPStatus(tt.code))
		})
	}
}

func TestResponseWithCodeUsesMappedHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseWithCode(c, CodeValidationFailed, gin.H{"field": "title"}, "invalid")

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var body Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, CodeValidationFailed, body.Code)
	require.Equal(t, "invalid", body.Message)
	require.NotNil(t, body.Data)
}
