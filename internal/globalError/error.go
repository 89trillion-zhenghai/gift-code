
package globalError

type GlobalError struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err GlobalError) Error() string {
	return err.Message
}
