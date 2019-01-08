rule "namespace_specification" {
  description = "Check namespace name is not empty"

  ignore_cases = []

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

  ignore_cases = [
    "${is_irregular_namespace_pattern()}",
  ]

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

  ignore_cases = [
    "${!is_irregular_namespace_pattern()}",
  ]

  expressions = [
    "${contains(lookuplist(var.namespace_name_map, jsonpath("metadata.namespace")), get_service_id_with_env(filename))}",
  ]

  report {
    level   = "ERROR"
    message = "${format("Namespace name %q is invalid", jsonpath("metadata.namespace"))}"
  }
}

rule "extension" {
  description = "Acceptable yaml file extensions are limtited"

  ignore_cases = []

  expressions = [
    "${ext(filename) == ".yaml" || ext(filename) == ".yaml.enc"}",
  ]

  report {
    level   = "ERROR"
    message = "File extension should be yaml or yaml.enc"
  }
}
