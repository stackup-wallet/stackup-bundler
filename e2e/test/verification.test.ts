import { ethers } from "ethers";
import { Client, Presets } from "userop";
import { errorCodes } from "../src/errors";
import { TestAccount } from "../src/testAccount";
import config from "../config";

describe("During the verification phase", () => {
  let client: Client;
  let acc: TestAccount;
  beforeAll(async () => {
    client = await Client.init(config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
    client.waitTimeoutMs = 2000;
    client.waitIntervalMs = 100;
    acc = await TestAccount.init(config.testAccount, config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
  });

  describe("With increasing call stack size", () => {
    describe("With zero sibling stacks", () => {
      [0, 2, 4, 6, 8, 10].forEach((depth) => {
        test(`Sender can run verification with ${depth} recursive calls`, async () => {
          const response = await client.sendUserOperation(
            acc.recursiveCall(depth, 0, 0)
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });
  });

  describe("With sender dependency on callGasLimit", () => {
    [0, 1, 2, 3, 4, 5].forEach((times) => {
      test(`Sender can run validation with non-simulated code that uses ${times} storage writes`, async () => {
        const response = await client.sendUserOperation(
          acc.forceValidationOOG(times)
        );
        const event = await response.wait();

        expect(event?.args.success).toBe(true);
      });
    });
  });

  describe("With paymaster dependency on callGasLimit", () => {
    [0, 1, 2, 3, 4, 5].forEach((times) => {
      test(`Paymaster can run postOp with non-simulated code that uses ${times} storage writes`, async () => {
        const response = await client.sendUserOperation(
          acc.forcePostOpValidationOOG(times)
        );
        const event = await response.wait();

        expect(event?.args.success).toBe(true);
      });
    });
  });

  describe("With state overrides", () => {
    test("Sender with zero funds can successfully estimate UserOperation gas", async () => {
      expect.assertions(4);
      const randAcc = await Presets.Builder.SimpleAccount.init(
        new ethers.Wallet(ethers.utils.randomBytes(32)),
        config.nodeUrl,
        {
          overrideBundlerRpc: config.bundlerUrl,
        }
      );
      randAcc
        .setPreVerificationGas(0)
        .setVerificationGasLimit(0)
        .setCallGasLimit(0);

      try {
        await client.sendUserOperation(
          randAcc.execute(randAcc.getSender(), ethers.constants.Zero, "0x")
        );
      } catch (error: any) {
        expect(error?.error.code).toBe(errorCodes.rejectedByEpOrAccount);
      }

      await client.sendUserOperation(
        randAcc.execute(randAcc.getSender(), ethers.constants.Zero, "0x"),
        {
          dryRun: true,
          stateOverrides: {
            [randAcc.getSender()]: {
              balance: ethers.constants.MaxUint256.toHexString(),
            },
          },
          onBuild(op) {
            expect(ethers.BigNumber.from(op.preVerificationGas).gte(0)).toBe(
              true
            );
            expect(ethers.BigNumber.from(op.verificationGasLimit).gte(0)).toBe(
              true
            );
            expect(ethers.BigNumber.from(op.callGasLimit).gte(0)).toBe(true);
          },
        }
      );
    });
  });
});
