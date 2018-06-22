package saxlike

import (
	"encoding/xml"
	"io"
)

//SAX-like XML Parser
type Parser struct {
	*xml.Decoder
	handler Handler
}

//Create a New Parser
func NewParser(reader io.Reader, handler Handler) *Parser {
	decoder := xml.NewDecoder(reader)
	return &Parser{decoder, handler}
}

//SetHTMLMode make Parser can parse invalid HTML
func (p *Parser) SetHTMLMode() {
	p.Strict = false
	p.AutoClose = xml.HTMLAutoClose
	p.Entity = xml.HTMLEntity
}

//Parse calls handler's methods
//when the parser encount a start-element,a end-element, a comment and so on.
func (p *Parser) Parse() (err error) {
	p.handler.StartDocument()

	for {
		token, err := p.Token()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			p.handler.StartElement(&t)
		case xml.EndElement:
			p.handler.EndElement(&t)
		case xml.CharData:
			p.handler.CharData(&t)
		case xml.Comment:
			p.handler.Comment(&t)
		case xml.ProcInst:
			p.handler.ProcInst(&t)
		case xml.Directive:
			p.handler.Directive(&t)
		default:
			panic("unknown xml token.")
		}
	}

	p.handler.EndDocument()
	return
}

//Create a parser and parse
func Parse(reader io.Reader, handler Handler, htmlMode bool) error {
	decoder := xml.NewDecoder(reader)
	parser := &Parser{decoder, handler}
	if htmlMode {
		parser.SetHTMLMode()
	}
	return parser.Parse()
}
