package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadConf(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	y := &Config{}
	err = yaml.Unmarshal([]byte(buf), &y)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return y, err
}

func CheckConf(conf *Config) (*Config, error) {
	existChannels := make(map[string]bool)
	for _, v := range conf.Channels {
		if v.Name == "" {
			return nil, fmt.Errorf("Channels: Channel name cannot be empty.")
		}
		if existChannels[v.Name] {
			return nil, fmt.Errorf("Channels: Duplicate channel name '%v'.", v.Name)
		} else {
			existChannels[v.Name] = true
		}
		reg, err := regexp.Compile("[^a-z]+")
		if err != nil {
			return nil, err
		}
		proccessedName := reg.ReplaceAllString(v.Name, "")
		if v.Name != proccessedName {
			return nil, fmt.Errorf("Channels: Channel name '%v' should only comprise lowercase alphabet.", v.Name)
		}
		if err := checkPolicy("Channel", v.Policies); err != nil {
			return nil, fmt.Errorf("Channels '%v': %w", v.Name, err)
		}
		if err := checkPolicyName(v.Name, v.Policies, conf.Organizations); err != nil {
			return nil, fmt.Errorf("Channels '%v': %w", v.Name, err)
		}
	}

	existChaincodes := make(map[string]bool)
	for _, v := range conf.Chaincodes {
		if v.Name == "" {
			return nil, fmt.Errorf("Chaincode name cannot be empty.")
		}
		if existChaincodes[v.Name] {
			return nil, fmt.Errorf("Chaincode: Duplicate chaincode name '%v'.", v.Name)
		} else {
			existChaincodes[v.Name] = true
		}
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return nil, err
		}
		proccessedName := reg.ReplaceAllString(v.Name, "")
		if v.Name != proccessedName {
			return nil, fmt.Errorf("Chaincode: Chaincode name '%v' should only comprise number, alphabet.", v.Name)
		}
		if v.Language == "" { //RE
			return nil, fmt.Errorf("Chaincode language cannot be empty.")
		}
		if v.Path == "" {
			return nil, fmt.Errorf("Chaincode path cannot be empty.")
		}
	}

	existOrganizations := make(map[string]bool)
	for _, v := range conf.Organizations {
		if v.Name == "" {
			return nil, fmt.Errorf("Organizations: Organization name cannot be empty.")
		}
		if existOrganizations[v.Name] {
			return nil, fmt.Errorf("Organizations: Duplicate organization name '%v'.", v.Name)
		} else {
			existOrganizations[v.Name] = true
		}
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return nil, err
		}
		proccessedName := reg.ReplaceAllString(v.Name, "")
		if v.Name != proccessedName {
			return nil, fmt.Errorf("Organizations: Organization name '%v' should only comprise number, alphabet.", v.Name)
		}
		if v.Type == "" {
			return nil, fmt.Errorf("Organization '%v': Organization type cannot be empty.", v.Name)
		}
		if v.Type != "peerOrg" && v.Type != "orderOrg" {
			return nil, fmt.Errorf("Organization '%v': Organization type cannot be '%v'. Organization type should be one of 'peerOrg' or 'orderOrg'.", v.Name, v.Type)
		}
		if v.Type == "orderOrg" {
			// if v.ConsensusType == "" {
			// 	return fmt.Errorf("Organizations: Consensus type of order organization cannot be empty.")
			// }
			// if v.ConsensusType != "Raft" && v.ConsensusType != "Solo" && v.ConsensusType != "Kafa"  {
			// 	return fmt.Errorf("Organizations: Consensus type of order organization should be one of 'Raft', Solo', and 'Kafka'.")
			// }
			if v.BatchTimeout == 0 {
				return nil, fmt.Errorf("Organization '%v': BatchTimeout cannot be empty or '0'.", v.Name)
			}
			if v.BatchSize.AbsoluteMaxBytes == 0 {
				return nil, fmt.Errorf("Organization '%v': AbsoluteMaxBytes cannot be empty or '0'.", v.Name)
			}
			if v.BatchSize.MaxMessageCount == 0 {
				return nil, fmt.Errorf("Organization '%v': MaxMessageCount cannot be empty or '0'.", v.Name)
			}
			if v.BatchSize.PreferredMaxBytes == 0 {
				return nil, fmt.Errorf("Organization '%v': PreferredMaxBytes cannot be empty or '0'.", v.Name)
			}
			if err := checkPeer(v.Type, "orderingPeer", v.Peers); err != nil {
				return nil, fmt.Errorf("Organizations '%v': %w", v.Name, err)
			}
			if v.Peers == nil {
				return nil, fmt.Errorf("Organization should have at least one peer.")
			}
		}
		if v.Type == "peerOrg" {
			if err := checkPeer(v.Type, "committingPeer", v.CommittingPeers); err != nil {
				return nil, fmt.Errorf("Organizations '%v': %w", v.Name, err)
			}
			if err := checkPeer(v.Type, "endorsingPeers", v.EndorsingPeers); err != nil {
				return nil, fmt.Errorf("Organizations '%v': %w", v.Name, err)
			}
			if v.CommittingPeers == nil && v.EndorsingPeers == nil {
				return nil, fmt.Errorf("Organization should have at least one peer.")
			}
			if err := checkPolicy("Organization", v.Policies); err != nil {
				return nil, fmt.Errorf("Organizations '%v': %w", v.Name, err)
			}
			if err := checkPolicyRole(v.Policies, v.Name); err != nil {
				return nil, fmt.Errorf("Organizations '%v': %w", v.Name, err)
			}
		}
	}

	existCertificateAuthorities := make(map[string]bool)
	existCAOrganizations := make(map[string]bool)
	for _, v := range conf.CertificateAuthorities {
		if v.Name == "" {
			return nil, fmt.Errorf("CertificateAuthorities: Certificate authority name cannot be empty.")
		}
		if existCertificateAuthorities[v.Name] {
			return nil, fmt.Errorf("CertificateAuthorities: Duplicate Certificate Authority name '%v'.", v.Name)
		} else {
			existCertificateAuthorities[v.Name] = true
		}
		if domainArr := strings.Split(v.Name, "."); len(domainArr) != 4 {
			return nil, fmt.Errorf("CertificateAuthorities: CA name '%v' should contains 4 sections. It should be like <CA>.<ORG>.<SLD>.<TLD>. ex. \"ca.org1.example.com.\"", v.Name)
		}
		reg, err := regexp.Compile("[^a-z0-9.]+")
		if err != nil {
			return nil, err
		}
		proccessedName := reg.ReplaceAllString(v.Name, "")
		if v.Name != proccessedName {
			return nil, fmt.Errorf("CertificateAuthorities: CA name '%v' should only comprise number, alphabet and dot.", v.Name)
		}
		if v.Address == "" {
			return nil, fmt.Errorf("CertificateAuthorities: The address of certificate authority, '%v', cannot be empty.", v.Name)
		}
		ip := strings.Split(v.Address, ".")
		if len(ip) != 4 {
			return nil, fmt.Errorf("CertificateAuthorities: The IP address, '%v', of certificate authority, '%v', should be like '127.0.0.1'.", v.Address, v.Name)
		}
		for _, value := range ip {
			ipNumber, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("CertificateAuthorities: Invalid IP address, '%v', of certificate authority, '%v'.", v.Address, v.Name)
			}
			if ipNumber > 255 || ipNumber < 0 {
				return nil, fmt.Errorf("CertificateAuthorities: The range of IP address, '%v', of certificate authority, '%v', is not between '0.0.0.0' and '255.255.255.255'.", v.Address, v.Name)
			}
		}
		if v.Port == "" {
			return nil, fmt.Errorf("CertificateAuthorities: The port of certificate authority, '%v', cannot be empty.", v.Name)
		}
		port, err := strconv.Atoi(v.Port)
		if err != nil {
			return nil, fmt.Errorf("CertificateAuthorities: The port, '%v', of certificate authority, '%v', is not a number.", v.Port, v.Name)
		}
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("CertificateAuthorities: The port, '%v', of certificate authority, '%v', is not between 1 and 65535.", v.Port, v.Name)
		}
		for _, org := range v.Organizations {
			if !existOrganizations[org] {
				return nil, fmt.Errorf("CertificateAuthorities: Organization, '%v', in certificate authority, '%v', had not founded in 'Organizations' section.", org, v.Name)
			}
			if existCAOrganizations[org] {
				return nil, fmt.Errorf("CertificateAuthorities: Duplicate organization name, '%v', in certificate authority, '%v'.", org, v.Name)
			} else {
				existCAOrganizations[org] = true
			}
		}
	}

	return conf, nil
}

