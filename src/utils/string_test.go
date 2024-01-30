package utils

import "testing"

func TestCreateFileName(t *testing.T) {
	type args struct {
		format         string
		fileOriginName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{format: "%f", fileOriginName: "temp.jpg"},
			want: "temp.jpg",
		},
		{
			name: "2",
			args: args{format: "%d-%f", fileOriginName: "temp.jpg"},
			want: "temp.jpg",
		},
		{
			name: "3",
			args: args{format: "%r-%d-%f", fileOriginName: "temp.jpg"},
			want: "temp.jpg",
		},
		{
			name: "4",
			args: args{format: "", fileOriginName: "temp.jpg"},
			want: "temp.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFileName(tt.args.format, tt.args.fileOriginName)
			t.Log(got)
		})
	}
}
