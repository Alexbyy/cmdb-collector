package manager

import (
	"fmt"
	"testing"
)

func TestReadFromJson(t *testing.T) {

	src := "../../config.json"
	res, err := ReadFromJson(src)
	fmt.Printf("结果：%v\n", res)
	fmt.Printf("错误： %v\n", err)
	fmt.Printf("job_create_pod:>>>>%v\n", (*res).AssociConfig["job_create_pod"]["job"])

}
