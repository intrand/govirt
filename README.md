# About
govirt allows you to perform basic tasks in oVirt on the command line.

I highly recommend starting with `--help` before following the example below. There are very simple things you'll want to know how to set as you go, such as:
  * the oVirt engine API URL
  * engine authentication
  * template name
  * template version
  * storage domain name

# Example of building a 5-node kubernetes cluster using k3s

0. clean up previous tests (if any):

    ```sh
    # clear known_hosts
    for i in k3s{0..4} k3s{0..4}.domain 192.168.0.20{0..4}; do \
        ssh-keygen -R ${i}; \
    done;
    ```

    ```sh
    # force-stop any existing VMs with the same names
    for i in {0..4}; do \
        govirt vm stop \
            --force \
            --name "k3s${i}"; \
        done;

    # delete any existing VMs with the same names
    for i in {0..4}; do \
        govirt vm rm \
            --yes \
            --name "k3s${i}"; \
        done;
    ```

1. create new VMs:

    1. clone the VM from template:

        ```sh
        for i in {0..4}; do \
            govirt vm create \
                --name "k3s${i}" \
                --memory 4 \
                --cpu 2; \
            done;
        ```

    2. prepare cloud-init script for the new VM to run at bootup:

        ```yaml
        fqdn: k3s0.domain

        write_files:
        - path: /etc/cloud/cloud.cfg.d/99-custom-networking.cfg
        permissions: '0644'
        content: |
            network: {config: disabled}
        - path: /etc/netplan/config.yaml
        permissions: '0644'
        content: |
          network:
            version: 2
            renderer: networkd
            ethernets:
              enp1s0:
                optional: yes
                dhcp4: no
                dhcp6: no
                addresses:
                  - 192.168.0.200/24
                gateway4: 192.168.0.1
                nameservers:
                  addresses:
                    - 1.1.1.1
                    - 1.0.0.1
                    - 8.8.8.8

        runcmd:
        - date > /opt/.creation
        - rm /etc/netplan/50-cloud-init.yaml
        - netplan generate
        - netplan apply
        - dpkg-reconfigure openssh-server

        users:
        - default
        - name: root
        ssh_authorized_keys:
        - "entire public key string goes here exactly as it would appear in authorized_keys"
        ```

    3. start VM with cloud-init script:

        ```sh
        for i in {0..4}; do \
            govirt vm start \
                --name "k3s${i}" \
                --init \
                --script "${HOME}/cloud-init/k3s${i}.yml"; \
            done;
        ```

2. setup ansible controller:

    ```sh
    # wait for VMs to come up with correct hostname, ip address, etc and accept host ssh key
    for i in k3s{0..4}.domain; do \
        ssh -o StrictHostKeyChecking=accept-new "${i}" -t hostname -s; \
    done;
    ```

3. config management:

    ```sh
    # deploy config management and the k3s cluster
    ansible-playbook -i ~/ansible/inventory/kubernetes.yml -e reboot=yes ~/ansible/playbooks/k3s.yml

    # copy the kubeconfig from first node to controller node to make life easy
    rsync -avhP k3s0:"~/.kube/config" ~/.kube/config && sed -i 's|https://127.0.0.1|https://k3s0.domain|g' ~/.kube/config
    ```

4. verify cluster:

    ```sh
    kubectl get nodes -o wide

    kubectl get pods -A
    ```
