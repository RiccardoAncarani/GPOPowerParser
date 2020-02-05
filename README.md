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