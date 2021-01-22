{{/* vim: set filetype=mustache: */}}

{{/*
Returns the secret name for the Secret containing the TLS certificate and key.
Uses `ingress.tls.secretName` first and falls back to `global.ingress.tls.secretName`
if there is a shared tls secret for all ingresses.
*/}}
{{- define "webservice.tlsSecret" -}}
{{- $defaultName := (dict "secretName" "") -}}
{{- if .Values.global.ingress.configureCertmanager -}}
{{- $_ := set $defaultName "secretName" (printf "%s-gitlab-tls" .Release.Name) -}}
{{- else -}}
{{- $_ := set $defaultName "secretName" (include "gitlab.wildcard-self-signed-cert-name" .) -}}
{{- end -}}
{{- pluck "secretName" .Values.ingress.tls .Values.global.ingress.tls $defaultName | first -}}
{{- end -}}

{{/*
Returns the workhorse image repository depending on the value of global.edition.

Used to switch the deployment from Enterprise Edition (default) to Community
Edition. If global.edition=ce, returns the Community Edition image repository
set in the Gitlab values.yaml, otherwise returns the Enterprise Edition
image repository.
*/}}
{{- define "workhorse.repository" -}}
{{- if eq "ce" .Values.global.edition -}}
{{ index .Values "global" "communityImages" "workhorse" "repository" }}
{{- else -}}
{{ index .Values "global" "enterpriseImages" "workhorse" "repository" }}
{{- end -}}
{{- end -}}

{{/*
Returns the webservice image depending on the value of global.edition.

Used to switch the deployment from Enterprise Edition (default) to Community
Edition. If global.edition=ce, returns the Community Edition image repository
set in the Gitlab values.yaml, otherwise returns the Enterprise Edition
image repository.
*/}}
{{- define "webservice.image" -}}
{{ coalesce .Values.image.repository (include "image.repository" .) }}:{{ coalesce .Values.image.tag (include "gitlab.versionTag" . ) }}
{{- end -}}

{{/*
Returns ERB section for Workhorse direct object storage configuration.

If Minio in use, set AWS and keys.
If consolidated object storage is in use, read the connection YAML
  If provider is AWS, render enabled as true.
*/}}
{{- define "workhorse.object_storage.config" -}}
<%
  require 'yaml'

  supported_providers = ['AWS']
  provider = ''
  aws_access_key_id = ''
  aws_secret_access_key = ''

  if File.exists? '/etc/gitlab/minio/accesskey'
    provider = 'AWS'
    aws_access_key_id = File.read('/etc/gitlab/minio/accesskey').strip.dump[1..-2]
    aws_secret_access_key = File.read('/etc/gitlab/minio/secretkey').strip.dump[1..-2]
  end

  if File.exists? '/etc/gitlab/objectstorage/object_store'
    connection = YAML.safe_load(File.read('/etc/gitlab/objectstorage/object_store'))
    provider = connection['provider']
    if connection.has_key? 'aws_access_key_id'
      aws_access_key_id = connection['aws_access_key_id']
      aws_secret_access_key = connection['aws_secret_access_key']
    end
  end

  if supported_providers.include? provider
%>
[object_storage]
provider = "<%= provider %>"
<%   if provider.eql? 'AWS' %>
# AWS / S3 object storage configuration.
[object_storage.s3]
# access/secret can be blank!
aws_access_key_id = "<%= aws_access_key_id %>"
aws_secret_access_key = "<%= aws_secret_access_key %>"
<%
    end
  end
%>
{{- end -}}