package yaml

import "testing"

func TestConfigProvider(t *testing.T) {
	cfgProv, err := NewConfigProvider("../gamedays/")
	if err != nil {
		panic(err)
	}
	for _, gd := range cfgProv.GamdeDays {
		t.Log(gd.Bucket)
	}
}
