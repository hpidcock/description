// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package description

import (
	"time"

	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
)

type RemoteApplicationSerializationSuite struct {
	SliceSerializationSuite
}

var _ = gc.Suite(&RemoteApplicationSerializationSuite{})

func (s *RemoteApplicationSerializationSuite) SetUpTest(c *gc.C) {
	s.SliceSerializationSuite.SetUpTest(c)
	s.importName = "remote applications"
	s.sliceName = "remote-applications"
	s.importFunc = func(m map[string]interface{}) (interface{}, error) {
		return importRemoteApplications(m)
	}
	s.testFields = func(m map[string]interface{}) {
		m["remote-applications"] = []interface{}{}
	}
}

func minimalRemoteApplicationMap() map[interface{}]interface{} {
	m := minimalRemoteApplicationMapWithoutStatus()
	m["status"] = map[interface{}]interface{}{
		"version": 2,
		"status": map[interface{}]interface{}{
			"value":   "running",
			"message": "monkey & bear",
			"data": map[interface{}]interface{}{
				"after": "the curtain",
			},
			"updated":  "2016-01-28T11:50:00Z",
			"neverset": false,
		},
	}
	return m
}

func minimalRemoteApplicationMapWithoutStatus() map[interface{}]interface{} {
	return map[interface{}]interface{}{
		"name":              "civil-wars",
		"offer-uuid":        "offer-uuid",
		"url":               "http://a.url",
		"source-model-uuid": "abcd-1234",
		"is-consumer-proxy": true,
		"consume-version":   666,
		"endpoints": map[interface{}]interface{}{
			"version": 1,
			"endpoints": []interface{}{map[interface{}]interface{}{
				"name":      "lana",
				"role":      "provider",
				"interface": "mysql",
			}},
		},
		"spaces": map[interface{}]interface{}{
			"version": 1,
			"spaces": []interface{}{map[interface{}]interface{}{
				"cloud-type":  "gce",
				"name":        "private",
				"provider-id": "juju-space-private",
				"provider-attributes": map[interface{}]interface{}{
					"project": "gothic",
				},
				"subnets": map[interface{}]interface{}{
					"version": 3,
					"subnets": []interface{}{map[interface{}]interface{}{
						"cidr":                "2.3.4.0/24",
						"subnet-id":           "",
						"is-public":           false,
						"space-id":            "",
						"space-name":          "",
						"vlan-tag":            0,
						"provider-id":         "juju-subnet-1",
						"availability-zones":  []interface{}{"az1", "az2"},
						"provider-space-id":   "juju-space-private",
						"provider-network-id": "network-1",
					}},
				},
			}},
		},
		"bindings": map[interface{}]interface{}{
			"lana": "private",
		},
	}
}

func minimalRemoteApplication() *remoteApplication {
	a := minimalRemoteApplicationWithoutStatus()
	a.SetStatus(StatusArgs{
		Value:   "running",
		Message: "monkey & bear",
		Data: map[string]interface{}{
			"after": "the curtain",
		},
		Updated: time.Date(2016, 1, 28, 11, 50, 0, 0, time.UTC),
	})
	return a
}

func minimalRemoteApplicationWithoutStatus() *remoteApplication {
	a := newRemoteApplication(RemoteApplicationArgs{
		Tag:             names.NewApplicationTag("civil-wars"),
		OfferUUID:       "offer-uuid",
		URL:             "http://a.url",
		SourceModel:     names.NewModelTag("abcd-1234"),
		IsConsumerProxy: true,
		ConsumeVersion:  666,
		Bindings:        map[string]string{"lana": "private"},
	})
	a.AddEndpoint(RemoteEndpointArgs{
		Name:      "lana",
		Role:      "provider",
		Interface: "mysql",
	})
	space := a.AddSpace(RemoteSpaceArgs{
		CloudType:  "gce",
		Name:       "private",
		ProviderId: "juju-space-private",
		ProviderAttributes: map[string]interface{}{
			"project": "gothic",
		},
	})
	space.AddSubnet(SubnetArgs{
		CIDR:              "2.3.4.0/24",
		ProviderId:        "juju-subnet-1",
		AvailabilityZones: []string{"az1", "az2"},
		ProviderSpaceId:   "juju-space-private",
		ProviderNetworkId: "network-1",
	})
	return a
}

func (*RemoteApplicationSerializationSuite) TestNew(c *gc.C) {
	r := minimalRemoteApplication()
	c.Check(r.Tag(), gc.Equals, names.NewApplicationTag("civil-wars"))
	c.Check(r.Name(), gc.Equals, "civil-wars")
	c.Check(r.OfferUUID(), gc.Equals, "offer-uuid")
	c.Check(r.URL(), gc.Equals, "http://a.url")
	c.Check(r.SourceModelTag(), gc.Equals, names.NewModelTag("abcd-1234"))
	c.Check(r.IsConsumerProxy(), jc.IsTrue)
	c.Check(r.Status(), gc.DeepEquals, &status{
		Version: 2,
		StatusPoint_: StatusPoint_{
			Value_:   "running",
			Message_: "monkey & bear",
			Data_: map[string]interface{}{
				"after": "the curtain",
			},
			Updated_: time.Date(2016, 1, 28, 11, 50, 0, 0, time.UTC),
		},
	})
	ep := r.Endpoints()
	c.Assert(ep, gc.HasLen, 1)
	c.Check(ep[0].Name(), gc.Equals, "lana")
	sp := r.Spaces()
	c.Assert(sp, gc.HasLen, 1)
	c.Check(sp[0].Name(), gc.Equals, "private")
	c.Check(r.Bindings(), gc.DeepEquals, map[string]string{"lana": "private"})
}

