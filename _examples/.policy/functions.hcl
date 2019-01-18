function "remove_ext" {
  params = [file]
  result = replace(basename(file), ext(file), "")
}

function "white" {
  params = [file]
  result = color(file, "white")
}

function "is_irregular_namespace_pattern" {
  params = []
  result = contains(var.special_cases, get_service_id_with_env(filename))
}

function "kind" {
  params = [name]
  result = jsonpath("kind") == title(name)
}
