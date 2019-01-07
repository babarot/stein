function "get_service_name" {
  params = [file]
  result = basename(dirname(dirname(dirname(file))))
}

function "get_env" {
  params = [file]
  result = basename(dirname(dirname(file)))
}

function "get_service_id_with_env" {
  params = [file]
  result = format("%s-%s", get_service_name(file), lookup(var.shortened_environment, get_env(file)))
}

