package configtx

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
	"strconv"
	"strings"
)

func ConvertConf(conf *config.Config) (Config, error) {
	capability := Capability{V2_0: true}
	capabilities := Capabilities{
		Channel:     capability,
		Orderer:     capability,
		Application: capability,
	}

	app := Application{}

	orgs := []Organization{}
	peerOrgs := []Organization{}
	orderOrgs := []Organization{}

	orderer := Orderer{
		OrdererType: "etcdraft",
		Policies: Policies{
			Readers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			Writers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			Admins: Policy{
				Type: "ImplicitMeta",
				Rule: "MAJORITY Admins",
			},
			BlockValidation: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
		},
		Organizations: []string{""},
	}

	channel := Channel{
		Policies: Policies{
			Readers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			Writers: Policy{
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			Admins: Policy{
				Type: "ImplicitMeta",
				Rule: "MAJORITY Admins",
			},
		},
		Capabilities: capabilities.Channel,
	}

	for _, channel := range conf.Channels {
		policies, err := extractPolicies(channel.Policies)
		if err != nil {
			return Config{}, fmt.Errorf("error occur in extracting policies: %v", err)
		}
		app.Policies = policies
		app.Capabilities = capabilities.Application
		app.Organizations = []Organization{{}}
	}

	for _, org := range conf.Organizations {
		policies, err := extractPolicies(org.Policies)
		if err != nil {
			return Config{}, fmt.Errorf("error occur in extracting policies: %v", err)
		}
		if org.Type == "orderOrg" {
			ordererEndpoints := []string{}
			consenters := []Consenter{}
			for _, peer := range org.Peers {
				ordererEndpoints = append(ordererEndpoints, peer.Name+":"+peer.Port)
				orghost, err := utils.ExtractHost(ordererEndpoints[0], 2)
				if err != nil {
					return Config{}, fmt.Errorf("error occur in extracting host: %v", err)
				}
				TLSCertPath := "../organizations/ordererOrganizations/" + orghost + "/orderers/" + peer.Name + "/tls/server.crt"
				peerPort, _ := strconv.Atoi(peer.Port)
				consenters = append(consenters, Consenter{
					Host:          peer.Name,
					Port:          peerPort,
					ClientTLSCert: TLSCertPath,
					ServerTLSCert: TLSCertPath,
				})
			}
			orghost, err := utils.ExtractHost(ordererEndpoints[0], 2)
			if err != nil {
				return Config{}, fmt.Errorf("error occur in extracting host: %v", err)
			}
			orderOrgs = append(orderOrgs, Organization{
				Name:             org.Name,
				ID:               org.Name,
				MSPDir:           "../organizations/ordererOrganizations/" + orghost + "/msp",
				Policies:         policies,
				OrdererEndpoints: ordererEndpoints,
			})

			orderer.Addresses = ordererEndpoints
			orderer.EtcdRaft.Consenters = consenters
			orderer.BatchTimeout = strconv.Itoa(org.BatchTimeout) + "s"
			orderer.BatchSize = BatchSize{
				MaxMessageCount:   org.BatchSize.MaxMessageCount,
				AbsoluteMaxBytes:  strconv.Itoa(org.BatchSize.AbsoluteMaxBytes) + "MB",
				PreferredMaxBytes: strconv.Itoa(org.BatchSize.PreferredMaxBytes) + "KB",
			}
		} else {
			endpoints := []string{}
			for _, peer := range org.CommittingPeers {
				endpoints = append(endpoints, peer.Name+":"+peer.Port)
			}
			for _, peer := range org.EndorsingPeers {
				endpoints = append(endpoints, peer.Name+":"+peer.Port)
			}
			orghost, err := utils.ExtractHost(endpoints[0], 2)
			if err != nil {
				return Config{}, fmt.Errorf("error occur in extracting host: %v", err)
			}
			peerOrgs = append(peerOrgs, Organization{
				Name:     org.Name,
				ID:       org.Name,
				MSPDir:   "../organizations/peerOrganizations/" + orghost + "/msp",
				Policies: policies,
			})
		}
	}
	orgs = append(orgs, peerOrgs...)
	orgs = append(orgs, orderOrgs...)
	profiles := Profiles{
		OrdererGenesis: OrdererGenesis{
			Policies:     channel.Policies,
			Capabilities: channel.Capabilities,
			Orderer: OrdererOrdererGenesis{
				OrdererType:   orderer.OrdererType,
				Addresses:     orderer.Addresses,
				EtcdRaft:      orderer.EtcdRaft,
				BatchTimeout:  orderer.BatchTimeout,
				BatchSize:     orderer.BatchSize,
				Policies:      orderer.Policies,
				Organizations: orderOrgs,
				Capabilities:  capabilities.Orderer,
			},
			Consortiums: Consortiums{
				SampleConsortium: SampleConsortium{
					Organizations: peerOrgs,
				},
			},
		},
		OrgsChannel: OrgsChannel{
			Consortium:   "SampleConsortium",
			Policies:     channel.Policies,
			Capabilities: channel.Capabilities,
			Application: Application{
				Organizations: peerOrgs,
				Policies:      app.Policies,
				Capabilities:  capabilities.Application,
			},
		},
	}

	configtx := Config{
		Organizations: orgs,
		Capabilities:  capabilities,
		Application:   app,
		Orderer:       orderer,
		Channel:       channel,
		Profiles:      profiles,
	}
	return configtx, nil
}

func extractPolicies(rpolicies []config.Policy) (Policies, error) {
	delimiter := func(r rune) bool {
		return r == '(' || r == ')' || r == ','
	}
	wpolicies := Policies{}
	for _, policy := range rpolicies {
		ptype := "ImplicitMeta"
		if strings.HasPrefix(strings.ToUpper(policy.Policy), "OR") || strings.HasPrefix(strings.ToUpper(policy.Policy), "AND") || strings.HasPrefix(strings.ToUpper(policy.Policy), "OUTOF") {
			ptype = "Signature"
		}
		policyRule := policy.Policy
		if ptype == "Signature" {
			policies := strings.FieldsFunc(policy.Policy, delimiter)
			for i, v := range policies {
				v = strings.TrimSpace(v)
				if i == 0 {
					policyRule = strings.ToUpper(v) + "("
				} else if i == 1 && strings.ToUpper(policies[0]) == "OUTOF" {
					policyRule = policyRule + v + ", "
				} else if i == len(policies)-1 {
					policyRule = policyRule + "'" + v + "')"
				} else {
					policyRule = policyRule + "'" + v + "', "
				}
			}
		}
		// else {
		// 	policyRule = "\"" + policyRule + "\""
		// }
		p := Policy{
			Type: ptype,
			Rule: policyRule,
		}
		if policy.Name == "Readers" {
			wpolicies.Readers = p
		} else if policy.Name == "Writers" {
			wpolicies.Writers = p
		} else if policy.Name == "Admins" {
			wpolicies.Admins = p
		} else if policy.Name == "Endorsement" {
			wpolicies.Endorsement = p
		} else if policy.Name == "LifecycleEndorsement" {
			wpolicies.LifecycleEndorsement = p
		} else if policy.Name == "BlockValidation" {
			wpolicies.BlockValidation = p
		} else {
			return wpolicies, fmt.Errorf("invalid policy when extracting policy")
		}
	}
	return wpolicies, nil
}
