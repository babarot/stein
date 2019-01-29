rule "filename" {
  description = "Check filename is the same as metadata.name"

  conditions = [
    "${jsonpath("metadata.name") == remove_ext(filename)}",
  ]

  report {
    level   = "ERROR"
    message = "${format("Filename should be %s.yaml (metadata.name + .yaml)", jsonpath("metadata.name"))}"
  }
}

rule "resource_per_file" {
  description = ""

  conditions = [
    "${wc(grep("^kind: ", file(filename))) == 0}",
  ]

  report {
    level   = "ERROR"
    message = "Only 1 resource should be defined in a YAML file"
  }
}

rule "yaml_separator" {
  description = "Do not use YAML separator"

  conditions = [
    "${length(grep("^---", file(filename))) == 0}",
  ]

  report {
    level   = "WARN"
    message = "YAML separator \"---\" should be removed"
  }
}
