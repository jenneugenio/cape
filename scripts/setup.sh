set -e
set -x

export CAPE_NAME=cape-user
export CAPE_EMAIL=cape_user@mycape.com
export CAPE_PASSWORD=capecape

if ! cape setup local http://localhost:8080; then
    cape config clusters remove -y local
    cape setup local http://localhost:8080
fi

output=$(cape services create --type data-connector --endpoint https://localhost:8081 service:dc@my-cape.com |  grep Token | awk '{print $2}')
kubectl create secret generic connector-secret --from-literal=token=$output

output=$(cape services create --type worker service:worker@my-cape.com | grep Token | awk '{print $2}')
kubectl create secret generic worker-secret --from-literal=token=$output

cape sources add --link service:dc@my-cape.com transactions postgres://postgres:dev@postgres-customer-postgresql:5434/customer

cape policies attach --from-file examples/allow-specific-fields.yaml allow-specific-fields global

cape pull transactions "SELECT * FROM transactions"