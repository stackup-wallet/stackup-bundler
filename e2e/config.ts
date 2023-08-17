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
  testGas: "0xdea3d7b717F3A6E634876092470efc281f9b3de7",
};

export default config;
