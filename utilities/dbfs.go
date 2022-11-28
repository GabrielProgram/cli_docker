package utilities

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/dbfs"
)

// move to go sdk / replace with utility function once
// https://github.com/databricks/databricks-sdk-go/issues/57 is Done
// Tracked in https://github.com/databricks/bricks/issues/25
func CreateDbfsFile(ctx context.Context,
	w *databricks.WorkspaceClient,
	path string,
	contents []byte,
	overwrite bool,
) error {
	// see https://docs.databricks.com/dev-tools/api/latest/dbfs.html#add-block
	const WRITE_BYTE_CHUNK_SIZE = 1e6
	createResponse, err := w.Dbfs.Create(ctx,
		dbfs.Create{
			Overwrite: overwrite,
			Path:      path,
		},
	)
	if err != nil {
		return err
	}
	handle := createResponse.Handle
	buffer := bytes.NewBuffer(contents)
	for {
		byteChunk := buffer.Next(WRITE_BYTE_CHUNK_SIZE)
		if len(byteChunk) == 0 {
			break
		}
		b64Data := base64.StdEncoding.EncodeToString(byteChunk)
		err := w.Dbfs.AddBlock(ctx,
			dbfs.AddBlock{
				Data:   b64Data,
				Handle: handle,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot add block: %w", err)
		}
	}
	err = w.Dbfs.Close(ctx,
		dbfs.Close{
			Handle: handle,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot close handle: %w", err)
	}
	return nil
}

func ReadDbfsFile(ctx context.Context,
	w *databricks.WorkspaceClient,
	path string,
) (content []byte, err error) {
	// see https://docs.databricks.com/dev-tools/api/latest/dbfs.html#read
	const READ_BYTE_CHUNK_SIZE = 1e6
	fetchLoop := true
	offSet := 0
	length := int(READ_BYTE_CHUNK_SIZE)
	for fetchLoop {
		dbfsReadReponse, err := w.Dbfs.Read(ctx,
			dbfs.ReadRequest{
				Path:   path,
				Offset: offSet,
				Length: length,
			},
		)
		if err != nil {
			return content, fmt.Errorf("cannot read %s: %w", path, err)
		}
		if dbfsReadReponse.BytesRead == 0 || dbfsReadReponse.BytesRead < int64(length) {
			fetchLoop = false
		}
		decodedBytes, err := base64.StdEncoding.DecodeString(dbfsReadReponse.Data)
		if err != nil {
			return content, err
		}
		content = append(content, decodedBytes...)
		offSet += length
	}
	return content, err
}
