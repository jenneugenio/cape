{{- define "cape.pod" }}
{{- if .Values.schedulerName }}
schedulerName: "{{ .Values.schedulerName }}"
{{- end }}
serviceAccountName: {{ template "cape.serviceAccountName" . }}
{{- if .Values.securityContext }}
securityContext:
{{ toYaml .Values.securityContext | indent 2 }}
{{- end }}
{{- if .Values.priorityClassName }}
priorityClassName: {{ .Values.priorityClassName }}
{{- end }}
{{- if .Values.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.image.pullSecrets }}
  - name {{ . }}
{{- end }}
{{- end }}
containers:
- name: {{ .Chart.Name }}
  image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
  imagePullPolicy: {{ .Values.image.pullPolicy }}
  command: ["cape"]
  args: 
  - "coordinator"
  - "start"
  - "--file"
  - "/etc/coordinator/coordinator-config.yaml"
  - "--instance-id"
  - "{{ .Values.instance_id }}"
  ports:
  - containerPort: {{ .Values.service.port }}
  volumeMounts:
  - name: config
    mountPath: "/etc/coordinator"
    readOnly: true
{{ if eq .Values.includeUI true }}
- name: {{ .Chart.Name }}-ui
  image: "{{ .Values.uiImage.repository }}:{{ .Values.uiImage.tag }}"
  imagePullPolicy: {{ .Values.uiImage.pullPolicy }}
{{- end }}
volumes:
- name: config
  secret:
    secretName: {{ template "cape.configSecret" . }}
{{- end }}
