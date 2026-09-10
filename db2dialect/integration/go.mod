module github.com/ibmdb/go_ibm_db/db2dialect/integration

go 1.24.0

require (
	github.com/ibmdb/go_ibm_db v0.0.0
	github.com/ibmdb/go_ibm_db/db2dialect v0.0.0
	github.com/uptrace/bun v1.2.18
)

require (
	github.com/ibmruntimes/go-recordio/v2 v2.0.0-20240416213906-ae0ad556db70 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

replace github.com/ibmdb/go_ibm_db => ../..

replace github.com/ibmdb/go_ibm_db/db2dialect => ..
