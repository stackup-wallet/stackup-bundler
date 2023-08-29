import { Client } from "userop";
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
});
