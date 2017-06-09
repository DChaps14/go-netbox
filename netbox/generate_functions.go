// Copyright 2017 The go-netbox Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
	"time"
)

func main() {
	typeName := flag.String("type-name", "Example", "Name of the type to use (e.g. TenantGroup).")
	serviceName := flag.String("service-name", "ExampleService", "Name of the service to create (e.g. TenantGroupsService).")
	endpoint := flag.String("endpoint", "tenancy", "Name of the endpoint (e.g. dcim, ipam, tenancy).")
	service := flag.String("service", "example", "Name of the service below endpoint (e.g. tenant-groups).")
	updateTypeName := flag.String("update-type-name", "", "Name of the type to use for creates and updates, to change the marshal behavior. Default typeName.")
	withoutListOpts := flag.Bool("without-list-opts", false, "Disable list options for this endpoint.")

	flag.Parse()

	if *updateTypeName == "" {
		*updateTypeName = *typeName
	}

	b := &bytes.Buffer{}
	functionsTemplate.Execute(b, struct {
		Timestamp      time.Time
		TypeName       string
		UpdateTypeName string
		ServiceName    string
		Endpoint       string
		Service        string
		ListOpts       bool
	}{
		Timestamp:      time.Now(),
		TypeName:       *typeName,
		UpdateTypeName: *updateTypeName,
		ServiceName:    *serviceName,
		Endpoint:       *endpoint,
		Service:        *service,
		ListOpts:       !*withoutListOpts,
	})

	// go fmt
	res, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fmt.Sprintf("%s_%s.go", *endpoint, *service))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

var functionsTemplate = template.Must(template.New("").Parse(`// Copyright 2017 The go-netbox Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by generate_functions.go. DO NOT EDIT.

package netbox

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// {{ .ServiceName }} is used in a Client to access NetBox's {{ .Endpoint }}/{{ .Service }} API methods.
type {{ .ServiceName }} struct {
	c *Client
}

// Get retrieves an {{ .TypeName }} object from NetBox by its ID.
func (s *{{ .ServiceName }}) Get(id int) (*{{ .TypeName }}, error) {
	req, err := s.c.NewRequest(
		http.MethodGet,
		fmt.Sprintf("api/{{ .Endpoint }}/{{ .Service }}/%d/", id),
		nil,
	)
	if err != nil {
		return nil, err
	}

	t := new({{ .TypeName }})
	err = s.c.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// List returns a Page associated with an NetBox API Endpoint.
{{ if .ListOpts -}}
func (s *{{ .ServiceName }}) List(options *List{{ .TypeName }}Options) *Page {
	return NewPage(s.c, "api/{{ .Endpoint }}/{{ .Service }}/", options)
}
{{ else -}}
func (s *{{ .ServiceName }}) List() *Page {
	return NewPage(s.c, "api/{{ .Endpoint }}/{{ .Service }}/", nil)
}
{{ end }}

// Extract retrives a list of {{ .TypeName }} objects from page.
func (s *{{ .ServiceName }}) Extract(page *Page) ([]*{{ .TypeName }}, error) {
	if err := page.Err(); err != nil {
		return nil, err
	}

	var groups []*{{ .TypeName }}
	if err := json.Unmarshal(page.data.Results, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// Create creates a new {{ .TypeName }} object in NetBox and returns the ID of the new object.
func (s *{{ .ServiceName }}) Create(data *{{ .TypeName }}) (int, error) {
	req, err := s.c.NewJSONRequest(http.MethodPost, "api/{{ .Endpoint }}/{{ .Service }}/", nil, data)
	if err != nil {
		return 0, err
	}

	g := new({{ .UpdateTypeName }})
	err = s.c.Do(req, g)
	if err != nil {
		return 0, err
	}
	return g.ID, nil
}

// Update changes an existing {{ .TypeName }} object in NetBox, and returns the ID of the new object.
func (s *{{ .ServiceName }}) Update(data *{{ .TypeName }}) (int, error) {
	req, err := s.c.NewJSONRequest(
		http.MethodPatch,
		fmt.Sprintf("api/{{ .Endpoint }}/{{ .Service }}/%d/", data.ID),
		nil,
		data,
	)
	if err != nil {
		return 0, err
	}

	// g is just used to verify correct api result.
	// data is not changed, because the g is not the full representation that one would
	// get with Get. But if the response was unmarshaled into {{ .UpdateTypeName }} correctly,
	// everything went fine, and we do not need to update data.
	g := new({{ .UpdateTypeName }})
	err = s.c.Do(req, g)
	if err != nil {
		return 0, err
	}
	return g.ID, nil
}

// Delete deletes an existing {{ .TypeName }} object from NetBox.
func (s *{{ .ServiceName }}) Delete(data *{{ .TypeName }}) error {
	req, err := s.c.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("api/{{ .Endpoint }}/{{ .Service }}/%d/", data.ID),
		nil,
	)
	if err != nil {
		return err
	}

	return s.c.Do(req, nil)
}
`))