func checkPeer(orgType string, peerType string, peers []Peer) error {
	existPeers := make(map[string]bool)
	for _, peer := range peers {
		if peer.Name == "" {
			return fmt.Errorf("Peer name cannot be empty.")
		}
		if existPeers[peer.Name] {
			return fmt.Errorf("Duplicate peer name '%v'.", peer.Name)
		} else {
			existPeers[peer.Name] = true
		}
		reg, err := regexp.Compile("[^a-z0-9.]+")
		if err != nil {
			return err
		}
		proccessedName := reg.ReplaceAllString(peer.Name, "")
		if peer.Name != proccessedName {
			return fmt.Errorf("Peer name '%v' should only comprise number, alphabet and dot.", peer.Name)
		}
		if peer.Address == "" {
			return fmt.Errorf("Peer address cannot be empty.")
		}
		ip := strings.Split(peer.Address, ".")
		if len(ip) != 4 {
			return fmt.Errorf("The address of peer '%v' should be like '127.0.0.1'.", peer.Address)
		}
		for _, v := range ip {
			ipNumber, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("Invalid IP address '%v'.", peer.Address)
			}
			if ipNumber > 255 || ipNumber < 0 {
				return fmt.Errorf("The range of IP address (peer address) '%v' is not between '0.0.0.0' and '255.255.255.255'.", peer.Address)
			}
		}
		if peer.Port == "" {
			return fmt.Errorf("Peer port cannot be empty.")
		}
		port, err := strconv.Atoi(peer.Port)
		if err != nil {
			return fmt.Errorf("The port of peer '%v' is not a number.", peer.Port)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("The port of peer '%v' is not between 1 and 65535.", peer.Port)
		}
		if orgType == "peerOrg" {
			if domainArr := strings.Split(peer.Name, "."); len(domainArr) != 4 {
				return fmt.Errorf("Peer name '%v' should contains 4 sections. It should be like <PEER>.<ORG>.<SLD>.<TLD>. ex. \"peer0.org1.example.com.\"", peer.Name)
			}
			if peer.DBPort == "" {
				return fmt.Errorf("DB port cannot be empty.")
			}
			dbPort, err := strconv.Atoi(peer.DBPort)
			if err != nil {
				return fmt.Errorf("The port of database '%v' is not a number.", peer.DBPort)
			}
			if dbPort < 1 || dbPort > 65535 {
				return fmt.Errorf("The port of database '%v' is not between 1 and 65535.", peer.DBPort)
			}
			if port == dbPort {
				return fmt.Errorf("Duplicate port '%v'.", peer.Port)
			}
		}
	}
	return nil
}

