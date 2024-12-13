package smugmug

import (
	"testing"
	"text/template"
)

func Test_albumImage_buildFilename(t *testing.T) {
	type fields struct {
		FileName         string
		ImageKey         string
		ArchivedMD5      string
		UploadKey        string
		DateTimeOriginal string
	}

	f := fields{
		FileName:         "ti_7095_R.jpg",
		ImageKey:         "ImageKeyValue",
		ArchivedMD5:      "ArchivedMD5Value",
		UploadKey:        "UploadKeyValue",
		DateTimeOriginal: "2012-06-06T21:08:48+00:00",
	}

	f_wo_date := fields{
		FileName:    "ti_7095_R.jpg",
		ImageKey:    "ImageKeyValue",
		ArchivedMD5: "ArchivedMD5Value",
		UploadKey:   "UploadKeyValue",
	}

	tests := []struct {
		name         string
		fields       fields
		filenameConf string
		want         string
		wantErr      bool
	}{
		{
			name:         "filename",
			fields:       f,
			filenameConf: "{{.FileName}}",
			want:         "ti_7095_R.jpg",
			wantErr:      false,
		},
		{
			name: "with apostrophe",
			fields: fields{
				FileName: "FileName'Value",
			},
			filenameConf: "{{.FileName}}",
			want:         "FileName'Value",
			wantErr:      false,
		},
		{
			name:         "empty",
			fields:       f,
			filenameConf: "",
			want:         "",
			wantErr:      true,
		},
		{
			name:         "wrong",
			fields:       f,
			filenameConf: "{{.WrongKey}}",
			want:         "",
			wantErr:      true,
		},
		{
			name:         "wrong with extra chars",
			fields:       f,
			filenameConf: "{{.WrongKey}}-",
			want:         "-",
			wantErr:      true,
		},
		{
			name:         "complex",
			fields:       f,
			filenameConf: "{{.ImageKey}}-{{.FileName}}",
			want:         "ImageKeyValue-ti_7095_R.jpg",
			wantErr:      false,
		},
		{
			name:         "all",
			fields:       f,
			filenameConf: "prefix-{{.UploadKey}}/{{.ImageKey}}-{{.FileName}}_{{.ArchivedMD5}}",
			want:         "prefix-UploadKeyValue/ImageKeyValue-ti_7095_R.jpg_ArchivedMD5Value",
			wantErr:      false,
		},
		{
			name:         "date and time",
			fields:       f,
			filenameConf: "prefix-{{.Date}}/{{.Time}}",
			want:         "prefix-2012-06-06/21_08_48",
			wantErr:      false,
		},
		{
			name:         "date and time but empty",
			fields:       f_wo_date,
			filenameConf: "prefix-{{.UploadKey}}/{{.ImageKey}}-{{.Date}}/{{.Time}}",
			want:         "prefix-UploadKeyValue/ImageKeyValue-/",
			wantErr:      false,
		},
		{
			name:         "Filename with FileNameNoExt and Extension",
			fields:       f_wo_date,
			filenameConf: "prefix-{{.FileName}}_{{.FileNameNoExt}}_{{.Extension}}",
			want:         "prefix-ti_7095_R.jpg_ti_7095_R_.jpg",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &albumImage{
				FileName:         tt.fields.FileName,
				ImageKey:         tt.fields.ImageKey,
				ArchivedMD5:      tt.fields.ArchivedMD5,
				UploadKey:        tt.fields.UploadKey,
				DateTimeOriginal: tt.fields.DateTimeOriginal,
			}

			tmpl, err := template.New("image_filename").Option("missingkey=error").Parse(tt.filenameConf)
			if err != nil {
				t.Fatal(err)
			}
			err = a.buildFilename(tmpl)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error: %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && a.Name() != tt.want {
				t.Fatalf("want: %s, got: %s", tt.want, a.Name())
			}
		})
	}
}
