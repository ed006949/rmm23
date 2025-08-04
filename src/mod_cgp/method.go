package mod_cgp

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/structs"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
)

func (r *Token) command(payload string) (outbound []string, err error) {
	var (
		request  *http.Request
		response *http.Response
		interim  = *r.URL
		buffer   = new(bytes.Buffer)
	)

	interim.RawQuery = url.PathEscape(payload)
	switch request, err = http.NewRequest(http.MethodGet, interim.String(), nil); {
	case err != nil:
		return nil, err
	}

	// request.SetBasicAuth(r.Username, r.BindPassword)
	switch response, err = http.DefaultClient.Do(request); {
	case err != nil:
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	switch {
	case response.StatusCode != http.StatusOK:
		l.Z{l.E: mod_errors.EINVALRESPONSE, l.M: response.Body}.Error()

		return nil, mod_errors.EINVALRESPONSE
	}

	switch _, err = buffer.ReadFrom(response.Body); {
	case err != nil:
		return nil, err
	}

	for _, b := range strings.Fields(buffer.String()) {
		for _, d := range re_output_delim.Split(b, -1) {
			switch {
			case len(d) == 0:
				continue
			}

			outbound = append(outbound, d)
		}
	}

	return
}

// Command will execute only first command found.
func (r *Token) Command(inbound *Command) (outbound []string, err error) {
	var (
		o             = make(l.Z)
		payload       string
		emptyResponse bool // check if response must be empty
	)

	switch {
	case inbound != nil:
		o["server"] = r.Name
		payload += "command"
		payload += "="

		switch {
		case inbound.Domain_Administration != nil:
			switch {
			case inbound.GETDOMAINALIASES != nil:
				o["command"] = "GETDOMAINALIASES"
				o["domain"] = inbound.GETDOMAINALIASES.DomainName
				payload += inbound.GETDOMAINALIASES.compile()

			case inbound.UPDATEDOMAINSETTINGS != nil:
				o["command"] = "UPDATEDOMAINSETTINGS"
				o["domain"] = inbound.UPDATEDOMAINSETTINGS.DomainName
				emptyResponse = true
				payload += inbound.UPDATEDOMAINSETTINGS.compile()

				switch {
				case l.Run.DryRunValue():
					o["payloadLen"] = len(payload)
					payload = ""
				}

			default:
				return nil, mod_errors.EComSetDomAdm
			}

		case inbound.Domain_Set_Administration != nil:
			switch {
			case inbound.MAINDOMAINNAME != nil:
				o["command"] = "MAINDOMAINNAME"
				payload += inbound.MAINDOMAINNAME.compile()

			case inbound.LISTDOMAINS != nil:
				o["command"] = "LISTDOMAINS"
				payload += inbound.LISTDOMAINS.compile()

			default:
				return nil, mod_errors.EComSetDomSetAdm
			}

		default:
			return nil, mod_errors.EComSet
		}

		o[l.M] = "do"
		o.Debug()

		switch outbound, err = r.command(payload); {
		case err != nil:
			o[l.E] = err
			o.Error()

			return
		case emptyResponse && outbound != nil:
			o[l.E] = mod_errors.EINVALRESPONSE
			o.Warning()

			return outbound, mod_errors.EINVALRESPONSE
		default:
			o[l.M] = "done"
			o.Informational()

			return
		}

	default:
		return nil, mod_errors.ECom
	}
}

func (r *Command_Dictionary) compile() (outbound string) {
	outbound += "{"
	outbound += " "

	for a, b := range structs.Map(r) {
		outbound += a

		switch {
		case len(b.(string)) > 0:
			outbound += "="

			switch a {
			case "CAChain", "PrivateSecureKey", "SecureCertificate":
				outbound += "["
				outbound += b.(string)
				outbound += "]"
			default:
				outbound += b.(string)
			}
		}

		outbound += ";"
		outbound += " "
	}

	outbound += " "
	outbound += "}"

	return
}

func (r *UPDATEDOMAINSETTINGS) compile() (outbound string) {
	return outbound + "UPDATEDOMAINSETTINGS" + " " + r.DomainName + " " + r.NewSettings.compile()
}

func (r *GETDOMAINALIASES) compile() (outbound string) {
	return outbound + "GETDOMAINALIASES" + " " + r.DomainName
}

func (r *MAINDOMAINNAME) compile() (outbound string) { return outbound + "MAINDOMAINNAME" }
func (r *LISTDOMAINS) compile() (outbound string)    { return outbound + "LISTDOMAINS" }
