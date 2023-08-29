import { ethers } from "ethers";
import { Constants } from "userop";

export const fundIfRequired = async (
  provider: ethers.providers.JsonRpcProvider,
  token: ethers.Contract,
  bundler: string,
  account: string,
  testAccount: string
) => {
  const signer = provider.getSigner(0);
  const [
    bundlerBalance,
    accountBalance,
    testAccountBalance,
    accountTokenBalance,
  ] = await Promise.all([
    provider.getBalance(bundler),
    provider.getBalance(account),
    provider.getBalance(testAccount),
    token.balanceOf(account) as ethers.BigNumber,
  ]);

  if (bundlerBalance.eq(0)) {
    const response = await signer.sendTransaction({
      to: bundler,
      value: ethers.constants.WeiPerEther,
    });
    await response.wait();
    console.log("Funded Bundler with 1 ETH...");
  }

  if (accountBalance.eq(0)) {
    const response = await signer.sendTransaction({
      to: account,
      value: ethers.constants.WeiPerEther.mul(2),
    });
    await response.wait();
    console.log("Funded Account with 2 ETH...");
  }

  if (testAccountBalance.eq(0)) {
    const response = await signer.sendTransaction({
      to: testAccount,
      value: ethers.constants.WeiPerEther.mul(2),
    });
    await response.wait();
    console.log("Funded Test Account with 2 ETH...");
  }

  if (accountTokenBalance.eq(0)) {
    const response = await signer.sendTransaction({
      to: token.address,
      value: 0,
      data: token.interface.encodeFunctionData("mint", [
        account,
        ethers.utils.parseUnits("10", 6),
      ]),
    });
    await response.wait();
    console.log("Minted 10 Test Tokens for Account...");
  }
};

export const getCallGasLimitBenchmark = async (
  provider: ethers.providers.JsonRpcProvider,
  sender: string,
  callData: ethers.BytesLike
) => {
  return provider.estimateGas({
    from: Constants.ERC4337.EntryPoint,
    to: sender,
    data: callData,
  });
};