func checkPolicy(level string, policies []Policy) error {
	delimiter := func(r rune) bool {
		return r == '(' || r == ')' || r == ' '
	}
	if policies == nil {
		return fmt.Errorf("%s should have at least one policy.", level)
	}
	existPolicies := make(map[string]bool)
	for _, policy := range policies {
		if policy.Name == "" {
			return fmt.Errorf("Policy name cannot be empty.")
		}
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return err
		}
		proccessedName := reg.ReplaceAllString(policy.Name, "")
		if policy.Name != proccessedName {
			return fmt.Errorf("Policy name '%v' should only comprise number and alphabet.", policy.Name)
		}
		if existPolicies[policy.Name] {
			return fmt.Errorf("Duplicate policy name '%v'.", policy.Name)
		} else {
			existPolicies[policy.Name] = true
		}
		if policy.Policy == "" {
			return fmt.Errorf("Policy cannot be empty.")
		}
		policyUpperCase := strings.ToUpper(policy.Policy)
		if level == "Organization" {
			policyRule := strings.FieldsFunc(policyUpperCase, delimiter)
			if !(policyRule[0] == "OUTOF" || policyRule[0] == "AND" || policyRule[0] == "OR") {
				return fmt.Errorf("The policy is '%v' while the beginning of organization policy should be 'OutOf', 'OR', or 'AND'.", policy.Policy)
			}
		} else if level == "Channel" {
			policyRule := strings.FieldsFunc(policyUpperCase, delimiter)
			if !(policyRule[0] == "OUTOF" || policyRule[0] == "AND" || policyRule[0] == "OR" || policyRule[0] == "ANY" || policyRule[0] == "MAJORITY" || policyRule[0] == "ALL") {
				return fmt.Errorf("The channel policy is '%v'. The beginning of channel policy should be 'OutOf', 'OR', 'AND', 'ANY, 'MAJORITY', or 'ALL'.", policy.Policy)
			}
		} else {
			return fmt.Errorf("'level' should be 'Organization' or 'Channel'.")
		}
	}
	return nil
}

