# Concepts

Stein is a tool like a linter for config files (such as JSON, YAML, HCL) with a customizable rule set.

```console
$ stein apply manifests/**/*
manifests/microservices/x-gateway-jp/development/Deployment/redis-master.yaml
  No violated rules

manifests/microservices/x-echo-jp/development/Deployment/test.yml
  [ERROR]  rule.one_resource_per_one_file  Only 1 resource should be defined in a YAML file
  [WARN ]  rule.yaml_separator             YAML separator "---" should be removed
  [ERROR]  rule.filename_extension         Filename extension should be yaml or yaml.enc

manifests/microservices/x-echo-jp/development/Service/service.yaml
  [ERROR]  rule.namespace_is_specified     Namespace is not specified
  [ERROR]  rule.metadata_name_is_correct   "service.yaml" should be "redis-master"+.yaml

4 error(s), 1 warn(s)
```

Motivation...bla bla bla
