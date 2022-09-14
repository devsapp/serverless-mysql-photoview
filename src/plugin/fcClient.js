const FC = require("@alicloud/fc2");

class FunctionHelper {
  constructor(config) {
    this.client = new FC( config.accountID, {
      accessKeyID: config.accessKeyId,
      accessKeySecret: config.accessKeySecret,
      securityToken: config.securityToken,
      region: config.regionId,
      timeout: 600 * 1000,
    });
  }

  async invoke(service, functionName, payload) {
    try {
      const response = await this.client.invokeFunction(
        service,
        functionName,
          JSON.stringify(payload)
      );
      return response.data;
    } catch (e) {
      throw new Error(e);
    }
  }
}

module.exports = { FunctionHelper };
