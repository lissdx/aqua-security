package validator

type ValidatorFn func(string) (bool, error)
type ValidatorGoFn func(interface{}) (bool, error)
type ValidatorFnMap map[string]ValidatorFn
type ValidatorGoFnMap map[string]ValidatorGoFn

type ValidatorFactory interface {
	GetValidatorFn(string) (ValidatorFn, error)
}

type ValidatorGoFactory interface {
	GetGoValidatorFn(string) (ValidatorGoFn, error)
}
