set -e

SHARED_SECRET="testing123"
POST_DATA='{"action":"completed","number":"54"}'
HMAC=$(echo -n "${POST_DATA}" | openssl dgst -sha256 -hmac "${SHARED_SECRET}" | cut -d ' ' -f 2)
HMAC_HEADER="x-hub-signature-256: sha256=${HMAC}"
CT_HEADER="Content-Type: application/json; charset=UTF-8"

curl -X POST -H "${CT_HEADER}" -H "${HMAC_HEADER}" -d "${POST_DATA}" $FUNCTION_URL

