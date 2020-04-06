package app

type Configuration struct {
	Verbose      bool
	DocumentRoot string `mapstructure:"document_root"`
	MetadataRoot string `mapstructure:"metadata_root"`
	TikaVersion  string `mapstructure:"tika_version"`
}
