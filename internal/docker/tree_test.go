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
		want   string
	}{
		{
			name:   "get string of only root node",
			fields: fields{"/", true, nil},
			want:   "/\n",
		},
		{
			name:   "get string of /etc/file",
			fields: fields{"/", true, []*fileTreeNode{&etcNode}},
			want:   "/\n└── etc/\n    └── file\n",
		},
		{
			name:   "get string of /etc/file + /other_file",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &otherFileNode}},
			want:   "/\n├── etc/\n│   └── file\n└── other_file\n",
		},
		{
			name:   "get string with symlink",
			fields: fields{"/", true, []*fileTreeNode{&etcNode, &binNodeWithSymlink}},
			want:   "/\n├── etc/\n│   └── file\n└── bin/\n    ├── file\n    └── link -> /tmp/file\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &fileTreeNode{
				Name:     tt.fields.Name,
				IsDir:    tt.fields.IsDir,
				Children: tt.fields.Children,
			}
			if got := n.getString("", true, true, true); got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeFileTrees(t *testing.T) {
	singleFileTree := &fileTreeNode{"file", "", false, nil}

	etcWithFile := &fileTreeNode{"etc", "", true, []*fileTreeNode{singleFileTree}}
	rootWithEtcTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithFile}}

	varWithFile := &fileTreeNode{"var", "", true, []*fileTreeNode{singleFileTree}}
	rootWithVarTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{varWithFile}}

	deleteSingleFileTree := &fileTreeNode{".wh.file", "", false, nil}
	etcWithDeleteFile := &fileTreeNode{"etc", "", true, []*fileTreeNode{deleteSingleFileTree}}
	rootWithEtcWithDeleteFileTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithDeleteFile}}

	rootWithEtcWithDeleteFileAndAddVarFileTreeNode := &fileTreeNode{"/", "", true, []*fileTreeNode{etcWithDeleteFile, varWithFile}}

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
