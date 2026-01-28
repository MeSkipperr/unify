package mtr

func applyDefaults(cfg *Config) {
	if cfg.Protocol == "" {
		cfg.Protocol = ProtocolICMP
	}

	if cfg.Count == 0 {
		cfg.Count = 10
	}

	// DNS default ON
	if !cfg.UseDNS {
		// false = user explicitly disable
	} else {
		cfg.UseDNS = true
	}

	// JSON default ON
	if !cfg.JSON {
		cfg.JSON = true
	}
}
