package iwrapper

func Convert(results []*ParseResult) ([]*GenerateConfig, error) {
	generateConfigs := make([]*GenerateConfig, 0, len(results))
	for _, result := range results {
		funcName := result.FuncName
		if funcName == "" {
			funcName = result.StructName + "Wrapper"
		}

		wrappedInterfaces := make([]*Interface, 0, len(result.RequiredInterfaces)+len(result.OptionalInterfaces))
		wrappedInterfaces = append(wrappedInterfaces, result.RequiredInterfaces...)
		wrappedInterfaces = append(wrappedInterfaces, result.OptionalInterfaces...)

		generateConfigs = append(generateConfigs, &GenerateConfig{
			FuncName:           funcName,
			RequireInterface:   NewAnonymousInterface(result.RequiredInterfaces),
			WrappedInterface:   NewNamedInterface(result.StructName, wrappedInterfaces, true),
			OptionalInterfaces: result.OptionalInterfaces,
		})
	}

	return generateConfigs, nil
}
