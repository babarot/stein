rule "replicas" {
  description = ""

  precondition {
    cases = ["${kind("Deployment")}"]
  }

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

  precondition {
    cases = ["${kind("Deployment")}"]
  }

  expressions = [
    "${jsonpath("spec.template.spec.containers.#.name") != ""}",
  ]

  report {
    level   = "ERROR"
    message = "hoge"
  }
}
