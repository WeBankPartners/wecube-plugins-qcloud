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
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="{{PORTBINDINGS}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/qcloud/logs:/home/app/qcloud/logs" envVariables=""/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="vpc" targetPackage="wecmdb" targetEntity="network_zone">
            <interface action="create" path="/qcloud/v1/vpc/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.key_name" required="Y">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.network_segment>wecmdb:network_segment.code" required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection" targetPackage="wecmdb" targetEntity="network_zone_link">
            <interface action="create" path="/qcloud/v1/peering-connection/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_1>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.asset_code" required="Y">peer_vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.user_id" required="Y">peer_uin</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.netband_width" required="Y">bandwidth</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_1>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="wecmdb" targetEntity="service">
            <interface action="create" path="/qcloud/v1/security-group/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.key_name" required="Y">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.description" required="Y">description</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.security_group_asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.security_group_asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id"  required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter"  required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.security_group_asset_code"  required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="invoke">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.security_group_asset_code" required="N">security_group_id</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="security_policy_inbound" required="Y">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.unit>wecmdb:unit~(unit)wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.service_type" required="Y">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.service_port" required="Y">policy_port</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="security_policy_action_permit" required="Y">policy_action</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.description" required="N">policy_description</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.security_group_asset_code" required="Y">security_group_id</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="security_policy_inbound" required="Y">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.unit>wecmdb:unit~(unit)wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.service_type" required="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.service>wecmdb:service.service_port" required="N">policy_port</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="security_policy_action_permit" required="Y">policy_action</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:invoke.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-table" targetPackage="wecmdb" targetEntity="resource_set">
            <interface action="create" path="/qcloud/v1/route-table/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-table/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="associate-subnet" path="/qcloud/v1/route-table/associate-subnet">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code" required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet" targetPackage="wecmdb" targetEntity="resource_set">
            <interface action="create" path="/qcloud/v1/subnet/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.internet_ip_segment>wecmdb:network_segment.code" required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.network_segment>wecmdb:network_segment.code" required="Y">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code">route_table_id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.subnet_asset_code" required="Y">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.routing_table_asset_code" required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_set.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="create" path="/qcloud/v1/vm/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.code" required="N">instance_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="Y">instance_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_system" required="Y">image_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_instance_spec" required="Y">host_type</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.storage" required="Y">system_disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.charge_type" required="Y">instance_charge_type</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.billing_cycle" required="Y">instance_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="N">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="N">project_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.cpu">cpu</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.memory">memory</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password">password</parameter>
                    <parameter datatype="string" mappingType="context">instance_state</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group-to-vm" path="/qcloud/v1/vm/bind-security-groups">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance~(resource_instance)wecmdb:business_app_instance.unit>wecmdb:unit.security_group_asset_code" required="Y">security_group_ids</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage" targetPackage="wecmdb" targetEntity="block_storage">
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
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
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
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>           
	    <interface action="buy-and-mount-cbs-disk" path="/qcloud/v1/cbs/create-mount">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.block_storage_type" required="Y">disk_type</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.disk_size" required="Y">disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.key_name" required="N">disk_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.key_name" required="Y">disk_charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.billing_cycle" required="Y">disk_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.asset_code" required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.id" required="Y">instance_guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.file_system" required="Y">file_system_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point" required="Y">mount_dir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.name">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_code">disk_id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
		        </outputParameters>
            </interface>
            <interface action="umount-destroy-cbs-disk" path="/qcloud/v1/cbs/umount-terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_code" required="Y">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.name" required="Y">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point" required="Y">mount_dir</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.asset_code" required="Y">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.id" required="Y">instance_guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway" targetPackage="wecmdb" targetEntity="network_zone_link">
            <interface action="create" path="/qcloud/v1/nat-gateway/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.max_concurrent" required="Y">max_concurrent</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.netband_width" required="Y">bandwidth</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.internet_ip>wecmdb:ip_address.code">eip</parameter>
                    <parameter datatype="string" mappingType="context">eip_id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.asset_code" required="Y">id</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="create" path="/qcloud/v1/mysql/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_system" required="Y">engine_version</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.memory" required="Y">memory</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.storage" required="Y">volume</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.key_name" required="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.charge_type" required="Y">charge_type</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.billing_cycle" required="Y">charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name" required="Y">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="N">password</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="mysql_character_set" required="Y">character_set</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="mysql_lower_case_table_names" required="Y">lower_case_table_names</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code">private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.login_port">private_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password">password</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql/restart">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mariadb" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="create" path="/qcloud/v1/mariadb/create">
            <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name" required="Y">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="N">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.code" required="Y">zones</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.code" required="Y">node_count</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.memory" required="Y">memory_size</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.storage" required="Y">storage_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.billing_cycle" required="Y">charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_system" required="Y">db_version</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="mysql_character_set" required="Y">character_set</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="mysql_lower_case_table_names" required="Y">lower_case_table_names</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.login_port">private_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code">private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password">password</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
		</outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy" targetPackage="wecmdb" targetEntity="routing_rule">
            <interface action="create" path="/qcloud/v1/route-policy/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.asset_code" required="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.resource_set>wecmdb:resource_set.routing_table_asset_code" required="Y">route_table_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.dest_network_segment>wecmdb:network_segment.code" required="Y">dest_cidr</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.network_zone_link_design>wecmdb:network_zone_link.network_zone_link_type" required="Y">gateway_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.network_zone_link_design>wecmdb:network_zone_link.asset_code" required="Y">gateway_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.description" required="N">desc</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.asset_code">id</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.asset_code" required="Y">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.resource_set>wecmdb:resource_set.routing_table_asset_code" required="Y">route_table_id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="create" path="/qcloud/v1/redis/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_instance_type" required="Y">type_id</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.memory" required="Y">mem_size</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.billing_cycle" required="Y">period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.charge_type" required="Y">billing_mode</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code">vip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.login_port">port</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
		</outputParameters>
            </interface>
        </plugin>
        <plugin name="clb" targetPackage="wecmdb" targetEntity="service">
            <interface action="create" path="/qcloud/v1/clb/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:service.key_name" required="N">name</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:service.service_type" required="Y">type</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.asset_code" required="Y">vpc_id</parameter>
                    <parameter datatype="number" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.subnet_asset_code" required="Y">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.asset_code" required="N">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.asset_code">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.service_ip>wecmdb:ip_address.code">vip</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.asset_code" required="Y">id</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
		            <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target" targetPackage="wecmdb" targetEntity="service">
            <interface action="add-backtarget" path="/qcloud/v1/clb-target/add-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.asset_code" required="Y">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.service_port" required="Y">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.service_type" required="Y">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.unit>wecmdb:unit~(unit)wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.asset_code" required="Y">host_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.unit>wecmdb:unit~(unit)wecmdb:business_app_instance.port" required="Y">host_port</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
		            <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
            <interface action="del-backtarget" path="/qcloud/v1/clb-target/del-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.resource_set>wecmdb:resource_set.business_zone>wecmdb:business_zone.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.auth_parameter" required="Y">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.asset_code" required="Y">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.service_port" required="Y">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.service_type" required="Y">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.removed_business_app_instance>wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.asset_code" required="Y">host_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.removed_business_app_instance>wecmdb:business_app_instance.port" required="Y">host_port</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:service.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>code</parameter>
                    <parameter datatype="string" mappingType='context'>msg</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
