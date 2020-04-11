package app

type Configuration struct {
	Verbose        bool
	DocumentRoot   string `mapstructure:"document_root"`
	MetadataDB     string `mapstructure:"metadata_db"`
	TikaVersion    string `mapstructure:"tika_version"`
	ClassifierData string `mapstructure:"classifier_data"`
}