func (*RemoteApplicationSerializationSuite) TestNewWithoutStatus(c *gc.C) {
	r := minimalRemoteApplicationWithoutStatus()
	c.Check(r.Tag(), gc.Equals, names.NewApplicationTag("civil-wars"))
	c.Check(r.Name(), gc.Equals, "civil-wars")
	c.Check(r.OfferUUID(), gc.Equals, "offer-uuid")
	c.Check(r.URL(), gc.Equals, "http://a.url")
	c.Check(r.SourceModelTag(), gc.Equals, names.NewModelTag("abcd-1234"))
	c.Check(r.IsConsumerProxy(), jc.IsTrue)
	c.Check(r.Status(), gc.IsNil)
	ep := r.Endpoints()
	c.Assert(ep, gc.HasLen, 1)
	c.Check(ep[0].Name(), gc.Equals, "lana")
	sp := r.Spaces()
	c.Assert(sp, gc.HasLen, 1)
	c.Check(sp[0].Name(), gc.Equals, "private")
	c.Check(r.Bindings(), gc.DeepEquals, map[string]string{"lana": "private"})
}

func (*RemoteApplicationSerializationSuite) TestBadSchema1(c *gc.C) {
	container := map[string]interface{}{
		"version":             1,
		"remote-applications": []interface{}{1234},
	}
	_, err := importRemoteApplications(container)
	c.Assert(err, gc.ErrorMatches, `remote applications version schema check failed: remote-applications\[0\]: expected map, got int\(1234\)`)
}

func (*RemoteApplicationSerializationSuite) TestBadSchema2(c *gc.C) {
	m := minimalRemoteApplicationMap()
	m["is-consumer-proxy"] = "blah"
	container := map[string]interface{}{
		"version":             1,
		"remote-applications": []interface{}{m},
	}
	_, err := importRemoteApplications(container)
	c.Assert(err, gc.ErrorMatches, `remote application 0 v1 schema check failed: is-consumer-proxy: expected bool, got string\("blah"\)`)
}

func (s *RemoteApplicationSerializationSuite) TestBadEndpoints(c *gc.C) {
	m := minimalRemoteApplicationMap()
	m["endpoints"] = map[interface{}]interface{}{
		"version": 1,
		"bishop":  "otter-trouserpress",
	}
	container := map[string]interface{}{
		"version":             1,
		"remote-applications": []interface{}{m},
	}
	_, err := importRemoteApplications(container)
	c.Assert(err, gc.ErrorMatches, `remote application 0: remote endpoints version schema check failed: endpoints: expected list, got nothing`)
}

func (*RemoteApplicationSerializationSuite) TestMinimalMatches(c *gc.C) {
	bytes, err := yaml.Marshal(minimalRemoteApplication())
	c.Assert(err, jc.ErrorIsNil)

	var source map[interface{}]interface{}
	err = yaml.Unmarshal(bytes, &source)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(source, jc.DeepEquals, minimalRemoteApplicationMap())
}

func (*RemoteApplicationSerializationSuite) TestMinimalMatchesWithoutStatus(c *gc.C) {
	bytes, err := yaml.Marshal(minimalRemoteApplicationWithoutStatus())
	c.Assert(err, jc.ErrorIsNil)

	var source map[interface{}]interface{}
	err = yaml.Unmarshal(bytes, &source)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(source, jc.DeepEquals, minimalRemoteApplicationMapWithoutStatus())
}

func (s *RemoteApplicationSerializationSuite) TestRoundTripVersion1(c *gc.C) {
	rIn := minimalRemoteApplication()
	rIn.ConsumeVersion_ = 0
	rOut := s.exportImport(c, 1, rIn)
	rIn.ConsumeVersion_ = 1
	c.Assert(rOut, jc.DeepEquals, rIn)
}

func (s *RemoteApplicationSerializationSuite) TestRoundTripVersion2(c *gc.C) {
	rIn := minimalRemoteApplication()
	rIn.ConsumeVersion_ = 0
	rIn.Macaroon_ = "mac"
	rOut := s.exportImport(c, 2, rIn)
	rIn.ConsumeVersion_ = 1
	c.Assert(rOut, jc.DeepEquals, rIn)
}

func (s *RemoteApplicationSerializationSuite) TestRoundTripVersion3(c *gc.C) {
	rIn := minimalRemoteApplication()
	rIn.Macaroon_ = "mac"
	rOut := s.exportImport(c, 3, rIn)
	c.Assert(rOut, jc.DeepEquals, rIn)
}

func (s *RemoteApplicationSerializationSuite) TestRoundTripWithoutStatus(c *gc.C) {
	rIn := minimalRemoteApplicationWithoutStatus()
	rOut := s.exportImport(c, 3, rIn)
	c.Assert(rOut, jc.DeepEquals, rIn)
}

func (s *RemoteApplicationSerializationSuite) exportImport(
	c *gc.C, version int, app *remoteApplication,
) *remoteApplication {
	applicationsIn := &remoteApplications{
		Version:            version,
		RemoteApplications: []*remoteApplication{app},
	}
	bytes, err := yaml.Marshal(applicationsIn)
	c.Assert(err, jc.ErrorIsNil)

	var source map[string]interface{}
	err = yaml.Unmarshal(bytes, &source)
	c.Assert(err, jc.ErrorIsNil)

	applicationsOut, err := importRemoteApplications(source)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(applicationsOut, gc.HasLen, 1)
	return applicationsOut[0]
}
