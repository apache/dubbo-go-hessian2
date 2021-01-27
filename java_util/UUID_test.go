/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_util

import "testing"

func TestUUID_ToString(t *testing.T) {
	type fields struct {
		MostSigBits  int64
		LeastSigBits int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{name: "one", fields: fields{
			MostSigBits:  int64(459021424248441700),
			LeastSigBits: int64(-7160773830801198154),
		}, want: "065ec58d-a89f-4b64-9c9f-d223ea2e73b6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid := UUID{
				MostSigBits:  tt.fields.MostSigBits,
				LeastSigBits: tt.fields.LeastSigBits,
			}
			if got := uuid.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
