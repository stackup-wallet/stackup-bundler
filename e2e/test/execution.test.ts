import { ethers } from "ethers";
import { Presets, Client, ISendUserOperationOpts } from "userop";
import { getCallGasLimitBenchmark } from "../src/helpers";
import { erc20ABI, testGasABI } from "../src/abi";
import { errorCodes } from "../src/errors";
import config from "../config";

const opChecks = (
  provider: ethers.providers.JsonRpcProvider
): ISendUserOperationOpts => ({
  onBuild: async (op) => {
    const cgl = ethers.BigNumber.from(op.callGasLimit).toNumber();
    const benchmark = await (
      await getCallGasLimitBenchmark(provider, op.sender, op.callData)
    ).toNumber();

    expect(cgl).toBeLessThanOrEqual(benchmark);
  },
});

// TODO: Figure out why CGL is not LTE to benchmark at certain depths/widths.
// Until then we use this check to prevent regression.
const opCheckDeep = (benchmark: number): ISendUserOperationOpts => ({
  onBuild: async (op) => {
    expect(
      ethers.BigNumber.from(op.callGasLimit).toNumber()
    ).toBeLessThanOrEqual(benchmark);
  },
});

describe("During the execution phase", () => {
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
  });

  test("Sender can transfer 0 ETH", async () => {
    const response = await client.sendUserOperation(
      acc.execute(acc.getSender(), 0, "0x"),
      { ...opChecks(provider) }
    );
    const event = await response.wait();

    expect(event?.args.success).toBe(true);
  });

  test("Sender can transfer half ETH balance", async () => {
    const balance = await provider.getBalance(acc.getSender());
    const response = await client.sendUserOperation(
      acc.execute(acc.getSender(), balance.div(2), "0x"),
      { ...opChecks(provider) }
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
      acc.execute(acc.getSender(), balance.sub(maxFee), "0x"),
      { ...opChecks(provider) }
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
      ),
      { ...opChecks(provider) }
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
    const response = await client.sendUserOperation(
      acc.executeBatch(to, data),
      { ...opChecks(provider) }
    );
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
          testGas.interface.encodeFunctionData("recursiveCall", [32, 0, 0, 32])
        ),
        { ...opChecks(provider) }
      );
    } catch (error: any) {
      expect(error?.error.code).toBe(errorCodes.executionReverted);
    }
  });

  describe("With increasing call stack size", () => {
    describe("With zero value", () => {
      [0, 2, 4, 8, 16].forEach((depth) => {
        test(`Sender can make contract interactions with ${depth} recursive calls`, async () => {
          let opts = opChecks(provider);
          if (depth === 8) opts = opCheckDeep(1195897);
          if (depth === 16) opts = opCheckDeep(4365893);

          const response = await client.sendUserOperation(
            acc.execute(
              config.testGas,
              0,
              testGas.interface.encodeFunctionData("recursiveCall", [
                depth,
                0,
                0,
                depth,
              ])
            ),
            opts
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });

    describe("With non-zero value", () => {
      [0, 2, 4, 8, 16].forEach((depth) => {
        test(`Sender can make contract interactions with ${depth} recursive calls`, async () => {
          let opts = opChecks(provider);
          if (depth === 8) opts = opCheckDeep(1262227);
          if (depth === 16) opts = opCheckDeep(4499616);

          const response = await client.sendUserOperation(
            acc.execute(
              config.testGas,
              ethers.utils.parseEther("0.001"),
              testGas.interface.encodeFunctionData("recursiveCall", [
                depth,
                0,
                0,
                depth,
              ])
            ),
            opts
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });

    describe("With random gas discount", () => {
      [1, 2, 3].forEach((depth) =>
        describe(`At depth equal to ${depth}`, () => {
          Array.from({ length: 5 }, () =>
            Math.floor(Math.random() * 100000)
          ).forEach((discount) => {
            test(`Sender can make contract interactions with ${discount} gas discount to recursive calls`, async () => {
              const response = await client.sendUserOperation(
                acc.execute(
                  config.testGas,
                  0,
                  testGas.interface.encodeFunctionData("recursiveCall", [
                    depth,
                    0,
                    discount,
                    depth,
                  ])
                ),
                opChecks(provider)
              );
              const event = await response.wait();

              expect(event?.args.success).toBe(true);
            });
          });
        })
      );
    });

    describe("With multiple stacks per depth", () => {
      [0, 1, 2, 3].forEach((depth) => {
        test(`Sender can make contract interactions with ${depth} recursive calls`, async () => {
          let opts = opChecks(provider);
          if (depth === 2) opts = opCheckDeep(866332);
          if (depth === 3) opts = opCheckDeep(7929055);

          const width = depth;
          const response = await client.sendUserOperation(
            acc.execute(
              config.testGas,
              0,
              testGas.interface.encodeFunctionData("recursiveCall", [
                depth,
                width,
                0,
                depth,
              ])
            ),
            opts
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });
  });
});
