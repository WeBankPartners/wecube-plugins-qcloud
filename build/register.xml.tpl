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
            <systemParameter name="QCLOUD_MYSQL_BACKUP_TYPE" scopeType="global" defaultValue="logical"/>
            <systemParameter name="QCLOUD_MYSQL_CHARACTER_SET" scopeType="global" defaultValue="UTF8"/>
            <systemParameter name="QCLOUD_MYSQL_LOWER_CASE_TABLE_NAMES" scopeType="global" defaultValue="0"/>
            <systemParameter name="QCLOUD_API_SECRET" scopeType="global" defaultValue="SecretID=XXXX;SecretKey=XXXX"/>
            <systemParameter name="QCLOUD_DELETE_LB_LISTENER" scopeType="global" defaultValue="Y"/>
            <systemParameter name="QCLOUD_NOT_DELETE_LB_LISTENER" scopeType="global" defaultValue="N"/>
            <systemParameter name="QCLOUD_SECURITY_POLICY_ACTION_EGRESS" scopeType="global" defaultValue="egress"/>
            <systemParameter name="QCLOUD_SECURITY_POLICY_ACTION_INGRESS" scopeType="global" defaultValue="ingress"/>
            <systemParameter name="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION" scopeType="global" defaultValue="accept"/>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="{{PORTBINDINGS}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/qcloud/logs:/home/app/qcloud/logs" envVariables="http_proxy={{HTTP_PROXY}},https_proxy={{HTTPS_PROXY}},HTTP_PROXY={{HTTP_PROXY}},HTTPS_PROXY={{HTTPS_PROXY}}"/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="vpc" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/vpc/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/peering-connection/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_uin</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">bandwidth</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/nat-gateway/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">max_concurrent</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">bandwidth</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">eip_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/subnet/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-table" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/route-table/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-table/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="associate-subnet" path="/qcloud/v1/route-table/associate-subnet" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/route-policy/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">dest_cidr</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">gateway_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">gateway_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">desc</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/security-group/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">description</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/vm/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_family</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">image_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">system_disk_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_charge_period</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_private_ip</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">project_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">cpu</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">memory</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_private_ip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    	        </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group" path="/qcloud/v1/vm/add-security-groups" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove-security-group" path="/qcloud/v1/vm/remove-security-groups" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
    	    <interface action="buy-and-mount-cbs-disk" path="/qcloud/v1/cbs/create-mount" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_size</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_charge_period</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">file_system_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mount_dir</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">volume_name</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">disk_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		    </outputParameters>
            </interface>
            <interface action="umount-destroy-cbs-disk" path="/qcloud/v1/cbs/umount-terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">volume_name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mount_dir</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/mysql/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">engine_version</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">memory_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_role</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">master_region</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">master_instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">volume_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">charge_period</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_CHARACTER_SET">character_set</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_LOWER_CASE_TABLE_NAMES">lower_case_table_names</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">private_ip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">private_port</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		 </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql/restart" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-security-group" path="/qcloud/v1/mysql/bind-security-group" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-instance-backup" path="/qcloud/v1/mysql/create-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_BACKUP_TYPE">backup_method</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">backup_database</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">backup_table</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">backup_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>

                </outputParameters>
            </interface>
             <interface action="delete-instance-backup" path="/qcloud/v1/mysql/delete-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">backup_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mariadb" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/mariadb/create" filterRule="">
            <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">zones</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">node_count</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">memory_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">storage_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">charge_period</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">db_version</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_CHARACTER_SET">character_set</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_LOWER_CASE_TABLE_NAMES">lower_case_table_names</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">private_port</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">private_ip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user_name</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		    </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/redis/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">type_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">mem_size</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">period</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">billing_mode</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">instance_name</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		    </outputParameters>
            </interface>
            <interface action="delete" path="/qcloud/v1/redis/delete" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/clb/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">vip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		 </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
    		        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="add-backtarget" path="/qcloud/v1/clb-target/add-backtarget" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">lb_id</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">lb_port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">protocol</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host_ids</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host_ports</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">listener_id</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="del-backtarget" path="/qcloud/v1/clb-target/del-backtarget" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">provider_params</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">lb_id</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">lb_port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">protocol</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host_ids</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host_ports</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">location</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">delete_listener</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
        </plugin>

        <!-- 最佳实践 -->
         <plugin name="vpc" targetPackage="wecmdb" targetEntity="network_segment" registerName="network_segment" targetEntityFilterRule="{network_segment_usage eq 'VPC'}">
            <interface action="create" path="/qcloud/v1/vpc/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vpc/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.vpc_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="peering-connection" targetPackage="wecmdb" targetEntity="network_link" registerName="network_link" targetEntityFilterRule="{code eq 'peer'}">
            <interface action="create" path="/qcloud/v1/peering-connection/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_link.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.key_name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.NONE">peer_provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_1>wecmdb:network_segment.vpc_asset_id">peer_vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_1>wecmdb:network_segment.data_center>wecmdb:data_center.cloud_uid">peer_uin</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.netband_width">bandwidth</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_1>wecmdb:network_segment.data_center>wecmdb:data_center.location">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/peering-connection/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_link.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_link.NONE">peer_provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_1>wecmdb:network_segment.data_center>wecmdb:data_center.location">peer_location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="nat-gateway" targetPackage="wecmdb" targetEntity="network_link" registerName="network_link" targetEntityFilterRule="{code eq 'nat'}">
            <interface action="create" path="/qcloud/v1/nat-gateway/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_link.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.key_name">name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.max_concurrent">max_concurrent</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.netband_width">bandwidth</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.nat_ip_address>wecmdb:ip_address.code">eip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.nat_ip_address>wecmdb:ip_address.asset_id">eip_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/nat-gateway/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity"  mappingEntityExpression="wecmdb:network_link.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="wecmdb:network_link.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity"  mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.network_segment_2>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_link.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="subnet" targetPackage="wecmdb" targetEntity="network_segment" registerName="network_segment" targetEntityFilterRule="{network_segment_usage eq 'SUBNET'}">
            <interface action="create" path="/qcloud/v1/subnet/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}{private_route_table eq 'N'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-with-routetable" path="/qcloud/v1/subnet/create-with-routetable" filterRule="{state_code eq 'created'}{fixed_date is NULL}{private_route_table eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.code">cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/subnet/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}{private_route_table eq 'N'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate-with-routetable" path="/qcloud/v1/subnet/terminate-with-routetable" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}{private_route_table eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-table" targetPackage="wecmdb" targetEntity="network_segment" registerName="subnet" targetEntityFilterRule="{network_segment_usage eq 'SUBNET'}">
            <interface action="create" path="/qcloud/v1/route-table/create" filterRule="{fixed_date is NULL}{private_route_table eq 'Y'}{route_table_asset_id eq ''}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-table/terminate" filterRule="{fixed_date is NULL}{private_route_table eq 'N'}{route_table_asset_id neq ''}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="associate-subnet" path="/qcloud/v1/route-table/associate-subnet" filterRule="{fixed_date is NULL}{private_route_table eq 'Y'}{route_table_asset_id eq ''}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.subnet_asset_id">subnet_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="route-policy" targetPackage="wecmdb" targetEntity="route" registerName="route" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/route-policy/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:route.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.dest_network_segment>wecmdb:network_segment.code">dest_cidr</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.network_link>wecmdb:network_link.network_link_type>wecmdb:network_link_type.code">gateway_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.network_link>wecmdb:network_link.asset_id">gateway_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.description">desc</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
            </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/route-policy/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:route.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.route_table_asset_id">route_table_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:route.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:routing_rule.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="wecmdb" targetEntity="network_segment" registerName="vpc" targetEntityFilterRule="{network_segment_usage eq 'VPC'}">
            <interface action="create" path="/qcloud/v1/security-group/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}{private_security_group eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">description</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id" >id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="wecmdb" targetEntity="network_segment" registerName="subnet" targetEntityFilterRule="{network_segment_usage eq 'SUBNET'}">
            <interface action="create" path="/qcloud/v1/security-group/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}{private_security_group eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.name">description</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}{private_security_group eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.security_group_asset_id" >id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:network_segment.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-group" targetPackage="wecmdb" targetEntity="unit" registerName="unit" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/security-group/create" filterRule="{fixed_date is NULL}{white_list_type neq 'N'}{security_group_asset_id eq ''}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:unit.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.key_name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.key_name">description</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.security_group_asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.resource_set>wecmdb:resource_set.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.security_group_asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/security-group/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}{white_list_type neq 'N'}{security_group_asset_id neq ''}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:unit.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.security_group_asset_id" >id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.resource_set>wecmdb:resource_set.network_zone>wecmdb:network_zone.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:unit.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="default_security_policy" registerName="default" targetEntityFilterRule="">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_type">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.policy_network_segment>wecmdb:network_segment.code">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_action">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.description">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_type">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.policy_network_segment>wecmdb:network_segment.code">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.security_policy_action">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.owner_network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:default_security_policy.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="invoke" registerName="egress_cache" targetEntityFilterRule="{invoked_resource_type eq 'CACHE'}">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:cache_instance.cache_resource_instance>wecmdb:cache_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:cache_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.description">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:cache_instance.cache_resource_instance>wecmdb:cache_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:cache_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="invoke" registerName="egress_rdb" targetEntityFilterRule="{invoked_resource_type eq 'RDB'}">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:rdb_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.description">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:rdb_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="invoke" registerName="egress_lb" targetEntityFilterRule="{invoked_resource_type eq 'LB'}">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:lb_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.description">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:lb_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="security-policy" targetPackage="wecmdb" targetEntity="invoke" registerName="egress_app" targetEntityFilterRule="{invoked_resource_type eq 'HOST'}">
            <interface action="create-policies" path="/qcloud/v1/security-policy/create-policies" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.app_resource_instance>wecmdb:app_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.description">policy_description</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-policies" path="/qcloud/v1/security-policy/delete-policies" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:invoke.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.security_group_asset_id">security_group_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_SECURITY_POLICY_ACTION_EGRESS">policy_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.app_resource_instance>wecmdb:app_resource_instance.ip_address">policy_cidr_block</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit.protocol">policy_protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port">policy_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DEFAULT_SECURITY_POLICY_ACTION">policy_action</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.invoke_unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:invoke.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/vm/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.subnet_asset_id">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.key_name">instance_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">instance_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">instance_family</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.resource_instance_system>wecmdb:resource_instance_system.code">image_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.resource_instance_spec>wecmdb:resource_instance_spec.code">host_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.storage">system_disk_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.charge_type>wecmdb:charge_type.code">instance_charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.billing_cycle">instance_charge_period</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">instance_private_ip</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">project_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.cpu">cpu</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.memory">memory</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">instance_state</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.ip_address">instance_private_ip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    	    </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/vm/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="start" path="/qcloud/v1/vm/start" filterRule="{state_code eq 'startup'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="stop" path="/qcloud/v1/vm/stop" filterRule="{state_code eq 'stoped'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-vpc-security-group" path="/qcloud/v1/vm/add-security-groups" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-subnet-security-group" path="/qcloud/v1/vm/add-security-groups" filterRule="{state_code eq 'created'}{fixed_date is NULL}{subnet_security_group eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove-vpc-security-group" path="/qcloud/v1/vm/remove-security-groups" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove-subnet-security-group" path="/qcloud/v1/vm/remove-security-groups" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}{subnet_security_group eq 'Y'}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="vm" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy" targetEntityFilterRule="">
            <interface action="bind_sg_app_created" path="/qcloud/v1/vm/add-security-groups" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove-sg-app-deleted" path="/qcloud/v1/vm/remove-security-groups" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="storage" targetPackage="wecmdb" targetEntity="block_storage" registerName="block_storage" targetEntityFilterRule="">
    	    <interface action="buy-and-mount-cbs-disk" path="/qcloud/v1/cbs/create-mount" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.storage_type>wecmdb:storage_type.code">disk_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.disk_size">disk_size</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.name">disk_name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.charge_type">disk_charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.billing_cycle">disk_charge_period</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.guid">instance_guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.file_system">file_system_type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point">mount_dir</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.code">volume_name</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id">disk_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		    </outputParameters>
            </interface>
            <interface action="umount-destroy-cbs-disk" path="/qcloud/v1/cbs/umount-terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:block_storage..NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.code">volume_name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.mount_point">mount_dir</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.asset_id">instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.guid">instance_guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.host_resource_instance>wecmdb:host_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:block_storage.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql" targetPackage="wecmdb" targetEntity="rdb_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/mysql/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.resource_instance_system>wecmdb:resource_instance_system.code">engine_version</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.resource_instance_spec>wecmdb:resource_instance_spec.code">memory_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.cluster_node_type>wecmdb:cluster_node_type.code">instance_role</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">master_region</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">master_instance_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.storage">volume_size</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.subnet_asset_id">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.key_name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.charge_type>wecmdb:charge_type.code">charge_type</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.billing_cycle">charge_period</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name">user_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_CHARACTER_SET">character_set</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_LOWER_CASE_TABLE_NAMES">lower_case_table_names</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.ip_address">private_ip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port">private_port</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name">user_name</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		 </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/mysql/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="restart" path="/qcloud/v1/mysql/restart" filterRule="{state_code eq 'startup'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="bind-vpc-security-group" path="/qcloud/v1/mysql/bind-security-group" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-instance-backup" path="/qcloud/v1/mysql/create-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_BACKUP_TYPE">backup_method</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance~(rdb_resource_instance)wecmdb:rdb_instance.unit>wecmdb:unit.code">backup_database</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">backup_table</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.backup_asset_id">backup_id</parameter>
                </outputParameters>
            </interface>
             <interface action="delete-instance-backup" path="/qcloud/v1/mysql/delete-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.backup_asset_id">backup_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="database" targetEntityFilterRule="">
            <interface action="create-deploy-backup" path="/qcloud/v1/mysql/create-backup" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_BACKUP_TYPE">backup_method</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance~(rdb_resource_instance)wecmdb:rdb_instance.unit>wecmdb:unit.code">backup_database</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.NONE">backup_table</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.backup_asset_id">backup_id</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-deploy-backup" path="/qcloud/v1/mysql/delete-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.backup_asset_id">backup_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="create-regular-backup" path="/qcloud/v1/mysql/create-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_MYSQL_BACKUP_TYPE">backup_method</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance~(rdb_resource_instance)wecmdb:rdb_instance.unit>wecmdb:unit.code">backup_database</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.NONE">backup_table</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.regular_backup_asset_id">backup_id</parameter>
                </outputParameters>
            </interface>
            <interface action="delete-regular-backup" path="/qcloud/v1/mysql/delete-backup" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.asset_id">mysql_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance._backup_asset_id">backup_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location" >location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis" targetPackage="wecmdb" targetEntity="cache_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/redis/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED">seed</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.key_name">instance_name</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.resource_instance_type>wecmdb:resource_instance_type.code">type_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.resource_instance_spec>wecmdb:resource_instance_spec.code">mem_size</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.billing_cycle">period</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.charge_type>wecmdb:charge_type.code">billing_mode</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.network_segment>wecmdb:network_segment.subnet_asset_id">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.security_group_asset_id">security_group_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">instance_name</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.ip_address">vip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.login_port">port</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.user_password">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
		        </outputParameters>
            </interface>
            <interface action="delete" path="/qcloud/v1/redis/delete" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:cache_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb" targetPackage="wecmdb" targetEntity="lb_resource_instance" registerName="resource" targetEntityFilterRule="">
            <interface action="create" path="/qcloud/v1/clb/create" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.key_name">name</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.resource_instance_type>wecmdb:resource_instance_type.code">type</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.f_network_segment>wecmdb:network_segment.vpc_asset_id">vpc_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.subnet_asset_id">subnet_id</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.ip_address">vip</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
    		    </outputParameters>
            </interface>
            <interface action="terminate" path="/qcloud/v1/clb/terminate" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.asset_id">id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_resource_instance.guid">guid</parameter>
    		        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target" targetPackage="wecmdb" targetEntity="lb_instance" registerName="whole" targetEntityFilterRule="">
            <interface action="add" path="/qcloud/v1/clb-target/add-backtarget" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id">lb_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port">lb_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.protocol">protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id">host_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port">host_ports</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_listener_asset_id">listener_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/qcloud/v1/clb-target/del-backtarget" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id">lb_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port">lb_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.protocol">protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.asset_id">host_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance.port">host_ports</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_DELETE_LB_LISTENER">delete_listener</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="clb-target" targetPackage="wecmdb" targetEntity="lb_instance" registerName="target" targetEntityFilterRule="">
            <interface action="add" path="/qcloud/v1/clb-target/add-backtarget" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id">lb_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port">lb_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.protocol">protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance{state_code eq 'created'}{fixed_date is NULL}.host_resource_instance>wecmdb:host_resource_instance.asset_id">host_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance{state_code eq 'created'}{fixed_date is NULL}.port">host_ports</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_listener_asset_id">listener_id</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/qcloud/v1/clb-target/del-backtarget" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.NONE">provider_params</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.asset_id">lb_id</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.port">lb_port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit.protocol">protocol</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance{state_code eq 'destroyed'}{fixed_date is NULL}.host_resource_instance>wecmdb:host_resource_instance.asset_id">host_ids</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.unit>wecmdb:unit~(invoke_unit)wecmdb:invoke.invoked_unit>wecmdb:unit~(unit)wecmdb:app_instance{state_code eq 'destroyed'}{fixed_date is NULL}.port">host_ports</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="QCLOUD_API_SECRET">api_secret</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.lb_resource_instance>wecmdb:lb_resource_instance.network_segment>wecmdb:network_segment.data_center>wecmdb:data_center.location">location</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="QCLOUD_NOT_DELETE_LB_LISTENER">delete_listener</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:lb_instance.guid">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
