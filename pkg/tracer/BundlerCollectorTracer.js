// This is the same BundlerCollectorTracer from github.com/eth-infinitism/bundler transpiled down to ES5.

var tracer = {
  callsFromEntryPoint: [],
  currentLevel: null,
  keccak: [],
  calls: [],
  logs: [],
  debug: [],
  lastOp: "",
  lastThreeOpcodes: [],
  // event sent after all validations are done: keccak("BeforeExecution()")
  stopCollectingTopic:
    "bb47ee3e183a558b1a2ff0874b079f3fc5478b7454eacf2bfc5af2ff5878f972",
  stopCollecting: false,
  topLevelCallCounter: 0,

  fault: function fault(log, _db) {
    this.debug.push(
      "fault depth=",
      log.getDepth(),
      " gas=",
      log.getGas(),
      " cost=",
      log.getCost(),
      " err=",
      log.getError()
    );
  },

  result: function result(_ctx, _db) {
    return {
      callsFromEntryPoint: this.callsFromEntryPoint,
      keccak: this.keccak,
      logs: this.logs,
      calls: this.calls,
      // for internal debugging.
      debug: this.debug,
    };
  },

  enter: function enter(frame) {
    if (this.stopCollecting) {
      return;
    }
    // this.debug.push('enter gas=', frame.getGas(), ' type=', frame.getType(), ' to=', toHex(frame.getTo()), ' in=', toHex(frame.getInput()).slice(0, 500))
    this.calls.push({
      type: frame.getType(),
      from: toHex(frame.getFrom()),
      to: toHex(frame.getTo()),
      method: toHex(frame.getInput()).slice(0, 10),
      gas: frame.getGas(),
      value: frame.getValue(),
    });
  },
  exit: function exit(frame) {
    if (this.stopCollecting) {
      return;
    }
    this.calls.push({
      type: frame.getError() != null ? "REVERT" : "RETURN",
      gasUsed: frame.getGasUsed(),
      data: toHex(frame.getOutput()).slice(0, 4000),
    });
  },

  // increment the "key" in the list. if the key is not defined yet, then set it to "1"
  countSlot: function countSlot(list, key) {
    if (!list[key]) list[key] = 0;
    list[key] += 1;
  },
  step: function step(log, db) {
    if (this.stopCollecting) {
      return;
    }
    var opcode = log.op.toString();

    var stackSize = log.stack.length();
    var stackTop3 = [];
    for (var i = 0; i < 3 && i < stackSize; i++) {
      stackTop3.push(log.stack.peek(i));
    }
    this.lastThreeOpcodes.push({
      opcode: opcode,
      stackTop3: stackTop3,
    });
    if (this.lastThreeOpcodes.length > 3) {
      this.lastThreeOpcodes.shift();
    }
    // this.debug.push(this.lastOp + '-' + opcode + '-' + log.getDepth() + '-' + log.getGas() + '-' + log.getCost())
    if (
      log.getGas() < log.getCost() ||
      // special rule for SSTORE with gas metering
      (opcode === "SSTORE" && log.getGas() < 2300)
    ) {
      this.currentLevel.oog = true;
    }

    if (opcode === "REVERT" || opcode === "RETURN") {
      if (log.getDepth() === 1) {
        // exit() is not called on top-level return/revent, so we reconstruct it
        // from opcode
        var ofs = parseInt(log.stack.peek(0).toString());
        var len = parseInt(log.stack.peek(1).toString());
        var data = toHex(log.memory.slice(ofs, ofs + len)).slice(0, 4000);
        // this.debug.push(opcode + ' ' + data)
        this.calls.push({
          type: opcode,
          gasUsed: 0,
          data: data,
        });
      }
      // NOTE: flushing all history after RETURN
      this.lastThreeOpcodes = [];
    }

    if (log.getDepth() === 1) {
      if (opcode === "CALL" || opcode === "STATICCALL") {
        // stack.peek(0) - gas
        var addr = toAddress(log.stack.peek(1).toString(16));
        var topLevelTargetAddress = toHex(addr);
        // stack.peek(2) - value
        var _ofs = parseInt(log.stack.peek(3).toString());
        // stack.peek(4) - len
        var topLevelMethodSig = toHex(log.memory.slice(_ofs, _ofs + 4));

        this.currentLevel = this.callsFromEntryPoint[this.topLevelCallCounter] =
          {
            topLevelMethodSig: topLevelMethodSig,
            topLevelTargetAddress: topLevelTargetAddress,
            access: {},
            opcodes: {},
            extCodeAccessInfo: {},
            contractSize: {},
          };
        this.topLevelCallCounter++;
      } else if (opcode === "LOG1") {
        // ignore log data ofs, len
        var topic = log.stack.peek(2).toString(16);
        if (topic === this.stopCollectingTopic) {
          this.stopCollecting = true;
        }
      }
      this.lastOp = "";
      return;
    }

    var lastOpInfo = this.lastThreeOpcodes[this.lastThreeOpcodes.length - 2];
    // store all addresses touched by EXTCODE* opcodes
    if (lastOpInfo && lastOpInfo.opcode.match(/^(EXT.*)$/) != null) {
      var _addr = toAddress(lastOpInfo.stackTop3[0].toString(16));
      var addrHex = toHex(_addr);
      var last3opcodesString = this.lastThreeOpcodes
        .map(function (x) {
          return x.opcode;
        })
        .join(",");
      // only store the last EXTCODE* opcode per address - could even be a boolean for our current use-case
      // [OP-051]
      if (last3opcodesString.match(/^(\w+),EXTCODESIZE,ISZERO$/) == null) {
        this.currentLevel.extCodeAccessInfo[addrHex] = opcode;
        // this.debug.push(`potentially illegal EXTCODESIZE without ISZERO for ${addrHex}`)
      } else {
        // this.debug.push(`safe EXTCODESIZE with ISZERO for ${addrHex}`)
      }
    }

    // not using 'isPrecompiled' to only allow the ones defined by the ERC-4337 as stateless precompiles
    // [OP-062]
    var isAllowedPrecompiled = function isAllowedPrecompiled(address) {
      var addrHex = toHex(address);
      var addressInt = parseInt(addrHex);
      // this.debug.push(`isPrecompiled address=${addrHex} addressInt=${addressInt}`)
      return addressInt > 0 && addressInt < 10;
    };
    // [OP-041]
    if (
      opcode.match(/^(EXT.*|CALL|CALLCODE|DELEGATECALL|STATICCALL)$/) != null
    ) {
      var idx = opcode.startsWith("EXT") ? 0 : 1;
      var _addr2 = toAddress(log.stack.peek(idx).toString(16));
      var _addrHex = toHex(_addr2);
      // this.debug.push('op=' + opcode + ' last=' + this.lastOp + ' stacksize=' + log.stack.length() + ' addr=' + addrHex)
      if (
        this.currentLevel.contractSize[_addrHex] == null &&
        !isAllowedPrecompiled(_addr2)
      ) {
        this.currentLevel.contractSize[_addrHex] = {
          contractSize: db.getCode(_addr2).length,
          opcode: opcode,
        };
      }
    }

    // [OP-012]
    if (this.lastOp === "GAS" && !opcode.includes("CALL")) {
      // count "GAS" opcode only if not followed by "CALL"
      this.countSlot(this.currentLevel.opcodes, "GAS");
    }
    if (opcode !== "GAS") {
      // ignore "unimportant" opcodes:
      if (
        opcode.match(
          /^(DUP\d+|PUSH\d+|SWAP\d+|POP|ADD|SUB|MUL|DIV|EQ|LTE?|S?GTE?|SLT|SH[LR]|AND|OR|NOT|ISZERO)$/
        ) == null
      ) {
        this.countSlot(this.currentLevel.opcodes, opcode);
      }
    }
    this.lastOp = opcode;

    if (opcode === "SLOAD" || opcode === "SSTORE") {
      var slot = toWord(log.stack.peek(0).toString(16));
      var slotHex = toHex(slot);
      var _addr3 = log.contract.getAddress();
      var _addrHex2 = toHex(_addr3);
      var access = this.currentLevel.access[_addrHex2];
      if (access == null) {
        access = {
          reads: {},
          writes: {},
        };
        this.currentLevel.access[_addrHex2] = access;
      }
      if (opcode === "SLOAD") {
        // read slot values before this UserOp was created
        // (so saving it if it was written before the first read)
        if (access.reads[slotHex] == null && access.writes[slotHex] == null) {
          access.reads[slotHex] = toHex(db.getState(_addr3, slot));
        }
      } else {
        this.countSlot(access.writes, slotHex);
      }
    }

    if (opcode === "KECCAK256") {
      // collect keccak on 64-byte blocks
      var _ofs2 = parseInt(log.stack.peek(0).toString());
      var _len = parseInt(log.stack.peek(1).toString());
      // currently, solidity uses only 2-word (6-byte) for a key. this might change..
      // still, no need to return too much
      if (_len > 20 && _len < 512) {
        // if (len === 64) {
        this.keccak.push(toHex(log.memory.slice(_ofs2, _ofs2 + _len)));
      }
    } else if (opcode.startsWith("LOG")) {
      var count = parseInt(opcode.substring(3));
      var _ofs3 = parseInt(log.stack.peek(0).toString());
      var _len2 = parseInt(log.stack.peek(1).toString());
      var topics = [];
      for (var _i = 0; _i < count; _i++) {
        // eslint-disable-next-line @typescript-eslint/restrict-plus-operands
        topics.push("0x" + log.stack.peek(2 + _i).toString(16));
      }
      var _data = toHex(log.memory.slice(_ofs3, _ofs3 + _len2));
      this.logs.push({
        topics: topics,
        data: _data,
      });
    }
  },
};
