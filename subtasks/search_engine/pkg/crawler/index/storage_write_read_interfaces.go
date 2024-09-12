// reverse index package (storage and functionality)
// here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted index

// implementation of interfaces to save/load into/from file  (PROTOBUF! solution)

package index

import (
	"io"
	"search_engine/pkg/crawler/indexpb"

	"google.golang.org/protobuf/proto"
)

func (s *Service) Write(p []byte) (n int, err error) {
	// Десериализация данных в структуру
	pbs := &indexpb.Service{}
	if err := proto.Unmarshal(p, pbs); err != nil {
		return 0, err
	}

	s.FromProto(pbs)
	return len(p), nil
}

func (s *Service) Read(p []byte) (n int, err error) {
	data, err := proto.Marshal(s.ToProto())
	if err != nil {
		return 0, err
	}

	n = copy(p, data)
	if n < len(data) {
		return n, nil
	}

	return n, io.EOF
}
