interface IConfig {
  signingKey: string;
  nodeUrl: string;
  bundlerUrl: string;
  testERC20Token: string;
  testGas: string;
  testAccount: string;
  testPaymaster: string;
}

const config: IConfig = {
  // This is for testing only. DO NOT use in production.
  signingKey:
    "c6cbc5ffad570fdad0544d1b5358a36edeb98d163b6567912ac4754e144d4edb",
  nodeUrl: "http://localhost:8545",
  bundlerUrl: "http://localhost:4337",

  // https://github.com/stackup-wallet/contracts/blob/main/contracts/test
  testERC20Token: "0x3870419Ba2BBf0127060bCB37f69A1b1C090992B",
  testGas: "0x3Eb396057D1eaB0aB9000c11c4E7F11e7974934f",
  testAccount: "0x6306eBAaa03a9F6DcA1d036c152EDeCF809147E3",
  testPaymaster: "0xfC7D123030203aaebcB0C72B1412Faa090b166f5",
};

export default config;
