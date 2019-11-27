<?xml version="1.0" encoding="UTF-8"?>
<package name="qcloud" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
    </packageDependencies>

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单 -->
    <menus>
    </menus>

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系 -->
    <dataModel>
    </dataModel>

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数 -->
    <systemParameters>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="{{PORTBINDINGS}}" volumeBindings="/etc/localtime:/etc/localtime,{{base_mount_path}}/qcloud/logs:/home/app/qcloud/logs" envVariables=""/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="vpc">
            <interface action="create" path="/qcloud/v1/vpc/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection">
            <interface action="create" path="/qcloud/v1/peering-connection/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">peer_vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">peer_uin</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">bandwidth</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group">
            <interface action="create" path="/qcloud/v1/security-group/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">description</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="create-policies" path="/qcloud/v1/security-group/create-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">description</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">policy_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">policy_protocol</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">policy_port</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">policy_action</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_description</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-group/delete-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_port</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">policy_action</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-table">
            <interface action="create" path="/qcloud/v1/route-table/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-table/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="associate-subnet" path="/qcloud/v1/route-table/associate-subnet">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet">
            <interface action="create" path="/qcloud/v1/subnet/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>route_table_id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm">
            <interface action="create" path="/qcloud/v1/vm/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">instance_name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">image_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host_type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">system_disk_size</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_charge_type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">instance_charge_period</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">project_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>cpu</parameter>
                    <parameter datatype="string" mappingType='context'>memory</parameter>
                    <parameter datatype="string" mappingType='context'>password</parameter>
                    <parameter datatype="string" mappingType='context'>instance_state</parameter>
                    <parameter datatype="string" mappingType='context'>instance_private_ip</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="bind security group to vm" path="/qcloud/v1/vm/bind-security-groups">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">security_group_ids</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage">
            <interface action="create" path="/qcloud/v1/storage/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">disk_size</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">disk_name</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_charge_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_charge_period</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/storage/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="buy cbs disk and mount" path="/qcloud/v1/cbs/create-mount">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">disk_size</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">disk_name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_charge_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">disk_charge_period</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">file_system_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">mount_dir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>volume_name</parameter>
                    <parameter datatype="string" mappingType='context'>disk_id</parameter>
                </outputParameters>
            </interface>
            <interface action="umount and destroy cbs disk" path="/qcloud/v1/cbs/umount-terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">volume_name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">mount_dir</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">instance_guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway">
            <interface action="create" path="/qcloud/v1/nat-gateway/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">max_concurrent</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">bandwidth</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>eip</parameter>
                    <parameter datatype="string" mappingType='context'>eip_id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-vm">
            <interface action="create" path="/qcloud/v1/mysql-vm/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">engine_version</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">memory</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">volume</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">charge_type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">charge_period</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">character_set</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lower_case_table_names</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>private_ip</parameter>
                    <parameter datatype="string" mappingType='context'>private_port</parameter>
                    <parameter datatype="string" mappingType='context'>user_name</parameter>
                    <parameter datatype="string" mappingType='context'>password</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql-vm/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql-vm/restart">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mariadb">
            <interface action="create" path="/qcloud/v1/mariadb/create">
            <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">user_name</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">zones</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">node_count</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">memory_size</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">storage_size</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">charge_period</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">db_version</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">character_set</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lower_case_table_names</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>private_port</parameter>
                    <parameter datatype="string" mappingType='context'>private_ip</parameter>
                    <parameter datatype="string" mappingType='context'>user_name</parameter>
                    <parameter datatype="string" mappingType='context'>password</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy">
            <interface action="create" path="/qcloud/v1/route-policy/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">route_table_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">dest_cidr</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">gateway_type</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">gateway_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">desc</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis">
            <interface action="create" path="/qcloud/v1/redis/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">type_id</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">mem_size</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">goods_num</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">period</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">billing_mode</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb">
            <interface action="create" path="/qcloud/v1/clb/create">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="number" mappingType='entity' required="N">name</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">type</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">vpc_id</parameter>
                    <parameter datatype="number" mappingType='entity' required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>id</parameter>
                    <parameter datatype="string" mappingType='context'>vip</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="add-backtarget" path="/qcloud/v1/clb/add-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lb_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lb_port</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">protocol</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host_port</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
            <interface action="del-backtarget" path="/qcloud/v1/clb/del-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lb_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">lb_port</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">protocol</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host_id</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host_port</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>