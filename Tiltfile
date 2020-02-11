db_name='privacy-db'
db_pw='dev'

install_db = 'helm install ' + db_name + ' stable/postgresql --set postgresqlPassword=' + db_pw
delete_db = 'helm delete ' + db_name

local_resource('create db', cmd=install_db, trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False)
local_resource('delete db', cmd=delete_db, trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False)

k8s_yaml(helm('charts/connector', values=['charts/connector/values.yaml']))
k8s_yaml(helm('charts/controller', values=['charts/connector/values.yaml']))

docker_build('dropoutlabs/privacyai:latest', '.')
