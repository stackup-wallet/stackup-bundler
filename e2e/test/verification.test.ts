import { Client } from "userop";
import { TestAccount } from "../src/testAccount";
import config from "../config";

describe("During the verification phase", () => {
  let client: Client;
  let testAcc: TestAccount;
  beforeAll(async () => {
    client = await Client.init(config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
    client.waitTimeoutMs = 2000;
    client.waitIntervalMs = 100;
    testAcc = await TestAccount.init(config.testAccount, config.nodeUrl, {
      overrideBundlerRpc: config.bundlerUrl,
    });
  });

  describe("With increasing call stack size", () => {
    describe("With zero sibling stacks", () => {
      [0, 2, 4, 8].forEach((depth) => {
        test(`Sender can run verification with ${depth} recursive calls`, async () => {
          const response = await client.sendUserOperation(
            testAcc.recursiveCall(depth, 0, 0)
          );
          const event = await response.wait();

          expect(event?.args.success).toBe(true);
        });
      });
    });
  });
});
