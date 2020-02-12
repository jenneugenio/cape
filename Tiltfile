db_name='privacy-db'
db_pw='dev'

start_db = 'helm upgrade --install ' + db_name + ' stable/postgresql --set postgresqlPassword=' + db_pw

# Deleting the helm chart doesn't delete the PVC that the helm chart creates
# So, we explicitly delete it here with kubectl
delete_db_chart = 'helm delete ' + db_name
delete_pvc = 'kubectl delete pvc/data-' + db_name + '-postgresql-0'
delete_cmd = ' && '.join([delete_db_chart, delete_pvc])

local_resource('create db', cmd=start_db)
local_resource('delete db', cmd=delete_cmd, trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False)

k8s_yaml(helm('charts/connector', values=['charts/connector/values.yaml']))
k8s_yaml(helm('charts/controller', values=['charts/connector/values.yaml']))

docker_build('dropoutlabs/privacyai:latest', '.')