func checkPolicyRole(policies []Policy, orgName string) error {
	delimiter := func(r rune) bool {
		return r == '(' || r == ')' || r == ',' || r == ' '
	}
	for _, policy := range policies {
		roles := strings.FieldsFunc(policy.Policy, delimiter)
		startIndex := 1
		if strings.ToUpper(roles[0]) == "OUTOF" {
			startIndex = 2
			roleNumber, err := strconv.Atoi(roles[1])
			if err != nil {
				return fmt.Errorf("'%v' used by rule 'OutOf' should be a number.", roles[1])
			}
			if roleNumber > len(roles)-2 {
				return fmt.Errorf("The number of role '%v' used by policy '%v' with rule 'OutOf' is larger than the number of role '%v'", roles[1], policy.Name, len(roles)-2)
			}
			if roleNumber < 1 {
				return fmt.Errorf("The number of role '%v' used by policy '%v' with rule 'OutOf' is less than 1. It should be larger or equal to '1'", roles[1], policy.Name)
			}
		}
		for i := startIndex; i < len(roles); i++ {
			org := strings.Split(roles[i], ".")[0]
			if org != orgName {
				return fmt.Errorf("'%v' is mismatch with '%v'. The organization name of policy is mismatch with organization name.", org, orgName)
			}
		}
	}
	return nil
}

func checkPolicyName(channelName string, policies []Policy, orgs []Organization) error {
	delimiter := func(r rune) bool {
		return r == '(' || r == ')' || r == ',' || r == ' '
	}
	for _, policy := range policies {
		channelPolicy := strings.FieldsFunc(policy.Policy, delimiter)
		channelPolicyRule := strings.ToUpper(channelPolicy[0])
		if channelPolicyRule == "OUTOF" || channelPolicyRule == "AND" || channelPolicyRule == "OR" {
			startIndex := 1
			if channelPolicyRule == "OUTOF" {
				startIndex = 2
				roleNumber, err := strconv.Atoi(channelPolicy[1])
				if err != nil {
					return fmt.Errorf("The first element '%v' used by rule 'OutOf' is not a number.", channelPolicy[1])
				}
				if roleNumber > len(channelPolicy)-2 {
					return fmt.Errorf("The number of role '%v' used by policy '%v' with rule 'OutOf' is larger than the number of role '%v'. It should less than or equal to '%v'", channelPolicy[1], policy.Name, len(channelPolicy)-2, len(channelPolicy)-2)
				}
				if roleNumber < 1 {
					return fmt.Errorf("The number of role '%v' used by policy '%v' with rule 'OutOf' is less than 1. It should be larger or equal to '1'", channelPolicy[1], policy.Name)
				}
			}
			orgsInChannel := []string{}
			for _, org := range orgs {
				for _, channel := range org.Channels {
					if channel == channelName {
						orgsInChannel = append(orgsInChannel, org.Name)
						break
					}
				}
			}
			exists := make(map[string]bool)
			for _, org := range orgsInChannel {
				exists[org] = true
			}
			for i := startIndex; i < len(channelPolicy); i++ {
				org := strings.Split(channelPolicy[i], ".")[0]
				if !exists[org] {
					return fmt.Errorf("The organization name, '%v', in channel policy is mismatch with organization which connected to channel.", org)
				}
			}
		} else if channelPolicyRule == "ALL" || channelPolicyRule == "MAJORITY" || channelPolicyRule == "ANY" {
			if len(channelPolicy) != 2 {
				return fmt.Errorf("The Channel Policy should be comprised of '%v' and 'one' organization policy name.", channelPolicyRule)
			}
			isMatch := false
			channelPolicyName := channelPolicy[1]
			orgsInChannel := []Organization{}
			for _, org := range orgs {
				for _, channel := range org.Channels {
					if channel == channelName {
						orgsInChannel = append(orgsInChannel, org)
						break
					}
				}
			}
			for _, org := range orgsInChannel {
				for _, orgPolicy := range org.Policies {
					if channelPolicyName == orgPolicy.Name {
						isMatch = true
						break
					}
				}
			}
			if !isMatch {
				return fmt.Errorf("The channel policy, '%v', is not match to any organization policies.", policy.Name)
			}
		} else {
			return fmt.Errorf("'channelPolicyRule' should be 'ANY', 'MAJORITY', 'ALL', 'AND', 'OR', or 'OutOf'.")
		}
	}
	return nil
}
