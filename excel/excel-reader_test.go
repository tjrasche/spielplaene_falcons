package excel

import (
	"github.com/xuri/excelize/v2"
	"os"
	"testing"
)

func TestReader_readCurrentExcelFiles(t *testing.T) {
	type fields struct {
		cfBucket          string
		cfAccId           string
		cfAccessKeyId     string
		cfAccessKeySecret string
	}
	tests := []struct {
		name     string
		fields   fields
		wantFile *excelize.File
		wantErr  bool
	}{
		{name: "Basic Read Test", fields: fields{
			cfAccId:           os.Getenv("CF_ACCOUNT_ID"),
			cfAccessKeyId:     os.Getenv("CF_ACCESSKEY_ID"),
			cfAccessKeySecret: os.Getenv("CF_ACCESS_KEY_SECRET"),
		}, wantFile: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reader{
				cfBucket:          tt.fields.cfBucket,
				cfAccId:           tt.fields.cfAccId,
				cfAccessKeyId:     tt.fields.cfAccessKeyId,
				cfAccessKeySecret: tt.fields.cfAccessKeySecret,
			}
			gotFile, err := r.readCurrentExcelFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("readCurrentExcelFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFile == nil {
				t.Errorf("need file here!")
			}
		})
	}
}
