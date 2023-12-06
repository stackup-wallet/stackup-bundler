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
  testGas: "0x450d8479B0ceF1e6933DED809e12845aF413A50D",
  testAccount: "0x6D7d359cE9e60dDa36EE712cE9B5947B4C72F862",
  testPaymaster: "0xa9C7F67D5Be8A805dC80f06E49BDe939384E300b",
};

export default config;
