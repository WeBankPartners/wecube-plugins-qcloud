<?xml version="1.0" encoding="UTF-8"?>

<package name="wecube-plugins-qcloud" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包
    <packageDependencies>
        <packageDependency name='xxx' version='1.0'/>
        <packageDependency name='xxx233' version='1.5'/>
    </packageDependencies> -->

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单
    <menus>
        <menu code='JOBS_SERVICE_CATALOG_MANAGEMENT' cat='JOBS' displayName="Servive Catalog Management">/service-catalog</menu>
        <menu code='JOBS_TASK_MANAGEMENT' cat='JOBS' displayName="Task Management">/task-management</menu>
    </menus> -->

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系
    <dataModel>
        <entity name="service_catalogue" displayName="服务目录" description="服务目录模型">
            <attribute name="id" datatype="int" description="唯一ID"/>
            <attribute name="name" datatype="string" description="名字"/>
            <attribute name="status" datatype="string" description="状态"/>
        </entity>
    </dataModel> -->

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数
    <systemParameters>
        <systemParameter name="xxx" defaultValue='xxxx' scopeType='global'/>
        <systemParameter name="xxx" defaultValue='xxxx' scopeType='plugin-package'/>
    </systemParameters> -->

    <!-- 5.权限设定
    <authorities>
        <authority systemRoleName="admin" >
            <menu code="JOBS_SERVICE_CATALOG_MANAGEMENT" />
            <menu code="JOBS_TASK_MANAGEMENT" />
        </authority >
        <authority systemRoleName="wecube_operator" >
            <menu code="JOBS_TASK_MANAGEMENT" />
        </authority >
    </authorities> -->

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{IMAGENAME}}" portBindings="{{PORTBINDINGS}}" volumeBindings="/etc/localtime:/etc/localtime,/home/app/wecube-plugins-qcloud/logs:/home/app/wecube-plugins-qcloud/logs" envVariables=""/>
        <!-- <mysql schema="service_management" initFileName="init.sql" upgradeFileName="upgrade.sql"/>
        <s3 bucketName="service_management"/> -->
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="vpc">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/vpc/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">cidr_block</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/vpc/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="peering-connection">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/peering-connection/create">
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
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/peering-connection/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">peer_provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="security-group">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/security-group/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">description</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/security-group/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="create-policies" path="/wecube-plugins-qcloud/v1/qcloud/security-group/create-policies">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">description</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">policy_type</parameter>
                    <parameter datatype="string">policy_cidr_block</parameter>
                    <parameter datatype="string">policy_protocol</parameter>
                    <parameter datatype="string">policy_port</parameter>
                    <parameter datatype="string">policy_action</parameter>
                    <parameter datatype="string">policy_description</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="delete-policies" path="/wecube-plugins-qcloud/v1/qcloud/security-group/delete-policies">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">policy_type</parameter>
                    <parameter datatype="string">policy_cidr_block</parameter>
                    <parameter datatype="string">policy_protocol</parameter>
                    <parameter datatype="string">policy_port</parameter>
                    <parameter datatype="string">policy_action</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="route-table">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/route-table/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/route-table/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="associate-subnet" path="/wecube-plugins-qcloud/v1/qcloud/route-table/associate-subnet">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">subnet_id</parameter>
                    <parameter datatype="string">route_table_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="subnet">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/subnet/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">cidr_block</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/subnet/create-with-routetable">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">cidr_block</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">route_table_id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/subnet/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/subnet/terminate-with-routetable">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">route_table_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="vm">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/vm/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                    <parameter datatype="string">subnet_id</parameter>
                    <parameter datatype="string">instance_name</parameter>
                    <parameter datatype="string">instance_type</parameter>
                    <parameter datatype="string">image_id</parameter>
                    <parameter datatype="string">host_type</parameter>
                    <parameter datatype="number">system_disk_size</parameter>
                    <parameter datatype="string">instance_charge_type</parameter>
                    <parameter datatype="number">instance_charge_period</parameter>
                    <parameter datatype="string">instance_private_ip</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">project_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">cpu</parameter>
                    <parameter datatype="string">memory</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">instance_state</parameter>
                    <parameter datatype="string">instance_private_ip</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/vm/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="start" path="/wecube-plugins-qcloud/v1/qcloud/vm/start">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="stop" path="/wecube-plugins-qcloud/v1/qcloud/vm/stop">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="bind security group to vm" path="/wecube-plugins-qcloud/v1/qcloud/vm/bind-security-groups">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">instance_id</parameter>
                    <parameter datatype="string">security_group_ids</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="storage">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/storage/create">
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
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/storage/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="buy cbs disk and mount" path="/wecube-plugins-qcloud/v1/qcloud/cbs/create-mount">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">disk_type</parameter>
                    <parameter datatype="number">disk_size</parameter>
                    <parameter datatype="string">disk_name</parameter>
                    <parameter datatype="string">disk_charge_type</parameter>
                    <parameter datatype="string">disk_charge_period</parameter>
                    <parameter datatype="string">instance_id</parameter>
                    <parameter datatype="string">instance_guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">file_system_type</parameter>
                    <parameter datatype="string">mount_dir</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">volume_name</parameter>
                    <parameter datatype="string">disk_id</parameter>
                </output-parameters>
            </interface>
            <interface name="umount and destroy cbs disk" path="/wecube-plugins-qcloud/v1/qcloud/cbs/umount-terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">volume_name</parameter>
                    <parameter datatype="string">mount_dir</parameter>
                    <parameter datatype="string">instance_id</parameter>
                    <parameter datatype="string">instance_guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">password</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>

        </plugin>
        <plugin name="nat-gateway">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/nat-gateway/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                    <parameter datatype="number">max_concurrent</parameter>
                    <parameter datatype="number">bandwidth</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">eip</parameter>
                    <parameter datatype="string">eip_id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/nat-gateway/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="mysql-vm">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/mysql-vm/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">engine_version</parameter>
                    <parameter datatype="number">memory</parameter>
                    <parameter datatype="number">volume</parameter>
                    <parameter datatype="string">vpc_id</parameter>
                    <parameter datatype="string">subnet_id</parameter>
                    <parameter datatype="string">name</parameter>
                    <parameter datatype="string">charge_type</parameter>
                    <parameter datatype="number">charge_period</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">character_set</parameter>
                    <parameter datatype="string">lower_case_table_names</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">private_ip</parameter>
                    <parameter datatype="string">private_port</parameter>
                    <parameter datatype="string">user_name</parameter>
                    <parameter datatype="string">password</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/mysql-vm/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="restart" path="/wecube-plugins-qcloud/v1/qcloud/mysql-vm/restart">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="mariadb">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/mariadb/create">
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
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">private_port</parameter>
                    <parameter datatype="string">private_ip</parameter>
                    <parameter datatype="string">user_name</parameter>
                    <parameter datatype="string">password</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="route-policy">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/route-policy/create">
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
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/route-policy/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">route_table_id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="redis">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/redis/create">
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
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="clb">
            <interface name="create" path="/wecube-plugins-qcloud/v1/qcloud/clb/create">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="number">name</parameter>
                    <parameter datatype="number">type</parameter>
                    <parameter datatype="number">vpc_id</parameter>
                    <parameter datatype="number">subnet_id</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">id</parameter>
                    <parameter datatype="string">vip</parameter>
                </output-parameters>
            </interface>
            <interface name="terminate" path="/wecube-plugins-qcloud/v1/qcloud/clb/terminate">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">id</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="add-backtarget" path="/wecube-plugins-qcloud/v1/qcloud/clb/add-backtarget">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">lb_id</parameter>
                    <parameter datatype="string">lb_port</parameter>
                    <parameter datatype="string">protocol</parameter>
                    <parameter datatype="string">host_id</parameter>
                    <parameter datatype="string">host_port</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="del-backtarget" path="/wecube-plugins-qcloud/v1/qcloud/clb/del-backtarget">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">provider_params</parameter>
                    <parameter datatype="string">lb_id</parameter>
                    <parameter datatype="string">lb_port</parameter>
                    <parameter datatype="string">protocol</parameter>
                    <parameter datatype="string">host_id</parameter>
                    <parameter datatype="string">host_port</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
    </plugins>
</package>