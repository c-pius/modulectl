package filegenerator

import (
	"fmt"

	"github.com/kyma-project/modulectl/internal/scaffold/common/errors"
	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/tools/io"
)

type FileReader interface {
	FileExists(path string) (bool, error)
}

type FileGenerator interface {
	GenerateFile(out io.Out, path string, args types.KeyValueArgs) error
}

type ReuseFileGeneratorService struct {
	kind          string
	fileReader    FileReader
	fileGenerator FileGenerator
}

func NewReuseFileGeneratorService(
	kind string,
	fileSystem FileReader,
	fileGenerator FileGenerator,
) (*ReuseFileGeneratorService, error) {
	if kind == "" {
		return nil, fmt.Errorf("%w: kind must not be empty", errors.ErrInvalidArg)
	}

	if fileSystem == nil {
		return nil, fmt.Errorf("%w: fileSystem must not be nil", errors.ErrInvalidArg)

	}

	if fileGenerator == nil {
		return nil, fmt.Errorf("%w: fileGenerator must not be nil", errors.ErrInvalidArg)
	}

	return &ReuseFileGeneratorService{
		kind:          kind,
		fileReader:    fileSystem,
		fileGenerator: fileGenerator,
	}, nil
}

func (s *ReuseFileGeneratorService) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	fileExists, err := s.fileReader.FileExists(path)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrCheckingFileExistence, path, err)
	}

	if fileExists {
		out.Write(fmt.Sprintf("The %s file already exists, reusing: %s\n", s.kind, path))
		return nil
	}

	return s.fileGenerator.GenerateFile(out, path, args)
}
