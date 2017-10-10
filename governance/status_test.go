// +build unit

package governance

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetContainerStatus(t *testing.T) {
	c1 := docker.APIContainers{
		ID:      "93f4354c97",
		Image:   "mycompany/x86/netspeed5:v2.5",
		Command: "/bin/sh -c 'python netspeed_edge.py --verbose --mqtt --policy'",
		Created: 1507728202,
		State:   "running",
		Status:  "Up Less than a second",
		Names:   []string{"/aaaa-netspeed5"},
		Labels: map[string]string{
			"network.bluehorizon.colonus.agreement_id": "aaaa",
			"network.bluehorizon.colonus.service_name": "netspeed5"},
	}
	c2 := docker.APIContainers{
		ID:      "73f4354c98",
		Image:   "mycompany/x86/test:v1.0",
		Command: "/bin/sh",
		Created: 1507728356,
		State:   "running",
		Status:  "Up 10 seconds",
		Names:   []string{"/aaaa-test"},
		Labels: map[string]string{
			"network.bluehorizon.colonus.agreement_id": "aaaa",
			"network.bluehorizon.colonus.service_name": "test"},
	}
	c3 := docker.APIContainers{
		ID:      "dec4ea9e5",
		Image:   "mycompany/x86/location:2.0.6",
		Command: "/bin/sh -c /start.sh",
		Created: 1507728184,
		State:   "running",
		Status:  "Up 18 second",
		Names:   []string{"/bbbb-location"},
		Labels: map[string]string{
			"network.bluehorizon.colonus.agreement_id": "bbbb",
			"network.bluehorizon.colonus.service_name": "location"},
	}
	c4 := docker.APIContainers{
		ID:      "dec4ea9e5",
		Image:   "mycompany/x86/gps:2.0.6",
		Command: "/bin/sh -c /start.sh",
		Created: 1507728188,
		State:   "running",
		Status:  "Up 18 second",
		Names:   []string{"/bluehorizon.network-microservices-gps_2.0.3_52df00-gps"},
		Labels: map[string]string{
			"network.bluehorizon.colonus.agreement_id":   "bluehorizon.network-microservices-gps_2.0.3_52df00",
			"network.bluehorizon.colonus.infrastructure": ""},
	}
	containers := []docker.APIContainers{c1, c2, c3, c4}

	agreementId := "aaaa"

	// test fail with a wrong deployment string
	deployment := "{\"services\":{\"netspeed5\":{st\":{\"image\":\"mycompany/x86/test:v1.0\"}}}"

	status, err := GetContainerStatus(deployment, agreementId, false, containers)

	assert.Error(t, err, "Error should be returned. ")

	// test workload containers succeeded
	deployment = "{\"services\":{\"netspeed5\":{\"image\":\"mycompany/x86/netspeed5:v2.5\",\"environment\":[\"FOO=bar\"]}, \"test\":{\"image\":\"mycompany/x86/test:v1.0\"}}}"
	exp_status := []ContainerStatus{ContainerStatus{Name: "/aaaa-netspeed5", Image: "mycompany/x86/netspeed5:v2.5", Created: 1507728202, State: "running"},
		{Name: "/aaaa-test", Image: "mycompany/x86/test:v1.0", Created: 1507728356, State: "running"}}

	status, err = GetContainerStatus(deployment, agreementId, false, containers)

	assert.Nil(t, err)
	assert.Equal(t, exp_status, status, "The elements should be the same.")

	// test with containers that does not have the agreement
	containers = []docker.APIContainers{c3, c4}
	deployment = "{\"services\":{\"netspeed5\":{\"image\":\"mycompany/x86/netspeed5:v2.5\",\"environment\":[\"FOO=bar\"]}, \"test\":{\"image\":\"mycompany/x86/test:v1.0\"}}}"
	exp_status = []ContainerStatus{ContainerStatus{Name: "netspeed5", Image: "mycompany/x86/netspeed5:v2.5", Created: 0, State: "not started"},
		{Name: "test", Image: "mycompany/x86/test:v1.0", Created: 0, State: "not started"}}

	status, err = GetContainerStatus(deployment, agreementId, false, containers)

	assert.Nil(t, err)
	assert.Equal(t, exp_status, status, "The elements should be the same.")

	// test with empty containers
	deployment = "{\"services\":{\"netspeed5\":{\"image\":\"mycompany/x86/netspeed5:v2.5\",\"environment\":[\"FOO=bar\"]}, \"test\":{\"image\":\"mycompany/x86/test:v1.0\"}}}"
	exp_status = []ContainerStatus{ContainerStatus{Name: "netspeed5", Image: "mycompany/x86/netspeed5:v2.5", Created: 0, State: "not started"},
		{Name: "test", Image: "mycompany/x86/test:v1.0", Created: 0, State: "not started"}}

	status, err = GetContainerStatus(deployment, agreementId, false, make([]docker.APIContainers, 0))

	assert.Nil(t, err)
	assert.Equal(t, exp_status, status, "The elements should be the same.")

	// test microservice containers succeeded
	key := "bluehorizon.network-microservices-gps_2.0.3_52df00"
	deployment = "{\"services\":{\"gps\":{\"image\":\"mycompany/x86/gps:2.0.6\"}}}"
	exp_status = []ContainerStatus{ContainerStatus{Name: "/bluehorizon.network-microservices-gps_2.0.3_52df00-gps", Image: "mycompany/x86/gps:2.0.6", Created: 1507728188, State: "running"}}
	containers = []docker.APIContainers{c1, c2, c3, c4}

	status, err = GetContainerStatus(deployment, key, true, containers)

	assert.Nil(t, err)
	assert.Equal(t, exp_status, status, "The elements should be the same.")

}
