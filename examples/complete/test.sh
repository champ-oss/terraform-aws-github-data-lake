set -e

SHARED_SECRET="testing123"
POST_DATA='{"action":"completed","number":"54"}'
HMAC=$(echo -n "${POST_DATA}" | openssl dgst -sha256 -hmac "${SHARED_SECRET}" | cut -d ' ' -f 2)
HMAC_HEADER="x-hub-signature-256: sha256=${HMAC}"
CT_HEADER="Content-Type: application/json; charset=UTF-8"

echo -e "\nTesting two HTTP POST curl requests..."
curl -s -X POST -H "${CT_HEADER}" -H "${HMAC_HEADER}" -d "${POST_DATA}" $FUNCTION_URL
curl -s -X POST -H "${CT_HEADER}" -H "${HMAC_HEADER}" -d "${POST_DATA}" $FUNCTION_URL

echo -e "\nWaiting for S3 objects..."
sleep 120
aws s3 ls s3://$BUCKET --recursive | grep .gz

echo -e "\nTest sending an HTTP POST request with an invalid secret"
curl -s -X POST -H "${CT_HEADER}" -H "x-hub-signature-256: sha256=invalid" -d "${POST_DATA}" $FUNCTION_URL | grep "Internal Server Error"

echo -e "\nStarting Athena query..."
aws athena start-query-execution --query-execution-context "Database=${DATABASE}" --result-configuration "OutputLocation=s3://${BUCKET}/" --query-string "SELECT * FROM \""${TABLE}"\""