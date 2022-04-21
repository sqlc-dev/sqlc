package golang

import "testing"

func TestStructName(t *testing.T) {

	testCases := []struct {
		name          string
		snakeCaseName string
		rename        string
		out           string
	}{
		{
			name:          "Rename Not Empty",
			snakeCaseName: "my_value",
			rename:        "MyValue",
			out:           "MyValue",
		},
		{
			name:          "ID to Upper Case",
			snakeCaseName: "id",
			rename:        "",
			out:           "ID",
		},
		{
			name:          "Camel Case",
			snakeCaseName: "user_id",
			rename:        "",
			out:           "UserID",
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			out := StructName(tc.snakeCaseName, tc.rename)
			if out != tc.out {
				t.Errorf("StructName(%s,%s) = %s; want %s", tc.snakeCaseName, tc.rename, out, tc.out)
			}
		})
	}
}
