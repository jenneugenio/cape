package coordinator

type SchemaOptions struct {
	BlobPath string `json:"blob_path"`
}

type SourceOptions struct {
	WithSchema    bool
	SchemaOptions *SchemaOptions
}
