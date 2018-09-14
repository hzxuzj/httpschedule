package cluster

import (
	"fmt"
)

func ParameterNotNULL(args string) error {
	return fmt.Errorf("Parameter %s cann't null", args)
}
