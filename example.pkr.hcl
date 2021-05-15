data "sshkey-sshkey" "sshkey" {
}

source "null" "example" {
  ssh_host = "127.0.0.1"
  ssh_username = "foo"
  ssh_private_key_file = data.sshkey-sshkey.sshkey.private_key_path
}

build {
  sources = ["sources.null.example"]
}
