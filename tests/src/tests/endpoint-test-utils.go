package tests

import (
	"fmt"

	schema "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
)

var Exposures = [...]schema.EndpointExposure{schema.PublicEndpointExposure, schema.InternalEndpointExposure, schema.NoneEndpointExposure}

// Get a random exposure value
func GetRandomExposure() schema.EndpointExposure {
	return Exposures[GetRandomNumber(len(Exposures))-1]
}

//var Protocols = [...]schema.EndpointProtocol{schema.HTTPEndpointProtocol, schema.HTTPSEndpointProtocol, schema.WSEndpointProtocol, schema.WSSEndpointProtocol, schema.TCPEndpointProtocol, schema.UDPEndpointProtocol}
var Protocols = [...]schema.EndpointProtocol{schema.HTTPEndpointProtocol, schema.WSEndpointProtocol, schema.TCPEndpointProtocol, schema.UDPEndpointProtocol}

// Get a random protocol value
func GetRandomProtocol() schema.EndpointProtocol {
	return Protocols[GetRandomNumber(len(Protocols))-1]
}

// Create one or more endpoints in a schema structure
func CreateEndpoints() []schema.Endpoint {

	numEndpoints := GetRandomNumber(5)
	endpoints := make([]schema.Endpoint, numEndpoints)

	for i := 0; i < numEndpoints; i++ {

		endpoint := schema.Endpoint{}

		endpoint.Name = GetRandomString(GetRandomNumber(15)+5, false)
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d name  : %s", i, endpoint.Name))

		endpoint.TargetPort = GetRandomNumber(9999)
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d targetPort: %d", i, endpoint.TargetPort))

		if GetBinaryDecision() {
			endpoint.Exposure = GetRandomExposure()
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d exposure: %s", i, endpoint.Exposure))
		}

		if GetBinaryDecision() {
			endpoint.Protocol = GetRandomProtocol()
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d protocol: %s", i, endpoint.Protocol))
		}

		endpoint.Secure = GetBinaryDecision()
		LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d secure: %t", i, endpoint.Secure))

		if GetBinaryDecision() {
			endpoint.Path = "/Path_" + GetRandomString(GetRandomNumber(10)+3, false)
			LogInfoMessage(fmt.Sprintf("   ....... add endpoint %d path: %s", i, endpoint.Path))
		}

		endpoints[i] = endpoint

	}

	return endpoints
}
