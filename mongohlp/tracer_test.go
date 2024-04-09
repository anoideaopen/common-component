package mongohlp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTraceMongoCmd(t *testing.T) {
	const findCmd = `{"find": "balshistory","filter": {"channel": "ch1","address": "addr1","token": "token1","allowed": "allowed1","orderedID": {"$lt": "3"}},"limit": {"$numberLong":"1"},"singleBatch": true,"sort": {"orderedID": {"$numberInt":"-1"}},"lsid": {"id": {"$binary":{"base64":"Yfn/oN3XT7+agzwk7cQVWQ==","subType":"04"}}},"$clusterTime": {"clusterTime": {"$timestamp":{"t":"1663388074","i":"20"}},"signature": {"hash": {"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId": {"$numberLong":"0"}}},"$db": "hlf-cb-tests","$readPreference": {"mode": "primaryPreferred"}}`
	const aggregateCmd = `{"aggregate": "totalbalsdate","pipeline": [{"$match": {"time": {"$lte": {"$numberLong":"2"}},"token": {"$in": ["token1"]}}},{"$group": {"_id": "$grOrderedID","address": {"$first": "$$ROOT.address"},"targetCh": {"$first": "$$ROOT.targetCh"},"token": {"$first": "$$ROOT.token"}}},{"$sort": {"_id": {"$numberInt":"1"}}},{"$limit": {"$numberLong":"2"}},{"$project": {"_id": "$_id","address": "$address","channel": "$targetCh","token": "$token"}},{"$lookup": {"from": "totalbalsdate","let": {"address": "$address","channel": "$channel","token": "$token"},"as": "items","pipeline": [{"$match": {"$expr": {"$and": [{"$eq": ["$address","$$address"]},{"$eq": ["$targetCh","$$channel"]},{"$eq": ["$token","$$token"]},{"$lte": ["$time",{"$numberLong":"2"}]}]}}},{"$sort": {"time": {"$numberInt":"-1"}}},{"$group": {"_id": {"channel": "$channel","targetCh": "$targetCh","token": "$token","address": "$address"},"mmm": {"$first": "$$ROOT"}}},{"$replaceRoot": {"newRoot": "$mmm"}}]}},{"$sort": {"_id": {"$numberInt":"1"}}}],"cursor": {},"lsid": {"id": {"$binary":{"base64":"gAweYmacS6m8VVD2yL2u5g==","subType":"04"}}},"$clusterTime": {"clusterTime": {"$timestamp":{"t":"1663381898","i":"29"}},"signature": {"hash": {"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId": {"$numberLong":"0"}}},"$db": "hlf-cb-tests","$readPreference": {"mode": "primaryPreferred"}}`

	trStr, err := TraceMongoCmd(findCmd)
	require.NoError(t, err)
	log.Println(trStr)

	trStr, err = TraceMongoCmd(aggregateCmd)
	require.NoError(t, err)
	log.Println(trStr)
}
