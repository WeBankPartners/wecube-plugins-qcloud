<?xml version="1.0" encoding="UTF-8"?>
<package name="qcloud" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
        <packageDependency name="wecmdb" version="v1.4.0"/>
    </packageDependencies>

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单 -->
    <menus>
    </menus>

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系 -->
    <dataModel>
    </dataModel>

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数 -->
    <systemParameters>
        <systemParameter name="SECURITY_POLICY_ACTION_ACCEPT" scopeType="global" defaultValue="ACCEPT"/>
        <systemParameter name="SECURITY_POLICY_TYPE_INGRESS" scopeType="global" defaultValue="ingress"/>
        <systemParameter name="SECURITY_POLICY_TYPE_EGRESS" scopeType="global" defaultValue="egress"/>
        <systemParameter name="MYSQL_BACKUP_TYPE_LOGICAL" scopeType="global" defaultValue="logical"/>
        <systemParameter name="MYSQL_CHARACTER_SET" scopeType="global" defaultValue="UTF8"/>
        <systemParameter name="MYSQL_LOWER_CASE_TABLE_NAMES" scopeType="global" defaultValue="0"/>
	    <systemParameter name="QCLOUD_API_SECRET" scopeType="global" defaultValue="SecretID=XXXX;SecretKey=XXXX"/>
        <systemParameter name="QCLOUD_UID" scopeType="global" defaultValue="XXXX"/>
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
        <plugin name="vpc">
            <interface action="create" path="/qcloud/v1/vpc/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection">
            <interface action="create" path="/qcloud/v1/peering-connection/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">peer_vpc_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_UID" required="Y" sensitiveData="N">peer_uin</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">bandwidth</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group">
            <interface action="create" path="/qcloud/v1/security-group/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">description</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression=""  required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression=""  required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression=""  required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">security_group_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_action</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">policy_description</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">security_group_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">policy_action</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-table">
            <interface action="create" path="/qcloud/v1/route-table/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-table/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="associate-subnet" path="/qcloud/v1/route-table/associate-subnet">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet">
            <interface action="create" path="/qcloud/v1/subnet/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm">
            <interface action="create" path="/qcloud/v1/vm/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">instance_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">image_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">host_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">system_disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">instance_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">project_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">cpu</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">memory</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">instance_state</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group-to-vm" path="/qcloud/v1/vm/bind-security-groups">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">security_group_ids</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage">
	        <interface action="buy-and-mount-cbs-disk" path="/qcloud/v1/cbs/create-mount">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">disk_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">disk_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">disk_charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">disk_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">file_system_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mount_dir</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">disk_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		        </outputParameters>
            </interface>
            <interface action="umount-destroy-cbs-disk" path="/qcloud/v1/cbs/umount-terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mount_dir</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway">
            <interface action="create" path="/qcloud/v1/nat-gateway/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">max_concurrent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">bandwidth</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">eip</parameter>
                    <parameter datatype="string" mappingType="context">eip_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql">
            <interface action="create" path="/qcloud/v1/mysql/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">engine_version</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">memory_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">instance_role</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">master_region</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">master_instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">volume_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_CHARACTER_SET" required="Y" sensitiveData="N">character_set</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_LOWER_CASE_TABLE_NAMES" required="Y" sensitiveData="N">lower_case_table_names</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">private_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql/restart">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group" path="/qcloud/v1/mysql/bind-security-group">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">security_group_ids</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-instance-backup" path="/qcloud/v1/mysql/create-backup">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_BACKUP_TYPE_LOGICAL" required="Y" sensitiveData="N">backup_method</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">backup_database</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">backup_table</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">backup_id</parameter>
                </outputParameters>
            </interface>
             <interface action="delete-instance-backup" path="/qcloud/v1/mysql/delete-backup">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">backup_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mariadb">
            <interface action="create" path="/qcloud/v1/mariadb/create">
            <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">zones</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">node_count</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">memory_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">storage_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">db_version</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_CHARACTER_SET" required="Y" sensitiveData="N">character_set</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_LOWER_CASE_TABLE_NAMES" required="Y" sensitiveData="N">lower_case_table_names</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">private_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		        </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy">
            <interface action="create" path="/qcloud/v1/route-policy/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">dest_cidr</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">gateway_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">gateway_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">desc</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis">
            <interface action="create" path="/qcloud/v1/redis/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">type_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">mem_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">billing_mode</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" required="N" sensitiveData="N">security_group_ids</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">vip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		        </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb">
            <interface action="create" path="/qcloud/v1/clb/create">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">vip</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target">
            <interface action="add-backtarget" path="/qcloud/v1/clb-target/add-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">host_ids</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">host_ports</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">listener_id</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="del-backtarget" path="/qcloud/v1/clb-target/del-backtarget">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">host_ids</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">host_ports</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>




        <!-- 最佳实践 -->
         <plugin name="vpc" targetPackage="wecmdb" targetEntity="network_segment" registerName="network_segment" targetEntityFilterRule="{network_segment_usage eq 'VPC'}">
            <interface action="create" path="/qcloud/v1/vpc/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection" targetPackage="wecmdb" targetEntity="network_zone_link" registerName="network_zone_link" targetEntityFilterRule="{network_zone_link_type eq 'PEERCONNECTION'}">
            <interface action="create" path="/qcloud/v1/peering-connection/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.key_name" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_1>wecmdb:network_zone.network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">peer_vpc_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_UID" required="Y" sensitiveData="N">peer_uin</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.netband_width" required="Y" sensitiveData="N">bandwidth</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_1>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">peer_provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_1>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="wecmdb" targetEntity="network_segment" registerName="vpc" targetEntityFilterRule="{network_segment_usage eq 'VPC'}">
            <interface action="create" path="/qcloud/v1/security-group/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.key_name" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.description" required="Y" sensitiveData="N">description</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid"  required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE"  required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id"  required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="default_security_policy" registerName="vpc" targetEntityFilterRule="">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.security_group_asset_id" required="N" sensitiveData="N">security_group_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_type" required="Y" sensitiveData="N">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.policy_network_segment>wecmdb:network_segment.code" required="Y" sensitiveData="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.protocol" required="Y" sensitiveData="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.port" required="Y" sensitiveData="N">policy_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_action" required="Y" sensitiveData="N">policy_action</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.description" required="N" sensitiveData="N">policy_description</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.security_group_asset_id" required="N" sensitiveData="N">security_group_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_type" required="Y" sensitiveData="N">policy_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.policy_network_segment>wecmdb:network_segment.code" required="Y" sensitiveData="N">policy_cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.protocol" required="Y" sensitiveData="N">policy_protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.port" required="Y" sensitiveData="N">policy_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_action" required="Y" sensitiveData="N">policy_action</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet" targetPackage="wecmdb" targetEntity="network_segment" registerName="network_segment" targetEntityFilterRule="{network_segment_usage eq 'SUBNET'}">
            <interface action="create" path="/qcloud/v1/subnet/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}{private_route_table eq 'N'}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable" filterRule="{state_code eq 'created'}{fixed_date eq ''}{private_route_table eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code" required="Y" sensitiveData="N">cidr_block</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}{private_route_table eq 'N'}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}{private_route_table eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/vm/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.subnet_asset_id" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.name" required="N" sensitiveData="N">instance_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="Y" sensitiveData="N">instance_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.resource_instance_system>wecmdb:resource_instance_system.code" required="Y" sensitiveData="N">image_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.resource_instance_spec>wecmdb:resource_instance_spec.code" required="Y" sensitiveData="N">host_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.storage" required="Y" sensitiveData="N">system_disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.charge_type" required="Y" sensitiveData="N">instance_charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.billing_cycle" required="N" sensitiveData="N">instance_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="N" sensitiveData="N">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">project_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.cpu">cpu</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.memory">memory</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">instance_state</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">instance_private_ip</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start" filterRule="{state_code eq 'startup'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop" filterRule="{state_code eq 'stoped'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group-to-vm" path="/qcloud/v1/vm/bind-security-groups" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id" required="Y" sensitiveData="N">security_group_ids</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage" targetPackage="wecmdb" targetEntity="block_storage" registerName="block_storage" targetEntityFilterRule="">
	        <interface action="buy-and-mount-cbs-disk" path="/qcloud/v1/cbs/create-mount" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.block_storage_type" required="Y" sensitiveData="N">disk_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.disk_size" required="Y" sensitiveData="N">disk_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.key_name" required="N" sensitiveData="N">disk_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.charge_type" required="Y" sensitiveData="N">disk_charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.billing_cycle" required="N" sensitiveData="N">disk_charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">instance_guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.file_system" required="Y" sensitiveData="N">file_system_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point" required="Y" sensitiveData="N">mount_dir</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.name">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id">disk_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		        </outputParameters>
            </interface>
            <interface action="umount-destroy-cbs-disk" path="/qcloud/v1/cbs/umount-terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage..NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.name" required="Y" sensitiveData="N">volume_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point" required="Y" sensitiveData="N">mount_dir</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.guid" required="Y" sensitiveData="N">instance_guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway" targetPackage="wecmdb" targetEntity="network_zone_link" registerName="network_zone_link" targetEntityFilterRule="{network_zone_link_type eq 'NAT'}">
            <interface action="create" path="/qcloud/v1/nat-gateway/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.key_name" required="Y" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.max_concurrent" required="Y" sensitiveData="N">max_concurrent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.netband_width" required="Y" sensitiveData="N">bandwidth</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="context">eip</parameter>
                    <parameter datatype="string" mappingType="context">eip_id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity"  mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.network_zone_2>wecmdb:network_zone.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:network_zone_link.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql" targetPackage="wecmdb" targetEntity="rdb_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/mysql/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y" sensitiveData="N">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.resource_instance_type>wecmdb:resource_instance_type.code" required="Y" sensitiveData="N">engine_version</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.resource_instance_spec>wecmdb:resource_instance_spec.code" required="Y" sensitiveData="N">memory_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.cluster_node_type>wecmdb:cluster_node_type.code" required="Y" sensitiveData="N">instance_role</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">master_region</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">master_instance_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.storage" required="Y" sensitiveData="N">volume_size</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.subnet_asset_id" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.key_name" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.charge_type" required="Y" sensitiveData="N">charge_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.billing_cycle" required="N" sensitiveData="N">charge_period</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name" required="Y" sensitiveData="N">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_CHARACTER_SET" required="Y" sensitiveData="N">character_set</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_LOWER_CASE_TABLE_NAMES" required="Y" sensitiveData="N">lower_case_table_names</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">private_ip</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port">private_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name">user_name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql/terminate">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql/restart" filterRule="{state_code eq 'startup'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group" path="/qcloud/v1/mysql/bind-security-group" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id" required="Y" sensitiveData="N">security_group_ids</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-instance-backup" path="/qcloud/v1/mysql/create-backup">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_BACKUP_TYPE_LOGICAL" required="Y" sensitiveData="N">backup_method</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance~(rdb_resource_instance)wecmdb:rdb_instance.unit>wecmdb:unit.code" required="Y" sensitiveData="N">backup_database</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">backup_table</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location"  required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.backup_asset_id">backup_id</parameter>
                </outputParameters>
            </interface>
             <interface action="delete-instance-backup" path="/qcloud/v1/mysql/delete-backup">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id" required="Y" sensitiveData="N">mysql_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.backup_asset_id" required="Y" sensitiveData="N">backup_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location"  required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy" targetPackage="wecmdb" targetEntity="route" registerName="route" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/route-policy/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.route_table_asset_id" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.dest_network_segment>wecmdb:network_segment.code" required="Y" sensitiveData="N">dest_cidr</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.network_zone_link>wecmdb:network_zone_link.network_zone_link_design>wecmdb:network_zone_link_design.network_zone_link_type" required="Y" sensitiveData="N">gateway_type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.network_zone_link>wecmdb:network_zone_link.asset_id" required="Y" sensitiveData="N">gateway_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.description" required="N" sensitiveData="N">desc</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.route_table_asset_id" required="Y" sensitiveData="N">route_table_id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb" targetPackage="wecmdb" targetEntity="lb_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/clb/create" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.name" required="N" sensitiveData="N">name</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.resource_instance_type>wecmdb:resource_instance_type.code" required="Y" sensitiveData="N">type</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id" required="Y" sensitiveData="N">vpc_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.subnet_asset_id" required="Y" sensitiveData="N">subnet_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id" required="N" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.code">vip</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
		     </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id" required="Y" sensitiveData="N">id</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target" targetPackage="wecmdb" targetEntity="lb_instance" registerName="app" targetEntityFilterRule="">
            <interface action="add-backtarget" path="/qcloud/v1/clb-target/add-backtarget" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id" required="Y" sensitiveData="N">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port" required="Y" sensitiveData="N">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.unit_design>wecmdb:unit_design.protocol" required="Y" sensitiveData="N">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">host_ids</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port" required="Y" sensitiveData="N">host_ports</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_listener_asset_id">listener_id</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="del-backtarget" path="/qcloud/v1/clb-target/del-backtarget" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid" required="Y" sensitiveData="N">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE" required="N" sensitiveData="N">provider_params</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id" required="Y" sensitiveData="N">lb_id</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port" required="Y" sensitiveData="N">lb_port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.unit_design>wecmdb:unit_design.protocol" required="Y" sensitiveData="N">protocol</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id" required="Y" sensitiveData="N">host_ids</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port" required="Y" sensitiveData="N">host_ports</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET" required="Y" sensitiveData="Y">api_secret</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.intranet_ip>wecmdb:ip_address.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" required="Y" sensitiveData="N">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
