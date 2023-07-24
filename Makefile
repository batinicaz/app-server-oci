.PHONY: all
all: create-vm add-public-key run-playbook

ANSIBLE_DIR := packer/ansible
CURRENT_DIR := $(shell basename "$(CURDIR)")
ROLES_FILE := $(ANSIBLE_DIR)/requirements.yml
VM_NAME ?= $(CURRENT_DIR)

KEY_FILE := ~/.ssh/$(CURRENT_DIR)
$(KEY_FILE):
	@echo "Generating private key..."
	ssh-keygen -t rsa -b 4096 -f $(KEY_FILE) -N ''

create-vm: $(KEY_FILE)
	@multipass info $(VM_NAME) || multipass launch --name $(VM_NAME) --memory 1g --cpus 1 --disk 10G lts

add-public-key: create-vm
	@echo "Adding public key to authorized_users file..."
	multipass exec $(VM_NAME) -- bash -c "echo '$(shell ssh-keygen -y -f $(KEY_FILE))' >> ~/.ssh/authorized_keys"

run-playbook:
	ansible-galaxy install -r $(ROLES_FILE); \
	ANSIBLE_HOST_KEY_CHECKING=false \
    ansible-playbook -i "$(shell multipass info "$(VM_NAME)" | grep "^IPv4:" | sed "s/IPv4:[ ]*//")," \
    --extra-vars "local_dev=true" \
    --private-key $(KEY_FILE) \
    --user ubuntu \
    $(ANSIBLE_DIR)/playbook.yml

clean:
	multipass delete --purge $(VM_NAME)
