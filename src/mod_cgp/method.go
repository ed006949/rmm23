package mod_cgp

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/structs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
		log.Warn().Err(mod_errors.EINVALRESPONSE).Msgf("response.StatusCode %v", response.StatusCode)

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
		logEvent      = new(zerolog.Event)
		payload       string
		emptyResponse bool // check if response must be empty
	)

	switch {
	case inbound != nil:
		logEvent.Str("server", r.Name)

		payload += "command"
		payload += "="

		switch {
		case inbound.Domain_Administration != nil:
			switch {
			case inbound.GETDOMAINALIASES != nil:
				logEvent.Str("command", "GETDOMAINALIASES")
				logEvent.Str("domain", inbound.GETDOMAINALIASES.DomainName)
				payload += inbound.GETDOMAINALIASES.compile()

			case inbound.UPDATEDOMAINSETTINGS != nil:
				logEvent.Str("command", "UPDATEDOMAINSETTINGS")
				logEvent.Str("domain", inbound.UPDATEDOMAINSETTINGS.DomainName)

				emptyResponse = true
				payload += inbound.UPDATEDOMAINSETTINGS.compile()

				switch {
				case l.Run.DryRunValue():
					logEvent.Int("payloadLen", len(payload))
					payload = ""
				}

			default:
				return nil, mod_errors.EComSetDomAdm
			}

		case inbound.Domain_Set_Administration != nil:
			switch {
			case inbound.MAINDOMAINNAME != nil:
				logEvent.Str("command", "MAINDOMAINNAME")

				payload += inbound.MAINDOMAINNAME.compile()

			case inbound.LISTDOMAINS != nil:
				logEvent.Str("command", "LISTDOMAINS")

				payload += inbound.LISTDOMAINS.compile()

			default:
				return nil, mod_errors.EComSetDomSetAdm
			}

		default:
			return nil, mod_errors.EComSet
		}

		log.Debug().Str("server", r.Name).Int("payloadLen", len(payload)).Msg("do")

		switch outbound, err = r.command(payload); {
		case err != nil:
			log.Error().Str("server", r.Name).Int("payloadLen", len(payload)).Err(err).Send()

			return
		case emptyResponse && outbound != nil:
			log.Warn().Str("server", r.Name).Int("payloadLen", len(payload)).Err(mod_errors.EINVALRESPONSE).Send()

			return outbound, mod_errors.EINVALRESPONSE
		default:
			log.Info().Str("server", r.Name).Int("payloadLen", len(payload)).Msg("done")

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
