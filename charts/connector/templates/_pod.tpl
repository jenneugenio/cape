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
    args: ["connector", "start"]
    ports:
      - containerPort: {{ .Values.service.port }}
    env:
      - name: CAPE_COORDINATOR_URL
        value: {{ .Values.coordinator_url }}
      - name: CAPE_TOKEN
        valueFrom:
          secretKeyRef:
            name: {{  default (include "cape.fullname" .) .Values.token_secret }}
            key: token
{{- end }}
