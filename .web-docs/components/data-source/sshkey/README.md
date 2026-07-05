Type: `sshkey`

Data source used to generate SSH keys

## Parameters and output

### Optional

  - `name` (string) - Key name, **must be unique** across `sshkey` datasource instances. Defaults to `packer`.
  - `type` (string) - Key type, must be either `rsa` or `ed25519`. Defaults to `rsa`.

## Output data

  - `public_key` (string) - SSH public key in "ssh-rsa ..." format
  - `private_key_path` (string) - Path to SSH private key

## Notes

  - Private key is cached in `PACKER_CACHE_DIR` (by default `packer_cache` directory is used). If you delete cached private key it will be regenerated on the next run.
  - Packer 1.7.3 or later is required

## Key lifecycle and cleanup

By default generated private keys stay in the Packer cache. If you run Packer again on the same machine with the same `name` and `type`, the plugin sees the same cache path and loads the existing key instead of generating a new one.

Temporary SSH keys are not automatically removed at the end of the build because Packer does not give data sources an end-of-build cleanup step. So cleanup has to live in the Packer template, where Packer already has a local execution step.

For successful builds, add a final `shell-local` post-processor inside the `build` block:

```hcl
post-processor "shell-local" {
  inline = [
    "rm -f -- '${data.sshkey.install.private_key_path}'",
  ]
}
```

For provisioning failures, add the same command as an `error-cleanup-provisioner` inside the `build` block:

```hcl
error-cleanup-provisioner "shell-local" {
  inline = [
    "rm -f -- '${data.sshkey.install.private_key_path}'",
  ]
}
```
