import { ethers } from "ethers";
import { Presets, Client } from "userop";
import { fundIfRequired } from "../src/helpers";
import { erc20ABI, testGasABI } from "../src/abi";
import { errorCodes } from "../src/errors";
import config from "../config";

describe("Without Paymaster", () => {
  const provider = new ethers.providers.JsonRpcProvider(config.nodeUrl);
  const signer = new ethers.Wallet(config.signingKey);
  const testToken = new ethers.Contract(
    config.testERC20Token,
    erc20ABI,
    provider
  );
  const testGas = new ethers.Contract(config.testGas, testGasABI, provider);
  let client: Client;
  let acc: Presets.Builder.SimpleAccount;
  beforeAll(async () => {
    client = await Client.init(config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
    client.waitTimeoutMs = 2000;
    client.waitIntervalMs = 100;
    acc = await Presets.Builder.SimpleAccount.init(signer, config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
    await fundIfRequired(
      provider,
      testToken,
      await signer.getAddress(),
      acc.getSender()
    );
  });

  test("Sender can transfer 0 ETH", async () => {
    const response = await client.sendUserOperation(
      acc.execute(acc.getSender(), 0, "0x")
    );
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender can transfer half ETH balance", async () => {
    const balance = await provider.getBalance(acc.getSender());
    const response = await client.sendUserOperation(
      acc.execute(acc.getSender(), balance.div(2), "0x")
    );
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender can transfer max ETH balance minus gas fee", async () => {
    const balance = await provider.getBalance(acc.getSender());
    const op = await client.buildUserOperation(
      acc.execute(acc.getSender(), balance.div(2), "0x")
    );
    const maxFee = ethers.BigNumber.from(op.maxFeePerGas).mul(
      ethers.BigNumber.from(op.preVerificationGas)
        .add(op.callGasLimit)
        .add(ethers.BigNumber.from(op.verificationGasLimit).mul(3))
    );
    const response = await client.sendUserOperation(
      acc.execute(acc.getSender(), balance.sub(maxFee), "0x")
    );
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender can transfer max balance of ERC20 token", async () => {
    const balance = await testToken.balanceOf(acc.getSender());
    const response = await client.sendUserOperation(
      acc.execute(
        config.testERC20Token,
        0,
        testToken.interface.encodeFunctionData("transfer", [
          acc.getSender(),
          balance,
        ])
      )
    );
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender can batch 30 ERC20 token transfers", async () => {
    const balance = await testToken.balanceOf(acc.getSender());
    const to: Array<string> = [];
    const data: Array<ethers.BytesLike> = [];
    for (let i = 0; i < 30; i++) {
      to.push(config.testERC20Token);
      data.push(
        testToken.interface.encodeFunctionData("transfer", [
          acc.getSender(),
          balance,
        ])
      );
    }
    const response = await client.sendUserOperation(acc.executeBatch(to, data));
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender cannot exceed the max batch gas limit", async () => {
    expect.assertions(1);
    try {
      await client.sendUserOperation(
        acc.execute(
          config.testGas,
          0,
          testGas.interface.encodeFunctionData("recursiveCall", [32, 32])
        )
      );
    } catch (error: any) {
      expect(error?.error.code).toBe(errorCodes.invalidFields);
    }
  });

  describe("With increasing call stack size", () => {
    describe("With zero value calls", () => {
      [0, 2, 4, 8, 16].forEach((depth) => {
        test(`Sender can make contract interactions with ${depth} recursive calls`, async () => {
          const response = await client.sendUserOperation(
            acc.execute(
              config.testGas,
              0,
              testGas.interface.encodeFunctionData("recursiveCall", [
                depth,
                depth,
              ])
            )
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });

    describe("With non-zero value calls", () => {
      [0, 2, 4, 8, 16].forEach((depth) => {
        test(`Sender can make contract interactions with ${depth} recursive calls`, async () => {
          const response = await client.sendUserOperation(
            acc.execute(
              config.testGas,
              ethers.utils.parseEther("0.001"),
              testGas.interface.encodeFunctionData("recursiveCall", [
                depth,
                depth,
              ])
            )
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });
  });
});
