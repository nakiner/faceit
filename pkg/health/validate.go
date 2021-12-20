package health

type validator interface {
	Validate() error
}

func validate(req interface{}) error {
	if val, ok := interface{}(req).(validator); ok {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
}
