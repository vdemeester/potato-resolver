/*
Copyright 2022 The Tekton Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"
	"regexp"

	"github.com/tektoncd/pipeline/pkg/apis/config"
	"github.com/tektoncd/pipeline/pkg/apis/version"
	"knative.dev/pkg/apis"
)

// ResultNameFormat Constant used to define the the regex Result.Name should follow
const ResultNameFormat = `^([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$`

var resultNameFormatRegex = regexp.MustCompile(ResultNameFormat)

// Validate implements apis.Validatable
func (tr TaskResult) Validate(ctx context.Context) (errs *apis.FieldError) {
	if !resultNameFormatRegex.MatchString(tr.Name) {
		return apis.ErrInvalidKeyName(tr.Name, "name", fmt.Sprintf("Name must consist of alphanumeric characters, '-', '_', and must start and end with an alphanumeric character (e.g. 'MyName',  or 'my-name',  or 'my_name', regex used for validation is '%s')", ResultNameFormat))
	}
	// Array and Object are alpha features
	if tr.Type == ResultsTypeArray || tr.Type == ResultsTypeObject {
		return errs.Also(version.ValidateEnabledAPIFields(ctx, "results type", config.AlphaAPIFields))
	}

	// Resources created before the result. Type was introduced may not have Type set
	// and should be considered valid
	if tr.Type == "" {
		return nil
	}

	// By default the result type is string
	if tr.Type != ResultsTypeString {
		return apis.ErrInvalidValue(tr.Type, "type", fmt.Sprintf("type must be string"))
	}

	return nil
}
