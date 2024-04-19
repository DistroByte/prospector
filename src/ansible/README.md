# Ansible

## Usage

### `apt-update.yaml`

| Variable            | Description                                | Default |
| ------------------- | ------------------------------------------ | ------- |
| `upgrade`           | upgrade packages                           | `false` |
| `packages`          | install packages                           | `[]`    |
| `check_hashicorp`   | check if hashicorp packages can be updated | `false` |
| `upgrade_hashicorp` | upgrade hashicorp packages                 | `false` |

#### Example

```bash
ansible-playbook -i hosts playbooks/apt-update.yaml
ansible-playbook -i hosts playbooks/apt-update.yaml -e "upgrade=true"
```

### `install-hashicorp.yaml`

#### Example

```bash
ansible-playbook -i hosts playbooks/install-hashicorp.yaml
```

### `configure-nomad.yaml`

Make sure to set the relevant variables in `hosts` before running this playbook.

#### Example

```bash
ansible-playbook -i hosts playbooks/configure-nomad.yaml
```