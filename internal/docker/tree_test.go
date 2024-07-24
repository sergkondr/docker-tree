package docker

import (
	"reflect"
	"testing"
)

func Test_fileTreeNode_String(t *testing.T) {
	fileNode := fileTreeNode{"file", "", false, nil}
	fileSymlinkNode := fileTreeNode{"link", "/tmp/file", false, nil}
	otherFileNode := fileTreeNode{"other_file", "", false, nil}
	etcNode := fileTreeNode{"etc", "", true, []*fileTreeNode{&fileNode}}
	binNodeWithSymlink := fileTreeNode{"bin", "", true, []*fileTreeNode{&fileNode, &fileSymlinkNode}}

	type fields struct {
		Name     string
		IsDir    bool
		Children []*fileTreeNode
	}
	tests := []struct {
		name   string
		fields fields
		opts   getStringOpts
		want   string
	}{
		{
			name:   "get string of only root node",
			fields: fields{"/", true, nil},
			opts:   getStringOpts{showLinks: false, depth: 99999},
			want:   "/\n",
		},
		{
			name:   "get string of /etc/file",
			fields: fields{"/", true, []*fileTreeNode{&etcNode}},
			opts:   getStringOpts{showLinks: false, depth: 99999},
			want:   "/\n└── etc/\n    └── file\n",
		},
		{
			name:   "get string of /etc/file + /other_file",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &otherFileNode}},
			opts:   getStringOpts{showLinks: false, depth: 99999},
			want:   "/\n├── etc/\n│   └── file\n└── other_file\n",
		},
		{
			name:   "get string with symlink",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &binNodeWithSymlink}},
			opts:   getStringOpts{showLinks: true, depth: 99999},
			want:   "/\n├── etc/\n│   └── file\n└── bin/\n    ├── file\n    └── link -> /tmp/file\n",
		},
		{
			name:   "get string with depth = 1",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &binNodeWithSymlink}},
			opts:   getStringOpts{showLinks: false, depth: 2}, // we use depth == 2 because we want it to handle root + one more level of nesting
			want:   "/\n├── etc/\n└── bin/\n",
		},
		{
			name:   "get string with depth = 2",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &binNodeWithSymlink}},
			opts:   getStringOpts{showLinks: false, depth: 3},
			want:   "/\n├── etc/\n│   └── file\n└── bin/\n    ├── file\n    └── link\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &fileTreeNode{
				Name:     tt.fields.Name,
				IsDir:    tt.fields.IsDir,
				Children: tt.fields.Children,
			}

			if got := n.getString("", tt.opts, true, true); got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeFileTrees(t *testing.T) {
	singleFileTree := &fileTreeNode{"file", "", false, nil}
	whiteOutDirNode := &fileTreeNode{".wh..wh..opq", "", false, nil}

	etcWithFile := &fileTreeNode{"etc", "", true, []*fileTreeNode{singleFileTree}}
	rootWithEtcTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithFile}}

	varWithFile := &fileTreeNode{"var", "", true, []*fileTreeNode{singleFileTree}}
	rootWithVarTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{varWithFile}}

	deleteSingleFileTree := &fileTreeNode{".wh.file", "", false, nil}
	etcWithDeleteFile := &fileTreeNode{"etc", "", true, []*fileTreeNode{deleteSingleFileTree}}
	rootWithEtcWithDeleteFileTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithDeleteFile}}

	rootWithEtcWithDeleteFileAndAddVarFileTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithDeleteFile, varWithFile}}

	varWithWhFile := &fileTreeNode{"var", "", true, []*fileTreeNode{whiteOutDirNode}}
	varWithoutFile := &fileTreeNode{"var", "", true, []*fileTreeNode{}}
	rootWithVarWH := &fileTreeNode{"/", "", true, []*fileTreeNode{varWithWhFile}}

	type args struct {
		original *fileTreeNode
		updated  *fileTreeNode
	}
	tests := []struct {
		name    string
		args    args
		want    *fileTreeNode
		wantErr bool
	}{
		{
			name: "original is nil",
			args: args{
				original: nil,
				updated:  singleFileTree,
			},
			want:    singleFileTree,
			wantErr: false,
		},
		{
			name: "add /var/file to /etc/file",
			args: args{
				original: rootWithEtcTreeNode,
				updated:  rootWithVarTreeNode,
			},
			want:    &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithFile, varWithFile}},
			wantErr: false,
		},
		{
			name: "delete /etc/file",
			args: args{
				original: rootWithEtcTreeNode,
				updated:  rootWithEtcWithDeleteFileTreeNode,
			},
			want:    &fileTreeNode{"/", "", true, []*fileTreeNode{{"etc", "", true, []*fileTreeNode{}}}},
			wantErr: false,
		},
		{
			name: "delete /etc/file and add /var/file",
			args: args{
				original: rootWithEtcTreeNode,
				updated:  rootWithEtcWithDeleteFileAndAddVarFileTreeNode,
			},
			want:    &fileTreeNode{"/", "", true, []*fileTreeNode{{"etc", "", true, []*fileTreeNode{}}, varWithFile}},
			wantErr: false,
		},
		{
			name: "white out /var/",
			args: args{
				original: rootWithVarTreeNode, // /var/file
				updated:  rootWithVarWH,
			},
			want:    varWithoutFile,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeFileTrees(tt.args.original, tt.args.updated)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeFileTrees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeFileTrees() got:\n%v, want:\n%v", got, tt.want)
			}
		})
	}
}
