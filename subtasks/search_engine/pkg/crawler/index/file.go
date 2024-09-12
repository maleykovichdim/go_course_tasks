// Reverse index package (storage and functionality)
// Here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted index

// Function - to load/save data to/from file on disk

package index

import (
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

// Name/path of data file
var internalStorageFileName string = "storage.pb"

// Open for read
func GetFileReader() (io.ReadCloser, error) {
	if _, err := os.Stat(internalStorageFileName); os.IsNotExist(err) {
		log.Error().Err(err).Msg("Storage file does not exist.")
		return nil, errors.New("storage file does not exist")
	}

	file, err := os.Open(internalStorageFileName)
	if err != nil {
		log.Error().Err(err).Msg("Error opening storage file.")
		return nil, errors.New("storage file open error")
	}

	log.Info().Msg("Storage file opened for reading.")
	return file, nil
}

// Open for write
func GetFileWriter() (io.WriteCloser, error) {
	file, err := os.OpenFile(internalStorageFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Error creating storage file.")
		return nil, errors.New("storage file was not created")
	}

	log.Info().Msg("Storage file opened for writing.")
	return file, nil
}

// Write to file
func WriteDataToStorage(data io.Reader) error {
	wr, err := GetFileWriter()
	if err != nil {
		return err
	}
	defer wr.Close()

	_, err = io.Copy(wr, data)
	if err != nil {
		log.Error().Err(err).Msg("Error writing data to storage file.")
		return err
	}
	log.Info().Msg("Data successfully written to storage file.")
	return nil
}

// Read from file
func ReadDataFromStorage(data io.Writer) error {
	rd, err := GetFileReader()
	if err != nil {
		return err
	}
	defer rd.Close()

	_, err = io.Copy(data, rd)
	if err != nil {
		log.Error().Err(err).Msg("Error reading data from storage file.")
		return err
	}
	log.Info().Msg("Data successfully read from storage file.")
	return nil
}
