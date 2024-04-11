package processor

type Name string

// Processor common interface
// to run and stop certain enrichment service
type Processor interface {
	Run()
	Stop()
}

// Nameable returns certain name
type Nameable interface {
	Name() Name
}

// NameableProcessor like  Processor just
// has the certain name
type NameableProcessor interface {
	Nameable
	Processor
}
