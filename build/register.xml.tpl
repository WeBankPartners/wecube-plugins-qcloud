<?xml version="1.0" encoding="UTF-8"?>
<package name="qcloud-resource-management" version="{{PLUGIN_VERSION}}">
    <docker-image-file>wecube-plugins-qcloud.tar</docker-image-file>
    <docker-image-repository>wecube-plugins-qcloud</docker-image-repository>
    <docker-image-tag>{{IMAGE_TAG}}</docker-image-tag>
    <container-port>8081</container-port>
    <container-config-directory>/home/app/wecube-plugins-qcloud/conf</container-config-directory>
    <container-log-directory>/home/app/wecube-plugins-qcloud/log</container-log-directory>
    <container-start-param>-v /etc/localtime:/etc/localtime -v /home/app/wecube-plugins-qcloud/logs:/home/app/wecube-plugins-qcloud/logs</container-start-param>
    <plugin id="vpc" name="Vpc Management" >
        <interface name="create" path="/v1/qcloud/vpc/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">cidr_block</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/vpc/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="peering-connection" name="Peer Connection Management">
        <interface name="create" path="/v1/qcloud/peering-connection/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">peer_provider_params</parameter>
                <parameter datatype="string">peer_vpc_id</parameter>
                <parameter datatype="string">peer_uin</parameter>
                <parameter datatype="string">bandwidth</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/peering-connection/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">peer_provider_params</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="security-group" name="Security Group Management">
        <interface name="create" path="/v1/qcloud/security-group/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">description</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/security-group/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="create-policies" path="/v1/qcloud/security-group/create-policies">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">description</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">rule_type</parameter>
                <parameter datatype="string">rule_cidr_ip</parameter>
                <parameter datatype="string">rule_ip_protocol</parameter>
                <parameter datatype="string">rule_port_range</parameter>
                <parameter datatype="string">rule_policy</parameter>
                <parameter datatype="string">rule_description</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="delete-policies" path="/v1/qcloud/security-group/delete-policies">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">description</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">rule_type</parameter>
                <parameter datatype="string">rule_cidr_ip</parameter>
                <parameter datatype="string">rule_ip_protocol</parameter>
                <parameter datatype="string">rule_port_range</parameter>
                <parameter datatype="string">rule_policy</parameter>
                <parameter datatype="string">rule_description</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="route-table" name="Route Table Management">
        <interface name="create" path="/v1/qcloud/route-table/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">vpc_id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/route-table/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="associate-subnet" path="/v1/qcloud/route-table/associate-subnet">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">subnet_id</parameter>
                <parameter datatype="string">route_table_id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="subnet" name="Subnet Management">
        <interface name="create" path="/v1/qcloud/subnet/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">cidr_block</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/subnet/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="vm" name="Virtual Machine Management">
        <interface name="create" path="/v1/qcloud/vm/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                 <parameter datatype="string">seed</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">subnet_id</parameter>
                <parameter datatype="string">instance_name</parameter>
                <parameter datatype="string">instance_type</parameter>
                <parameter datatype="string">image_id</parameter>
                <parameter datatype="number">system_disk_size</parameter>
                <parameter datatype="string">instance_charge_type</parameter>
                <parameter datatype="number">instance_charge_period</parameter>
                <parameter datatype="string">instance_private_ip</parameter>
                 <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">cpu</parameter>
                <parameter datatype="string">memory</parameter>
                <parameter datatype="string">password</parameter>
                <parameter datatype="string">instance_state</parameter>
                <parameter datatype="string">instance_private_ip</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/vm/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="start" path="/v1/qcloud/vm/start">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="stop" path="/v1/qcloud/vm/stop">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="storage" name="Storage Management">
        <interface name="create" path="/v1/qcloud/storage/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">disk_type</parameter>
                <parameter datatype="number">disk_size</parameter>
                <parameter datatype="string">disk_name</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">disk_charge_type</parameter>
                <parameter datatype="string">disk_charge_period</parameter>
                <parameter datatype="string">instance_id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/storage/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="nat-gateway" name="Nat Gateway Management">
        <interface name="create" path="/v1/qcloud/nat-gateway/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="number">max_concurrent</parameter>
                <parameter datatype="number">bandwidth</parameter>
                <parameter datatype="string">assigned_eip_set</parameter>
                <parameter datatype="number">auto_alloc_eip_num</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/nat-gateway/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="mysql-vm" name="Mysql Management">
        <interface name="create" path="/v1/qcloud/mysql-vm/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">engine_version</parameter>
                <parameter datatype="number">memory</parameter>
                <parameter datatype="number">volume</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">subnet_id</parameter>
                <parameter datatype="string">name</parameter>
                <parameter datatype="number">count</parameter>
                <parameter datatype="string">charge_type</parameter>
                <parameter datatype="number">charge_period</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/mysql-vm/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="restart" path="/v1/qcloud/mysql-vm/restart">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="mariadb" name="Mariadb Management">
        <interface name="create" path="/v1/qcloud/mariadb/create">
         <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">seed</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">user_name</parameter>
                <parameter datatype="string">zones</parameter>
                <parameter datatype="number">node_count</parameter>
                <parameter datatype="number">memory_size</parameter>
                <parameter datatype="number">storage_size</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">subnet_id</parameter>
                <parameter datatype="number">charge_period</parameter>
                <parameter datatype="string">db_version</parameter>
                <parameter datatype="string">character_set</parameter>
                <parameter datatype="string">lower_case_table_names</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="number">private_port</parameter>
                <parameter datatype="string">private_ip</parameter>
                <parameter datatype="string">user_name</parameter>
                <parameter datatype="string">password</parameter>
             </output-parameters>
        </interface>
    </plugin>

    <plugin id="route-policy" name="Route Policy Management">
        <interface name="create" path="/v1/qcloud/route-policy/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">route_table_id</parameter>
                <parameter datatype="string">dest_cidr</parameter>
                <parameter datatype="string">gateway_type</parameter>
                <parameter datatype="string">gateway_id</parameter>
                <parameter datatype="string">desc</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
        <interface name="terminate" path="/v1/qcloud/route-policy/terminate">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="string">id</parameter>
                <parameter datatype="string">route_table_id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
            </output-parameters>
        </interface>
    </plugin>

    <plugin id="redis" name="Redis Management">
        <interface name="create" path="/v1/qcloud/redis/create">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">provider_params</parameter>
                <parameter datatype="number">type_id</parameter>
                <parameter datatype="number">mem_size</parameter>
                <parameter datatype="number">goods_num</parameter>
                <parameter datatype="number">period</parameter>
                <parameter datatype="string">password</parameter>
                <parameter datatype="number">billing_mode</parameter>
                <parameter datatype="string">vpc_id</parameter>
                <parameter datatype="string">subnet_id</parameter>
                <parameter datatype="string">id</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">request_id</parameter>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">id</parameter>
            </output-parameters>
        </interface>
    </plugin>
</package>
