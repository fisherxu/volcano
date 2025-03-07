/*
Copyright 2019 The Volcano Authors.

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

package options

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestAddFlags(t *testing.T) {
	fs := pflag.NewFlagSet("addflagstest", pflag.ContinueOnError)
	s := NewServerOption()
	s.AddFlags(fs)

	args := []string{
		"--master=127.0.0.1",
		"--kube-api-burst=200",
	}
	fs.Parse(args)

	// This is a snapshot of expected options parsed by args.
	expected := &ServerOption{
		Master:        "127.0.0.1",
		KubeAPIQPS:    defaultQPS,
		KubeAPIBurst:  200,
		PrintVersion:  false,
		WorkerThreads: defaultWorkers,
		SchedulerName: defaultSchedulerName,
	}

	if !reflect.DeepEqual(expected, s) {
		t.Errorf("Got different run options than expected.\nGot: %+v\nExpected: %+v\n", s, expected)
	}

	err := s.CheckOptionOrDie()
	if err != nil {
		t.Errorf("expected nil but got %v\n", err)
	}

}
