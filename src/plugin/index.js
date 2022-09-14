const core = require('@serverless-devs/core');
const {validate} = require('./validate');
const {FunctionHelper} = require('./fcClient');
const {RdsHelper} = require('./rdsClient');

const {lodash, Logger} = core;
const logger = new Logger('fc-resource-creator');

/**
 * Plugin 插件入口
 * @param inputs 组件的入口参数
 * @param args 插件的自定义参数 {url: 请求的url, method: 请求方式(默认head), interval: 请求的频率（默认2m）}
 * @return inputs
 */

module.exports = async function index(inputs, args = {}) {
    logger.debug(`inputs params: ${JSON.stringify(inputs)}`);
    logger.debug(`args params: ${JSON.stringify(args)}`);
    const params = {inputs, args};


    validate(params);
    const serviceName = args.service_name;
    const functionName = args.function_name;
    const payload = {
        // invoke_type：0为创建资源，1为删除资源，此处只需要创建。
        invoke_type: 0,
        oss_config: {
            oss_bucket: args.oss_bucket,
            oss_object_name: args.oss_object_name,
            oss_region: inputs.props.region
        }
    }

    const config = {
        endpoint: `${inputs.credentials.AccountID}.${args.function_region}.fc.aliyuncs.com`,
        accessKeyId: inputs.credentials && inputs.credentials.AccessKeyID,
        accessKeySecret: inputs.credentials && inputs.credentials.AccessKeySecret,
        securityToken: inputs.credentials && inputs.credentials.SecurityToken,
        readTimeout: 1200000,
        regionId: args.function_region,
        accountID: inputs.credentials.AccountID
    };


    const fcClient = new FunctionHelper(config);
    let body;
    await logger.task('Finish', [
        {
            title: 'Invoke resource creator function',
            id: 'invoking function',
            task: async () => {
                body = await fcClient.invoke(
                    serviceName,
                    functionName,
                    payload
                );
            },
        },
    ]);

    const Output = JSON.parse(body.toString());
    if (lodash.get(Output, 'status') != 'SUCCESS') {
        logger.error(`Create resource error, operations: ${JSON.stringify(Output)}`);
        throw new Error(`Create resource error, operations: ${JSON.stringify(Output)}`);
    }
    const result = lodash.get(Output, 'result');
    const terraformOut = JSON.parse(result);
    const resourceConfig = lodash.get(terraformOut, 'outputs');
    logger.debug(`instance info: ${resourceConfig}`)


    const vpc = lodash.get(resourceConfig, 'VPC_ID.value')
    const vSwitch = lodash.get(resourceConfig, 'VSWITCH_ID.value')
    const securityGroup = lodash.get(resourceConfig, 'SECURITY_GROUP_ID.value')

    const dbId = lodash.get(resourceConfig, 'DB_ID.value')
    const resourceGroupId = lodash.get(resourceConfig, 'RESOURCE_ID.value')


    const rdsHelper = new RdsHelper(inputs.credentials)
    const rdsParams = {
        dbId, resourceGroupId,
    }

    let secretArn = await rdsHelper.checkSecret(rdsParams);
    if (secretArn == null) {
        secretArn = await rdsHelper.createSecret(rdsParams)
    }
    logger.debug(`secretArn: ${secretArn}`)
    const resourceArn = `acs:rds:cn-hangzhou:${inputs.credentials.AccountID}:dbinstance/${dbId}`
    logger.debug(`resourceArn: ${resourceArn}`)


    inputs = lodash.merge(inputs, {
        props: {
            service: {
                vpcConfig: {
                    vpcId: vpc,
                    securityGroupId: securityGroup,
                    vswitchIds: [vSwitch]
                }
            },
            function: {
                environmentVariables: {
                    ENDPOINT: args.endpoint,
                    DATABASE: args.databaseName,
                    RESOURCE_ARN: resourceArn,
                    SECRET_ARN: secretArn,
                    PHOTOVIEW_DEVELOPMENT_MODE: 1
                },
            },
        },
    });

    return inputs;
};
