const core = require('@alicloud/pop-core');
const {lodash, Logger} = require("@serverless-devs/core");
const logger = new Logger('fc-resource-creator');

class RdsHelper {
    constructor(config) {
        this.client = new core({
            accessKeyId: config.AccessKeyID,
            accessKeySecret: config.AccessKeySecret,
            securityToken: config.SecurityToken, // use STS Token
            endpoint: 'https://rds.aliyuncs.com',
            apiVersion: '2014-08-15'
        });

    }

    async checkSecret(rdsParams) {
        try {
            let result = null
            const params = {
                "PageNumber": 1,
                "PageSize": 1000,
                "RegionId": "cn-hangzhou",
                "Engine": "MySQL",
                "DbInstanceId": rdsParams.dbId
            }

            const requestOption = {
                method: 'POST',
                formatParams: false,

            };

            const response = await this.client.request('DescribeSecrets', params, requestOption)
            logger.debug(JSON.stringify(response))
            const secretsArray = lodash.get(response, 'Secrets')
            await secretsArray.forEach(function(secret) {
               if( secret.SecretName === `photoview`){
                   result = secret.SecretArn
               }
            });
            return result
        } catch (e) {
            throw new Error(e);
        }
    }

    async createSecret(rdsParams) {
        try {
            const params = {
                "RegionId": "cn-hangzhou",
                "Username": "photoview",
                "Password": "Photoview2022",
                "DbInstanceId": rdsParams.dbId,
                "ResourceGroupId": rdsParams.resourceGroupId,
                "SecretName": "photoview",
                "Engine": "mysql"
            }
            const requestOption = {
                method: 'POST',
                formatParams: false,
                // 超时设置，仅对当前请求有效。该产品部分接口调用比较慢，请您适当调整超时时间。
                timeout: 50000,
            };
            const response = await this.client.request('CreateSecret', params, requestOption)
            logger.debug(JSON.stringify(response))
            return lodash.get(response, 'SecretArn')
        } catch (e) {
            throw new Error(e);
        }
    }
}

module.exports = {RdsHelper};
