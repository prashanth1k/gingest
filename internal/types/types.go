package types

// FileInfo represents information about a processed file
type FileInfo struct {
	RelativePath string
	AbsolutePath string
	Content      string
	IsBinary     bool
	Error        error
}

// Stats represents processing statistics
type Stats struct {
	NumFilesProcessed int
	NumDirsProcessed  int
	NumBinaryFiles    int
	NumSkippedFiles   int
	TotalContentBytes int64
	Source            string
	Branch            string
	AllPaths          []string // All file paths for tree generation
}
