rule "replicas" {
  description = ""

  ignore_cases = [
    "${jsonpath("kind") != "Deployment"}",
  ]

  expressions = [
    "${jsonpath("spec.replicas", 0) >= 1}",
  ]

  report {
    level   = "ERROR"
    message = "Too few replicas"
  }
}

rule "images" {
  description = ""

  ignore_cases = [
    "${jsonpath("kind") != "Deployment"}",
  ]

  expressions = [
    "${jsonpath("spec.template.spec.containers.#.name") != ""}",
  ]

  report {
    level   = "ERROR"
    message = "hoge"
  }
}
