/*
Copyright 2020 The Tekton Authors

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

package builder

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func TestParam(t *testing.T) {
	got := Param("foo", "bar")

	want := v1beta1.Param{
		Name: "foo",
		Value: v1beta1.ArrayOrString{
			Type:      v1beta1.ParamTypeString,
			StringVal: "bar",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("-want/+got: %s", diff)
	}
}
