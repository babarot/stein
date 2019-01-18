rule "namespace_specification" {
  description = "Check namespace name is not empty"

  expressions = [
    "${jsonpath("metadata.namespace") != ""}",
  ]

  report {
    level   = "ERROR"
    message = "Namespace is not specified"
  }
}

rule "namespace_name" {
  description = "Check namespace name is valid"

  depends_on = ["rule.namespace_specification"]

  precondition {
    cases = [
      "${!is_irregular_namespace_pattern()}",
    ]
  }

  expressions = [
    "${jsonpath("metadata.namespace") == get_service_id_with_env(filename)}",
  ]

  report {
    level   = "ERROR"
    message = "${format("Namespace name %q is invalid", jsonpath("metadata.namespace"))}"
  }
}

rule "namespace_name_irregular" {
  description = "Check namespace name is valid"

  depends_on = ["rule.namespace_specification"]

  precondition {
    cases = [
      "${is_irregular_namespace_pattern()}",
    ]
  }

  expressions = [
    "${contains(lookuplist(var.namespace_name_map, jsonpath("metadata.namespace")), get_service_id_with_env(filename))}",
  ]

  report {
    level   = "ERROR"
    message = "${format("This case is irregular pattern, so %q is invalid", jsonpath("metadata.namespace"))}"
  }
}

rule "extension" {
  description = "Acceptable yaml file extensions are limited"

  expressions = [
    "${ext(filename) == ".yaml" || ext(filename) == ".yaml.enc"}",
  ]

  report {
    level   = "ERROR"
    message = "File extension should be yaml or yaml.enc"
  }
}
