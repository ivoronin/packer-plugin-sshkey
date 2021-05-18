# Packer SSH key plugin
Packer plugin for auto-generating SSH keys.
Packer 1.7 and later is required.

# Documentation
## Usage example
```hcl
packer {
  required_plugins {
    sshkey = {
      version = ">= 0.0.2"
      source = "github.com/ivoronin/sshkey"
    }
  }
}

data "sshkey" "install" {
}

source "qemu" "install" {
  ssh_username              = "root"
  ssh_private_key_file      = data.sshkey.install.private_key_path
  ssh_clear_authorized_keys = true
  http_content = {
    "/install.conf" = templatefile("install.conf.pkrtpl", { "ssh_public_key" : data.sshkey.install.public_key })
  }
  <...>
}

build {
  sources = ["source.qemu.install"]
}
```
Run `packer init` to automatically download and install plugin.

## Arguments
  - `name` (default: "packer") - Key name, *should be unique across `sshkey` datasources*.

## Attributes
  - `public_key` - SSH public key in "ssh-rsa ..." format
  - `private_key_path` - Path to SSH private key

## Notes
  - Plugin generates 2048-bit RSA keys
  - Private key is cached in `PACKER_CACHE_DIR` (by default "packer_cache" directory is used). If you delete cached private key it will be regenerated on next run.
