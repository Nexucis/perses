// Copyright 2021 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"encoding/json"
	"fmt"
	"reflect"

	modelAPI "github.com/perses/perses/pkg/model/api"
	"github.com/perses/perses/pkg/model/api/v1/common"
)

func FilterDatasource[T DatasourceInterface](kind string, defaultDTS *bool, list []T) []T {
	result := make([]T, 0, len(list))
	for _, d := range list {
		if (len(kind) > 0 && kind != d.GetDTSSpec().Plugin.Kind) ||
			(defaultDTS != nil && *defaultDTS != d.GetDTSSpec().Default) {
			continue
		}
		result = append(result, d)
	}
	return result
}

type DatasourceInterface interface {
	GetMetadata() modelAPI.Metadata
	GetDTSSpec() DatasourceSpec
}

type DatasourceSpec struct {
	Display *common.Display `json:"display,omitempty" yaml:"display,omitempty"`
	Default bool            `json:"default" yaml:"default"`
	// Plugin will contain the datasource configuration.
	// The data typed is available in Cue.
	Plugin common.Plugin `json:"plugin" yaml:"plugin"`
}

// GlobalDatasource is the struct representing the datasource shared to everybody.
// Any Dashboard can reference it.
type GlobalDatasource struct {
	Kind     Kind           `json:"kind" yaml:"kind"`
	Metadata Metadata       `json:"metadata" yaml:"metadata"`
	Spec     DatasourceSpec `json:"spec" yaml:"spec"`
}

func (d *GlobalDatasource) UnmarshalJSON(data []byte) error {
	var tmp GlobalDatasource
	type plain GlobalDatasource
	if err := json.Unmarshal(data, (*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*d = tmp
	return nil
}

func (d *GlobalDatasource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp GlobalDatasource
	type plain GlobalDatasource
	if err := unmarshal((*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*d = tmp
	return nil
}

func (d *GlobalDatasource) validate() error {
	if d.Kind != KindGlobalDatasource {
		return fmt.Errorf("invalid kind: %q for a GlobalDatasource type", d.Kind)
	}
	if reflect.DeepEqual(d.Spec, DatasourceSpec{}) {
		return fmt.Errorf("spec cannot be empty")
	}
	return nil
}

func (d *GlobalDatasource) GetMetadata() modelAPI.Metadata {
	return &d.Metadata
}

func (d *GlobalDatasource) GetKind() string {
	return string(d.Kind)
}

func (d *GlobalDatasource) GetDTSSpec() DatasourceSpec {
	return d.Spec
}

func (d *GlobalDatasource) GetSpec() interface{} {
	return d.Spec
}

// Datasource will be the datasource you can define in your project/namespace
// This is a resource that won't be shared across projects.
// A Dashboard can use it only if it is in the same project.
type Datasource struct {
	Kind     Kind            `json:"kind" yaml:"kind"`
	Metadata ProjectMetadata `json:"metadata" yaml:"metadata"`
	Spec     DatasourceSpec  `json:"spec" yaml:"spec"`
}

func (d *Datasource) UnmarshalJSON(data []byte) error {
	var tmp Datasource
	type plain Datasource
	if err := json.Unmarshal(data, (*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*d = tmp
	return nil
}

func (d *Datasource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp Datasource
	type plain Datasource
	if err := unmarshal((*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*d = tmp
	return nil
}

func (d *Datasource) validate() error {
	if d.Kind != KindDatasource {
		return fmt.Errorf("invalid kind: %q for a Datasource type", d.Kind)
	}
	if reflect.DeepEqual(d.Spec, DatasourceSpec{}) {
		return fmt.Errorf("spec cannot be empty")
	}
	return nil
}

func (d *Datasource) GetMetadata() modelAPI.Metadata {
	return &d.Metadata
}

func (d *Datasource) GetKind() string {
	return string(d.Kind)
}

func (d *Datasource) GetDTSSpec() DatasourceSpec {
	return d.Spec
}

func (d *Datasource) GetSpec() interface{} {
	return d.Spec
}
