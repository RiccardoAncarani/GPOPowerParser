```
Get-DomainPolicyData -Policy All | exort-clixml gpo.xml
```

```
iconv -f utf-16 -t utf-8 gpo.xml >  utf8_corp_gpo.txt
```

```
brew install michael-simons/homebrew-seabolt/seabolt
git clone https://gitlab.com/riccardo.ancarani94/power-gpo-parser.git
cd power-gpo-parser
export GOPATH=$(pwd)
go get -v github.com/neo4j/neo4j-go-driver/neo4j
cd src/power-gpo-parser
go install
```

```
./power-gpo-parser --gpo utf8_corp_gpo.txt --bloodhound
```


Find all the users that can RDP into a machine where they have special privileges:
```
MATCH (u:User)-[:CanRDP]->(c:Computer) WITH u,c
OPTIONAL MATCH (u)-[:MemberOf*1..]->(g:Group)-[:CanRDP]->(c) WITH u,c
MATCH (u)-[:CanPrivesc]->(c) RETURN u.name, c.name
```


Find all the NTLM relay opportunities for computer accounts:
```
MATCH (u1:Computer)-[:AdminTo]->(c1:Computer {signing: false}) RETURN u1.name, c1.name
MATCH (u2)-[:MemberOf*1..]->(g:Group)-[:AdminTo]->(c2 {signing: false}) RETURN u2.name, c2.name
```