# ***RedHat has abandoned RHV and thusly oVirt is likely dead. This project will not be maintained going forward :(***


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

    1. clone each VM from a sealed template:

        ```sh
        for i in {0..4}; do \
            govirt vm create \
                --name "k3s${i}" \
                --memory 4 \
                --cpu 2; \
            done;
        ```

    2. prepare cloud-init scripts for each new VM to run at bootup:

        ```sh
        mkdir -p "${HOME}/cloud-init" && \
        for i in {0..4}; do \
        govirt cloud-init create \
            --fqdn "k3s${i}.lan" \
            --ssh-key "ssh-rsa some_ssh_public_key with_comment_if_you_like" > "${HOME}/cloud-init/k3s${i}.yml"; \
        done;

    3. start VMs with cloud-init script:

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
