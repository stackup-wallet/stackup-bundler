import { ethers } from "ethers";
import {
  Constants,
  Presets,
  UserOperationBuilder,
  BundlerJsonRpcProvider,
  IPresetBuilderOpts,
  UserOperationMiddlewareFn,
} from "userop";
import { EntryPoint, EntryPoint__factory } from "userop/dist/typechain";
import { testAccountABI } from "./abi";
import config from "../config";

const RECURSIVE_CALL_MODE = "0x0001";
const FORCE_VALIDATION_OOG_MODE = "0x0002";

export class TestAccount extends UserOperationBuilder {
  private provider: ethers.providers.JsonRpcProvider;
  private entryPoint: EntryPoint;
  private account: ethers.Contract;

  private constructor(
    address: string,
    rpcUrl: string,
    opts?: IPresetBuilderOpts
  ) {
    super();
    this.provider = new BundlerJsonRpcProvider(rpcUrl).setBundlerRpc(
      opts?.overrideBundlerRpc
    );
    this.entryPoint = EntryPoint__factory.connect(
      opts?.entryPoint || Constants.ERC4337.EntryPoint,
      this.provider
    );
    this.account = new ethers.Contract(address, testAccountABI, this.provider);
  }

  private resolveAccount: UserOperationMiddlewareFn = async (ctx) => {
    ctx.op.nonce = await this.entryPoint.getNonce(ctx.op.sender, 0);
  };

  public static async init(
    address: string,
    rpcUrl: string,
    opts?: IPresetBuilderOpts
  ): Promise<TestAccount> {
    const instance = new TestAccount(address, rpcUrl, opts);

    const base = instance
      .useDefaults({ sender: instance.account.address })
      .useMiddleware(instance.resolveAccount)
      .useMiddleware(Presets.Middleware.getGasPrice(instance.provider));

    return opts?.paymasterMiddleware
      ? base.useMiddleware(opts.paymasterMiddleware)
      : base.useMiddleware(
          Presets.Middleware.estimateUserOperationGas(instance.provider)
        );
  }

  recursiveCall(depth: number, width: number, discount: number) {
    return this.setCallData(
      ethers.utils.defaultAbiCoder.encode(
        ["uint256", "uint256", "uint256"],
        [depth, width, discount]
      )
    ).setSignature(RECURSIVE_CALL_MODE);
  }

  forceValidationOOG(wasteGasMultiplier: number) {
    return this.setCallData(
      ethers.utils.defaultAbiCoder.encode(["uint256"], [wasteGasMultiplier])
    ).setSignature(FORCE_VALIDATION_OOG_MODE);
  }

  forcePostOpValidationOOG(wasteGasMultiplier: number) {
    return this.setPaymasterAndData(
      ethers.utils.hexConcat([
        config.testPaymaster,
        ethers.utils.defaultAbiCoder.encode(["uint256"], [wasteGasMultiplier]),
      ])
    );
  }
}
