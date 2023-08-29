interface IConfig {
  signingKey: string;
  nodeUrl: string;
  bundlerUrl: string;
  testERC20Token: string;
  testGas: string;
  testAccount: string;
}

const config: IConfig = {
  // This is for testing only. DO NOT use in production.
  signingKey:
    "c6cbc5ffad570fdad0544d1b5358a36edeb98d163b6567912ac4754e144d4edb",
  nodeUrl: "http://localhost:8545",
  bundlerUrl: "http://localhost:4337",

  // https://github.com/stackup-wallet/contracts/blob/main/contracts/test
  testERC20Token: "0x3870419Ba2BBf0127060bCB37f69A1b1C090992B",
  testGas: "0xc2e76Ee793a194Dd930C18c4cDeC93E7C75d567C",
  testAccount: "0x3dFD39F2c17625b301ae0EF72B411D1de5211325",
};

export default config;
