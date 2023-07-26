interface IConfig {
  signingKey: string;
  nodeUrl: string;
  bundlerUrl: string;
  testERC20Token: string;
  testGas: string;
}

const config: IConfig = {
  // This is for testing only. DO NOT use in production.
  signingKey:
    "c6cbc5ffad570fdad0544d1b5358a36edeb98d163b6567912ac4754e144d4edb",
  nodeUrl: "http://localhost:8545",
  bundlerUrl: "http://localhost:4337",
  testERC20Token: "0x3870419Ba2BBf0127060bCB37f69A1b1C090992B",
  // https://github.com/stackup-wallet/contracts/blob/main/contracts/test/TestGas.sol
  testGas: "0x310de3671FCDA16eC24433e7610c2F48DEE7f8c2",
};

export default config;
