package smugmug

import (
	"html/template"
	"testing"
)

func Test_albumImage_buildFilename(t *testing.T) {
	type fields struct {
		FileName                      string
		ImageKey                      string
		ArchivedMD5                   string
		UploadKey                     string
		PreferredDisplayFileExtension string
	}

	f := fields{
		FileName:                      "FileNameValue",
		ImageKey:                      "ImageKeyValue",
		ArchivedMD5:                   "ArchivedMD5Value",
		UploadKey:                     "UploadKeyValue",
		PreferredDisplayFileExtension: "PreferredDisplayFileExtensionValue",
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
			want:         "FileNameValue",
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
			want:         "ImageKeyValue-FileNameValue",
			wantErr:      false,
		},
		{
			name:         "all",
			fields:       f,
			filenameConf: "prefix-{{.UploadKey}}/{{.ImageKey}}-{{.FileName}}_{{.ArchivedMD5}}",
			want:         "prefix-UploadKeyValue/ImageKeyValue-FileNameValue_ArchivedMD5Value",
			wantErr:      false,
		},
		{
			name: "extension manipulation - happy path",
			fields: fields{
				FileName:                      "FileNameValue.ExtenionInFileName",
				ImageKey:                      "ImageKeyValue",
				ArchivedMD5:                   "ArchivedMD5Value",
				UploadKey:                     "UploadKeyValue",
				PreferredDisplayFileExtension: "PreferredDisplayFileExtensionValue",
			},
			filenameConf: "{{.FileNameNoExt}}-{{.ImageKey}}.{{.Extension}}",
			want:         "FileNameValue-ImageKeyValue.ExtenionInFileName",
			wantErr:      false,
		},
		{
			name: "extension manipulation - no extension in FileName",
			fields: fields{
				FileName:                      "FileNameValue",
				ImageKey:                      "ImageKeyValue",
				ArchivedMD5:                   "ArchivedMD5Value",
				UploadKey:                     "UploadKeyValue",
				PreferredDisplayFileExtension: "PreferredDisplayFileExtensionValue",
			},
			filenameConf: "{{.FileNameNoExt}}-{{.ImageKey}}.{{.Extension}}",
			want:         "FileNameValue-ImageKeyValue.PreferredDisplayFileExtensionValue",
			wantErr:      false,
		},
		{
			name: "extension manipulation - last char is period",
			fields: fields{
				FileName:                      "FileNameValue.",
				ImageKey:                      "ImageKeyValue",
				ArchivedMD5:                   "ArchivedMD5Value",
				UploadKey:                     "UploadKeyValue",
				PreferredDisplayFileExtension: "PreferredDisplayFileExtensionValue",
			},
			filenameConf: "{{.FileNameNoExt}}-{{.ImageKey}}.{{.Extension}}",
			want:         "FileNameValue-ImageKeyValue.",
			wantErr:      false,
		},
		{
			name: "extension manipulation - first char is period",
			fields: fields{
				FileName:                      ".HiddenFiles",
				ImageKey:                      "ImageKeyValue",
				ArchivedMD5:                   "ArchivedMD5Value",
				UploadKey:                     "UploadKeyValue",
				PreferredDisplayFileExtension: "PreferredDisplayFileExtensionValue",
			},
			filenameConf: "-{{.ImageKey}}.{{.Extension}}",
			want:         "-ImageKeyValue.HiddenFiles",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &albumImage{
				FileName:                      tt.fields.FileName,
				ImageKey:                      tt.fields.ImageKey,
				ArchivedMD5:                   tt.fields.ArchivedMD5,
				UploadKey:                     tt.fields.UploadKey,
				PreferredDisplayFileExtension: tt.fields.PreferredDisplayFileExtension,
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
