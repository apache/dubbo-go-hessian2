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
