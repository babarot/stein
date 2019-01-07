config {
  # args = "_examples/manifests/microservices/*/*/*/*"

  # args = "_examples/manifests/**/*"

  report {
    format = "${format("[{{.Level}}]  {{.Rule}}  %s", color("{{.Message}}", "white"))}"

    # format = "{{.Level}}: {{.Rule}}: {{.Message}}"
    style = "console"
    color = true
  }
}
