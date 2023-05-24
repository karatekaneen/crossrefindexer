package elastic

type IndexSettings struct {
	Settings Settings `json:"settings,omitempty"`
	Mappings Mappings `json:"mappings,omitempty"`
}

type Index struct {
	NumberOfReplicas int    `json:"number_of_replicas,omitempty"`
	RefreshInterval  int    `json:"refresh_interval,omitempty"`
	Codec            string `json:"codec,omitempty"`
}

type Stop struct {
	Type      string `json:"type,omitempty"`
	Stopwords string `json:"stopwords,omitempty"`
}

type Filter struct {
	MyStop Stop `json:"my_stop,omitempty"`
}

type AnalyzerSettings struct {
	Type      string   `json:"type,omitempty"`
	Tokenizer string   `json:"tokenizer,omitempty"`
	Filter    []string `json:"filter,omitempty"`
}

type Analysis struct {
	Filter   Filter                      `json:"filter,omitempty"`
	Analyzer map[string]AnalyzerSettings `json:"analyzer,omitempty"`
}

type Settings struct {
	Index    Index    `json:"index,omitempty"`
	Analysis Analysis `json:"analysis,omitempty"`
}

type FieldSetting struct {
	Type     string `json:"type,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

type Mappings struct {
	Properties map[string]FieldSetting `json:"properties,omitempty"`
}

func DefaultSettings(modifiers ...func(*IndexSettings)) IndexSettings {
	settings := IndexSettings{
		Settings: Settings{
			Index: Index{
				NumberOfReplicas: 0,
				RefreshInterval:  -1,
				Codec:            "best_compression",
			},
			Analysis: Analysis{
				Filter: Filter{
					MyStop: Stop{
						Type:      "stop",
						Stopwords: "_english_",
					},
				},
				Analyzer: map[string]AnalyzerSettings{
					"case_insensitive_keyword": {
						Type:      "custom",
						Tokenizer: "keyword",
						Filter:    []string{"lowercase"},
					},
					"case_insensitive_folding_keyword": {
						Type:      "custom",
						Tokenizer: "keyword",
						Filter:    []string{"lowercase", "asciifolding"},
					},
					"case_insensitive_folding_text": {
						Type:      "custom",
						Tokenizer: "standard",
						Filter:    []string{"lowercase", "asciifolding"},
					},
					"case_insensitive_folding_text_stopwords": {
						Type:      "custom",
						Tokenizer: "standard",
						Filter:    []string{"lowercase", "asciifolding", "my_stop"},
					},
				},
			},
		},
		Mappings: Mappings{
			Properties: map[string]FieldSetting{
				"DOI": {
					Type:     "text",
					Analyzer: "case_insensitive_keyword",
				},
				"title": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_text_stopwords",
				},
				"first_author": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"author": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_text",
				},
				"first_page": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"journal": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_text_stopwords",
				},
				"abbreviated_journal": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"volume": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"issue": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"year": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_keyword",
				},
				"query": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_text",
				},
				"bibliographic": {
					Type:     "text",
					Analyzer: "case_insensitive_folding_text_stopwords",
				},
			},
		},
	}

	// Apply modifiers
	for _, modifier := range modifiers {
		modifier(&settings)
	}

	return settings
}
