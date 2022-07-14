package serverconfig

import (
	"bufio"
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func MakeDirsAndWriteConf(conf *config.Config) error {
	for _, ca := range conf.CertificateAuthorities {
		port, err := strconv.Atoi(ca.Port)
		if err != nil {
			return err
		}
		caHosts := strings.Split(ca.Name, ".")
		caHost := ""
		for i, v := range caHosts {
			if i == len(caHosts)-1 {
				caHost = caHost + v
			} else if i != 0 {
				caHost = caHost + v + "."
			}
		}
		for _, org := range ca.Organizations {
			err := os.RemoveAll("./organizations/fabric-ca/" + org)
			if err != nil {
				return err
			}
			err = os.Mkdir("./organizations/fabric-ca/"+org, 0750)
			if err != nil {
				return err
			}
			serverConfig := Config{
				Version:      "1.2.0",
				Port:         port,
				Debug:        false,
				CrlsizeLimit: 512000,
				TLS: TLS{
					Enabled: true,
					ClientAuth: ClientAuth{
						Type: "noclientcert",
					},
				},
				CA: CA{
					Name: ca.Name,
				},
				CRL: CRL{
					Expiry: "24h",
				},
				Registry: Registry{
					MaxEnrollments: -1,
					Identities: []Identity{
						{
							Name: "admin",
							Pass: "adminpw",
							Type: "client",
							Attrs: Attrs{
								HfRegistrarRoles:         "*",
								HfRegistrarDelegateRoles: "*",
								HfRevoker:                true,
								HfIntermediateCA:         true,
								HfGenCRL:                 true,
								HfRegistrarAttributes:    "*",
								HfAffiliationMgr:         true,
							},
						},
					},
				},
				DB: DB{
					Type:       "sqlite3",
					Datasource: "fabric-ca-server.db",
					TLS:        TLSDB{Enabled: false},
				},
				LDAP: LDAP{
					Enabled: false,
					URL:     "ldap://<adminDN>:<adminPassword>@<host>:<port>/<base>",
					Attribute: Attribute{
						Names: []string{}, //[]string{"uid", "member"},
					},
				},
				Affiliations: map[string]string{
					"org": "department1",
				},
				Signing: Signing{
					Default: Default{
						Usage:  []string{"digital signature"},
						Expiry: "8760h",
					},
					Profiles: Profile{
						CA: CAProfile{
							Usage:  []string{"cert sign", "crl sign"},
							Expiry: "43800h",
							CAConstraint: Caconstraint{
								ISCA:       true,
								MaxPathLen: 0,
							},
						},
						TLS: Default{
							Usage: []string{
								"signing",
								"key encipherment",
								"server auth",
								"client auth",
								"key agreement",
							},
							Expiry: "8760h",
						},
					},
				},
				CSR: CSR{
					CN: ca.Name,
					Names: []NamesCSR{
						{
							C:  "US",
							ST: "North Carolina",
							L:  "Raleigh",
							O:  caHost,
						},
					},
					Hosts: []string{"localhost", caHost},
					CA: CACSR{
						Expiry:     "131400h",
						PathLength: 1,
					},
				},
				BCCSP: BCCSP{
					Default: "SW",
					SW: SW{
						Hash:     "SHA2",
						Security: 256,
						FileKeyStore: FileKeyStore{
							Keystore: "msp/keystore",
						},
					},
				},
			}
			err = utils.WriteYaml(serverConfig, "./organizations/fabric-ca/"+org+"/fabric-ca-server-config.yaml")
			if err != nil {
				log.Fatal(err)
			}
			err = removeEmptyValue("./organizations/fabric-ca/" + org + "/fabric-ca-server-config.yaml")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func removeEmptyValue(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	rowNumber, targetRowIndex := 1, 0
	targetRow := []int{6, 7, 10, 13, 14, 15, 38, 40, 41, 46, 48, 49, 52, 54, 99, 100, 103, 104, 106, 107, 108, 110, 112, 113}

	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		row := scanner.Text()
		if targetRowIndex < len(targetRow) && rowNumber == targetRow[targetRowIndex] {
			row = strings.Split(row, ":")[0] + ":"
			targetRowIndex++
		}
		rowNumber++
		content = content + row + "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, []byte(content), 0666); err != nil {
		return err
	}
	return nil
}
