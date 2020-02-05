package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

/* TODO

- Always install elevated
- firewall

*/

var (
	driver  neo4j.Driver
	session neo4j.Session
	result  neo4j.Result
	err     error
)

type GPO []struct {
	Text string `xml:",chardata"`
	N    string `xml:"N,attr"`
}

var (
	// types of settings we'll look for
	settings = []string{
		"PrivilegeRights",
		"RegistryValues",
	}
	// registry keys keywords to look for
	registryKeys = []string{
		"RequireSecuritySignature",
		"EnableSecuritySignature",
		"LmCompatibilityLevel",
		"EnableLUA",
		"FilterAdministratorToken",
		"LocalAccountTokenFilterPolicy",
	}

	// dangerous privileges
	dangerousPrivileges = []string{
		"SeEnableDelegationPrivilege",
		"SeDebugPrivilege",
		"SeBackupPrivilege",
		"SeImpersonationPrivilege",
		"SeRestorePrivilege",
		"SeCreateToken",
		"SeTakeOwnership",
		"SeTcbPrivilege",
		"SeCreateToken",
		"SeLoadDriver",
		"SeAssignPrimaryToken",
	}
)

func main() {
	configFile := flag.String("gpo", "non-existent.xml", "The XML configuration exported by PowerView")
	bloodHound := flag.Bool("bloodhound", false, "Use BloodHound APIs to augment data")
	flag.Parse()

	fileContent, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	gpos := new(Objs)
	xml.Unmarshal(fileContent, &gpos)

	//connect to BH
	var driver neo4j.Driver
	var session neo4j.Session

	if *bloodHound {
		driver, err = connectToBloodHound(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PASSWORD"))
		if err != nil {
			panic(err)
		}
		defer driver.Close()
		if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
			panic(err)
		}
		defer session.Close()
	}

	for _, gpo := range gpos.Obj {
		for _, attr := range gpo.MS.Obj {
			if isIn(attr.N, settings) {
				// Permissions
				// cover the cases where the privilege is assigned to multiple SIDs
				for _, value := range attr.MS.Obj {
					if isIn(value.N, registryKeys) {
						fmt.Println("[+] Found GPO that sets: " + value.N)
						fmt.Println(gpo.MS.S)
						regVal := value.LST.S[1]
						fmt.Println("\tValue -> " + regVal)
						if *bloodHound {
							if regVal == "0" && strings.Contains(value.N, "Signature") {
								executeQuery("computersWithoutSMBSigning", "none", getGPOName(GPO(gpo.MS.S)), session)
							}
						}
					}
					if isIn(value.N, dangerousPrivileges) {
						fmt.Println("[+] Found GPO that assigns: " + value.N)
						fmt.Println(gpo.MS.S)
						for _, s := range value.LST.S {
							userSID := strings.Trim(s, "*")
							fmt.Println("\t[+] " + userSID)
							if *bloodHound {
								executeQuery("usersAdminViaURA", userSID, getGPOName(GPO(gpo.MS.S)), session)
							}
						}
					}
				}
				// cover the cases where the privilege is assigned to only one SID
				for _, s := range attr.MS.S {
					if isIn(s.N, dangerousPrivileges) {
						fmt.Println("[+] Found GPO that assigns: " + s.N)
						fmt.Println(gpo.MS.S)
						userSID := strings.Trim(s.Text, "*")
						fmt.Println("\t[+] " + userSID)
						if *bloodHound {
							executeQuery("usersAdminViaURA", userSID, getGPOName(GPO(gpo.MS.S)), session)
						}
					}
				}

			}
		}

	}
}

func isIn(str string, list []string) bool {
	for _, v := range list {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func connectToBloodHound(username, password string) (neo4j.Driver, error) {
	if driver, err = neo4j.NewDriver("bolt://"+os.Getenv("NEO4J_SEVER")+":7687", neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	}); err != nil {
		fmt.Println("Error while establishing graph connection")
	}

	return driver, nil
}

func getGPOName(str GPO) string {
	for _, c := range str {
		if c.N == "GPOName" {
			return strings.Trim(c.Text, "{}")
		}
	}
	return "place"
}

func executeQuery(queryType string, SID string, GPOName string, session neo4j.Session) {
	var computersByGPO = `MATCH (g:GPO {guid: $GpoName}) WITH g 
	OPTIONAL MATCH (g)-[r1:GpLink {enforced:false}]->(container1) WITH g,container1 
	OPTIONAL MATCH (g)-[r2:GpLink {enforced:true}]->(container2) WITH g,container1,container2 
	OPTIONAL MATCH p1 = (g)-[r1:GpLink]->(container1)-[r2:Contains*1..]->(n1:Computer) WHERE NONE(x in NODES(p1) WHERE x.blocksinheritance = true AND LABELS(x) = 'OU') WITH g,p1,container2,n1 
	OPTIONAL MATCH p2 = (g)-[r1:GpLink]->(container2)-[r2:Contains*1..]->(n2:Computer) WITH n1,n2
	MATCH (n1), (n2) WITH collect(n1) + collect(n2) AS computers WITH computers
	UNWIND computers as c WITH c
	SET c.signing = false
	RETURN count(c)
	`

	var usersByGPO = `MATCH (u {objectsid: $sid}) WITH u
	MATCH (g:GPO {guid: $GpoName}) WITH u,g 
	OPTIONAL MATCH (g)-[r1:GpLink {enforced:false}]->(container1) WITH u,g,container1 
	OPTIONAL MATCH (g)-[r2:GpLink {enforced:true}]->(container2) WITH u,g,container1,container2 
	OPTIONAL MATCH p1 = (g)-[r1:GpLink]->(container1)-[r2:Contains*1..]->(n1:Computer) WHERE NONE(x in NODES(p1) WHERE x.blocksinheritance = true AND LABELS(x) = 'OU') WITH u,g,p1,container2,n1 
	OPTIONAL MATCH p2 = (g)-[r1:GpLink]->(container2)-[r2:Contains*1..]->(n2:Computer) WITH u,n1,n2
	MATCH (n1), (n2) WITH collect(n1) + collect(n2) AS computers,u
	UNWIND computers as c WITH c,u
	CREATE (u)-[:CanPrivesc]->(c)
	RETURN count(c)
	`

	if queryType == "usersAdminViaURA" {
		result, err = session.Run(usersByGPO, map[string]interface{}{
			"sid":     SID,
			"GpoName": GPOName,
		})
	}

	if queryType == "computersWithoutSMBSigning" {
		result, err = session.Run(computersByGPO, map[string]interface{}{
			"GpoName": GPOName,
		})
	}
}
