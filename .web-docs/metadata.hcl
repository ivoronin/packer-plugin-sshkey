# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "SSH Key"
  description = "TODO"
  identifier = "packer/BrandonRomano/sshkey"
  component {
    type = "data-source"
    name = "SSH Key"
    slug = "sshkey"
  }
}
