// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package forwarder

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
)

// validationStatus represents the validation status
// known to the HTTP endpoint.
type validationStatus string

const (
	validValidationStatus        = validationStatus("valid")
	invalidValidationStatus      = validationStatus("invalid")
	notValidatedValidationStatus = validationStatus("not_validated")
)

func parseValidationStatus(v *bool) validationStatus {
	if v == nil {
		return notValidatedValidationStatus
	}
	if *v {
		return invalidValidationStatus
	}
	return validValidationStatus
}

var escapeQuotes = strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace

// CreateFormFile creates an [io.Writer] like [mime/multipart.Writer.CreateFromFile].
// This version allows to set the mime type, too.
func createFormFile(w *multipart.Writer, fieldname, filename, mimeType string) (io.Writer, error) {
	// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.20:src/mime/multipart/writer.go;l=140
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", mimeType)
	return w.CreatePart(h)
}

func buildRequest(
	doc []byte,
	filename *string,
	status validationStatus,
	url string,
	headers http.Header,
) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	var err error
	part := func(name, fname, mimeType, content string) {
		if err != nil {
			return
		}
		if fname == "" {
			err = writer.WriteField(name, content)
			return
		}
		var w io.Writer
		if w, err = createFormFile(writer, name, fname, mimeType); err == nil {
			_, err = w.Write([]byte(content))
		}
	}
	var fn string
	if filename != nil {
		fn = filepath.Base(*filename)
	} else {
		fn = "document.json"
	}
	part("advisory", fn, "application/json", string(doc))
	part("validation_status", "", "text/plain", string(status))

	if err := errors.Join(err, writer.Close()); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	contentType := writer.FormDataContentType()
	req.Header.Set("Content-Type", contentType)
	return req, nil
}
