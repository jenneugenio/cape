db_deployment_name='cape-db'
db_name='cape'
db_pw='dev'

override_pw = 'postgresqlPassword=' + db_pw
override_db = 'postgresqlDatabase=' + db_name

fetch_script = "helm fetch stable/postgresql --untar --untardir ./deploy/helm/"
local(fetch_script + " || true")
pg_yaml = helm('deploy/helm/postgresql', name=db_deployment_name, set=[override_pw, override_db], values=['./deploy/helm/postgresql/values.yaml'])
k8s_yaml(pg_yaml)
k8s_resource(db_deployment_name + '-postgresql', port_forwards=5432)

# Deleting the helm chart doesn't delete the PVC that the helm chart creates
# So, we explicitly delete it here with kubectl
# Also need to make sure don't delete the secret so remove that from the
# list of manifests
db_resources = [decode_yaml(yaml) for yaml in str(pg_yaml).split('\n---\n')]
without_secret = [r for r in db_resources if r["kind"] != 'Secret']
delete_db_chart = 'kubectl delete ' + ' '.join([(r["kind"] + "/" + r["metadata"]["name"]) for r in without_secret])
delete_pvc = 'kubectl delete pvc/data-' + db_deployment_name + '-postgresql-0'
delete_cmd = ' && '.join([delete_db_chart, delete_pvc])

local_resource('delete db', cmd=delete_cmd, trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False)

k8s_yaml(helm('charts/connector', values=['charts/local_values/connector_values.yaml']))
k8s_yaml(helm('charts/controller', values=['charts/local_values/controller_values.yaml']))

k8s_yaml('manifests/test_job.yaml')
k8s_resource("test", trigger_mode=TRIGGER_MODE_MANUAL)

docker_build('dropoutlabs/cape:latest', '.', dockerfile='dockerfiles/Dockerfile.base')
docker_build('dropoutlabs/cape-test:latest', '.', dockerfile='dockerfiles/Dockerfile.test')
docker_build('dropoutlabs/controller:latest', '.', dockerfile='dockerfiles/Dockerfile.controller')
docker_build('dropoutlabs/connector:latest', '.', dockerfile='dockerfiles/Dockerfile.connector')
docker_build('dropoutlabs/update:latest', '.', dockerfile='dockerfiles/Dockerfile.update')

k8s_resource('connector', port_forwards=8081, trigger_mode=TRIGGER_MODE_MANUAL)
k8s_resource('controller', port_forwards=8080, trigger_mode=TRIGGER_MODE_MANUAL)
