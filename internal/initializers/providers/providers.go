package providers

func Dependencies() []interface{} {
	return []interface{}{
		NewConfig,
		NewProcessor,
		NewLogger,
		NewStore,
	}
}
