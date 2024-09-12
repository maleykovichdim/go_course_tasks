// reverse index package (storage and functionality)
// here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted inde

//Format conversion: storage <-> protobuf, OVERHEAD! , because of PROTOBUF  (INCOMPATIBILITY crawler.Document and generated protobuf struct for it :(( )

package index

import (
	"search_engine/pkg/crawler"
	"search_engine/pkg/crawler/indexpb"
)

func ConvertDocuments(p *crawler.Document) *indexpb.Document {
	return &indexpb.Document{
		Id:    uint32(p.ID),
		Url:   p.URL,
		Title: p.Title,
		Body:  p.Body,
	}
}

func ConvertDocumentsR(p *indexpb.Document) *crawler.Document {
	return &crawler.Document{
		ID:    uint32(p.Id),
		URL:   p.Url,
		Title: p.Title,
		Body:  p.Body,
	}
}

func ConvertMapValue(value []uint32) *indexpb.Service_MapFieldEntry {
	a := indexpb.Service_MapFieldEntry{}
	a.Index = append(a.Index, value...)
	return &a
}

func (s *Service) ToProto() *indexpb.Service {
	pbs := indexpb.Service{}
	pbs.Counter = uint32(s.counter)
	pbs.Links = make([]*indexpb.Document, len(s.links))
	for i, link := range s.links {
		pbs.Links[i] = ConvertDocuments(&link)
	}
	pbs.Index = make(map[string]*indexpb.Service_MapFieldEntry)
	for key, value := range s.index {
		pbs.Index[key] = ConvertMapValue(value)
	}
	return &pbs
}

func (s *Service) FromProto(pbs *indexpb.Service) {

	s.counter = pbs.Counter
	s.links = make([]crawler.Document, len(pbs.Links))
	s.index = make(map[string][]uint32)

	// s.links =
	for i := range s.links {
		s.links[i] = *ConvertDocumentsR(pbs.Links[i])
	}

	for key, value := range pbs.Index {
		s.index[key] = (*value).Index
	}
}
