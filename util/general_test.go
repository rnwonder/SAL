package util

import (
	"github.com/rnwonder/SAL/data"
	"reflect"
	"testing"
	"time"
)

func TestMyCmpWorkAround(t *testing.T) {
	type args struct {
		value1 string
		value2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test MyCmpWorkAround",
			args: args{
				value1: "test",
				value2: "test",
			},
			want: "test",
		},
		{
			name: "Test MyCmpWorkAround",
			args: args{
				value1: "",
				value2: "test",
			},
			want: "test",
		},
		{
			name: "Test MyCmpWorkAround",
			args: args{
				value1: "test",
				value2: "",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MyCmpWorkAround(tt.args.value1, tt.args.value2); got != tt.want {
				t.Errorf("MyCmpWorkAround() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePageInfo(t *testing.T) {
	type args struct {
		page  string
		limit string
		total int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
		want2 int
		want3 int
		want4 int
	}{
		{
			name: "Test CalculatePageInfo",
			args: args{
				page:  "1",
				limit: "10",
				total: 100,
			},
			want:  0,
			want1: 10,
			want2: 10,
			want3: 10,
			want4: 1,
		},
		{
			name: "Test CalculatePageInfo",
			args: args{
				page:  "2",
				limit: "10",
				total: 100,
			},
			want:  10,
			want1: 20,
			want2: 10,
			want3: 10,
			want4: 2,
		},
		{
			name: "Test CalculatePageInfo",
			args: args{
				page:  "1",
				limit: "15",
				total: 36,
			},
			want:  0,
			want1: 15,
			want2: 3,
			want3: 15,
			want4: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4 := CalculatePageInfo(tt.args.page, tt.args.limit, tt.args.total)
			if got != tt.want {
				t.Errorf("CalculatePageInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CalculatePageInfo() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("CalculatePageInfo() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("CalculatePageInfo() got3 = %v, want %v", got3, tt.want3)
			}
			if got4 != tt.want4 {
				t.Errorf("CalculatePageInfo() got4 = %v, want %v", got4, tt.want4)
			}
		})
	}
}

func TestClientProductFormat(t *testing.T) {
	type args struct {
		product data.Product
	}
	tests := []struct {
		name string
		args args
		want data.ProductResponse
	}{
		{
			name: "Test ClientProductFormat",
			args: args{
				product: data.Product{
					Id:          "1",
					SkuId:       "1",
					Name:        "test",
					Description: "test",
					Price:       10,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want: data.ProductResponse{
				Id:          "1",
				Name:        "test",
				Description: "test",
				Price:       10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "Test ClientProductFormat",
			args: args{
				product: data.Product{
					Id:          "2",
					SkuId:       "2",
					Name:        "test2",
					Description: "test2",
					Price:       20,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want: data.ProductResponse{
				Id:          "2",
				Name:        "test2",
				Description: "test2",
				Price:       20,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientProductFormat(tt.args.product); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientProductFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test CompareHashAndPassword",
			args: args{
				hash:     HashPassword("test"),
				password: "test",
			},
			want: true,
		},
		{
			name: "Test CompareHashAndPassword",
			args: args{
				hash:     "test",
				password: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareHashAndPassword(tt.args.hash, tt.args.password); got != tt.want {
				t.Errorf("CompareHashAndPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test HashPassword",
			args: args{
				password: "test",
			},
			want: HashPassword("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashPassword(tt.args.password); got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextPage(t *testing.T) {
	type args struct {
		page       int
		totalPages int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test NextPage",
			args: args{
				page:       1,
				totalPages: 10,
			},
			want: "2",
		},
		{
			name: "Test NextPage",
			args: args{
				page:       10,
				totalPages: 10,
			},
			want: "10",
		},
		{
			name: "Test NextPage",
			args: args{
				page:       5,
				totalPages: 10,
			},
			want: "6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextPage(tt.args.page, tt.args.totalPages); got != tt.want {
				t.Errorf("NextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrevPage(t *testing.T) {
	type args struct {
		page int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test PrevPage",
			args: args{
				page: 1,
			},
			want: "1",
		},
		{
			name: "Test PrevPage",
			args: args{
				page: 10,
			},
			want: "9",
		},
		{
			name: "Test PrevPage",
			args: args{
				page: 5,
			},
			want: "4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrevPage(tt.args.page); got != tt.want {
				t.Errorf("PrevPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateUniqueEmail(t *testing.T) {
	type args struct {
		emailSet map[string]bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateUniqueEmail(tt.args.emailSet); got != tt.want {
				t.Errorf("generateUniqueEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateUniqueName(t *testing.T) {
	type args struct {
		nameSet map[string]bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test generateUniqueName",
			args: args{
				nameSet: make(map[string]bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateUniqueName(tt.args.nameSet)
			got1 := generateUniqueName(tt.args.nameSet)
			if got == got1 {
				t.Errorf("generateUniqueName() = %v, want %v", got, got1)
			}
		})
	}
}
