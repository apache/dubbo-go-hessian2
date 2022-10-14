package idl

import "testing"

func TestRun(t *testing.T) {
	Run(Opt{
		GroupID:    "com.amap.aos",
		ArtifactID: "amap-aos-hsf-demo-api",
		OutType:    outTypeInit,
	})
}

func TestRun2(t *testing.T) {
	Run(Opt{
		IDLPath: "xxx",
		OutType: outTypeDeploy,
	})
}

func TestRun3(t *testing.T) {
	Run(Opt{
		IDLPath: "xxx",
		OutType: outTypeProvider,
	})
}
